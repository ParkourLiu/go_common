package cdoris

import (
	"fmt"
	"sync"
	"time"
)

// DorisBuffer 本地批量缓冲器
type DorisBuffer struct {
	mu            sync.Mutex
	dbName        string
	tableName     string
	maxSize       int           // 触发刷新的字节阈值
	flushInterval time.Duration // 超时刷新间隔
	buf           []byte
	lastFlushTS   time.Time
	client        *DorisClient
	useUpdate     bool // true=PutJson(支持更新unique key) false=PutJsonNoUpdate
}

// NewDorisBuffer 创建本地缓冲实例
// maxSize: 缓冲区最大字节数，到达立即刷新
// flushInterval: 最长等待多久强制刷新（避免数据长期滞留）
// useUpdate: 是否开启unique_key更新模式
func NewDorisBuffer(dc *DorisClient, db, table string, maxSize int, flushInterval time.Duration, useUpdate bool) *DorisBuffer {
	b := &DorisBuffer{
		client:        dc,
		dbName:        db,
		tableName:     table,
		maxSize:       maxSize,
		flushInterval: flushInterval,
		useUpdate:     useUpdate,
		lastFlushTS:   time.Now(),
	}
	return b
}

// Append 自动兼容两种输入：
// 1. 单条对象：[]byte(`{"id":1}`)
// 2. 对象数组：[]byte(`[{"id":1},{"id":2}]`)
func (d *DorisBuffer) Append(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	var body []byte
	switch data[0] {
	case '{':
		body = data
	case '[':
		// 数组，剔除首尾 []
		if len(data) < 2 {
			return nil
		}
		body = data[1 : len(data)-1]
		if len(body) == 0 {
			return nil
		}
	default:
		return fmt.Errorf("unsupported json prefix, first char: %c", data[0])
	}

	// 写入缓冲区
	if len(d.buf) == 0 {
		d.buf = append(d.buf, body...)
	} else {
		d.buf = append(d.buf, ',')
		d.buf = append(d.buf, body...)
	}

	return d.checkFlushLocked()
}

// checkFlushLocked 判断是否达到大小/时间阈值，满足则刷新
func (d *DorisBuffer) checkFlushLocked() error {
	if len(d.buf) >= d.maxSize {
		return d.flushLocked()
	}
	if time.Since(d.lastFlushTS) >= d.flushInterval {
		return d.flushLocked()
	}
	return nil
}

// Flush 强制刷新缓冲区所有数据
func (d *DorisBuffer) Flush() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.flushLocked()
}

// flushLocked 内部刷新逻辑（持有锁调用）
func (d *DorisBuffer) flushLocked() error {
	if len(d.buf) == 0 {
		return nil
	}
	// 组装成完整数组 [xxx]
	payload := make([]byte, 0, len(d.buf)+2)
	payload = append(payload, '[')
	payload = append(payload, d.buf...)
	payload = append(payload, ']')

	// 清空缓冲区
	d.buf = d.buf[:0]
	d.lastFlushTS = time.Now()

	var err error
	if d.useUpdate {
		err = d.client.PutJson(d.dbName, d.tableName, payload)
	} else {
		err = d.client.PutJsonNoUpdate(d.dbName, d.tableName, payload)
	}
	if err != nil {
		return fmt.Errorf("doris buffer flush failed: %w, payload_len=%d", err, len(payload))
	}
	return nil
}

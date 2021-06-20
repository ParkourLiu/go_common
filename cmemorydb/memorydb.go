package cmemorydb

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"go_common/clogs"
	"go_common/cmethod"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type MemoryDB struct {
	M   map[string]map[string]bool `json:"m"`
	Ttl map[string]int64           `json:"ttl"`
	s   *sync.RWMutex

	dataPath string //./expiredMap.db
	sigs     chan os.Signal
	exitFlag bool
	log      *clogs.Log
}

func json2struct(strByte []byte) (dat *MemoryDB, err error) {
	if len(strByte) == 0 {
		dat = &MemoryDB{M: map[string]map[string]bool{}, Ttl: map[string]int64{}}
		return
	}
	err = json.Unmarshal(strByte, &dat)
	return
}
func NewMemoryDB(dataPath string, log *clogs.Log) (e *MemoryDB, err error) {
	if log == nil {
		return e, errors.New("clogs.Log is nil")
	}
	e = &MemoryDB{M: map[string]map[string]bool{}, Ttl: map[string]int64{}}
	t1 := time.Now().Unix()
	if cmethod.IsExist(dataPath) { //文件是否存在，存在就读取存档数据
		dataBytes, err := ioutil.ReadFile(dataPath)
		if err != nil {
			return e, err
		}
		e, err = json2struct(dataBytes) //根据文件初始化
		if err != nil {
			return e, errors.New("db file error:" + err.Error())
		}
	}
	t2 := time.Now().Unix()
	log.Info("MemoryDB:", len(e.M), "time:", t2-t1)
	e.s = &sync.RWMutex{}
	e.dataPath = dataPath
	e.log = log
	e.sigs = make(chan os.Signal, 1)
	signal.Notify(e.sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT, syscall.SIGBUS, syscall.SIGFPE, syscall.SIGKILL, syscall.SIGSEGV, syscall.SIGPIPE, syscall.SIGALRM, syscall.SIGTERM) //linux Signal
	go func() {
		e.ttlCheck() //启动检查过期
		for {
			time.Sleep(70 * time.Minute)
			e.ttlCheck()
		}
	}()
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			e.rdb()
		}
	}()
	go func() {
		<-e.sigs
		e.exitFlag = false
		e.rdb()
	exitAgen:
		if !e.exitFlag {
			time.Sleep(time.Second)
			goto exitAgen
		}
		os.Exit(0)
	}()
	return
}

//redis set集合，同一个key  只有一个超时时间  可以存不同value,不设超时即永不超时
func (e *MemoryDB) Sadd(key string, value string, seconds ...int64) {
	e.s.Lock()
	defer e.s.Unlock()
	for _, v := range seconds {
		expiredtime := time.Now().Unix() + v
		e.Ttl[key] = expiredtime
		break
	}

	if _, ok := e.M[key]; !ok {
		e.M[key] = map[string]bool{}
	}
	e.M[key][value] = true

}

func (e *MemoryDB) Del(key string) {
	e.s.Lock()
	defer e.s.Unlock()

	delete(e.M, key)
	delete(e.Ttl, key)
}

func (e *MemoryDB) Smembers(key string) (map[string]bool, bool) {
	e.s.RLock()
	defer e.s.RUnlock()

	if data, ok := e.M[key]; ok { //key存在
		dataFlag := false
		if ttl, ok := e.Ttl[key]; ok { //有过期时间
			if ttl > time.Now().Unix() { //没过期
				dataFlag = true
			} else { //已过期
				delete(e.M, key)
				delete(e.Ttl, key)
			}
		} else { //没过期时间，永不过期
			dataFlag = true
		}

		if dataFlag { //有数据
			datas := map[string]bool{} //深度拷贝
			for k, v := range data {
				datas[k] = v
			}
			return datas, true
		}
	}

	return nil, false
}

func (e *MemoryDB) ttlCheck() {
	e.s.Lock()
	defer e.s.Unlock()
	nowUnix := time.Now().Unix()
	for k, v := range e.Ttl {
		if v <= nowUnix { //过期
			delete(e.M, k)
			delete(e.Ttl, k)
		}
	}
}

func (e *MemoryDB) rdb() { //持久化
	e.s.RLock()
	defer e.s.RUnlock()

	t1 := time.Now().Unix()
	mBuf := bytes.Buffer{}
	tBuf := bytes.Buffer{}
	mBuf.WriteString(`{"m":{`)
	tBuf.WriteString(`},"ttl":{`)
	for k, mm := range e.M {
		mBuf.WriteString(`"`)
		mBuf.WriteString(k)
		mBuf.WriteString(`":{`)
		for k1, _ := range mm {
			mBuf.WriteString(`"`)
			mBuf.WriteString(k1)
			mBuf.WriteString(`":true,`) //存在末尾逗号
		}
		mBuf.WriteString(`},`) //存在末尾逗号

		tBuf.WriteString(`"`)
		tBuf.WriteString(k)
		tBuf.WriteString(`":`)
		tBuf.WriteString(strconv.FormatInt(e.Ttl[k], 10))
		tBuf.WriteString(`,`) //存在末尾逗号
	}
	tBuf.WriteString(`}}`)
	mBuf.Write(tBuf.Bytes())
	go e.rdbfork(bytes.Replace(mBuf.Bytes(), []byte{44, 125}, []byte{125}, -1))

	t2 := time.Now().Unix()
	e.log.Info("forkTime:", t2-t1)

}

func (e *MemoryDB) rdbfork(jsonBytes []byte) {
	defer func() { e.exitFlag = true }() //linux Signal exit flag
	t1 := time.Now().Unix()
	pathTemp := e.dataPath + ".tmp"
	f, err := os.Create(pathTemp)
	if err != nil {
		e.log.Error(err)
		return
	}
	_, err = f.Write(jsonBytes)
	f.Close()
	if err != nil {
		e.log.Error(err)
		return
	}
	os.Remove(e.dataPath)
	os.Rename(pathTemp, e.dataPath)
	t2 := time.Now().Unix()
	e.log.Info("rdbTime:", t2-t1)
}

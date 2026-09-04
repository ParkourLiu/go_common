package clogs

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Log struct {
	log         *logs.BeeLogger
	errLogCache map[string]printInfo
	l           *sync.Mutex
}
type printInfo struct {
	lastTime int64
	count    int
}

//log.Emergency("log3--->Emergency")
//	log.Alert("log3--->Alert")       //1
//	log.Critical("log3--->Critical") //2
//	log.Error("log3--->Error")       //3
//	log.Warn("log3--->Warning")      //4
//	log.Notice("log3--->Notice")     //5
//	log.Info("log3--->Info")         //6
//	log.Debug("log3--->Debug")       //7
func NewLog(logLevel string, Async bool) *Log {
	log := logs.NewLogger()
	if Async {
		log.Async() //设置异步
	}
	level, err := strconv.Atoi(logLevel)
	if err != nil {
		level = 7
	}
	log.SetLevel(level)
	log.EnableFuncCallDepth(true) //设置打印行号
	log.SetLogFuncCallDepth(3)    //设置打印深度

	log.SetLogger(logs.AdapterConsole)
	return &Log{log: log, errLogCache: map[string]printInfo{}, l: &sync.Mutex{}}
}

func NewLogDance(logLevel *int, Async bool) (log *Log) {
	log = &Log{errLogCache: map[string]printInfo{}, l: &sync.Mutex{}}
	log.log = logs.NewLogger()
	if Async {
		log.log.Async() //设置异步
	}
	log.log.SetLevel(*logLevel)       //设置级别
	log.log.EnableFuncCallDepth(true) //设置打印行号
	log.log.SetLogFuncCallDepth(3)    //设置打印深度
	log.log.SetLogger(logs.AdapterConsole)

	go func() {
		for {
			time.Sleep(3 * time.Second)
			if log.log.GetLevel() != *logLevel {
				log.log.SetLevel(*logLevel)
				log.Info("update log level", *logLevel)
			}
		}
	}()
	return
}

//logName   ./logs/test.log
func NewFileLog(logName string, logLevel string, Async bool) *Log {
	log := logs.NewLogger()
	if Async {
		log.Async() //设置异步
	}
	level, err := strconv.Atoi(logLevel)
	if err != nil {
		level = 7
	}
	log.SetLevel(level)
	log.EnableFuncCallDepth(true) //设置打印行号
	log.SetLogFuncCallDepth(3)    //设置打印深度
	//`{"level":7,"filename":"test.log","separate":["error", "warning", "info", "debug"]}`
	log.SetLogger(logs.AdapterMultiFile, `{"filename":"`+logName+`","separate":["error","warning","info","debug"]}`)
	return &Log{log: log, errLogCache: map[string]printInfo{}, l: &sync.Mutex{}}
}

func NewFileLogDance(logName string, logLevel *string, Async bool) *Log {
	log := logs.NewLogger()
	if Async {
		log.Async() //设置异步
	}
	log.EnableFuncCallDepth(true) //设置打印行号
	log.SetLogFuncCallDepth(3)    //设置打印深度
	//`{"level":7,"filename":"test.log","separate":["error", "warning", "info", "debug"]}`
	log.SetLogger(logs.AdapterMultiFile, `{"filename":"`+logName+`","separate":["error","warning","info","debug"]}`)
	go func() {
		for {
			level, err := strconv.Atoi(*logLevel)
			if err != nil {
				level = 7
			}
			log.SetLevel(level)
			time.Sleep(3 * time.Second)
		}
	}()
	return &Log{log: log, errLogCache: map[string]printInfo{}, l: &sync.Mutex{}}
}
func (l *Log) Debug(v ...interface{}) {
	l.log.Debug(strings.Repeat("%+v ", len(v)), v...)
}
func (l *Log) Warn(v ...interface{}) {
	l.log.Warn(strings.Repeat("%+v ", len(v)), v...)
}
func (l *Log) Info(v ...interface{}) {
	l.log.Info(strings.Repeat("%+v ", len(v)), v...)
}
func (l *Log) Error(v ...interface{}) {
	l.log.Error(strings.Repeat("%+v ", len(v)), v...)
}
func (l *Log) ErrorMerge(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	k := fmt.Sprintf("%s%d", file, line)
	now := time.Now().Unix()

	l.l.Lock() //加锁
	p, ok := l.errLogCache[k]
	if !ok { //第一次报错需要打印并记录
		l.errLogCache[k] = printInfo{lastTime: now, count: 1}
		l.l.Unlock() //解锁
		l.log.Error(strings.Repeat("%+v ", len(v)), v...)
		return
	}
	p.count = p.count + 1    //累加错误次数
	if now-p.lastTime < 60 { //距离上次打印后1分钟没到不重复打印
		l.errLogCache[k] = p //记录
		l.l.Unlock()         //解锁
		if p.count <= 3 {    //前3次打印，之后不打印
			l.log.Error(strings.Repeat("%+v ", len(v)), v...)
		}
		return
	}
	p.lastTime = now     //重置最后打印时间
	l.errLogCache[k] = p //记录
	l.l.Unlock()         //解锁

	v = append([]interface{}{fmt.Sprintf("[%d]", p.count)}, v...) //报错次数加入到错误信息头部，从第一次打印到现在的所有错误次数
	l.log.Error(strings.Repeat("%+v ", len(v)), v...)
}

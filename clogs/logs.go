package clogs

import (
	"github.com/astaxie/beego/logs"
	"strings"
)

type Log struct {
	log *logs.BeeLogger
}

//log.Emergency("log3--->Emergency")
//	log.Alert("log3--->Alert")       //1
//	log.Critical("log3--->Critical") //2
//	log.Error("log3--->Error")       //3
//	log.Warn("log3--->Warning")      //4
//	log.Notice("log3--->Notice")     //5
//	log.Info("log3--->Info")         //6
//	log.Debug("log3--->Debug")       //7
func NewLog(logLevel string) *Log {
	log := logs.NewLogger()
	log.EnableFuncCallDepth(true) //设置打印行号
	log.SetLogFuncCallDepth(3)    //设置打印深度
	log.SetLogger(logs.AdapterConsole, `{"level":`+logLevel+`}`)
	return &Log{log: log}
}
func (l *Log) Debug(v ...interface{}) {
	l.log.Debug(strings.Repeat("%v ", len(v)), v...)
}
func (l *Log) Warn(v ...interface{}) {
	l.log.Warn(strings.Repeat("%v ", len(v)), v...)
}

func (l *Log) Info(v ...interface{}) {
	l.log.Info(strings.Repeat("%v ", len(v)), v...)
}
func (l *Log) Error(v ...interface{}) {
	l.log.Error(strings.Repeat("%v ", len(v)), v...)
}

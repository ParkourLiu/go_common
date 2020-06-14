package clogs_test

import (
	"go_common/clogs"
	"testing"
)

var log *clogs.Log

func init() {
	//log.Emergency("log3--->Emergency")
	//	log.Alert("log3--->Alert")       //1
	//	log.Critical("log3--->Critical") //2
	//	log.Error("log3--->Error")       //3
	//	log.Warn("log3--->Warning")      //4
	//	log.Notice("log3--->Notice")     //5
	//	log.Info("log3--->Info")         //6
	//	log.Debug("log3--->Debug")       //7
	log = clogs.NewLog("7")
}

func Test_log(t *testing.T) {
	log.Debug("Debug")
	log.Warn("Warn")
	log.Info("Info")
	log.Error("Error")
}

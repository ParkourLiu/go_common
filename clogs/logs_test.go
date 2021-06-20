package clogs_test

import (
	"go_common/clogs"
	"testing"
	"time"
)

var log *clogs.Log

func init() {
	log = clogs.NewFileLog("aaa.log", "7")
}

func Test_log(t *testing.T) {
	log.Debug("Debug")
	log.Warn("Warn")
	log.Info("Info")
	log.Error("Error")
	time.Sleep(5 * time.Second)
}

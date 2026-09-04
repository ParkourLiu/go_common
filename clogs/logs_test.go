package clogs_test

import (
	"go_common/clogs"
	"testing"
	"time"
)

var log *clogs.Log
var level = "7"

func init() {
	log = clogs.NewLog(level, false)
}

func Test_log(t *testing.T) {
	log.Debug("Debug")
	log.Warn("Warn")
	log.Info("Info")
	log.Error("Error")
	for i := 0; i < 9; i++ {
		log.ErrorMerge("ErrorMerge1")
		if i == 3 || i == 6 {
			time.Sleep(time.Minute + time.Second + time.Second)
		}
	}
}

package clogs_test

import (
	"fmt"
	"go_common/clogs"
	"testing"
	"time"
)

var log *clogs.Log
var level = "7"

func init() {
	log = clogs.NewLog(&level, false)
}

func Test_log(t *testing.T) {
	log.Debug("Debug")
	log.Warn("Warn")
	log.Info("Info")
	log.Error("Error")
	level = "6"
	time.Sleep(5 * time.Second)
	fmt.Println("-------------------------")
	log.Debug("Debug")
	log.Warn("Warn")
	log.Info("Info")
	log.Error("Error")
}

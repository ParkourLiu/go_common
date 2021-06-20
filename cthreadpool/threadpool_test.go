package cthreadpool_test

import (
	"fmt"
	"go_common/cthreadpool"
	"testing"
)

var threadPllClient *cthreadpool.ThreadPllClient

func init() {
	threadPllClient = cthreadpool.NewThreadPllClient(3)
}

func Test_ThreadPllClient_Run(t *testing.T) {
	for i := 0; i < 1000; i++ {
		threadPllClient.Run(func() {
			fmt.Println(111)
		})
	}
}

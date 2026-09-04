package cidutil_test

import (
	"fmt"
	"go_common/cidutil"
	"testing"
)

var (
	idClient *cidutil.IdWorker
)

func init() {
	idClient, _ = cidutil.NewIdWorker(1) //传入nodeID
}

func Test_IdWorker_NextId(t *testing.T) {
	nextId := fmt.Sprint(idClient.NextId())
	fmt.Println(nextId)
}

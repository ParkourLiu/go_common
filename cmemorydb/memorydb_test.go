package cmemorydb_test

import (
	"go_common/clogs"
	"go_common/cmemorydb"
	"os"
	"testing"
	"time"
)

var (
	log      *clogs.Log
	memorydb *cmemorydb.MemoryDB
)

func init() {
	var err error
	log = clogs.NewLog("7", true)
	memorydb, err = cmemorydb.NewMemoryDB("./data.db", log)
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
}
func TestMemoryDb(t *testing.T) {
	memorydb.Sadd("1", "1")
	memorydb.Sadd("2", "2", 5)
	log.Info(memorydb.Smembers("1"))
	log.Info(memorydb.Smembers("2"))
	time.Sleep(6 * time.Second)
	log.Info(memorydb.Smembers("1"))
	log.Info(memorydb.Smembers("2"))

}

func TestMemoryDB_Del(t *testing.T) {
	memorydb.Sadd("a", "a")
	memorydb.Smembers("a")
}

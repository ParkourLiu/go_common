package cmysql_test

import (
	"fmt"
	"go_common/clogs"
	"go_common/cmysql"
	"testing"
	"time"
)

var mysqlClient *cmysql.MysqlClient
var log = clogs.NewLog("7", false)

func init() {
	mysqlClient = cmysql.NewMysqlClient(&cmysql.MysqlInfo{
		UserName:     "root",
		Password:     "root",
		IP:           "127.0.0.1",
		Port:         "3306",
		DatabaseName: "test",
		MaxIdleConns: 1000,
		Log:          log,
		ConnArgs: map[string]string{
			"parseTime": "true",
		},
	})
}

func TestMysqlClient_Count(t *testing.T) {
	i, err := mysqlClient.Count(&cmysql.Stmt{Sql: "SELECT count(*) FROM aaa WHERE id=2", Args: []interface{}{}})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(i)
}
func TestMysqlClient_SearchOneRow(t *testing.T) {
	i, err := mysqlClient.SearchOneRow(&cmysql.Stmt{Sql: "SELECT * FROM aaa WHERE id=2", Args: []interface{}{}})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(i)
}

func TestMysqlClient_SearchMutiRows(t *testing.T) {
	for i := 0; i < 100000; i++ {
		a, err := mysqlClient.SearchMutiRows(&cmysql.Stmt{Sql: "SELECT * FROM aaa", Args: []interface{}{}})
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(a)
	}

}

func TestMysqlClient_Execute(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		err := mysqlClient.Execute(&cmysql.Stmt{Sql: "INSERT INTO aaa(NAME) VALUES(?)", Args: []interface{}{"握草"}})
		if err != nil {
			t.Error(err)
			return
		}
	}
}
func TestMysqlClient_Close(t *testing.T) {
	mysqlClient.Close()
	time.Sleep(10 * time.Second)
}

func TestMysqlClient_ExecuteByTransaction(t *testing.T) {
	sqls := []cmysql.Stmt{
		{Sql: "INSERT INTO aaa(NAME) VALUES(?)", Args: []interface{}{"握草1"}},
		{Sql: "INSERT INTO aaa(NAME) VALUES(?)", Args: []interface{}{"握草2"}},
		{Sql: "INSERT INTO aaa(NAME) VALUES(?)", Args: []interface{}{"握草3"}},
		{Sql: "INSERT INTO aaa(NAME) VALUES(?)", Args: []interface{}{"握草3"}},
	}
	err := mysqlClient.ExecuteByTransaction(sqls)
	if err != nil {
		t.Error(err)
	}
}

func TestMysqlClient_HandTransaction(t *testing.T) {
	tx, err := mysqlClient.GetTransaction()
	if err != nil {
		t.Error(err)
		return
	}
	err = mysqlClient.AddTransactionSql(tx, &cmysql.Stmt{Sql: "INSERT INTO aaa(NAME) VALUES(?)", Args: []interface{}{"握草1"}})
	if err != nil {
		t.Error(err)
		return
	}
	err = mysqlClient.AddTransactionSql(tx, &cmysql.Stmt{Sql: "INSERT INTO aaa(NAME) VALUES(?)", Args: []interface{}{"握草2"}})
	if err != nil {
		t.Error(err)
		return
	}
	err = mysqlClient.AddTransactionSql(tx, &cmysql.Stmt{Sql: "INSERT INTO aaa(NAME) VALUES(?)", Args: []interface{}{"握草3"}})
	if err != nil {
		t.Error(err)
		return
	}
	mysqlClient.Commit(tx)
}

//关闭mysql，然后再打开，查看mysql是否重连
func TestCheckConn(t *testing.T) {
	for {
		i, err := mysqlClient.SearchOneRow(&cmysql.Stmt{Sql: "SELECT * FROM aaa WHERE id=2", Args: []interface{}{}})
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Println(i)
		time.Sleep(5 * time.Second)
	}
}

//关闭mysql，然后再打开，查看mysql是否重连
func TestMysqlClient_SearchFormat(t *testing.T) {
	type User struct {
		Id   int       `json:"iii,omitempty"`
		Name string    `json:"name,omitempty"`
		Ct   time.Time `json:"ct,omitempty"`
	}
	var count int
	err := mysqlClient.SearchFormat(&cmysql.Stmt{Sql: "SELECT count(1) FROM user where id>1", Args: []interface{}{}}, &count)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(count)

	var us []User
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: "SELECT * FROM user where id>1", Args: []interface{}{}}, &us)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(us)

	var u User
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: "SELECT * FROM user where id=1", Args: []interface{}{}}, &u)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(u)

	var ms []map[string]interface{}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: "SELECT * FROM user where id>1", Args: []interface{}{}}, &ms)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(ms)

	var m map[string]string
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: "SELECT * FROM user where id=1", Args: []interface{}{}}, &m)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(m)

	var ss []string
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: "SELECT name FROM user where id>1", Args: []interface{}{}}, &ss)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(ss)
}

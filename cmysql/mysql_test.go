package cmysql_test

import (
	"fmt"
	"go_common/cmysql"
	"testing"
	"time"
)

var mysqlClient *cmysql.MysqlClient

func init() {
	mysqlClient = cmysql.NewMysqlClient(&cmysql.MysqlInfo{
		UserName:     "root",
		Password:     "root",
		IP:           "127.0.0.1",
		Port:         "3306",
		DatabaseName: "aaa",
		MaxIdleConns: 1000,
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

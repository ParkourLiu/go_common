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

//func init() {
//	mysqlClient = cmysql.NewMysqlClient(&cmysql.MysqlInfo{
//		UserName:     "root",
//		Password:     "root",
//		IP:           "127.0.0.1",
//		Port:         "3306",
//		DatabaseName: "test",
//		MaxIdleConns: 1000,
//		Log:          log,
//	})
//}

func TestMysqlClient_Count(t *testing.T) {
	var f = func() {
		var data []map[string]string
		var i int
		//err := mysqlClient.SelectAndCountAut(&cmysql.Stmt{Sql: "SELECT * FROM auth"}, &data, &i)
		i, err := mysqlClient.Count(&cmysql.Stmt{Sql: "SELECT count(*) FROM auth a"}, true)
		fmt.Println(data, i, err)
		time.Sleep(time.Second * 1)
	}
	go f()
	go f()
	f()
	time.Sleep(time.Hour)
}

func TestMysqlClient_SearchOneRow(t *testing.T) {
	i, err := mysqlClient.SearchOneRow(&cmysql.Stmt{Sql: "SELECT * FROM auth WHERE auid=?", Args: []interface{}{"au1"}})
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

func TestMysqlClient_SearchFormats(t *testing.T) {
	type User struct {
		Uid  string
		Name string
	}
	type User1 struct {
		Rids string
	}
	var u, u1 User
	err := mysqlClient.SearchFormats(&[]cmysql.Stmt{{Sql: "SELECT * FROM user where uid='u1';", Args: []interface{}{}}, {Sql: "SELECT * FROM user where uid='u2';", Args: []interface{}{}}}, &u, &u1)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(u, u1)
}

func TestMysqlClient_SelectAndCount(t *testing.T) {
	type User struct {
		Name string `json:"name,omitempty"`
	}
	var u []User
	var count int
	err := mysqlClient.SelectAndCount(&cmysql.Stmt{Sql: "SELECT * FROM user1;"}, &cmysql.Stmt{Sql: "SELECT count(*) FROM user;"}, &u, &count)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(u, count, err)
}

func TestMysqlClient_SelectAndCountAut(t *testing.T) {
	type User struct {
		Name string `json:"name,omitempty"`
	}
	var u []User
	var count int
	//stmt := &cmysql.Stmt{Sql: "SELECT u.* FROM `user` u;"}
	stmt := &cmysql.Stmt{Sql: "SELECT u.* FROM `user` u order by name limit ?,?;", Args: []interface{}{0, 1}}
	//stmt := &cmysql.Stmt{Sql: "SELECT u.*,(select user_id from `test` where `user_id`=?)as t FROM `user` u;", Args: []interface{}{"aaa"}}
	//stmt := &cmysql.Stmt{Sql: "SELECT u.*,(select user_id from `test` where `user_id`=?)as t,(select user_id from `test` where `user_id`=?)as t2 FROM `user` u;", Args: []interface{}{"aaa", "bbbb"}}
	//stmt := &cmysql.Stmt{Sql: "SELECT u.*,(select user_id from `test` where `user_id`=?)as t,(select user_id from `test` where `user_id`=?)as t2 FROM `user` u limit ?,?;", Args: []interface{}{"aaa", "bbbb", 0, 1}}
	//stmt := &cmysql.Stmt{Sql: "SELECT u.*,(select user_id from `test` where `user_id`=?)as t,(select user_id from `test` where `user_id`=?)as t2 FROM `user` u order by u.name desc", Args: []interface{}{"aaa", "bbbb"}}
	//stmt := &cmysql.Stmt{Sql: "SELECT u.*,(select user_id from `test` where `user_id`=?)as t,(select user_id from `test` where `user_id`=?)as t2 FROM `user` u order by u.name limit ?,?;", Args: []interface{}{"aaa", "bbbb", 0, 1}}
	//stmt := &cmysql.Stmt{Sql: "SELECT * FROM `user` WHERE `uid` =(SELECT`user_id`FROM`test`.`test` WHERE user_id=? LIMIT ?,?);", Args: []interface{}{"u1", 0, 1}}
	err := mysqlClient.SelectAndCountAut(stmt, &u, &count)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(u, count, err)
}
func TestMysqlClient_SelectAndCountAut2(t *testing.T) {

	var u []map[string]string
	var count int
	sql := "(SELECT `called`, `province`, `city`, `districts`, `type`, `level`, `sms`, `ct`, `ut` FROM `voip0s`.`voip0s_20220701`  where `level`=?) " +
		"UNION ALL " +
		"(SELECT `called`, `province`, `city`, `districts`, `type`, `level`, `sms`, `ct`, `ut` FROM `voip0s`.`voip0s_20220701` )" +
		"LIMIT ?,?;"
	stmt := &cmysql.Stmt{Sql: sql, Args: []interface{}{1, 0, 3}}
	err := mysqlClient.SelectAndCountAut(stmt, &u, &count)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(u, count, err)
}

func TestMysqlClient_SelectAndCountAut3(t *testing.T) {
	var u []map[string]string
	var count int
	sql := " (SELECT `called`, `province`, `city`, `districts`, `type`, `level`, `sms`, `ct`, `ut` FROM `voip0s`.`voip0s_20220705` WHERE (city IN (?,?) OR FIND_IN_SET(?,districts) OR FIND_IN_SET(?,districts) ))ORDER BY `level` DESC ,ct DESC LIMIT ?,?"
	stmt := &cmysql.Stmt{Sql: sql, Args: []interface{}{"邯郸", "台州区县", "邯郸", "台州区县", 0, 10}}
	err := mysqlClient.SelectAndCountAut(stmt, &u, &count)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(u, count, err)
}

func TestCount2(t *testing.T) {
	sql := "select count(1) from `role` where `rid`=?"
	stmt := &cmysql.Stmt{Sql: sql, Args: []interface{}{"r1"}}
	for {
		c, err := mysqlClient.Count(stmt)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info(c)
		time.Sleep(time.Second)
		break
	}

}

func TestNewProxyMysqlClient(t *testing.T) {
	pm, err := cmysql.NewMysqlClient2(&cmysql.MysqlInfo{
		UserName:     "root",
		Password:     "DHJIsd3nh@#9877",
		IP:           "8.218.242.69",
		Port:         "3306",
		DatabaseName: "meiqiu",
		ProxyInfo: &cmysql.ProxyInfo{
			Addres:       "8.222.244.146:34235",
			AuthUser:     "dante",
			AuthPassword: "mmxe&765**v",
		},
	})
	if err != nil {
		log.Error(err)
		return
	}
	a, err := pm.SearchMutiRows(&cmysql.Stmt{Sql: "show tables;"})
	log.Info(a, err)

}

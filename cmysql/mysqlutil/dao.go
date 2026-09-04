package main

import (
	"bytes"
	"fmt"
	"go_common/cmysql"
)

//查询库里面有多少个表
func s_tables(dbName string) (tableNames []string, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `table_name` FROM `information_schema`.`TABLES` WHERE `TABLE_SCHEMA`=?;")
	sqlArgs := []interface{}{dbName}
	//return mysqlClient.SearchMutiRows(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &tableNames)
	return
}

type Field struct {
	Field string
	Key   string //主键的值为PRI
	Type  string //int(11) varchar(20) datetime
}

//查询表里面有多少个字段
func s_fields(tableName string) (fields []Field, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString(fmt.Sprintf("SHOW COLUMNS FROM `%s`", tableName)) //这里不能当作参数传递
	sqlArgs := []interface{}{}
	//return mysqlClient.SearchMutiRows(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &fields)
	return
}

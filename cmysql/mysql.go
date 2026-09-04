package cmysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-sqlite3"
	"go_common/clogs"
	"golang.org/x/net/proxy"
	"net"
	"strings"
	"sync"
	"time"
)

type MysqlClient struct {
	db        *sql.DB
	mysqlInfo *MysqlInfo
	log       *clogs.Log

	stmtMap          *stmtMap
	l                *sync.Mutex
	StrictParameters bool
}

type stmtMap struct {
	l           *sync.RWMutex
	stmtLiveMap map[string]*stmtLive
}
type stmtLive struct {
	stmt     *sql.Stmt
	lastTime int64
}

func newStmtMap() (s *stmtMap) {
	s = &stmtMap{
		l:           &sync.RWMutex{},
		stmtLiveMap: map[string]*stmtLive{},
	}
	go func() {
		for {
			time.Sleep(time.Hour)
			s.ttl()
		}
	}()
	return
}

func (s *stmtMap) set(sql string, stmt *sql.Stmt) {
	stmtLive := &stmtLive{
		stmt:     stmt,
		lastTime: time.Now().Unix(),
	}
	s.l.Lock()
	defer s.l.Unlock()
	s.stmtLiveMap[sql] = stmtLive
}
func (s *stmtMap) get(sql string) (stmt *sql.Stmt, ok bool) {
	now := time.Now().Unix()
	s.l.RLock()
	stmtLive, ok := s.stmtLiveMap[sql]
	s.l.RUnlock()
	if ok {
		stmt = stmtLive.stmt
		stmtLive.lastTime = now //记录此sql最后一次使用的时间
	}
	return
}
func (s *stmtMap) ttl() {
	now := time.Now().Unix()
	s.l.Lock()
	defer s.l.Unlock()
	for sqlStr, v := range s.stmtLiveMap {
		if now-v.lastTime > 600 { //超过10分钟未使用的，则关闭并删除此预处理
			delete(s.stmtLiveMap, sqlStr)
			if v.stmt != nil {
				_ = v.stmt.Close()
			}
		}
	}
}

type MysqlInfo struct {
	UserName     string
	Password     string
	IP           string
	Port         string
	DatabaseName string
	//MaxOpenConns int //用于设置最大打开的连接数，默认值为0表示不限制。
	MaxIdleConns     int               //用于设置闲置的连接数，默认值为0表示不保留空闲连接。但是在远程连接中，0会因为并发报错
	ConnArgs         map[string]string //连接参数
	Log              *clogs.Log
	ProxyInfo        *ProxyInfo
	StrictParameters bool //严格参数模式
}

// 代理信息的结构体
type ProxyInfo struct {
	Addres       string //"127.0.0.1:2008"
	AuthUser     string //代理用户
	AuthPassword string //代理密码
}
type Stmt struct {
	Sql  string
	Args []interface{}
}

func NewClickHouseClient(mysqlInfo *MysqlInfo) *MysqlClient {
	db, _ := sql.Open("clickhouse", fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&database=%s", mysqlInfo.IP, mysqlInfo.Port, mysqlInfo.UserName, mysqlInfo.Password, mysqlInfo.DatabaseName))
	err := db.Ping()
	if err != nil {
		if mysqlInfo.Log != nil {
			mysqlInfo.Log.Error(err)
		}
		panic(err)
	}
	if mysqlInfo.MaxIdleConns < 30 {
		mysqlInfo.MaxIdleConns = 30
	}
	//db.SetMaxOpenConns(mysqlInfo.MaxOpenConns) //用于设置最大打开的连接数，默认值为0表示不限制。
	db.SetMaxIdleConns(mysqlInfo.MaxIdleConns) //用于设置闲置的连接数，默认值为0表示不保留空闲连接,
	m := &MysqlClient{
		db:        db,
		mysqlInfo: mysqlInfo,
		log:       mysqlInfo.Log,
		stmtMap:   newStmtMap(),
		l:         &sync.Mutex{},
	}
	return m
}
func NewSqliteClient(dbPath string) *MysqlClient {
	db, _ := sql.Open("sqlite3", dbPath)
	m := &MysqlClient{
		db:      db,
		stmtMap: newStmtMap(),
		l:       &sync.Mutex{},
	}
	return m
}
func NewMysqlClient(mysqlInfo *MysqlInfo) *MysqlClient {
	////uri: "root:zaq12wsx1@tcp(localhost:3306)/mm?charset=utf8"
	netType := "tcp"
	if mysqlInfo.ProxyInfo != nil {
		netType = netType + mysqlInfo.IP
		dialer, err := proxy.SOCKS5("tcp", mysqlInfo.ProxyInfo.Addres, &proxy.Auth{User: mysqlInfo.ProxyInfo.AuthUser, Password: mysqlInfo.ProxyInfo.AuthPassword}, proxy.Direct)
		if err != nil {
			panic(err)
		}
		mysql.RegisterDialContext(netType, func(ctx context.Context, addr string) (net.Conn, error) {
			return dialer.Dial("tcp", addr)
		})
	}
	uri := mysqlInfo.UserName + ":" + mysqlInfo.Password + "@" + netType + "(" + mysqlInfo.IP + ":" + mysqlInfo.Port + ")/" + mysqlInfo.DatabaseName + "?charset=utf8mb4&allowOldPasswords=1&allowAllFiles=true" //charset=utf8mb4&     allowOldPasswords=1是为了兼容老版本mysql
	for k, v := range mysqlInfo.ConnArgs {
		uri = fmt.Sprintf("%s&%s=%s", uri, k, v)
	}
	if mysqlInfo.Log != nil {
		mysqlInfo.Log.Info(uri)
	}
	db, _ := sql.Open("mysql", uri)
	err := db.Ping()
	if err != nil {
		if mysqlInfo.Log != nil {
			mysqlInfo.Log.Error(err)
		}
		panic(err)
	}
	if mysqlInfo.MaxIdleConns < 30 {
		mysqlInfo.MaxIdleConns = 30
	}
	//db.SetMaxOpenConns(mysqlInfo.MaxOpenConns) //用于设置最大打开的连接数，默认值为0表示不限制。
	db.SetMaxIdleConns(mysqlInfo.MaxIdleConns) //用于设置闲置的连接数，默认值为0表示不保留空闲连接,
	m := &MysqlClient{
		db:               db,
		mysqlInfo:        mysqlInfo,
		log:              mysqlInfo.Log,
		stmtMap:          newStmtMap(),
		l:                &sync.Mutex{},
		StrictParameters: mysqlInfo.StrictParameters,
	}
	return m
}

func NewMysqlClient2(mysqlInfo *MysqlInfo) (m *MysqlClient, err error) {
	netType := "tcp"
	if mysqlInfo.ProxyInfo != nil {
		netType = netType + mysqlInfo.IP
		dialer, err := proxy.SOCKS5("tcp", mysqlInfo.ProxyInfo.Addres, &proxy.Auth{User: mysqlInfo.ProxyInfo.AuthUser, Password: mysqlInfo.ProxyInfo.AuthPassword}, proxy.Direct)
		if err != nil {
			panic(err)
		}
		mysql.RegisterDialContext(netType, func(ctx context.Context, addr string) (net.Conn, error) {
			return dialer.Dial("tcp", addr)
		})
	}
	////uri: "root:zaq12wsx1@tcp(localhost:3306)/mm?charset=utf8"
	uri := mysqlInfo.UserName + ":" + mysqlInfo.Password + "@" + netType + "(" + mysqlInfo.IP + ":" + mysqlInfo.Port + ")/" + mysqlInfo.DatabaseName + "?charset=utf8mb4&allowOldPasswords=1" //charset=utf8mb4&     allowOldPasswords=1是为了兼容老版本mysql
	for k, v := range mysqlInfo.ConnArgs {
		uri = fmt.Sprintf("%s&%s=%s", uri, k, v)
	}
	db, _ := sql.Open("mysql", uri)
	err = db.Ping()
	if err != nil {
		return
	}
	if mysqlInfo.MaxIdleConns < 30 {
		mysqlInfo.MaxIdleConns = 30
	}
	//db.SetMaxOpenConns(mysqlInfo.MaxOpenConns) //用于设置最大打开的连接数，默认值为0表示不限制。
	db.SetMaxIdleConns(mysqlInfo.MaxIdleConns) //用于设置闲置的连接数，默认值为0表示不保留空闲连接,
	m = &MysqlClient{
		db:        db,
		mysqlInfo: mysqlInfo,
		log:       mysqlInfo.Log,
		stmtMap:   newStmtMap(),
		l:         &sync.Mutex{},
	}
	return
}

func NewMssqlClient(mysqlInfo *MysqlInfo) (m *MysqlClient, err error) {
	////uri: "root:zaq12wsx1@tcp(localhost:3306)/mm?charset=utf8"
	uri := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", mysqlInfo.UserName, mysqlInfo.Password, mysqlInfo.IP, mysqlInfo.Port, mysqlInfo.DatabaseName)
	for k, v := range mysqlInfo.ConnArgs {
		uri = fmt.Sprintf("%s&%s=%s", uri, k, v)
	}
	db, _ := sql.Open("sqlserver", uri)
	err = db.Ping()
	if err != nil {
		return
	}
	if mysqlInfo.MaxIdleConns < 30 {
		mysqlInfo.MaxIdleConns = 30
	}
	//db.SetMaxOpenConns(mysqlInfo.MaxOpenConns) //用于设置最大打开的连接数，默认值为0表示不限制。
	db.SetMaxIdleConns(mysqlInfo.MaxIdleConns) //用于设置闲置的连接数，默认值为0表示不保留空闲连接,
	m = &MysqlClient{
		db:        db,
		mysqlInfo: mysqlInfo,
		log:       mysqlInfo.Log,
		stmtMap:   newStmtMap(),
		l:         &sync.Mutex{},
	}
	return
}

func (c *MysqlClient) getStmtIns(sqlStr string) (stmtIns *sql.Stmt, err error) {
	c.l.Lock()
	defer c.l.Unlock()
	stmtIns, ok := c.stmtMap.get(sqlStr)
	if !ok {
		stmtIns, err = c.db.Prepare(sqlStr)
		if err != nil {
			return stmtIns, err
		}
		c.stmtMap.set(sqlStr, stmtIns)
		if c.log != nil {
			c.log.Info("reg sql:", sqlStr)
		}
	}
	return
}

func (c *MysqlClient) Close() {
	if c.db != nil {
		_ = c.db.Close()
	}
}

func (c *MysqlClient) Count(stmt *Stmt, useStmtCache ...bool) (count int, err error) {
	if len(useStmtCache) > 0 && useStmtCache[0] {
		stmtIns, err := c.getStmtIns(stmt.Sql)
		if err != nil {
			return 0, err
		}
		err = stmtIns.QueryRow(stmt.Args...).Scan(&count)
	} else {
		err = c.db.QueryRow(stmt.Sql, stmt.Args...).Scan(&count)
	}
	return
}
func (c *MysqlClient) ExecuteClickhouse(stmt *Stmt) (err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return
	}
	st, err := tx.Prepare(stmt.Sql)
	if err != nil {
		return
	}
	_, err = st.Exec(stmt.Args...)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

// 批量录入
func (c *MysqlClient) ExecuteClickhouseBatch(stmt *Stmt) (err error) {
	tx, err := c.db.Begin()
	if err != nil {
		return
	}
	st, err := tx.Prepare(stmt.Sql)
	if err != nil {
		return
	}
	argCount := strings.Count(stmt.Sql, "?")
	if len(stmt.Args)%argCount != 0 {
		return fmt.Errorf("arg count Not equal to value")
	}
	var args []interface{}
	for i := 0; i < len(stmt.Args); i++ {
		args = append(args, stmt.Args[i]) //拼接一条的参数value
		if i%argCount == argCount-1 {     //value数量等于参数数量
			_, err = st.Exec(args...)
			if err != nil {
				return
			}
			args = []interface{}{} //置空
		}
	}

	err = tx.Commit()
	return
}

func (c *MysqlClient) Execute(stmt *Stmt, useStmtCache ...bool) (err error) {
	if len(useStmtCache) > 0 && useStmtCache[0] {
		stmtIns, err := c.getStmtIns(stmt.Sql)
		if err != nil {
			return err
		}
		_, err = stmtIns.Exec(stmt.Args...)
	} else {
		_, err = c.db.Exec(stmt.Sql, stmt.Args...)
	}
	return
}

// 返回值：nil代表没有数据。数组元素为空字符串代表null
func (c *MysqlClient) SearchOneRow(stmt *Stmt, useStmtCache ...bool) (m map[string]string, err error) {
	if c.StrictParameters {
		err = formatStmt(stmt)
		if err != nil {
			return
		}
	}
	var rows *sql.Rows
	if len(useStmtCache) > 0 && useStmtCache[0] {
		stmtIns, err := c.getStmtIns(stmt.Sql)
		if err != nil {
			return nil, err
		}
		rows, err = stmtIns.Query(stmt.Args...)
	} else {
		rows, err = c.db.Query(stmt.Sql, stmt.Args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms, err := getRows(rows)
	if err != nil {
		return nil, err
	}
	if len(ms) > 1 {
		return nil, errors.New("Not only one row")
	} else if len(ms) == 1 {
		m = ms[0]
	}
	return
}

// 返回值：长度为0代表没有数据。数组元素为空字符串代表null
func (c *MysqlClient) SearchMutiRows(stmt *Stmt, useStmtCache ...bool) (ms []map[string]string, err error) {
	if c.StrictParameters {
		err = formatStmt(stmt)
		if err != nil {
			return
		}
	}
	var rows *sql.Rows
	if len(useStmtCache) > 0 && useStmtCache[0] {
		stmtIns, err := c.getStmtIns(stmt.Sql)
		if err != nil {
			return nil, err
		}
		rows, err = stmtIns.Query(stmt.Args...)
	} else {
		rows, err = c.db.Query(stmt.Sql, stmt.Args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ms, err = getRows(rows)
	return
}

func (c *MysqlClient) SearchMutiRowsAndCols(stmt *Stmt) (ms []map[string]string, cols []string, err error) {
	if c.StrictParameters {
		err = formatStmt(stmt)
		if err != nil {
			return
		}
	}
	rows, err := c.db.Query(stmt.Sql, stmt.Args...)
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		return
	}
	cols, err = rows.Columns()
	if err != nil {
		return
	}
	ms, err = getRows(rows)
	return
}

func getRows(rows *sql.Rows) ([]map[string]string, error) {
	results := make([]map[string]string, 0) //result
	if rows == nil {
		return nil, errors.New("rows is nil")
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var rawResult [][]byte
	var result map[string]string
	var dest []interface{}
	for rows.Next() {
		result = make(map[string]string, len(cols))
		dest = make([]interface{}, len(cols))
		rawResult = make([][]byte, len(cols))
		for i, _ := range rawResult {
			dest[i] = &rawResult[i]
		}

		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[cols[i]] = ""
			} else {
				result[cols[i]] = string(raw)
			}
		}
		results = append(results, result)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return results, nil
}

// 用事务批量执行sql命令
func (c *MysqlClient) ExecuteByTransaction(stmts []Stmt) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, stmt := range stmts {
		_, err = tx.Exec(stmt.Sql, stmt.Args...)
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

// 手动开启一个事务
func (c *MysqlClient) GetTransaction() (*sql.Tx, error) {
	return c.db.Begin()
}

// 添加要执行的sql
func (c *MysqlClient) AddTransactionSql(tx *sql.Tx, stmt *Stmt) error {
	_, err := tx.Exec(stmt.Sql, stmt.Args...)
	if err != nil {
		defer tx.Rollback()
	}
	return err
}

// 提交事务
func (c *MysqlClient) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

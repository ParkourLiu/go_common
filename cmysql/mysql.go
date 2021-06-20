package cmysql

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"go_common/clogs"
)

type MysqlClient struct {
	db        *sql.DB
	mysqlInfo *MysqlInfo
	//checkBreak bool //是否停止检查连通性的开关,   mysql库内部已实现断开自动重连机制，无需再实现此功能
	log *clogs.Log
}
type MysqlInfo struct {
	UserName     string
	Password     string
	IP           string
	Port         string
	DatabaseName string
	//MaxOpenConns int //用于设置最大打开的连接数，默认值为0表示不限制。
	MaxIdleConns int //用于设置闲置的连接数，默认值为0表示不保留空闲连接。但是在远程连接中，0会因为并发报错
	Log          *clogs.Log
}
type Stmt struct {
	Sql  string
	Args []interface{}
}

func NewMysqlClient(mysqlInfo *MysqlInfo) *MysqlClient {
	////uri: "root:zaq12wsx1@tcp(localhost:3306)/mm?charset=utf8"
	uri := mysqlInfo.UserName + ":" + mysqlInfo.Password + "@tcp(" + mysqlInfo.IP + ":" + mysqlInfo.Port + ")/" + mysqlInfo.DatabaseName + "?charset=utf8mb4&allowOldPasswords=1" //allowOldPasswords=1是为了兼容老版本mysql
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
		db:        db,
		mysqlInfo: mysqlInfo,
	}
	return m
}

func (c *MysqlClient) Close() {
	if c.db != nil {
		_ = c.db.Close()
	}
}

func (c *MysqlClient) Count(stmt *Stmt) (int, error) {
	stmtIns, err := c.db.Prepare(stmt.Sql)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()

	var count int
	err = stmtIns.QueryRow(stmt.Args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *MysqlClient) Execute(stmt *Stmt) error {
	stmtIns, err := c.db.Prepare(stmt.Sql)
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(stmt.Args...)
	if err != nil {
		return err
	}
	return err
}

//返回值：nil代表没有数据。数组元素为空字符串代表null
func (c *MysqlClient) SearchOneRow(stmt *Stmt) (map[string]string, error) {
	stmtIns, err := c.db.Prepare(stmt.Sql)
	if err != nil {
		return nil, err
	}
	defer stmtIns.Close()

	rows, err3 := stmtIns.Query(stmt.Args...)
	if err3 != nil {
		return nil, err3
	}
	defer rows.Close()

	results, err2 := getRows(rows)
	if err2 != nil {
		return nil, err2
	}
	if len(results) > 1 {
		return nil, errors.New("Not only one row")
	} else if len(results) == 1 {
		return results[0], nil
	} else {
		return nil, nil
	}
}

//返回值：长度为0代表没有数据。数组元素为空字符串代表null
func (c *MysqlClient) SearchMutiRows(stmt *Stmt) ([]map[string]string, error) {
	stmtIns, err := c.db.Prepare(stmt.Sql)
	if err != nil {
		return nil, err
	}
	defer stmtIns.Close()

	rows, err3 := stmtIns.Query(stmt.Args...)
	if err3 != nil {
		return nil, err3
	}
	defer rows.Close()

	results, err2 := getRows(rows)
	if err2 != nil {
		return nil, err2
	}
	return results, nil
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
		rawResult = make([][]byte, len(cols))
		result = make(map[string]string, len(cols))
		dest = make([]interface{}, len(cols))
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

//用事务批量执行sql命令
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

//手动开启一个事务
func (c *MysqlClient) GetTransaction() (*sql.Tx, error) {
	return c.db.Begin()
}

//添加要执行的sql
func (c *MysqlClient) AddTransactionSql(tx *sql.Tx, stmt *Stmt) error {
	_, err := tx.Exec(stmt.Sql, stmt.Args...)
	if err != nil {
		defer tx.Rollback()
	}
	return err
}

//提交事务
func (c *MysqlClient) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

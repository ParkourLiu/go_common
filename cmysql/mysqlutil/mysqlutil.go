package main

import (
	"fmt"
	"go_common/cconf"
	"go_common/clogs"
	"go_common/cmysql"
	"math"
	"os"
	"strings"
	"time"
)

var (
	log          = clogs.NewLog("7", false)
	mysqlClient  *cmysql.MysqlClient
	databaseName = cconf.String("mysqlDatabaseName") //mysql数据库名
	sqlFile      *os.File
)

func init() {
	//aiXin()
	mysqlClient = cmysql.NewMysqlClient(&cmysql.MysqlInfo{
		UserName:     cconf.String("mysqlUserName"),
		Password:     cconf.String("mysqlPassword"),
		IP:           cconf.String("mysqlIP"),
		Port:         cconf.String("mysqlPort"),
		DatabaseName: databaseName,
		MaxIdleConns: 30,
		//Log:          log,
	})
	var err error
	sqlFile, err = os.Create("./" + databaseName + ".sql")
	if err != nil {
		log.Error("创建sql文件失败：可能是权限不够", err)
		sleepOut()
	}
}

func main() {
	tableNameList, err := s_tables(databaseName)
	if err != nil {
		log.Error("查询数据库表格错误：", err)
		return
	}
	for _, tableName := range tableNameList {
		fieldList, _ := s_fields(tableName) //获取表格字段
		jointSql(tableName, fieldList)
	}
}

func jointSql(tableName string, fieldList []Field) {
	insertSqlStr := insertSql(tableName, fieldList)
	selectSqlStr := selectSql(tableName, fieldList)
	updateSqlStr := updateSql(tableName, fieldList)
	goStructStr := goStruct(tableName, fieldList)
	sqlFile.WriteString(fmt.Sprintf("#%s\n%s\n%s\n%s\n%s\n\n", tableName, insertSqlStr, selectSqlStr, updateSqlStr, goStructStr))
}

//INSERT INTO `test`.`user`(`id`,`name`,`ct`) VALUES(?,?,?) ON DUPLICATE KEY UPDATE name=?, ct=?;
func insertSql(tableName string, fieldList []Field) (sqlStr string) {
	fields, args, updateArgs := []string{}, []string{}, []string{}
	for _, v := range fieldList {
		fields = append(fields, fmt.Sprintf("`%s`", v.Field))
		args = append(args, "?")
		if v.Key != "PRI" { //不是主键才添加
			updateArgs = append(updateArgs, fmt.Sprintf("`%s`=?", v.Field))
		}
	}
	fieldsStr, argsStr, updateArgsStr := strings.Join(fields, ", "), strings.Join(args, ", "), strings.Join(updateArgs, ", ")
	sqlStr = fmt.Sprintf("INSERT INTO `%s`.`%s`(%s) VALUES(%s) ON DUPLICATE KEY UPDATE %s;", databaseName, tableName, fieldsStr, argsStr, updateArgsStr)
	return
}

//SELECT `id`, `name`, `ct` FROM `test`.`user` WHERE id=? AND name=?
func selectSql(tableName string, fieldList []Field) (sqlStr string) {
	fields, args := []string{}, []string{}
	for _, v := range fieldList {
		fields = append(fields, fmt.Sprintf("`%s`", v.Field))
		if v.Key == "PRI" {
			args = append(args, fmt.Sprintf("`%s`=?", v.Field))
		}
	}
	fieldsStr, argsStr := strings.Join(fields, ", "), strings.Join(args, " AND ")
	sqlStr = fmt.Sprintf("SELECT %s FROM `%s`.`%s` WHERE %s;", fieldsStr, databaseName, tableName, argsStr)
	return
}

//UPDATE `test`.`user` SET `id` = 'id', `name` = 'name', `ct` = 'ct' WHERE `id` = 'id' AND `name` = 'name';
func updateSql(tableName string, fieldList []Field) (sqlStr string) {
	fields, args := []string{}, []string{}
	for _, v := range fieldList {
		fields = append(fields, fmt.Sprintf("`%s`=?", v.Field))
		if v.Key == "PRI" {
			args = append(args, fmt.Sprintf("`%s`=?", v.Field))
		}
	}
	fieldsStr, argsStr := strings.Join(fields, ", "), strings.Join(args, " AND ")
	sqlStr = fmt.Sprintf("UPDATE `%s`.`%s` SET %s WHERE %s;", databaseName, tableName, fieldsStr, argsStr)
	return
}

//type User struct {
//	Id   int       `json:"iii,omitempty"`
//	Name string    `json:"name,omitempty"`
//	Ct   time.Time `json:"ct,omitempty"`
//}
func goStruct(tableName string, fieldList []Field) (structStr string) {
	fields := []string{}
	for _, v := range fieldList {
		fields = append(fields, fmt.Sprintf("	%s %s `json:\"%s,omitempty\"`", initialsToUpper(v.Field), typeFormat(v.Type), v.Field))
	}
	structName := initialsToUpper(tableName)
	fieldsStr := strings.Join(fields, "\n")
	structStr = fmt.Sprintf("type %s struct {\n%s\n}", structName, fieldsStr)

	return
}

// initialsToUpper 首字母转大写，不支持中文
func initialsToUpper(field string) string {
	return strings.ToUpper(field[:1]) + field[1:]
}

func typeFormat(mysqlType string) (goType string) {
	if strings.Contains(mysqlType, "int") {
		return "int"
	} else {
		return "string"
	}
}

func sleepOut() {
	for i := 10; i > 0; i-- {
		fmt.Println(i, "秒后退出。。。")
		time.Sleep(time.Second)
	}
	os.Exit(0)
}

func aiXin() {
	MYWORD := "AnDi love you"
	chars := strings.Split(MYWORD, " ")
	zoom := float64(1)
	tail := "--- 后羿vs丘比特"
	for _, char := range chars {
		allChar := make([]string, 0)

		for y := 12 * zoom; y > -12*zoom; y-- {
			lst := make([]string, 0)
			lstCon := ""
			for x := -30 * zoom; x < 30*zoom; x++ {
				x2 := float64(x)
				y2 := float64(y)
				formula := math.Pow(math.Pow(x2*0.04/zoom, 2)+math.Pow(y2*0.1/zoom, 2)-1, 3) - math.Pow(x2*0.04/zoom, 2)*math.Pow(y2*0.1/zoom, 3)
				if formula <= 0 {
					index := int(x) % len(char)
					if index >= 0 {
						lstCon += string(char[index])
					} else {
						lstCon += string(char[int(float64(len(char))-math.Abs(float64(index)))])
					}

				} else {
					lstCon += " "
				}
			}
			lst = append(lst, lstCon)
			allChar = append(allChar, lst...)
		}

		for _, text := range allChar {
			fmt.Printf("%s\n", text)
			time.Sleep(20 * time.Millisecond)
		}
	}
	fmt.Println("\t\t\t\t", tail)
}

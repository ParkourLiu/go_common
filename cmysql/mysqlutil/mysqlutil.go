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
		fmt.Println(tableName)
		fieldList, _ := s_fields(tableName) //获取表格字段
		jointSql(tableName, fieldList)
	}
}

func jointSql(tableName string, fieldList []Field) {
	insertSqlStr := insertSql(tableName, fieldList)
	selectSqlStr := selectSql(tableName, fieldList)
	updateSqlStr := updateSql(tableName, fieldList)
	goStructStr := goStruct(tableName, fieldList)
	insertsSqlStr := insertsSql(tableName, fieldList)
	SFuncStr := dynamicSelect(tableName, fieldList)
	IuFuncStr := dynamicDUPLICATE(tableName, fieldList)
	sqlFile.WriteString(fmt.Sprintf("#%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n\n", tableName, insertSqlStr, selectSqlStr, updateSqlStr, goStructStr, insertsSqlStr, SFuncStr, IuFuncStr))
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

func insertsSql(tableName string, fieldList []Field) (funcStr string) {
	initials := getInitialsToLower(tableName)
	structName := initialsToUpper(tableName)
	fields, args, valueStrs := []string{}, []string{}, []string{}
	for _, v := range fieldList {
		fields = append(fields, fmt.Sprintf("`%s`", v.Field))
		args = append(args, "?")
		valueStrs = append(valueStrs, fmt.Sprintf("%s.%s", initials, initialsToUpper(v.Field)))
	}
	funcStr = fmt.Sprintf("func I_%ss(%ss []%s) (err error) {\n", tableName, initials, structName) +
		"	sqlBuf := bytes.Buffer{}\n" +
		fmt.Sprintf("	sqlBuf.WriteString(\"INSERT IGNORE INTO `%s`.`%s`(%s) VALUES\")\n", databaseName, tableName, strings.Join(fields, ", ")) +
		"	sqlArgs := []interface{}{}\n" +
		fmt.Sprintf("	for i, %s := range %ss {\n", initials, initials) +
		fmt.Sprintf("		if i < len(%ss)-1 {\n", initials) +
		fmt.Sprintf("			sqlBuf.WriteString(\"(%s),\")\n", strings.Join(args, ", ")) +
		"		} else {\n" +
		fmt.Sprintf("			sqlBuf.WriteString(\"(%s);\")\n", strings.Join(args, ", ")) +
		"		}\n" +
		fmt.Sprintf("		sqlArgs = append(sqlArgs, %s)\n", strings.Join(valueStrs, ", ")) +
		"	}\n" +
		"	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})\n" +
		"}"
	return
}

func dynamicSelect(tableName string, fieldList []Field) (funcStr string) {
	initials := getInitialsToLower(tableName)
	structName := initialsToUpper(tableName)
	fields, dynamicStrs := []string{}, []string{}
	for _, v := range fieldList {
		fields = append(fields, fmt.Sprintf("`%s`", v.Field))
		valueInitStr := `""`
		if typeFormat(v.Type) == "int" {
			valueInitStr = `0`
		}
		dynamicStrs = append(dynamicStrs, fmt.Sprintf("	if %s.%s != %s {\n		sqlBuf.WriteString(\"AND `%s`=? \")\n		sqlArgs = append(sqlArgs, %s.%s)\n	}\n", initials, initialsToUpper(v.Field), valueInitStr, v.Field, initials, initialsToUpper(v.Field)))
	}
	funcStr = fmt.Sprintf("func S_%s(%s *%s) (%ss []%s, err error) {\n", tableName, initials, structName, initials, structName) +
		"	sqlBuf := bytes.Buffer{}\n" +
		fmt.Sprintf("	sqlBuf.WriteString(\"SELECT %s FROM `%s`.`%s` WHERE 1=1  \")\n", strings.Join(fields, ", "), databaseName, tableName) +
		"	sqlArgs := []interface{}{}\n" +
		strings.Join(dynamicStrs, "") +
		fmt.Sprintf("	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &%ss)\n", initials) +
		"	return\n" +
		"}"
	return
}

func dynamicDUPLICATE(tableName string, fieldList []Field) (funcStr string) {
	initials := getInitialsToLower(tableName)
	structName := initialsToUpper(tableName)
	fields, args, argsValues, dynamicStrs := []string{}, []string{}, []string{}, []string{}
	for _, v := range fieldList {
		fields = append(fields, fmt.Sprintf("`%s`", v.Field))
		args = append(args, "?")
		argsValues = append(argsValues, fmt.Sprintf("%s.%s", initials, initialsToUpper(v.Field)))
		valueInitStr := `""`
		if typeFormat(v.Type) == "int" {
			valueInitStr = `0`
		}
		dynamicStrs = append(dynamicStrs, fmt.Sprintf("	if %s.%s != %s {\n		sqlBuf.WriteString(\"`%s`=?,\")\n		sqlArgs = append(sqlArgs, %s.%s)\n	}\n", initials, initialsToUpper(v.Field), valueInitStr, v.Field, initials, initialsToUpper(v.Field)))
	}
	fieldsStr, argsStr := strings.Join(fields, ", "), strings.Join(args, ", ")
	funcStr = fmt.Sprintf("func IU_%s(%s *%s) (err error) {\n", tableName, initials, structName) +
		"	sqlBuf := bytes.Buffer{}\n" +
		fmt.Sprintf("	sqlBuf.WriteString(\"INSERT INTO `%s`.`%s`(%s) VALUES(%s) ON DUPLICATE KEY UPDATE \")\n", databaseName, tableName, fieldsStr, argsStr) +
		fmt.Sprintf("	sqlArgs := []interface{}{%s}\n", strings.Join(argsValues, ", ")) +
		fmt.Sprintf("%s\n请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString(\"`ut`=NOW();\")\n", strings.Join(dynamicStrs, "")) +
		"	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})\n" +
		"}"
	return
}

// initialsToUpper 首字母转大写，不支持中文
func initialsToUpper(field string) string {
	return strings.ToUpper(field[:1]) + field[1:]
}
func getInitialsToLower(field string) string {
	if field == "" {
		return ""
	}
	return strings.ToLower(field[:1])
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

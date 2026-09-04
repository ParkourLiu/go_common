package cmysql

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

func (c *MysqlClient) SelectAndCountAut(s *Stmt, T interface{}, count *int, useStmtCache ...bool) (err error) {
	err = formatStmt(s)
	if err != nil {
		return
	}
	s.Sql = deleteExtraSpace(s.Sql) //去除多余空格
	wait := &sync.WaitGroup{}
	var err1 error
	wait.Add(1)
	go func() {
		defer wait.Done()
		err1 = c.SearchFormat(s, T, useStmtCache...)
	}()

	s2, err := select2Count(*s)
	if err != nil {
		return
	}
	*count, err = c.Count(s2, useStmtCache...)
	if err != nil {
		return err
	}
	wait.Wait()
	if err1 != nil {
		return err1
	}
	return
}

//手动处理sql数量
func (c *MysqlClient) SelectAndCount(selectStmt, countStmt *Stmt, T interface{}, count *int, useStmtCache ...bool) (err error) {
	err = formatStmt(countStmt)
	if err != nil {
		return
	}
	wait := &sync.WaitGroup{}
	var err1 error
	wait.Add(1)
	go func() {
		defer wait.Done()
		err1 = c.SearchFormat(selectStmt, T, useStmtCache...)
	}()
	*count, err = c.Count(countStmt, useStmtCache...)
	if err != nil {
		return err
	}
	wait.Wait()
	if err1 != nil {
		return err1
	}
	return
}

func findFromStr(sql string) (offsize, endOffsize int, err error) {
	sql = strings.Replace(sql, " from", " FROM", -1)
	sql = strings.Replace(sql, "select ", "SELECT ", -1)
	sql = strings.Replace(sql, " order by ", " ORDER BY ", -1)
	sql = strings.Replace(sql, " limit ", " LIMIT ", -1)
	if !strings.HasPrefix(sql, "SELECT ") {
		return offsize, endOffsize, fmt.Errorf("This is an unrecognized select statement")
	}

	for offsize = len("SELECT "); ; {
		sIndex := strings.Index(sql[offsize:], "SELECT ")
		fIndex := strings.Index(sql[offsize:], " FROM")
		offsize = offsize + fIndex
		if sIndex > -1 && sIndex < fIndex {
			offsize = offsize + len(" FROM")
			continue
		} else {
			break
		}
	}
	endOffsize = len(sql)
	kIndex := strings.LastIndex(sql, ")")
	lIndex := strings.LastIndex(sql, " LIMIT ")
	if lIndex > kIndex {
		endOffsize = lIndex
	}
	oIndex := strings.LastIndex(sql, " ORDER BY ")
	if oIndex > kIndex {
		endOffsize = oIndex
	}
	return
}

func strCount(s, str string) (count int) {
	for _, v := range s {
		if string(v) == str {
			count++
		}
	}
	return
}

//删除字符串中的多余空格和换行，有多个空格和换行时，仅保留一个空格
func deleteExtraSpace(s string) string {
	s = strings.TrimSpace(s)
	reg, _ := regexp.Compile(`\s{2,}`)
	tmpstr := []byte(s)
	for {
		spc_index := reg.FindStringIndex(string(tmpstr))
		if spc_index == nil {
			break
		}
		tmpstr = append(tmpstr[:spc_index[0]+1], tmpstr[spc_index[1]:]...)
	}
	return string(tmpstr)
}

func select2Count(s Stmt) (*Stmt, error) {
	var err error
	if strings.Contains(s.Sql, "union all") && strings.HasPrefix(s.Sql, "(") {
		err = unionSelect2Count(&s, "union all")
	} else if strings.Contains(s.Sql, "UNION ALL") && strings.HasPrefix(s.Sql, "(") {
		err = unionSelect2Count(&s, "UNION ALL")
	} else if strings.Contains(s.Sql, "union") && strings.HasPrefix(s.Sql, "(") {
		err = unionSelect2Count(&s, "union")
	} else if strings.Contains(s.Sql, "UNION") && strings.HasPrefix(s.Sql, "(") {
		err = unionSelect2Count(&s, "UNION")
	} else {
		err = oneSelect2Count(&s)
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func unionSelect2Count(s *Stmt, uType string) (err error) {
	lastIndex := strings.LastIndex(s.Sql, ")")
	s.Sql = s.Sql[:lastIndex+1] //去除排序和limit
	sqlSplit := strings.Split(s.Sql, uType)
	newSqlBuf := bytes.Buffer{}
	newSqlBuf.WriteString("SELECT SUM(c) FROM(")
	newAgrs := []interface{}{}
	for i, oneSelectSql := range sqlSplit {
		oneSelectSql = cleanBracket(oneSelectSql) //去除括号
		offsize, endOffsize, err := findFromStr(oneSelectSql)
		if err != nil {
			return err
		}
		sc := strCount(oneSelectSql, "?")                  //单条sql的总参数数量
		tempAgrs := s.Args[len(newAgrs) : len(newAgrs)+sc] //获取其参数
		sc = strCount(oneSelectSql[:offsize], "?")
		tempAgrs = tempAgrs[sc:]
		sc = strCount(oneSelectSql[endOffsize:], "?")
		tempAgrs = tempAgrs[:len(tempAgrs)-sc]
		newSqlBuf.WriteString(fmt.Sprintf("(SELECT COUNT(*) c %s)", oneSelectSql[offsize:endOffsize]))
		if i < len(sqlSplit)-1 {
			newSqlBuf.WriteString(uType)
		}
		newAgrs = append(newAgrs, tempAgrs...)
	}
	newSqlBuf.WriteString(") a")
	fmt.Println(newSqlBuf.String(), newAgrs)
	s.Sql = newSqlBuf.String()
	s.Args = newAgrs
	return
}

func oneSelect2Count(s *Stmt) (err error) {
	s.Sql = cleanBracket(s.Sql)
	offsize, endOffsize, err := findFromStr(s.Sql)
	if err != nil {
		return
	}
	sc := strCount(s.Sql[:offsize], "?")
	s.Args = s.Args[sc:]
	sc = strCount(s.Sql[endOffsize:], "?")
	s.Args = s.Args[:len(s.Args)-sc]
	s.Sql = "SELECT COUNT(*) " + s.Sql[offsize:endOffsize]
	return
}

func cleanBracket(sql string) string {
	sql = strings.TrimSpace(sql)
	if !strings.HasPrefix(sql, "(") {
		return sql
	}
	bFlag := 1
	for i, v := range sql[1:] {
		if v == 40 { //(
			bFlag++
		} else if v == 41 { //)
			bFlag--
			if bFlag == 0 {
				sql = sql[1:i+1] + sql[i+2:]
				break
			}
		}
	}
	sql = strings.TrimSpace(sql)
	return sql
}

package cmysql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (c *MysqlClient) SearchFormat(stmt *Stmt, T interface{}) (err error) {
	stmtIns, err := c.db.Prepare(stmt.Sql)
	if err != nil {
		return
	}
	defer stmtIns.Close()
	rows, err := stmtIns.Query(stmt.Args...)
	if err != nil {
		return
	}
	defer rows.Close()
	err = scanAll(rows, T)
	if err != nil {
		return
	}
	return
}

func scanAll(rows *sql.Rows, T interface{}) (err error) {
	cols, err := rows.Columns() //列
	if err != nil {
		return err
	}
	colsLen := len(cols)
	if colsLen == 0 { //没数据直接返回
		return
	}
	mysqlFieldTemp := make([]string, colsLen) //mysql临时字段，全小写
	for i, col := range cols {
		mysqlFieldTemp[i] = strings.ToLower(col)
	}
	v := reflect.ValueOf(T) //反射元素内容
	if v.Kind() != reflect.Ptr {
		//fmt.Println("不是指针类型，没法进行修改操作")
		return fmt.Errorf("%v It is not a pointer type and cannot be modified", v.Kind())
	}
	if v.IsNil() {
		//fmt.Println("空反射")
		return fmt.Errorf("%v is Null pointer", v.Kind())
	}
	t := v.Type()
	tElem := t.Elem() // 获取指针所指向的元素
	switch tElem.Kind() {
	case reflect.Slice:
		err = scanSlice(v, tElem, rows, colsLen, mysqlFieldTemp, cols)
	case reflect.Struct:
		err = scanStruct(v, tElem, rows, colsLen, mysqlFieldTemp)
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		if colsLen > 1 { //查询结果不止一个字段
			return fmt.Errorf("There are multiple fields,Please receive in Struct or Map")
		}
		err = scanBaseType(v, tElem, rows)
	case reflect.Map:
		err = scanMap(v, tElem, rows, cols, colsLen)
	}
	if err != nil {
		return
	}
	return rows.Err()
}
func scanBaseType(v reflect.Value, tElem reflect.Type, rows *sql.Rows) (err error) {
	vElem := v.Elem()
	if rows.Next() {
		var dest interface{}
		err = rows.Scan(&dest)
		if err != nil {
			return err
		}
		err = colFormatField(tElem.Name(), dest, vElem, tElem.Kind()) //格式化每个字段
	}
	if rows.Next() {
		//第二次判断，如果有多条提醒用户改用数组形式接收
		return fmt.Errorf("There are multiple results,Please receive in array")
	}
	return
}
func scanMap(v reflect.Value, tElem reflect.Type, rows *sql.Rows, cols []string, colsLen int) (err error) {
	vElem := v.Elem()
	switch tElem.String() {
	case "map[string]string":
		if rows.Next() {
			rawResult := make([]interface{}, colsLen)
			result := make(map[string]string, colsLen)
			dest := make([]interface{}, colsLen)
			for i, _ := range rawResult {
				dest[i] = &rawResult[i]
			}
			err = rows.Scan(dest...)
			if err != nil {
				return
			}
			for i, raw := range rawResult {
				switch raw.(type) { //获取数据库中类型
				case int64:
					result[cols[i]] = strconv.FormatInt(raw.(int64), 10)
				case []uint8:
					result[cols[i]] = string(raw.([]byte))
				case time.Time:
					result[cols[i]] = raw.(time.Time).Format("2006-01-02 15:04:05")
				case nil:
					result[cols[i]] = ""
				}
			}
			vElem.Set(reflect.ValueOf(result))
		}
	case "map[string]interface {}":
		if rows.Next() {
			rawResult := make([]interface{}, colsLen)
			result := make(map[string]interface{}, colsLen)
			dest := make([]interface{}, colsLen)
			for i, _ := range rawResult {
				dest[i] = &rawResult[i]
			}
			err = rows.Scan(dest...)
			if err != nil {
				return
			}
			for i, raw := range rawResult {
				switch raw.(type) { //获取数据库中类型
				case []uint8:
					result[cols[i]] = string(raw.([]byte))
				default:
					result[cols[i]] = raw
				}
			}
			vElem.Set(reflect.ValueOf(result))
		}
	default:
		return fmt.Errorf("Only use map[string]string or map[string]interface{}")
	}
	if rows.Next() {
		//第二次判断，如果有多条提醒用户改用数组形式接收
		return fmt.Errorf("There are multiple results,Please receive in array")
	}
	return
}
func scanStruct(v reflect.Value, tElem reflect.Type, rows *sql.Rows, colsLen int, mysqlFieldTemp []string) (err error) {
	structFieldMap := map[string]string{}
	for i := 0; i < tElem.NumField(); i++ { //遍历所有字段
		key := tElem.Field(i)  // 从0开始获取所包含的字段
		if !key.IsExported() { //不是导出字段（字段首字母没有大写）就不管
			continue
		}
		structFieldMap[strings.ToLower(key.Name)] = key.Name //全部小写，方便比对字段
	}
	vElem := v.Elem()
	var dest []interface{}
	var destTemp []interface{}
	if rows.Next() {
		dest = make([]interface{}, colsLen)
		destTemp = make([]interface{}, colsLen)
		for i, _ := range destTemp {
			dest[i] = &destTemp[i]
		}
		err = rows.Scan(dest...)
		if err != nil {
			return err
		}
		for i, value := range destTemp { //遍历一行数据中的所有列
			fieldName, ok := structFieldMap[mysqlFieldTemp[i]] //获取此列名，是否在结构体中也存在此字段
			//fmt.Printf("%d:%T,%v,%s\n", i, value, value, fieldName)
			if !ok { //查出来的字段结构体中不存在
				continue
			}
			field := vElem.FieldByName(fieldName)                    //获取字段本身
			fieldType := field.Kind()                                //获取字段反射后的类型
			err = colFormatField(fieldName, value, field, fieldType) //格式化每个字段
			if err != nil {
				return
			}
		}
		//fmt.Printf("333 %+v\n", vp)
	}
	if rows.Next() {
		//第二次判断，如果有多条提醒用户改用数组形式接收
		return fmt.Errorf("There are multiple results,Please receive in array")
	}
	return
}

func colFormatField(fieldName string, value interface{}, field reflect.Value, fieldType reflect.Kind) (err error) {
	switch value.(type) { //获取数据库中类型
	case int64: //数据库中是int64
		switch fieldType {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(value.(int64))
		case reflect.String:
			field.SetString(strconv.FormatInt(value.(int64), 10))
		default:
			return fmt.Errorf("converting  %s (%q) to  %s", fieldName, value, field.Kind())
		}
	case []uint8: //数据库中是[]uint8
		valueStr := string(value.([]byte))
		switch fieldType {
		case reflect.String:
			field.SetString(valueStr)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f64, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return fmt.Errorf("converting  %s (%q) to %s: %v", fieldName, valueStr, field.Kind(), err)
			}
			field.SetInt(f64)
		case reflect.Float32, reflect.Float64:
			f64, err := strconv.ParseFloat(valueStr, field.Type().Bits())
			if err != nil {
				return fmt.Errorf("converting  %s (%q) to %s: %v", fieldName, valueStr, field.Kind(), err)
			}
			field.SetFloat(f64)
		default:
			if "time.Time" == field.Type().String() { //时间格式
				//如果你想把mysql中的datetime转换成Golang中"time.Time"类型，请在打开连接后拼接参数：parseTime=true and loc=Local
				return fmt.Errorf("If you want to convert  datetime to time.Time, please splice the parameter after the open connection: parseTime = true and loc=Local")
			} else {
				return fmt.Errorf("converting  %s (%q) to %s", fieldName, valueStr, field.Kind())
			}
		}
	case time.Time: //数据库中是time.Time
		if fieldType == reflect.String {
			valueTime := value.(time.Time)
			field.SetString(valueTime.Format("2006-01-02 15:04:05"))
		} else if fieldType == reflect.Int64 {
			valueTime := value.(time.Time)
			field.SetInt(valueTime.Unix()) //暂时格式化成秒级时间戳
		} else if "time.Time" == field.Type().String() { //时间格式
			field.Set(reflect.ValueOf(value))
		} else {
			return fmt.Errorf("converting  %s (%q) to %s", fieldName, value, field.Kind())
		}
	case nil:
		//数据库中是null,这里不用做任何操作
	default:
		return fmt.Errorf("Unknown MySQL type converting  %s (%q) to %s", fieldName, value, field.Kind())
	}
	return
}

func scanSlice(v reflect.Value, tElem reflect.Type, rows *sql.Rows, colsLen int, mysqlFieldTemp, cols []string) (err error) { //对切片的扫描
	ttElem := tElem.Elem()                //获取切片内的元素反射 main.User
	isPtr := ttElem.Kind() == reflect.Ptr //判断切片内的元素是否是指针类型
	if isPtr {
		ttElem = ttElem.Elem() //如果是指针就获取元素本身main.User
	}

	switch ttElem.Kind() { //判断元素的类型
	case reflect.Struct:
		err = scanSliceStruct(v, ttElem, rows, colsLen, mysqlFieldTemp, isPtr)
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		if colsLen > 1 { //查询结果不止一个字段
			return fmt.Errorf("There are multiple fields,Please receive in []Struct or []Map")
		}
		err = scanSliceBaseType(v, ttElem, rows, isPtr)
	case reflect.Map:
		err = scanSliceMap(v, ttElem, rows, cols, colsLen, isPtr)
	default:
		return fmt.Errorf("%v is Not a Can be formatted type", ttElem.Kind())
	}
	if err != nil {
		return
	}
	return rows.Err()
}
func scanSliceMap(v reflect.Value, ttElem reflect.Type, rows *sql.Rows, cols []string, colsLen int, isPtr bool) (err error) {
	vElem := v.Elem()
	switch ttElem.String() {
	case "map[string]string":
		for rows.Next() {
			rawResult := make([]interface{}, colsLen)
			result := make(map[string]string, colsLen)
			dest := make([]interface{}, colsLen)
			for i, _ := range rawResult {
				dest[i] = &rawResult[i]
			}
			err = rows.Scan(dest...)
			if err != nil {
				return
			}
			for i, raw := range rawResult {
				switch raw.(type) { //获取数据库中类型
				case int64:
					result[cols[i]] = strconv.FormatInt(raw.(int64), 10)
				case []uint8:
					result[cols[i]] = string(raw.([]byte))
				case time.Time:
					result[cols[i]] = raw.(time.Time).Format("2006-01-02 15:04:05")
				case nil:
					result[cols[i]] = ""
				}
			}
			vp := reflect.New(ttElem) ///New返回一个值，该值表示指向指定类型的新零值的指针。也就是说，返回值的类型是PointerTo（typ）
			vpElem := vp.Elem()       // 获取指针所指向的元素
			vpElem.Set(reflect.ValueOf(result))
			if isPtr {
				vElem.Set(reflect.Append(vElem, vp))
			} else {
				vElem.Set(reflect.Append(vElem, vp.Elem()))
			}
		}
	case "map[string]interface {}":
		for rows.Next() {
			rawResult := make([]interface{}, colsLen)
			result := make(map[string]interface{}, colsLen)
			dest := make([]interface{}, colsLen)
			for i, _ := range rawResult {
				dest[i] = &rawResult[i]
			}
			err = rows.Scan(dest...)
			if err != nil {
				return
			}
			for i, raw := range rawResult {
				switch raw.(type) { //获取数据库中类型
				case []uint8:
					result[cols[i]] = string(raw.([]byte))
				default:
					result[cols[i]] = raw
				}
			}
			vp := reflect.New(ttElem) ///New返回一个值，该值表示指向指定类型的新零值的指针。也就是说，返回值的类型是PointerTo（typ）
			vpElem := vp.Elem()       // 获取指针所指向的元素
			vpElem.Set(reflect.ValueOf(result))
			if isPtr {
				vElem.Set(reflect.Append(vElem, vp))
			} else {
				vElem.Set(reflect.Append(vElem, vp.Elem()))
			}
		}
	default:
		return fmt.Errorf("Only use []map[string]string or []map[string]interface{}")
	}
	return
}

func scanSliceBaseType(v reflect.Value, ttElem reflect.Type, rows *sql.Rows, isPtr bool) (err error) {
	vElem := v.Elem()
	for rows.Next() {
		var dest interface{}
		err = rows.Scan(&dest)
		if err != nil {
			return err
		}
		vp := reflect.New(ttElem)                                        ///New返回一个值，该值表示指向指定类型的新零值的指针。也就是说，返回值的类型是PointerTo（typ）
		vpElem := vp.Elem()                                              // 获取指针所指向的元素
		err = colFormatField(ttElem.Name(), dest, vpElem, ttElem.Kind()) //格式化每个字段
		if isPtr {
			vElem.Set(reflect.Append(vElem, vp))
		} else {
			vElem.Set(reflect.Append(vElem, vp.Elem()))
		}
	}
	return
}
func scanSliceStruct(v reflect.Value, ttElem reflect.Type, rows *sql.Rows, colsLen int, mysqlFieldTemp []string, isPtr bool) (err error) {
	structFieldMap := map[string]string{}
	for i := 0; i < ttElem.NumField(); i++ { //遍历所有字段
		key := ttElem.Field(i) // 从0开始获取所包含的字段
		if !key.IsExported() { //不是导出字段（字段首字母没有大写）就不管
			continue
		}
		structFieldMap[strings.ToLower(key.Name)] = key.Name //全部小写，方便比对字段
	}
	//fmt.Println("字段", structFieldMap)
	vElem := v.Elem() // 获取指针指向的元素本身，下面的内容append到此内存中
	var dest []interface{}
	var destTemp []interface{}
	for rows.Next() {
		dest = make([]interface{}, colsLen)
		destTemp = make([]interface{}, colsLen)
		for i, _ := range destTemp {
			dest[i] = &destTemp[i]
		}
		err = rows.Scan(dest...)
		if err != nil {
			return err
		}

		vp := reflect.New(ttElem) ///New返回一个值，该值表示指向指定类型的新零值的指针。也就是说，返回值的类型是PointerTo（typ）
		vpElem := vp.Elem()       // 获取指针所指向的元素

		for i, value := range destTemp { //遍历一行数据中的所有列
			fieldName, ok := structFieldMap[mysqlFieldTemp[i]] //获取此列名，是否在结构体中也存在此字段
			//fmt.Printf("%d:%T,%v,%s\n", i, value, value, fieldName)
			if !ok { //查出来的字段结构体中不存在
				continue
			}
			field := vpElem.FieldByName(fieldName)                   //获取字段本身
			fieldType := field.Kind()                                //获取字段反射后的类型
			err = colFormatField(fieldName, value, field, fieldType) //格式化每个字段
			if err != nil {
				return
			}
		}
		//fmt.Printf("333 %+v\n", vp)
		if isPtr {
			vElem.Set(reflect.Append(vElem, vp))
		} else {
			vElem.Set(reflect.Append(vElem, vp.Elem()))
		}
	}
	return
}

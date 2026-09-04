#login_log
INSERT INTO `lt`.`login_log`(`uid`, `login_id`, `uname`, `nickname`, `ip`, `expire`, `ct`, `ut`) VALUES(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `uid`=?, `uname`=?, `nickname`=?, `ip`=?, `expire`=?, `ct`=?, `ut`=?;
SELECT `uid`, `login_id`, `uname`, `nickname`, `ip`, `expire`, `ct`, `ut` FROM `lt`.`login_log` WHERE `login_id`=?;
UPDATE `lt`.`login_log` SET `uid`=?, `login_id`=?, `uname`=?, `nickname`=?, `ip`=?, `expire`=?, `ct`=?, `ut`=? WHERE `login_id`=?;
type Login_log struct {
	Uid int `json:"uid,omitempty"`
	Login_id string `json:"login_id,omitempty"`
	Uname string `json:"uname,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Ip string `json:"ip,omitempty"`
	Expire string `json:"expire,omitempty"`
	Ct string `json:"ct,omitempty"`
	Ut string `json:"ut,omitempty"`
}
func I_login_logs(ls []Login_log) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `lt`.`login_log`(`uid`, `login_id`, `uname`, `nickname`, `ip`, `expire`, `ct`, `ut`) VALUES")
	sqlArgs := []interface{}{}
	for i, l := range ls {
		if i < len(ls)-1 {
			sqlBuf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?);")
		}
		sqlArgs = append(sqlArgs, l.Uid, l.Login_id, l.Uname, l.Nickname, l.Ip, l.Expire, l.Ct, l.Ut)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_login_log(l *Login_log) (ls []Login_log, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `uid`, `login_id`, `uname`, `nickname`, `ip`, `expire`, `ct`, `ut` FROM `lt`.`login_log` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if l.Uid != 0 {
		sqlBuf.WriteString("AND `uid`=? ")
		sqlArgs = append(sqlArgs, l.Uid)
	}
	if l.Login_id != "" {
		sqlBuf.WriteString("AND `login_id`=? ")
		sqlArgs = append(sqlArgs, l.Login_id)
	}
	if l.Uname != "" {
		sqlBuf.WriteString("AND `uname`=? ")
		sqlArgs = append(sqlArgs, l.Uname)
	}
	if l.Nickname != "" {
		sqlBuf.WriteString("AND `nickname`=? ")
		sqlArgs = append(sqlArgs, l.Nickname)
	}
	if l.Ip != "" {
		sqlBuf.WriteString("AND `ip`=? ")
		sqlArgs = append(sqlArgs, l.Ip)
	}
	if l.Expire != "" {
		sqlBuf.WriteString("AND `expire`=? ")
		sqlArgs = append(sqlArgs, l.Expire)
	}
	if l.Ct != "" {
		sqlBuf.WriteString("AND `ct`=? ")
		sqlArgs = append(sqlArgs, l.Ct)
	}
	if l.Ut != "" {
		sqlBuf.WriteString("AND `ut`=? ")
		sqlArgs = append(sqlArgs, l.Ut)
	}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &ls)
	return
}
func IU_login_log(l *Login_log) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT INTO `lt`.`login_log`(`uid`, `login_id`, `uname`, `nickname`, `ip`, `expire`, `ct`, `ut`) VALUES(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{l.Uid, l.Login_id, l.Uname, l.Nickname, l.Ip, l.Expire, l.Ct, l.Ut}
	if l.Uid != 0 {
		sqlBuf.WriteString("`uid`=?,")
		sqlArgs = append(sqlArgs, l.Uid)
	}
	if l.Login_id != "" {
		sqlBuf.WriteString("`login_id`=?,")
		sqlArgs = append(sqlArgs, l.Login_id)
	}
	if l.Uname != "" {
		sqlBuf.WriteString("`uname`=?,")
		sqlArgs = append(sqlArgs, l.Uname)
	}
	if l.Nickname != "" {
		sqlBuf.WriteString("`nickname`=?,")
		sqlArgs = append(sqlArgs, l.Nickname)
	}
	if l.Ip != "" {
		sqlBuf.WriteString("`ip`=?,")
		sqlArgs = append(sqlArgs, l.Ip)
	}
	if l.Expire != "" {
		sqlBuf.WriteString("`expire`=?,")
		sqlArgs = append(sqlArgs, l.Expire)
	}
	if l.Ct != "" {
		sqlBuf.WriteString("`ct`=?,")
		sqlArgs = append(sqlArgs, l.Ct)
	}
	if l.Ut != "" {
		sqlBuf.WriteString("`ut`=?,")
		sqlArgs = append(sqlArgs, l.Ut)
	}

请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString("`ut`=NOW();")
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}

#uri
INSERT INTO `lt`.`uri`(`uri`, `remark`) VALUES(?, ?) ON DUPLICATE KEY UPDATE `remark`=?;
SELECT `uri`, `remark` FROM `lt`.`uri` WHERE `uri`=?;
UPDATE `lt`.`uri` SET `uri`=?, `remark`=? WHERE `uri`=?;
type Uri struct {
	Uri string `json:"uri,omitempty"`
	Remark string `json:"remark,omitempty"`
}
func I_uris(us []Uri) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `lt`.`uri`(`uri`, `remark`) VALUES")
	sqlArgs := []interface{}{}
	for i, u := range us {
		if i < len(us)-1 {
			sqlBuf.WriteString("(?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?);")
		}
		sqlArgs = append(sqlArgs, u.Uri, u.Remark)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_uri(u *Uri) (us []Uri, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `uri`, `remark` FROM `lt`.`uri` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if u.Uri != "" {
		sqlBuf.WriteString("AND `uri`=? ")
		sqlArgs = append(sqlArgs, u.Uri)
	}
	if u.Remark != "" {
		sqlBuf.WriteString("AND `remark`=? ")
		sqlArgs = append(sqlArgs, u.Remark)
	}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &us)
	return
}
func IU_uri(u *Uri) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT INTO `lt`.`uri`(`uri`, `remark`) VALUES(?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{u.Uri, u.Remark}
	if u.Uri != "" {
		sqlBuf.WriteString("`uri`=?,")
		sqlArgs = append(sqlArgs, u.Uri)
	}
	if u.Remark != "" {
		sqlBuf.WriteString("`remark`=?,")
		sqlArgs = append(sqlArgs, u.Remark)
	}

请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString("`ut`=NOW();")
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}

#user
INSERT INTO `lt`.`user`(`uid`, `uname`, `nickname`, `passwd`, `expire`, `puid`, `title`, `ct`, `ut`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `uname`=?, `nickname`=?, `passwd`=?, `expire`=?, `puid`=?, `title`=?, `ct`=?, `ut`=?;
SELECT `uid`, `uname`, `nickname`, `passwd`, `expire`, `puid`, `title`, `ct`, `ut` FROM `lt`.`user` WHERE `uid`=?;
UPDATE `lt`.`user` SET `uid`=?, `uname`=?, `nickname`=?, `passwd`=?, `expire`=?, `puid`=?, `title`=?, `ct`=?, `ut`=? WHERE `uid`=?;
type User struct {
	Uid int `json:"uid,omitempty"`
	Uname string `json:"uname,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Passwd string `json:"passwd,omitempty"`
	Expire string `json:"expire,omitempty"`
	Puid int `json:"puid,omitempty"`
	Title string `json:"title,omitempty"`
	Ct string `json:"ct,omitempty"`
	Ut string `json:"ut,omitempty"`
}
func I_users(us []User) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `lt`.`user`(`uid`, `uname`, `nickname`, `passwd`, `expire`, `puid`, `title`, `ct`, `ut`) VALUES")
	sqlArgs := []interface{}{}
	for i, u := range us {
		if i < len(us)-1 {
			sqlBuf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?);")
		}
		sqlArgs = append(sqlArgs, u.Uid, u.Uname, u.Nickname, u.Passwd, u.Expire, u.Puid, u.Title, u.Ct, u.Ut)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_user(u *User) (us []User, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `uid`, `uname`, `nickname`, `passwd`, `expire`, `puid`, `title`, `ct`, `ut` FROM `lt`.`user` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if u.Uid != 0 {
		sqlBuf.WriteString("AND `uid`=? ")
		sqlArgs = append(sqlArgs, u.Uid)
	}
	if u.Uname != "" {
		sqlBuf.WriteString("AND `uname`=? ")
		sqlArgs = append(sqlArgs, u.Uname)
	}
	if u.Nickname != "" {
		sqlBuf.WriteString("AND `nickname`=? ")
		sqlArgs = append(sqlArgs, u.Nickname)
	}
	if u.Passwd != "" {
		sqlBuf.WriteString("AND `passwd`=? ")
		sqlArgs = append(sqlArgs, u.Passwd)
	}
	if u.Expire != "" {
		sqlBuf.WriteString("AND `expire`=? ")
		sqlArgs = append(sqlArgs, u.Expire)
	}
	if u.Puid != 0 {
		sqlBuf.WriteString("AND `puid`=? ")
		sqlArgs = append(sqlArgs, u.Puid)
	}
	if u.Title != "" {
		sqlBuf.WriteString("AND `title`=? ")
		sqlArgs = append(sqlArgs, u.Title)
	}
	if u.Ct != "" {
		sqlBuf.WriteString("AND `ct`=? ")
		sqlArgs = append(sqlArgs, u.Ct)
	}
	if u.Ut != "" {
		sqlBuf.WriteString("AND `ut`=? ")
		sqlArgs = append(sqlArgs, u.Ut)
	}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &us)
	return
}
func IU_user(u *User) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT INTO `lt`.`user`(`uid`, `uname`, `nickname`, `passwd`, `expire`, `puid`, `title`, `ct`, `ut`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{u.Uid, u.Uname, u.Nickname, u.Passwd, u.Expire, u.Puid, u.Title, u.Ct, u.Ut}
	if u.Uid != 0 {
		sqlBuf.WriteString("`uid`=?,")
		sqlArgs = append(sqlArgs, u.Uid)
	}
	if u.Uname != "" {
		sqlBuf.WriteString("`uname`=?,")
		sqlArgs = append(sqlArgs, u.Uname)
	}
	if u.Nickname != "" {
		sqlBuf.WriteString("`nickname`=?,")
		sqlArgs = append(sqlArgs, u.Nickname)
	}
	if u.Passwd != "" {
		sqlBuf.WriteString("`passwd`=?,")
		sqlArgs = append(sqlArgs, u.Passwd)
	}
	if u.Expire != "" {
		sqlBuf.WriteString("`expire`=?,")
		sqlArgs = append(sqlArgs, u.Expire)
	}
	if u.Puid != 0 {
		sqlBuf.WriteString("`puid`=?,")
		sqlArgs = append(sqlArgs, u.Puid)
	}
	if u.Title != "" {
		sqlBuf.WriteString("`title`=?,")
		sqlArgs = append(sqlArgs, u.Title)
	}
	if u.Ct != "" {
		sqlBuf.WriteString("`ct`=?,")
		sqlArgs = append(sqlArgs, u.Ct)
	}
	if u.Ut != "" {
		sqlBuf.WriteString("`ut`=?,")
		sqlArgs = append(sqlArgs, u.Ut)
	}

请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString("`ut`=NOW();")
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}

#user_uri
INSERT INTO `lt`.`user_uri`(`uid`, `uri`, `ct`) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE `ct`=?;
SELECT `uid`, `uri`, `ct` FROM `lt`.`user_uri` WHERE `uid`=? AND `uri`=?;
UPDATE `lt`.`user_uri` SET `uid`=?, `uri`=?, `ct`=? WHERE `uid`=? AND `uri`=?;
type User_uri struct {
	Uid int `json:"uid,omitempty"`
	Uri string `json:"uri,omitempty"`
	Ct string `json:"ct,omitempty"`
}
func I_user_uris(us []User_uri) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `lt`.`user_uri`(`uid`, `uri`, `ct`) VALUES")
	sqlArgs := []interface{}{}
	for i, u := range us {
		if i < len(us)-1 {
			sqlBuf.WriteString("(?, ?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?, ?);")
		}
		sqlArgs = append(sqlArgs, u.Uid, u.Uri, u.Ct)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_user_uri(u *User_uri) (us []User_uri, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `uid`, `uri`, `ct` FROM `lt`.`user_uri` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if u.Uid != 0 {
		sqlBuf.WriteString("AND `uid`=? ")
		sqlArgs = append(sqlArgs, u.Uid)
	}
	if u.Uri != "" {
		sqlBuf.WriteString("AND `uri`=? ")
		sqlArgs = append(sqlArgs, u.Uri)
	}
	if u.Ct != "" {
		sqlBuf.WriteString("AND `ct`=? ")
		sqlArgs = append(sqlArgs, u.Ct)
	}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &us)
	return
}
func IU_user_uri(u *User_uri) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT INTO `lt`.`user_uri`(`uid`, `uri`, `ct`) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{u.Uid, u.Uri, u.Ct}
	if u.Uid != 0 {
		sqlBuf.WriteString("`uid`=?,")
		sqlArgs = append(sqlArgs, u.Uid)
	}
	if u.Uri != "" {
		sqlBuf.WriteString("`uri`=?,")
		sqlArgs = append(sqlArgs, u.Uri)
	}
	if u.Ct != "" {
		sqlBuf.WriteString("`ct`=?,")
		sqlArgs = append(sqlArgs, u.Ct)
	}

请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString("`ut`=NOW();")
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}


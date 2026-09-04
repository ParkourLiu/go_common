#client
INSERT INTO `cf`.`client`(`cid`, `v`, `ip`, `ipaddr`, `iszp`, `package`, `file_paths`, `remark`, `ht`, `ct`, `ut`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `v`=?, `ip`=?, `ipaddr`=?, `iszp`=?, `package`=?, `file_paths`=?, `remark`=?, `ht`=?, `ct`=?, `ut`=?;
SELECT `cid`, `v`, `ip`, `ipaddr`, `iszp`, `package`, `file_paths`, `remark`, `ht`, `ct`, `ut` FROM `cf`.`client` WHERE `cid`=?;
UPDATE `cf`.`client` SET `cid`=?, `v`=?, `ip`=?, `ipaddr`=?, `iszp`=?, `package`=?, `file_paths`=?, `remark`=?, `ht`=?, `ct`=?, `ut`=? WHERE `cid`=?;
type Client struct {
	Cid string `json:"cid,omitempty"`
	V int `json:"v,omitempty"`
	Ip string `json:"ip,omitempty"`
	Ipaddr string `json:"ipaddr,omitempty"`
	Iszp int `json:"iszp,omitempty"`
	Package string `json:"package,omitempty"`
	File_paths string `json:"file_paths,omitempty"`
	Remark string `json:"remark,omitempty"`
	Ht string `json:"ht,omitempty"`
	Ct string `json:"ct,omitempty"`
	Ut string `json:"ut,omitempty"`
}
func I_clients(cs []Client) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `cf`.`client`(`cid`, `v`, `ip`, `ipaddr`, `iszp`, `package`, `file_paths`, `remark`, `ht`, `ct`, `ut`) VALUES")
	sqlArgs := []interface{}{}
	for i, c := range cs {
		if i < len(cs)-1 {
			sqlBuf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
		}
		sqlArgs = append(sqlArgs, c.Cid, c.V, c.Ip, c.Ipaddr, c.Iszp, c.Package, c.File_paths, c.Remark, c.Ht, c.Ct, c.Ut)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_client(c *Client) (cs []Client, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `cid`, `v`, `ip`, `ipaddr`, `iszp`, `package`, `file_paths`, `remark`, `ht`, `ct`, `ut` FROM `cf`.`client` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if c.Cid != "" {
		sqlBuf.WriteString("AND `cid`=? ")
		sqlArgs = append(sqlArgs, c.Cid)
	}
	if c.V != 0 {
		sqlBuf.WriteString("AND `v`=? ")
		sqlArgs = append(sqlArgs, c.V)
	}
	if c.Ip != "" {
		sqlBuf.WriteString("AND `ip`=? ")
		sqlArgs = append(sqlArgs, c.Ip)
	}
	if c.Ipaddr != "" {
		sqlBuf.WriteString("AND `ipaddr`=? ")
		sqlArgs = append(sqlArgs, c.Ipaddr)
	}
	if c.Iszp != 0 {
		sqlBuf.WriteString("AND `iszp`=? ")
		sqlArgs = append(sqlArgs, c.Iszp)
	}
	if c.Package != "" {
		sqlBuf.WriteString("AND `package`=? ")
		sqlArgs = append(sqlArgs, c.Package)
	}
	if c.File_paths != "" {
		sqlBuf.WriteString("AND `file_paths`=? ")
		sqlArgs = append(sqlArgs, c.File_paths)
	}
	if c.Remark != "" {
		sqlBuf.WriteString("AND `remark`=? ")
		sqlArgs = append(sqlArgs, c.Remark)
	}
	if c.Ht != "" {
		sqlBuf.WriteString("AND `ht`=? ")
		sqlArgs = append(sqlArgs, c.Ht)
	}
	if c.Ct != "" {
		sqlBuf.WriteString("AND `ct`=? ")
		sqlArgs = append(sqlArgs, c.Ct)
	}
	if c.Ut != "" {
		sqlBuf.WriteString("AND `ut`=? ")
		sqlArgs = append(sqlArgs, c.Ut)
	}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &cs)
	return
}
func IU_client(c *Client) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT INTO `cf`.`client`(`cid`, `v`, `ip`, `ipaddr`, `iszp`, `package`, `file_paths`, `remark`, `ht`, `ct`, `ut`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{c.Cid, c.V, c.Ip, c.Ipaddr, c.Iszp, c.Package, c.File_paths, c.Remark, c.Ht, c.Ct, c.Ut}
	if c.Cid != "" {
		sqlBuf.WriteString("`cid`=?,")
		sqlArgs = append(sqlArgs, c.Cid)
	}
	if c.V != 0 {
		sqlBuf.WriteString("`v`=?,")
		sqlArgs = append(sqlArgs, c.V)
	}
	if c.Ip != "" {
		sqlBuf.WriteString("`ip`=?,")
		sqlArgs = append(sqlArgs, c.Ip)
	}
	if c.Ipaddr != "" {
		sqlBuf.WriteString("`ipaddr`=?,")
		sqlArgs = append(sqlArgs, c.Ipaddr)
	}
	if c.Iszp != 0 {
		sqlBuf.WriteString("`iszp`=?,")
		sqlArgs = append(sqlArgs, c.Iszp)
	}
	if c.Package != "" {
		sqlBuf.WriteString("`package`=?,")
		sqlArgs = append(sqlArgs, c.Package)
	}
	if c.File_paths != "" {
		sqlBuf.WriteString("`file_paths`=?,")
		sqlArgs = append(sqlArgs, c.File_paths)
	}
	if c.Remark != "" {
		sqlBuf.WriteString("`remark`=?,")
		sqlArgs = append(sqlArgs, c.Remark)
	}
	if c.Ht != "" {
		sqlBuf.WriteString("`ht`=?,")
		sqlArgs = append(sqlArgs, c.Ht)
	}
	if c.Ct != "" {
		sqlBuf.WriteString("`ct`=?,")
		sqlArgs = append(sqlArgs, c.Ct)
	}
	if c.Ut != "" {
		sqlBuf.WriteString("`ut`=?,")
		sqlArgs = append(sqlArgs, c.Ut)
	}

请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString("`ut`=NOW();")
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}

#client_history
INSERT INTO `cf`.`client_history`(`cid`, `ct`, `ip`, `ipaddr`) VALUES(?, ?, ?, ?) ON DUPLICATE KEY UPDATE `ip`=?, `ipaddr`=?;
SELECT `cid`, `ct`, `ip`, `ipaddr` FROM `cf`.`client_history` WHERE `cid`=? AND `ct`=?;
UPDATE `cf`.`client_history` SET `cid`=?, `ct`=?, `ip`=?, `ipaddr`=? WHERE `cid`=? AND `ct`=?;
type Client_history struct {
	Cid string `json:"cid,omitempty"`
	Ct string `json:"ct,omitempty"`
	Ip string `json:"ip,omitempty"`
	Ipaddr string `json:"ipaddr,omitempty"`
}
func I_client_historys(cs []Client_history) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `cf`.`client_history`(`cid`, `ct`, `ip`, `ipaddr`) VALUES")
	sqlArgs := []interface{}{}
	for i, c := range cs {
		if i < len(cs)-1 {
			sqlBuf.WriteString("(?, ?, ?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?, ?, ?);")
		}
		sqlArgs = append(sqlArgs, c.Cid, c.Ct, c.Ip, c.Ipaddr)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_client_history(c *Client_history) (cs []Client_history, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `cid`, `ct`, `ip`, `ipaddr` FROM `cf`.`client_history` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if c.Cid != "" {
		sqlBuf.WriteString("AND `cid`=? ")
		sqlArgs = append(sqlArgs, c.Cid)
	}
	if c.Ct != "" {
		sqlBuf.WriteString("AND `ct`=? ")
		sqlArgs = append(sqlArgs, c.Ct)
	}
	if c.Ip != "" {
		sqlBuf.WriteString("AND `ip`=? ")
		sqlArgs = append(sqlArgs, c.Ip)
	}
	if c.Ipaddr != "" {
		sqlBuf.WriteString("AND `ipaddr`=? ")
		sqlArgs = append(sqlArgs, c.Ipaddr)
	}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &cs)
	return
}
func IU_client_history(c *Client_history) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT INTO `cf`.`client_history`(`cid`, `ct`, `ip`, `ipaddr`) VALUES(?, ?, ?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{c.Cid, c.Ct, c.Ip, c.Ipaddr}
	if c.Cid != "" {
		sqlBuf.WriteString("`cid`=?,")
		sqlArgs = append(sqlArgs, c.Cid)
	}
	if c.Ct != "" {
		sqlBuf.WriteString("`ct`=?,")
		sqlArgs = append(sqlArgs, c.Ct)
	}
	if c.Ip != "" {
		sqlBuf.WriteString("`ip`=?,")
		sqlArgs = append(sqlArgs, c.Ip)
	}
	if c.Ipaddr != "" {
		sqlBuf.WriteString("`ipaddr`=?,")
		sqlArgs = append(sqlArgs, c.Ipaddr)
	}

请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString("`ut`=NOW();")
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}

#package
INSERT INTO `cf`.`package`(`name`, `url`, `open_path`, `exec_name`) VALUES(?, ?, ?, ?) ON DUPLICATE KEY UPDATE `url`=?, `open_path`=?, `exec_name`=?;
SELECT `name`, `url`, `open_path`, `exec_name` FROM `cf`.`package` WHERE `name`=?;
UPDATE `cf`.`package` SET `name`=?, `url`=?, `open_path`=?, `exec_name`=? WHERE `name`=?;
type Package struct {
	Name string `json:"name,omitempty"`
	Url string `json:"url,omitempty"`
	Open_path string `json:"open_path,omitempty"`
	Exec_name string `json:"exec_name,omitempty"`
}
func I_packages(ps []Package) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `cf`.`package`(`name`, `url`, `open_path`, `exec_name`) VALUES")
	sqlArgs := []interface{}{}
	for i, p := range ps {
		if i < len(ps)-1 {
			sqlBuf.WriteString("(?, ?, ?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?, ?, ?);")
		}
		sqlArgs = append(sqlArgs, p.Name, p.Url, p.Open_path, p.Exec_name)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_package(p *Package) (ps []Package, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `name`, `url`, `open_path`, `exec_name` FROM `cf`.`package` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if p.Name != "" {
		sqlBuf.WriteString("AND `name`=? ")
		sqlArgs = append(sqlArgs, p.Name)
	}
	if p.Url != "" {
		sqlBuf.WriteString("AND `url`=? ")
		sqlArgs = append(sqlArgs, p.Url)
	}
	if p.Open_path != "" {
		sqlBuf.WriteString("AND `open_path`=? ")
		sqlArgs = append(sqlArgs, p.Open_path)
	}
	if p.Exec_name != "" {
		sqlBuf.WriteString("AND `exec_name`=? ")
		sqlArgs = append(sqlArgs, p.Exec_name)
	}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &ps)
	return
}
func IU_package(p *Package) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT INTO `cf`.`package`(`name`, `url`, `open_path`, `exec_name`) VALUES(?, ?, ?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{p.Name, p.Url, p.Open_path, p.Exec_name}
	if p.Name != "" {
		sqlBuf.WriteString("`name`=?,")
		sqlArgs = append(sqlArgs, p.Name)
	}
	if p.Url != "" {
		sqlBuf.WriteString("`url`=?,")
		sqlArgs = append(sqlArgs, p.Url)
	}
	if p.Open_path != "" {
		sqlBuf.WriteString("`open_path`=?,")
		sqlArgs = append(sqlArgs, p.Open_path)
	}
	if p.Exec_name != "" {
		sqlBuf.WriteString("`exec_name`=?,")
		sqlArgs = append(sqlArgs, p.Exec_name)
	}

请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString("`ut`=NOW();")
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}

#package_history
INSERT INTO `cf`.`package_history`(`cid`, `package`, `ct`) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE `cid`=?, `package`=?, `ct`=?;
SELECT `cid`, `package`, `ct` FROM `cf`.`package_history` WHERE ;
UPDATE `cf`.`package_history` SET `cid`=?, `package`=?, `ct`=? WHERE ;
type Package_history struct {
	Cid string `json:"cid,omitempty"`
	Package string `json:"package,omitempty"`
	Ct string `json:"ct,omitempty"`
}
func I_package_historys(ps []Package_history) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `cf`.`package_history`(`cid`, `package`, `ct`) VALUES")
	sqlArgs := []interface{}{}
	for i, p := range ps {
		if i < len(ps)-1 {
			sqlBuf.WriteString("(?, ?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?, ?);")
		}
		sqlArgs = append(sqlArgs, p.Cid, p.Package, p.Ct)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_package_history(p *Package_history) (ps []Package_history, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `cid`, `package`, `ct` FROM `cf`.`package_history` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if p.Cid != "" {
		sqlBuf.WriteString("AND `cid`=? ")
		sqlArgs = append(sqlArgs, p.Cid)
	}
	if p.Package != "" {
		sqlBuf.WriteString("AND `package`=? ")
		sqlArgs = append(sqlArgs, p.Package)
	}
	if p.Ct != "" {
		sqlBuf.WriteString("AND `ct`=? ")
		sqlArgs = append(sqlArgs, p.Ct)
	}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &ps)
	return
}
func IU_package_history(p *Package_history) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT INTO `cf`.`package_history`(`cid`, `package`, `ct`) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{p.Cid, p.Package, p.Ct}
	if p.Cid != "" {
		sqlBuf.WriteString("`cid`=?,")
		sqlArgs = append(sqlArgs, p.Cid)
	}
	if p.Package != "" {
		sqlBuf.WriteString("`package`=?,")
		sqlArgs = append(sqlArgs, p.Package)
	}
	if p.Ct != "" {
		sqlBuf.WriteString("`ct`=?,")
		sqlArgs = append(sqlArgs, p.Ct)
	}

请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString("`ut`=NOW();")
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}

#zpkeyword
INSERT INTO `cf`.`zpkeyword`(`key`) VALUES(?) ON DUPLICATE KEY UPDATE ;
SELECT `key` FROM `cf`.`zpkeyword` WHERE `key`=?;
UPDATE `cf`.`zpkeyword` SET `key`=? WHERE `key`=?;
type Zpkeyword struct {
	Key string `json:"key,omitempty"`
}
func I_zpkeywords(zs []Zpkeyword) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `cf`.`zpkeyword`(`key`) VALUES")
	sqlArgs := []interface{}{}
	for i, z := range zs {
		if i < len(zs)-1 {
			sqlBuf.WriteString("(?),")
		} else {
			sqlBuf.WriteString("(?);")
		}
		sqlArgs = append(sqlArgs, z.Key)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_zpkeyword(z *Zpkeyword) (zs []Zpkeyword, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `key` FROM `cf`.`zpkeyword` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if z.Key != "" {
		sqlBuf.WriteString("AND `key`=? ")
		sqlArgs = append(sqlArgs, z.Key)
	}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &zs)
	return
}
func IU_zpkeyword(z *Zpkeyword) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT INTO `cf`.`zpkeyword`(`key`) VALUES(?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{z.Key}
	if z.Key != "" {
		sqlBuf.WriteString("`key`=?,")
		sqlArgs = append(sqlArgs, z.Key)
	}

请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString("`ut`=NOW();")
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}


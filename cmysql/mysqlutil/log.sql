#client
INSERT INTO `log`.`client`(`ip`, `ipaddr`, `log`, `remark`, `return`, `ht`, `ct`, `ut`) VALUES(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `ipaddr`=?, `log`=?, `remark`=?, `return`=?, `ht`=?, `ct`=?, `ut`=?;
SELECT `ip`, `ipaddr`, `log`, `remark`, `return`, `ht`, `ct`, `ut` FROM `log`.`client` WHERE `ip`=?;
UPDATE `log`.`client` SET `ip`=?, `ipaddr`=?, `log`=?, `remark`=?, `return`=?, `ht`=?, `ct`=?, `ut`=? WHERE `ip`=?;
type Client struct {
	Ip string `json:"ip,omitempty"`
	Ipaddr string `json:"ipaddr,omitempty"`
	Log string `json:"log,omitempty"`
	Remark string `json:"remark,omitempty"`
	Return string `json:"return,omitempty"`
	Ht string `json:"ht,omitempty"`
	Ct string `json:"ct,omitempty"`
	Ut string `json:"ut,omitempty"`
}
func I_clients(cs []Client) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `log`.`client`(`ip`, `ipaddr`, `log`, `remark`, `return`, `ht`, `ct`, `ut`) VALUES")
	sqlArgs := []interface{}{}
	for i, c := range cs {
		if i < len(cs)-1 {
			sqlBuf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?);")
		}
		sqlArgs = append(sqlArgs, c.Ip, c.Ipaddr, c.Log, c.Remark, c.Return, c.Ht, c.Ct, c.Ut)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_client(c *Client) (cs []Client, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `ip`, `ipaddr`, `log`, `remark`, `return`, `ht`, `ct`, `ut` FROM `log`.`client` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if c.Ip != "" {
		sqlBuf.WriteString("AND `ip`=? ")
		sqlArgs = append(sqlArgs, c.Ip)
	}
	if c.Ipaddr != "" {
		sqlBuf.WriteString("AND `ipaddr`=? ")
		sqlArgs = append(sqlArgs, c.Ipaddr)
	}
	if c.Log != "" {
		sqlBuf.WriteString("AND `log`=? ")
		sqlArgs = append(sqlArgs, c.Log)
	}
	if c.Remark != "" {
		sqlBuf.WriteString("AND `remark`=? ")
		sqlArgs = append(sqlArgs, c.Remark)
	}
	if c.Return != "" {
		sqlBuf.WriteString("AND `return`=? ")
		sqlArgs = append(sqlArgs, c.Return)
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
	sqlBuf.WriteString("INSERT INTO `log`.`client`(`ip`, `ipaddr`, `log`, `remark`, `return`, `ht`, `ct`, `ut`) VALUES(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{c.Ip, c.Ipaddr, c.Log, c.Remark, c.Return, c.Ht, c.Ct, c.Ut}
	if c.Ip != "" {
		sqlBuf.WriteString("`ip`=?,")
		sqlArgs = append(sqlArgs, c.Ip)
	}
	if c.Ipaddr != "" {
		sqlBuf.WriteString("`ipaddr`=?,")
		sqlArgs = append(sqlArgs, c.Ipaddr)
	}
	if c.Log != "" {
		sqlBuf.WriteString("`log`=?,")
		sqlArgs = append(sqlArgs, c.Log)
	}
	if c.Remark != "" {
		sqlBuf.WriteString("`remark`=?,")
		sqlArgs = append(sqlArgs, c.Remark)
	}
	if c.Return != "" {
		sqlBuf.WriteString("`return`=?,")
		sqlArgs = append(sqlArgs, c.Return)
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


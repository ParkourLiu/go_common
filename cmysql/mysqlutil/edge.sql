#client
INSERT INTO `edge`.`client`(`cid`, `sid`, `group`, `remark`, `serverflag`, `bit`, `pid`, `system`, `ip`, `addr`, `innerip`, `whoami`, `hostname`, `version`, `sleeptime`, `ht`, `ct`, `ut`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `cid`=?, `group`=?, `remark`=?, `serverflag`=?, `bit`=?, `pid`=?, `system`=?, `ip`=?, `addr`=?, `innerip`=?, `whoami`=?, `hostname`=?, `version`=?, `sleeptime`=?, `ht`=?, `ct`=?, `ut`=?;
SELECT `cid`, `sid`, `group`, `remark`, `serverflag`, `bit`, `pid`, `system`, `ip`, `addr`, `innerip`, `whoami`, `hostname`, `version`, `sleeptime`, `ht`, `ct`, `ut` FROM `edge`.`client` WHERE `sid`=?;
UPDATE `edge`.`client` SET `cid`=?, `sid`=?, `group`=?, `remark`=?, `serverflag`=?, `bit`=?, `pid`=?, `system`=?, `ip`=?, `addr`=?, `innerip`=?, `whoami`=?, `hostname`=?, `version`=?, `sleeptime`=?, `ht`=?, `ct`=?, `ut`=? WHERE `sid`=?;
type Client struct {
	Cid string `json:"cid,omitempty"`
	Sid string `json:"sid,omitempty"`
	Group string `json:"group,omitempty"`
	Remark string `json:"remark,omitempty"`
	Serverflag string `json:"serverflag,omitempty"`
	Bit int `json:"bit,omitempty"`
	Pid int `json:"pid,omitempty"`
	System int `json:"system,omitempty"`
	Ip string `json:"ip,omitempty"`
	Addr string `json:"addr,omitempty"`
	Innerip string `json:"innerip,omitempty"`
	Whoami string `json:"whoami,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Version int `json:"version,omitempty"`
	Sleeptime int `json:"sleeptime,omitempty"`
	Ht string `json:"ht,omitempty"`
	Ct string `json:"ct,omitempty"`
	Ut string `json:"ut,omitempty"`
}
func I_clients(cs []Client) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `edge`.`client`(`cid`, `sid`, `group`, `remark`, `serverflag`, `bit`, `pid`, `system`, `ip`, `addr`, `innerip`, `whoami`, `hostname`, `version`, `sleeptime`, `ht`, `ct`, `ut`) VALUES")
	sqlArgs := []interface{}{}
	for i, c := range cs {
		if i < len(cs)-1 {
			sqlBuf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);")
		}
		sqlArgs = append(sqlArgs, c.Cid, c.Sid, c.Group, c.Remark, c.Serverflag, c.Bit, c.Pid, c.System, c.Ip, c.Addr, c.Innerip, c.Whoami, c.Hostname, c.Version, c.Sleeptime, c.Ht, c.Ct, c.Ut)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_client(c *Client) (cs []Client, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `cid`, `sid`, `group`, `remark`, `serverflag`, `bit`, `pid`, `system`, `ip`, `addr`, `innerip`, `whoami`, `hostname`, `version`, `sleeptime`, `ht`, `ct`, `ut` FROM `edge`.`client` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if c.Cid != "" {
		sqlBuf.WriteString("AND `cid`=? ")
		sqlArgs = append(sqlArgs, c.Cid)
	}
	if c.Sid != "" {
		sqlBuf.WriteString("AND `sid`=? ")
		sqlArgs = append(sqlArgs, c.Sid)
	}
	if c.Group != "" {
		sqlBuf.WriteString("AND `group`=? ")
		sqlArgs = append(sqlArgs, c.Group)
	}
	if c.Remark != "" {
		sqlBuf.WriteString("AND `remark`=? ")
		sqlArgs = append(sqlArgs, c.Remark)
	}
	if c.Serverflag != "" {
		sqlBuf.WriteString("AND `serverflag`=? ")
		sqlArgs = append(sqlArgs, c.Serverflag)
	}
	if c.Bit != 0 {
		sqlBuf.WriteString("AND `bit`=? ")
		sqlArgs = append(sqlArgs, c.Bit)
	}
	if c.Pid != 0 {
		sqlBuf.WriteString("AND `pid`=? ")
		sqlArgs = append(sqlArgs, c.Pid)
	}
	if c.System != 0 {
		sqlBuf.WriteString("AND `system`=? ")
		sqlArgs = append(sqlArgs, c.System)
	}
	if c.Ip != "" {
		sqlBuf.WriteString("AND `ip`=? ")
		sqlArgs = append(sqlArgs, c.Ip)
	}
	if c.Addr != "" {
		sqlBuf.WriteString("AND `addr`=? ")
		sqlArgs = append(sqlArgs, c.Addr)
	}
	if c.Innerip != "" {
		sqlBuf.WriteString("AND `innerip`=? ")
		sqlArgs = append(sqlArgs, c.Innerip)
	}
	if c.Whoami != "" {
		sqlBuf.WriteString("AND `whoami`=? ")
		sqlArgs = append(sqlArgs, c.Whoami)
	}
	if c.Hostname != "" {
		sqlBuf.WriteString("AND `hostname`=? ")
		sqlArgs = append(sqlArgs, c.Hostname)
	}
	if c.Version != 0 {
		sqlBuf.WriteString("AND `version`=? ")
		sqlArgs = append(sqlArgs, c.Version)
	}
	if c.Sleeptime != 0 {
		sqlBuf.WriteString("AND `sleeptime`=? ")
		sqlArgs = append(sqlArgs, c.Sleeptime)
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
	sqlBuf.WriteString("INSERT INTO `edge`.`client`(`cid`, `sid`, `group`, `remark`, `serverflag`, `bit`, `pid`, `system`, `ip`, `addr`, `innerip`, `whoami`, `hostname`, `version`, `sleeptime`, `ht`, `ct`, `ut`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{c.Cid, c.Sid, c.Group, c.Remark, c.Serverflag, c.Bit, c.Pid, c.System, c.Ip, c.Addr, c.Innerip, c.Whoami, c.Hostname, c.Version, c.Sleeptime, c.Ht, c.Ct, c.Ut}
	if c.Cid != "" {
		sqlBuf.WriteString("`cid`=?,")
		sqlArgs = append(sqlArgs, c.Cid)
	}
	if c.Sid != "" {
		sqlBuf.WriteString("`sid`=?,")
		sqlArgs = append(sqlArgs, c.Sid)
	}
	if c.Group != "" {
		sqlBuf.WriteString("`group`=?,")
		sqlArgs = append(sqlArgs, c.Group)
	}
	if c.Remark != "" {
		sqlBuf.WriteString("`remark`=?,")
		sqlArgs = append(sqlArgs, c.Remark)
	}
	if c.Serverflag != "" {
		sqlBuf.WriteString("`serverflag`=?,")
		sqlArgs = append(sqlArgs, c.Serverflag)
	}
	if c.Bit != 0 {
		sqlBuf.WriteString("`bit`=?,")
		sqlArgs = append(sqlArgs, c.Bit)
	}
	if c.Pid != 0 {
		sqlBuf.WriteString("`pid`=?,")
		sqlArgs = append(sqlArgs, c.Pid)
	}
	if c.System != 0 {
		sqlBuf.WriteString("`system`=?,")
		sqlArgs = append(sqlArgs, c.System)
	}
	if c.Ip != "" {
		sqlBuf.WriteString("`ip`=?,")
		sqlArgs = append(sqlArgs, c.Ip)
	}
	if c.Addr != "" {
		sqlBuf.WriteString("`addr`=?,")
		sqlArgs = append(sqlArgs, c.Addr)
	}
	if c.Innerip != "" {
		sqlBuf.WriteString("`innerip`=?,")
		sqlArgs = append(sqlArgs, c.Innerip)
	}
	if c.Whoami != "" {
		sqlBuf.WriteString("`whoami`=?,")
		sqlArgs = append(sqlArgs, c.Whoami)
	}
	if c.Hostname != "" {
		sqlBuf.WriteString("`hostname`=?,")
		sqlArgs = append(sqlArgs, c.Hostname)
	}
	if c.Version != 0 {
		sqlBuf.WriteString("`version`=?,")
		sqlArgs = append(sqlArgs, c.Version)
	}
	if c.Sleeptime != 0 {
		sqlBuf.WriteString("`sleeptime`=?,")
		sqlArgs = append(sqlArgs, c.Sleeptime)
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

#cmd
INSERT INTO `edge`.`cmd`(`sid`, `cmd`, `ct`) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE `sid`=?, `cmd`=?, `ct`=?;
SELECT `sid`, `cmd`, `ct` FROM `edge`.`cmd` WHERE ;
UPDATE `edge`.`cmd` SET `sid`=?, `cmd`=?, `ct`=? WHERE ;
type Cmd struct {
	Sid string `json:"sid,omitempty"`
	Cmd string `json:"cmd,omitempty"`
	Ct string `json:"ct,omitempty"`
}
func I_cmds(cs []Cmd) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT IGNORE INTO `edge`.`cmd`(`sid`, `cmd`, `ct`) VALUES")
	sqlArgs := []interface{}{}
	for i, c := range cs {
		if i < len(cs)-1 {
			sqlBuf.WriteString("(?, ?, ?),")
		} else {
			sqlBuf.WriteString("(?, ?, ?);")
		}
		sqlArgs = append(sqlArgs, c.Sid, c.Cmd, c.Ct)
	}
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}
func S_cmd(c *Cmd) (cs []Cmd, err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("SELECT `sid`, `cmd`, `ct` FROM `edge`.`cmd` WHERE 1=1  ")
	sqlArgs := []interface{}{}
	if c.Sid != "" {
		sqlBuf.WriteString("AND `sid`=? ")
		sqlArgs = append(sqlArgs, c.Sid)
	}
	if c.Cmd != "" {
		sqlBuf.WriteString("AND `cmd`=? ")
		sqlArgs = append(sqlArgs, c.Cmd)
	}
	if c.Ct != "" {
		sqlBuf.WriteString("AND `ct`=? ")
		sqlArgs = append(sqlArgs, c.Ct)
	}
	err = mysqlClient.SearchFormat(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs}, &cs)
	return
}
func IU_cmd(c *Cmd) (err error) {
	sqlBuf := bytes.Buffer{}
	sqlBuf.WriteString("INSERT INTO `edge`.`cmd`(`sid`, `cmd`, `ct`) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE ")
	sqlArgs := []interface{}{c.Sid, c.Cmd, c.Ct}
	if c.Sid != "" {
		sqlBuf.WriteString("`sid`=?,")
		sqlArgs = append(sqlArgs, c.Sid)
	}
	if c.Cmd != "" {
		sqlBuf.WriteString("`cmd`=?,")
		sqlArgs = append(sqlArgs, c.Cmd)
	}
	if c.Ct != "" {
		sqlBuf.WriteString("`ct`=?,")
		sqlArgs = append(sqlArgs, c.Ct)
	}

请自行添加结束SQL,保证所有字段为空时不报错,例 sqlBuf.WriteString("`ut`=NOW();")
	return mysqlClient.Execute(&cmysql.Stmt{Sql: sqlBuf.String(), Args: sqlArgs})
}


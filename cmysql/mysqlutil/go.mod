module go_common/cmysql/mysqlutil

go 1.18

require (
	go_common/cconf v0.0.0
	go_common/clogs v0.0.0
	go_common/cmysql v0.0.0
)

require (
	github.com/astaxie/beego v1.12.3 // indirect
	github.com/denisenkom/go-mssqldb v0.12.3 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/larspensjo/config v0.0.0-20160228172812-b6db95dc6321 // indirect
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.12.0 // indirect
)

replace go_common/clogs => ..\..\clogs

replace go_common/cconf => ..\..\cconf

replace go_common/cmysql => ..\..\cmysql

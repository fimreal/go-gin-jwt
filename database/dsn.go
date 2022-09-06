// 配置数据库
package database

import (
	"os"

	"gorm.io/gorm"
)

type dbInfo struct {
	username string
	password string
	host     string
	database string
}

var db = dbInfo{
	os.Getenv("DBUSER"),
	os.Getenv("DBPWD"),
	os.Getenv("DBHOST"),
	os.Getenv("DBNAME"),
}

var DSN = db.username + ":" + db.password + "@tcp(" + db.host + ")/" + db.database

type DBORM struct {
	*gorm.DB
}

var ShowSQL = true

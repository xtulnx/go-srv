// Package accesskit 一些 db 封装，依赖:
//
//	gorm.io/gorm
//	gorm.io/gen
package accesskit

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
)

// NewDb 建立连接
//
//	conf 配置
//	tablePrefix 数据表名的前缀，可以为空
func NewDb(conf DbConfig, tablePrefix string) (*gorm.DB, error) {
	c := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		},
	}
	switch strings.ToLower(conf.Dialect) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(conf.DSN), c)
		return db, err
	case "sqlite", "sqlite3":
		db, err := gorm.Open(sqlite.Open(conf.DSN), c)
		return db, err
	default:
		return nil, fmt.Errorf("不支持的数据库: [%s]", conf.Dialect)
	}
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// SafeStr 限定字符串长度，避免字段溢出
func SafeStr(s string, size int) string {
	if size <= 1 || len(s) < size {
		return s
	}
	cc := []rune(s)
	if size <= 1 || len(cc) < size {
		return s
	}
	return string(cc[:size])
}

// SetDbLogger 打开调试日志
func SetDbLogger(db *gorm.DB, logSqlEnabled string) *gorm.DB {
	if logSqlEnabled == "1" || logSqlEnabled == "true" || logSqlEnabled == "info" {
		db = db.Debug()
	} else {
		var t logger.LogLevel
		switch logSqlEnabled {
		case "debug", "info":
			t = logger.Info
		case "warn":
			t = logger.Warn
		case "error":
			t = logger.Error
		}
		if t > 0 {
			db = db.Session(&gorm.Session{Logger: db.Logger.LogMode(t)})
		}
	}
	return db
}

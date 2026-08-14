package db

import (
	"server/pkg/global"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func gormSqlite() *gorm.DB {
	m := global.MPA_CONFIG.Sqlite
	if m.Dbname == "" {
		return nil
	}

	// sqliteConfig := sqlite.
	if db, err := gorm.Open(sqlite.Open(m.Dsn()), Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

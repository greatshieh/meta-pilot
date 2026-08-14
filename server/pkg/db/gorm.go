package db

import (
	"fmt"
	"log"
	"os"
	"server/pkg/global"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBBASE interface {
	GetLogMode() string
}

type DataBaseTable interface {
	TableName() string
}

type _gorm struct{}

var Gorm = new(_gorm)

// Config gorm 自定义配置

func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	_default := logger.New(newWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})
	var logMode DBBASE
	switch global.MPA_CONFIG.System.DbType {
	case "mysql":
		logMode = &global.MPA_CONFIG.Mysql
	case "pgsql":
		logMode = &global.MPA_CONFIG.Pgsql
	case "sqlite":
		logMode = &global.MPA_CONFIG.Sqlite
	default:
		logMode = &global.MPA_CONFIG.Mysql
	}

	switch logMode.GetLogMode() {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
	return config
}

// Gorm 初始化数据库并产生数据库全局变量
// Author SliverHorn
func NewGorm() *gorm.DB {
	switch global.MPA_CONFIG.System.DbType {
	case "mysql":
		return gormMysql()
	case "pgsql":
		return gormPgSql()
	case "sqlite":
		return gormSqlite()
	default:
		return gormMysql()
	}
}

// InitGorm 初始化数据库并注册到单例
func InitGorm() {
	db := NewGorm()
	InitDB(db)
	// 同时设置全局变量以保持向后兼容
	global.MPA_DB = db
}

// RegisterTables 注册数据库表专用
func RegisterTables(tables ...any) {
	db := global.MPA_DB

	if err := db.AutoMigrate(tables...); err != nil {
		global.MPA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(1)
	}

	global.MPA_LOG.Info("🚀 数据表注册成功")
}

type writer struct {
	logger.Writer
}

// newWriter writer 构造函数

func newWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf 格式化打印日志

func (w *writer) Printf(message string, data ...any) {
	var logZap bool
	switch global.MPA_CONFIG.System.DbType {
	case "mysql":
		logZap = global.MPA_CONFIG.Mysql.LogZap
	case "pgsql":
		logZap = global.MPA_CONFIG.Pgsql.LogZap
	}
	if logZap {
		global.MPA_LOG.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}

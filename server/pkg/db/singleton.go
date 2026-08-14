package db

import (
	"sync"

	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

// InitDB 初始化数据库连接（单例模式，只执行一次）
func InitDB(db *gorm.DB) {
	once.Do(func() {
		instance = db
	})
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return instance
}

// ResetDB 重置数据库连接（主要用于测试）
func ResetDB() {
	instance = nil
	once = sync.Once{}
}

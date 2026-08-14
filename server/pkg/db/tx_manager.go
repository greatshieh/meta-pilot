package db

import (
	"context"

	"gorm.io/gorm"
)

// Transaction 事务管理器，service 层只依赖这个接口
type Transaction interface {
	Transaction(fn func(ctx context.Context) error) error
	// todo 实现上下文事务
	// Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// gormTxManager 事务实现
type gormTxManager struct {
	db *gorm.DB
}

func NewGormTxManager(db *gorm.DB) Transaction {
	return &gormTxManager{db: db}
}

// Transaction 核心：开启事务
func (m *gormTxManager) Transaction(fn func(ctx context.Context) error) error {
	// 开启事务
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	ctx := context.Background()

	// 把 tx 注入 ctx → 所有DAO自动获取事务
	backgroundCtx := context.WithValue(ctx, "tx", tx)

	// 执行业务逻辑
	err := fn(backgroundCtx)
	if err != nil {
		tx.Rollback() // 失败回滚
		return err
	}

	return tx.Commit().Error // 成功提交
}

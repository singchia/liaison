package dao

import (
	"context"
	"fmt"
)

// TransactionFunc 事务函数类型
type TransactionFunc func(Dao) error

// WithTransaction 执行事务的辅助函数
func WithTransaction(dao Dao, fn TransactionFunc) error {
	txDao := dao.Begin()

	// 使用 defer 确保在 panic 时回滚
	defer func() {
		if r := recover(); r != nil {
			txDao.Rollback()
			panic(r) // 重新抛出 panic
		}
	}()

	// 执行事务函数
	if err := fn(txDao); err != nil {
		txDao.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	// 提交事务
	return txDao.Commit()
}

// WithTransactionContext 带上下文的事务执行函数
func WithTransactionContext(ctx context.Context, dao Dao, fn TransactionFunc) error {
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	txDao := dao.Begin()

	defer func() {
		if r := recover(); r != nil {
			txDao.Rollback()
			panic(r)
		}
	}()

	// 执行事务函数
	if err := fn(txDao); err != nil {
		txDao.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	// 再次检查上下文
	select {
	case <-ctx.Done():
		txDao.Rollback()
		return ctx.Err()
	default:
	}

	// 提交事务
	return txDao.Commit()
}

package db

import (
	"context"

	"gorm.io/gorm"

	"github.com/go-goim/core/pkg/db/mysql"
)

type mysqlTransactionCtxKey struct{}

// GetDBFromCtx try to get gorm.DB from context, if not found then return DB with context.Background
func GetDBFromCtx(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return mysql.GetDB().WithContext(context.Background())
	}

	v := ctx.Value(mysqlTransactionCtxKey{})
	if v == nil {
		return mysql.GetDB().WithContext(ctx)
	}

	// double check
	gdb, ok := v.(*gorm.DB)
	if !ok {
		// maybe set by others
		return mysql.GetDB().WithContext(ctx)
	}

	return gdb
}

// ctxWithGormDB return new context.Context contain value with gorm.DB
func ctxWithGormDB(ctx context.Context, tx *gorm.DB) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, mysqlTransactionCtxKey{}, tx.WithContext(ctx))
}

// Transaction get gorm.DB from ctx and run Transaction Operation with ctx
// How to use:
//	// 所有 db 操作的第一参数为 context.Context, 然后通过 ctx 读取 DB 对象
//	if err := db.Transaction(context.Background(), func(ctx context.Context) error {
//		if err := d.Create(ctx); err != nil {
//			return err
//		}
//
//		d.Name = "123"
//		return d.Update(ctx)
//	}); err != nil {
//		return
//	}
//
//  func (d *Domain) Create(ctx context.Context) error {
//	  return GetDBFromCtx(ctx).Create(d).Error
//  }
//
//  func (d *Domain) Update(ctx context.Context) error {
//	  return GetDBFromCtx(ctx).Updates(d).Error
//  }
func Transaction(ctx context.Context, f func(context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	gdb := GetDBFromCtx(ctx)
	return gdb.Transaction(func(tx *gorm.DB) error {
		return f(ctxWithGormDB(ctx, tx))
	})
}

// IsInTransaction return true if current operation within a transaction
func IsInTransaction(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	if ctx.Value(mysqlTransactionCtxKey{}) == nil {
		return false
	}

	return true
}

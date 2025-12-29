package util

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Trans struct {
	DB *gorm.DB
}

type TransFunc func(context.Context) error

func (a *Trans) Exec(ctx context.Context, fn TransFunc) error {
	if _, ok := FromTrans(ctx); ok {
		return fn(ctx)
	}
	return a.DB.Transaction(func(tx *gorm.DB) error {
		return fn(NewTrans(ctx, tx))
	})
}

func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	db := defDB
	if tdb, ok := FromTrans(ctx); ok {
		db = tdb
	}
	if FromRowLock(ctx) {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return db.WithContext(ctx)
}

func wrapQueryOptions(db *gorm.DB, opts QueryOptions) *gorm.DB {
	if len(opts.SelectFields) > 0 {
		db = db.Select(opts.SelectFields)
	}
	if len(opts.OmitFields) > 0 {
		db = db.Omit(opts.OmitFields...)
	}
	if len(opts.OrderFields) > 0 {
		db = db.Order(opts.OrderFields...)
	}
	return db
}

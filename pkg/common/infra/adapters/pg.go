package adapters

import (
	"context"

	"github.com/go-pg/pg/v9"

	"tiu-85/gorest-iman/pkg/common/infra/utils"
)

type db struct {
	*pg.DB
}

func (d *db) Stats() *pg.PoolStats {
	return d.DB.PoolStats()
}

func (d *db) WithForUpdate(_ ...string) DB {
	return d
}

func (d *db) WithForNoKeyUpdate(_ ...string) DB {
	return d
}

func (d *db) WithSkipLocked() DB {
	return d
}

func NewPG(conn *pg.DB) DB {
	return &db{
		DB: conn,
	}
}

func NewPGWithCtx(ctx context.Context, d *db) DB {
	return &db{
		DB: d.DB.WithContext(ctx),
	}
}

func (d *db) Save(obj interface{}) error {
	has, err := utils.ReflectHasID(obj)
	if err != nil {
		return err
	}
	if has {
		return d.Update(obj)
	}
	return d.Insert(obj)
}

func (d *db) Begin(ctx context.Context, options ...Option) (DB, error) {
	pgTx, err := d.DB.WithContext(ctx).Begin()
	if err != nil {
		return nil, err
	}

	var newTx DB = &tx{
		db: d,
		Tx: pgTx,
	}

	for _, option := range options {
		newTx = option(ctx, newTx)
	}

	return newTx, nil
}

func (d *db) Commit() error {
	return nil
}

func (d *db) Rollback() error {
	return nil
}

func (d *db) WithCtx(ctx context.Context) DB {
	return NewPGWithCtx(ctx, d)
}

package adapters

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"gorm.io/gorm"
)

type Option func(ctx context.Context, db DB) DB

type DB interface {
	Begin(ctx context.Context, options ...Option) (DB, error)
	Commit() error
	Rollback() error
	Close() error
	Model(...interface{}) *orm.Query
	Insert(...interface{}) error
	Update(interface{}) error
	Save(interface{}) error
	CreateTable(interface{}, *orm.CreateTableOptions) error
	DropTable(interface{}, *orm.DropTableOptions) error
	Exec(interface{}, ...interface{}) (pg.Result, error)
	WithCtx(context.Context) DB
	WithForUpdate(...string) DB
	WithForNoKeyUpdate(...string) DB
	WithSkipLocked() DB
	Stats() *pg.PoolStats
}

type GormDB interface {
	Open(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error)
}

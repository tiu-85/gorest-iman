package adapters

import (
	"fmt"
	"net"
	"sync"

	"github.com/go-pg/pg/v9"

	"tiu-85/gorest-iman/pkg/common/infra/values"
)

const (
	DefaultKind Kind = "default"
	PostKind    Kind = "post"

	queryStartTimeKey = "queryDuration"
)

type Kind string

type Pool interface {
	Get(Kind) (DB, error)
}

func NewPool(dbs map[string]*values.DbConfig, logger Logger) (Pool, error) {
	p := &pool{
		connections: make(map[Kind]DB),
		logger:      logger.Named("pool"),
		mtx:         new(sync.Mutex),
	}
	for key, cfg := range dbs {
		dbPool := pg.Connect(&pg.Options{
			Addr:                  net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port)),
			User:                  cfg.Username,
			Password:              cfg.Password,
			Database:              cfg.Database,
			PoolSize:              cfg.MaxOpenConns,
			MaxConnAge:            cfg.MaxConnLifetime,
			ReadTimeout:           cfg.ReadTimeout,
			WriteTimeout:          cfg.WriteTimeout,
			MaxRetries:            cfg.MaxRetries,
			MinIdleConns:          cfg.MinIdleConns,
			RetryStatementTimeout: true,
			ApplicationName:       key,
		})

		_, err := dbPool.Exec("SELECT 1")
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		dbPool.AddQueryHook(LogHook{Logger: logger.Named("query_logger")})
		p.connections[Kind(key)] = NewPG(dbPool)
	}

	return p, nil
}

type pool struct {
	connections map[Kind]DB
	mtx         *sync.Mutex
	logger      Logger
}

func (p *pool) Get(kind Kind) (DB, error) {
	logger := p.logger.With("method", "get")
	p.mtx.Lock()
	defer p.mtx.Unlock()
	conn, ok := p.connections[kind]
	if !ok {
		logger.Errorf("db kind=%s not found", kind)
		return nil, values.ErrInvalidKind
	}
	return conn, nil
}

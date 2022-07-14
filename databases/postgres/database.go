package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Db struct {
	pgpool *pgxpool.Pool
}

type Model interface {
	GetRowValues() RowValues
	GetManyRowValuesFlattened() RowValues
	GetManyRowValues() RowEntries
}

type RowValues []interface{}

type RowEntries struct {
	Rows []RowValues
}

var connStr string
var Pg Db

func (d *Db) pool(ctx context.Context) *pgxpool.Pool {
	if Pg.pgpool != nil {
		return Pg.pgpool
	}
	Pg.pgpool = d.InitPG(ctx, connStr)
	return d.InitPG(ctx, connStr)
}

func (d *Db) InitPG(ctx context.Context, connStr string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(err)
	}
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute

	c, err := pgxpool.Connect(ctx, config.ConnString())
	if err != nil {
		panic(err)
	}
	Pg.pgpool = c
	return Pg.pgpool
}

func (d *Db) QueryRow(ctx context.Context, query string) pgx.Row {
	return Pg.pool(ctx).QueryRow(ctx, query)
}

func (d *Db) Query(ctx context.Context, query string) (pgx.Rows, error) {
	return Pg.pool(ctx).Query(ctx, query)
}

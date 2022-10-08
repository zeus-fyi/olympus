package postgres_apps

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Db struct {
	Pgpool *pgxpool.Pool
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

var ConnStr string
var Pg Db

func (d *Db) InitPG(ctx context.Context, pgConnStr string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(pgConnStr)
	if err != nil {
		panic(err)
	}
	ConnStr = config.ConnString()
	c, err := pgxpool.Connect(ctx, ConnStr)
	if err != nil {
		panic(err)
	}
	Pg.Pgpool = c
	return Pg.Pgpool
}

func (d *Db) QueryRow(ctx context.Context, query string) pgx.Row {
	return Pg.Pgpool.QueryRow(ctx, query)
}

func (d *Db) Query(ctx context.Context, query string) (pgx.Rows, error) {
	return Pg.Pgpool.Query(ctx, query)
}

func (d *Db) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return Pg.Pgpool.Exec(ctx, query, args...)
}

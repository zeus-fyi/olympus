package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Db struct {
	pgpool *pgxpool.Pool
}

type Model interface {
	GetRowValues() RowValues
	GetManyRowValues() RowEntries
}

type RowValues []string

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
	c, err := pgxpool.Connect(ctx, connStr)
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

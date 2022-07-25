package postgres

import (
	"context"

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

var connStr string
var Pg Db

func (d *Db) InitPG(ctx context.Context, pgConnStr string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(pgConnStr)
	if err != nil {
		panic(err)
	}
	connStr = config.ConnString()
	c, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		panic(err)
	}
	Pg.Pgpool = c
	return Pg.Pgpool
}

func (d *Db) QueryRow(ctx context.Context, query string) pgx.Row {
	defer Pg.Pgpool.Close()
	return Pg.Pgpool.QueryRow(ctx, query)
}

func (d *Db) Query(ctx context.Context, query string) (pgx.Rows, error) {
	defer Pg.Pgpool.Close()
	return Pg.Pgpool.Query(ctx, query)
}

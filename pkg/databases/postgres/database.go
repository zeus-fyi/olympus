package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
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
	return Pg.Pgpool.QueryRow(ctx, query)
}

func (d *Db) Query(ctx context.Context, query string) (pgx.Rows, error) {
	return Pg.Pgpool.Query(ctx, query)
}

func (d *Db) PoolStats(ctx context.Context) *pgxpool.Stat {
	log.Ctx(ctx).Info().Msg("Getting Pool Stats")
	stats := Pg.Pgpool.Stat()
	return stats
}

func (d *Db) Ping(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("Pinging DB")
	err := Pg.Pgpool.Ping(ctx)
	log.Err(err).Msg("Pinging DB failed")
	return err
}

type ConfigChangePG struct {
	MinConn           *int32
	MaxConns          *int32
	MaxConnLifetime   *time.Duration
	HealthCheckPeriod *time.Duration
}

func UpdateConfigPG(ctx context.Context, cfg ConfigChangePG) error {
	log.Ctx(ctx).Debug().Msg("UpdateConfigPG")
	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(err)
	}

	if cfg.MinConn != nil {
		log.Info().Msgf("min conn updated. was %s, is now %s", dbConfig.MinConns, *cfg.MinConn)
		dbConfig.MinConns = *cfg.MinConn
	}

	if cfg.MaxConns != nil {
		log.Info().Msgf("max conn updated. was %s, is now %s", dbConfig.MaxConns, *cfg.MaxConns)
		dbConfig.MaxConns = *cfg.MaxConns
	}

	if cfg.MaxConnLifetime != nil {
		log.Info().Msgf("max conn lifetime updated. was %s, is now %s", dbConfig.MaxConnLifetime, *cfg.MaxConnLifetime)
		dbConfig.MaxConnLifetime = *cfg.MaxConnLifetime
	}

	if cfg.HealthCheckPeriod != nil {
		log.Info().Msgf("max conn lifetime updated. was %s, is now %s", dbConfig.HealthCheckPeriod, *cfg.HealthCheckPeriod)
		dbConfig.HealthCheckPeriod = *cfg.HealthCheckPeriod
	}

	connStr = dbConfig.ConnString()
	_ = Pg.InitPG(ctx, connStr)
	return nil
}

func ReadCfg(ctx context.Context) ConfigChangePG {
	log.Ctx(ctx).Debug().Msg("ReadCfg")
	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(err)
	}
	var cfg ConfigChangePG
	cfgCopy := dbConfig.Copy()
	if cfg.MinConn != nil {
		cfg.MinConn = &cfgCopy.MinConns
	}

	if cfg.MaxConns != nil {
		cfg.MaxConns = &cfgCopy.MaxConns
	}

	if cfg.MaxConnLifetime != nil {
		cfg.MaxConnLifetime = &cfgCopy.MaxConnLifetime
	}

	if cfg.HealthCheckPeriod != nil {
		cfg.HealthCheckPeriod = &cfgCopy.HealthCheckPeriod
	}
	return cfg
}

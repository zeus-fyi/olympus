package admin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type AdminDb struct {
	Pgpool *pgxpool.Pool
}

type ConfigChangePG struct {
	MinConn           *int32
	MaxConns          *int32
	MaxConnLifetime   *time.Duration
	HealthCheckPeriod *time.Duration
}

type ConfigReadPG struct {
	MinConn           int32
	MaxConns          int32
	MaxConnLifetime   time.Duration
	HealthCheckPeriod time.Duration
}

type PoolStats struct {
	TotalConns    int32
	IdleConns     int32
	AcquiredConns int32
	MaxConns      int32
}

func (a *AdminDb) PoolStats(ctx context.Context) PoolStats {
	log.Ctx(ctx).Info().Msg("Getting Pool Stats")

	var Pstat PoolStats
	stats := apps.Pg.Pgpool.Stat()
	if stats != nil {
		Pstat.TotalConns = stats.TotalConns()
		Pstat.IdleConns = stats.IdleConns()
		Pstat.AcquiredConns = stats.AcquiredConns()
		Pstat.MaxConns = stats.MaxConns()
	}

	return Pstat
}

func (a *AdminDb) Ping(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("Pinging DB")
	err := apps.Pg.Pgpool.Ping(ctx)
	log.Err(err).Msg("Pinging DB failed")
	return err
}

func UpdateConfigPG(ctx context.Context, cfg ConfigChangePG) error {
	log.Ctx(ctx).Debug().Msg("UpdateConfigPG")
	cfgCopy := apps.Pg.Pgpool.Config().Copy()
	if cfgCopy == nil {
		panic(errors.New("should be a connStr"))
	}
	dbConfig := *cfgCopy
	if cfg.MinConn != nil {
		log.Info().Msgf("min conn updated. was %d, is now %d", dbConfig.MinConns, *cfg.MinConn)
		dbConfig.MinConns = *cfg.MinConn
		dbConfig.ConnConfig.Config.RuntimeParams["pool_min_conns"] = fmt.Sprintf("%d", *cfg.MinConn)
	}

	if cfg.MaxConns != nil {
		log.Info().Msgf("max conn updated. was %d, is now %d", dbConfig.MaxConns, *cfg.MaxConns)
		dbConfig.MaxConns = *cfg.MaxConns
		dbConfig.ConnConfig.Config.RuntimeParams["pool_max_conns"] = fmt.Sprintf("%d", *cfg.MaxConns)
	}

	if cfg.MaxConnLifetime != nil {
		log.Info().Msgf("max conn lifetime updated. was %s, is now %s", dbConfig.MaxConnLifetime, *cfg.MaxConnLifetime)
		dbConfig.MaxConnLifetime = *cfg.MaxConnLifetime
	}

	if cfg.HealthCheckPeriod != nil {
		log.Info().Msgf("max conn lifetime updated. was %s, is now %s", dbConfig.HealthCheckPeriod, *cfg.HealthCheckPeriod)
		dbConfig.HealthCheckPeriod = *cfg.HealthCheckPeriod
	}

	apps.Pg.Pgpool.Close()
	apps.ConnStr = dbConfig.ConnString()
	apps.Pg.Pgpool = apps.Pg.InitPG(ctx, apps.ConnStr)
	return nil
}

func ReadCfg(ctx context.Context) ConfigReadPG {
	log.Ctx(ctx).Debug().Msg("ReadCfg")
	dbConf, err := pgxpool.ParseConfig(apps.ConnStr)
	if err != nil {
		panic(err)
	}
	if dbConf == nil {
		panic(errors.New("should be a connStr"))
	}
	dbConfig := *dbConf
	var cfg ConfigReadPG
	cfg.MinConn = dbConfig.MinConns
	cfg.MaxConns = dbConfig.MaxConns
	cfg.MaxConnLifetime = dbConfig.MaxConnLifetime
	cfg.HealthCheckPeriod = dbConfig.HealthCheckPeriod
	return cfg
}

func (a *AdminDb) FetchTableSize(ctx context.Context, tableName string) (string, error) {
	log.Ctx(ctx).Info().Msgf("FetchTableSize Table: %s", tableName)
	var tableSize string
	query := fmt.Sprintf(`SELECT pg_size_pretty(pg_total_relation_size('%s'))`, tableName)
	err := apps.Pg.Pgpool.QueryRow(ctx, query).Scan(&tableSize)
	log.Err(err).Msgf("FetchTableSize DB failed with response: %s", tableSize)
	return tableSize, err
}

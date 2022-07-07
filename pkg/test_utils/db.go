package test_utils

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"bitbucket.org/zeus/eth-indexer/databases/postgres"
)

func InitLocalTestDBConn() *pgxpool.Pool {
	conn, err := postgres.InitPG("TODO")
	if err != nil {
		panic(err)
	}
	return conn
}

package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// InitPG Make sure you defer dbClose
func InitPG(connStr string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

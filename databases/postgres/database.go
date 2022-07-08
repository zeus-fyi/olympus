package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// InitPG Make sure you defer dbClose
func InitPG(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

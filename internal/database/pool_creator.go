package database

import (
	"context"
	"github.com/dimitryshirokov/simple-app/internal/config"
	"github.com/dimitryshirokov/simple-app/internal/internal_error"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePool(ctx context.Context, c *config.Config) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(c.DbUrl)
	if err != nil {
		return nil, internal_error.NewError("can't parse database url", err, nil)
	}
	dbConfig.MinConns = int32(c.DbMinConnections)
	dbConfig.MaxConns = int32(c.DbMaxConnections)
	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, internal_error.NewError("can't create database connection pool", err, nil)
	}
	return pool, nil
}

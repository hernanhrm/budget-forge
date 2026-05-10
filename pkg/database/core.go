package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	shared_domain "github.com/hernanhrm/budget-forge/pkg/shared_domain"
)

const sqlPreviewMaxLen = 100

type Database struct {
	Pool   PoolInterface
	logger shared_domain.Logger
}

func NewConnection(ctx context.Context, connString string, log shared_domain.Logger) (*Database, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("db_pool_creation_failed: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db_ping_failed: %w", err)
	}

	db := &Database{
		Pool:   pool,
		logger: log.With("component", "database"),
	}

	poolConfig := pool.Config()
	db.logger.WithContext(ctx).Info("database connection established",
		"max_connections", poolConfig.MaxConns,
		"min_connections", poolConfig.MinConns,
	)

	return db, nil
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		db.logger.Info("database connection pool closed")
	}
}

func (db *Database) GetPool() PoolInterface {
	return db.Pool
}

func (db *Database) PgxPool() (*pgxpool.Pool, bool) {
	p, ok := db.Pool.(*pgxpool.Pool)
	return p, ok
}

func (db *Database) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	rows, err := db.Pool.Query(ctx, sql, args...)
	if err != nil {
		db.logger.WithContext(ctx).Error("query execution failed",
			"error", err,
			"sql_preview", truncateSQL(sql),
		)
		return nil, fmt.Errorf("db_query_failed: %w", err)
	}
	return rows, nil
}

func (db *Database) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}

func (db *Database) QueryRowScan(ctx context.Context, scanFunc func(row pgx.Row) error, sql string, args ...any) error {
	row := db.Pool.QueryRow(ctx, sql, args...)
	if err := scanFunc(row); err != nil {
		db.logger.WithContext(ctx).Error("query row scan failed",
			"error", err,
			"sql_preview", truncateSQL(sql),
		)
		return fmt.Errorf("db_query_row_scan_failed: %w", err)
	}

	return nil
}

func (db *Database) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	tag, err := db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		db.logger.WithContext(ctx).Error("exec operation failed",
			"error", err,
			"sql_preview", truncateSQL(sql),
		)
		return tag, fmt.Errorf("db_exec_failed: %w", err)
	}

	db.logger.WithContext(ctx).Debug("exec operation completed",
		"rows_affected", tag.RowsAffected(),
	)

	return tag, nil
}

func (db *Database) HealthCheck(ctx context.Context) error {
	if db.Pool == nil {
		return fmt.Errorf("db_health_check_failed: pool is nil")
	}
	if err := db.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("db_health_check_failed: %w", err)
	}
	return nil
}

func (db *Database) Shutdown(ctx context.Context) error {
	db.logger.WithContext(ctx).Info("shutting down database connection")
	if db.Pool != nil {
		db.Pool.Close()
	}
	return nil
}

func truncateSQL(sql string) string {
	if len(sql) <= sqlPreviewMaxLen {
		return sql
	}
	return sql[:sqlPreviewMaxLen] + "..."
}

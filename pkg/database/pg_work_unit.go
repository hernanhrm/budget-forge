package database

import (
	"context"
	"fmt"

	shared_domain "github.com/hernanhrm/budget-forge/pkg/shared_domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgWorkUnit struct {
	pool *pgxpool.Pool
}

func NewPgWorkUnit(pool *pgxpool.Pool) *PgWorkUnit {
	return &PgWorkUnit{pool: pool}
}

func (w *PgWorkUnit) Run(ctx context.Context, fn func(context.Context) error) (err error) {
	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	err = fn(contextWithTx(ctx, tx))
	return
}

var _ shared_domain.WorkUnit = (*PgWorkUnit)(nil)

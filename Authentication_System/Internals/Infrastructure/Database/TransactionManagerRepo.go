package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type typeKey string

const txKey typeKey = "tx"

type pgxTxManager struct {
	db *pgxpool.Pool
}

func NewPgxTransactionManager(db *pgxpool.Pool) *pgxTxManager {
	return &pgxTxManager{db: db}
}

func (m *pgxTxManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, txKey, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction err: %v, rollback err: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
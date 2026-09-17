package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresReceptionstRepository struct {
	db *pgxpool.Pool
}

func NewPostgreReceptionstRepo(db *pgxpool.Pool) *PostgresReceptionstRepository {
	return &PostgresReceptionstRepository{
		db: db,
	}
}

func (r *PostgresReceptionstRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM users 
			WHERE email = $1
		);
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}
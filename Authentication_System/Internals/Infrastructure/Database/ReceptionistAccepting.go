package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) CreateReceptionistUser(
	ctx context.Context,
	email string,
	hashedPassword string,
	fullName string,
	doctorID string,
) (string, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	userID := uuid.New().String()
	receptionistID := uuid.New().String()
	now := time.Now().UTC()

	insertUserQuery := `
		INSERT INTO users (
			id, name, email, password, role, gender, is_verified, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`

	_, err = tx.Exec(
		ctx,
		insertUserQuery,
		userID,       
		fullName,     
		email,       
		hashedPassword, 
		"RECEPTIONIST",
		"Male",         
		true,          
		now,         
		now,            
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert user record: %w", err)
	}

	insertReceptionistQuery := `
		INSERT INTO receptionists (
			id, user_id, doctor_id, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5);
	`

	_, err = tx.Exec(
		ctx,
		insertReceptionistQuery,
		receptionistID,
		userID,
		doctorID,
		now,
		now,
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert receptionist record: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return receptionistID, nil
}  
package database


import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)


type VerificationTokenModel struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type emailSaveRepository struct {
	db *pgxpool.Pool
}

func NewEmailSaveRepository(db *pgxpool.Pool) *emailSaveRepository {
	return &emailSaveRepository{db: db}
}

func (r *emailSaveRepository) SaveToken(ctx context.Context, userID string, token string, expiresAt time.Time) error {
	exec := getExecutor(ctx, r.db)

	query := `
		INSERT INTO email_verifications (id, user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := exec.Exec(ctx, query,
		uuid.New().String(),
		userID,
		token,
		expiresAt,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert verification token: %w", err)
	}

	return nil
}
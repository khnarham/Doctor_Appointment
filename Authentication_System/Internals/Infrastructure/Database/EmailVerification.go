package database

import (
	"context"
	"errors"
	"fmt"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"my-appointment-app/Authentication_System/Internals/Usecase" // Ensure path matches your actual package
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmailVerificationModel struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type emailVerificationRepository struct {
	db *pgxpool.Pool
}

func NewEmailVerificationRepository(db *pgxpool.Pool) *emailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

// 2. FindByToken DB Query
func (r *emailVerificationRepository) FindByToken(ctx context.Context, token string) (*usecase.EmailVerificationData, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM email_verifications
		WHERE token = $1
		LIMIT 1
	`

	row := exec.QueryRow(ctx, query, token)

	var model EmailVerificationModel
	err := row.Scan(
		&model.ID,
		&model.UserID,
		&model.Token,
		&model.ExpiresAt,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to scan email verification row: %w", err)
	}

	return &usecase.EmailVerificationData{
		ID:        model.ID,
		UserID:    model.UserID,
		Token:     model.Token,
		ExpiresAt: model.ExpiresAt,
	}, nil
}

func (r *emailVerificationRepository) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, name, email, password, role, gender, is_verified, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	row := exec.QueryRow(ctx, query, userID)
	return mapRowToUserVerificaiton(row)
}

func (r *emailVerificationRepository) Save(ctx context.Context, user *domain.User) error {
	exec := getExecutor(ctx, r.db)

	query := `
		INSERT INTO users (id, name, email, password, role, gender, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			password = EXCLUDED.password,
			role = EXCLUDED.role,
			gender = EXCLUDED.gender,
			is_verified = EXCLUDED.is_verified,
			updated_at = EXCLUDED.updated_at
	`

	_, err := exec.Exec(ctx, query,
		user.ID(),
		user.Name(),
		user.Email().Value(),
		user.PasswordHash(),
		string(user.Role()),
		string(user.Gender()),
		user.IsVerified(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save user record: %w", err)
	}

	return nil
}

func (r *emailVerificationRepository) Delete(ctx context.Context, tokenID string) error {
	exec := getExecutor(ctx, r.db)

	query := `
		DELETE FROM email_verifications
		WHERE id = $1
	`

	_, err := exec.Exec(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to delete verification token: %w", err)
	}

	return nil
}

func mapRowToUserVerificaiton(row pgx.Row) (*domain.User, error) {
	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Name,
		&userModel.Email,
		&userModel.Password,
		&userModel.Role,
		&userModel.Gender,
		&userModel.IsVerified,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan user row: %w", err)
	}

	email, err := domain.NewEmail(userModel.Email)
	if err != nil {
		return nil, err
	}

	return domain.ReconstituteUser(
		userModel.ID,
		userModel.Name,
		email,
		userModel.Password,
		userModel.IsVerified,
		domain.Role(userModel.Role),
		domain.Gender(userModel.Gender),
		userModel.CreatedAt,
		userModel.UpdatedAt,
	)
}
package database

import (
	"context"
	"errors"
	"fmt"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Role       string    `db:"role"`
	Gender     string    `db:"gender"`
	IsVerified bool      `db:"is_verified"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type userRegisterRepository struct {
	db *pgxpool.Pool
}

func NewUserRegisterRepository(db *pgxpool.Pool) *userRegisterRepository {
	return &userRegisterRepository{db: db}
}

func (r *userRegisterRepository) Save(ctx context.Context, user *domain.User) error {
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

func (r *userRegisterRepository) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, name, email, password, role, gender, is_verified, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	row := exec.QueryRow(ctx, query, email.Value())
	return mapRowToUser(row, email)
}

func mapRowToUser(row pgx.Row, email domain.Email) (*domain.User, error) {
	var model UserModel

	err := row.Scan(
		&model.ID,
		&model.Name,
		&model.Email,
		&model.Password,
		&model.Role,
		&model.Gender,
		&model.IsVerified,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan user row: %w", err)
	}

	userEntity, err := domain.ReconstituteUser(
		model.ID,
		model.Name,
		email,
		model.Password,
		model.IsVerified,
		domain.Role(model.Role),
		domain.Gender(model.Gender),
		model.CreatedAt,
		model.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstitute user entity: %w", err)
	}

	return userEntity, nil
}
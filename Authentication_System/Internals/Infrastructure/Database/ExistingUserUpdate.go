package database

import (
	"context"
	"fmt"
	domain "my-appointment-app/Authentication_System/Internals/Domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExistingUserRegistrationRepo struct {
	db *pgxpool.Pool
}

func NewExistingUserRegistrationRepo(db *pgxpool.Pool) *ExistingUserRegistrationRepo {
	return &ExistingUserRegistrationRepo{db: db}
}

func (r *ExistingUserRegistrationRepo) Update(ctx context.Context, user *domain.User) error {
	exec := getExecutor(ctx, r.db)

	query := `
		UPDATE users
		SET 
			name = $1,
			email = $2,
			password = $3,
			role = $4,
			gender = $5,
			is_verified = $6,
			updated_at = $7
		WHERE id = $8
	`

	_, err := exec.Exec(ctx, query,
		user.Name(),
		user.Email().Value(),
		user.PasswordHash(),
		string(user.Role()),
		string(user.Gender()),
		user.IsVerified(),
		user.UpdatedAt(),
		user.ID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update user record: %w", err)
	}

	return nil
}
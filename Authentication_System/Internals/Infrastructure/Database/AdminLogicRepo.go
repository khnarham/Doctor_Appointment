package database

import (
	"context"
	"errors"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresAdminRepository struct {
	db *pgxpool.Pool
}

func NewPostgresAdminRepository(db *pgxpool.Pool) *postgresAdminRepository {
	return &postgresAdminRepository{db: db}
}

func (r *postgresAdminRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, role, gender, is_verified, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`

	var (
		id           string
		name         string
		userEmail    string
		passwordHash string
		roleStr      string
		gender       string
		isVerified   bool
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&id,
		&name,
		&userEmail,
		&passwordHash,
		&roleStr,
		&gender,
		&isVerified,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
    
	newEmail , err := domain.NewEmail(userEmail);
	if err !=nil {
		return nil , err
	}
	
	// Reconstitute domain entity safely via factory function
	user, err := domain.ReconstituteUser(
		id,
		name,
		newEmail,
		passwordHash,
		isVerified,
		domain.Role(roleStr),
		domain.Gender(gender),
		createdAt,
		updatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
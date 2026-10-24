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

type DoctorInvitationModel struct {
	ID               string    `db:"id"`
	Email            string    `db:"email"`
	SpecializationID string    `db:"specialization_id"`
	Token            string    `db:"token"`
	IsUsed           bool      `db:"is_used"`
	ExpiresAt        time.Time `db:"expires_at"`
	CreatedAt        time.Time `db:"created_at"`
}

type doctorInvitationRepository struct {
	db *pgxpool.Pool
}

func NewDoctorInvitationRepository(db *pgxpool.Pool) *doctorInvitationRepository {
	return &doctorInvitationRepository{db: db}
}

func (r *doctorInvitationRepository) Save(ctx context.Context, invite *domain.DoctorInvitation) error {
	exec := getExecutor(ctx, r.db)

	query := `
		INSERT INTO doctor_invitations (id, email, specialization_id, token, is_used, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			specialization_id = EXCLUDED.specialization_id,
			token = EXCLUDED.token,
			is_used = EXCLUDED.is_used,
			expires_at = EXCLUDED.expires_at
	`

	_, err := exec.Exec(ctx, query,
		invite.ID(),
		invite.Email().Value(),
		invite.SpecializationID(),
		invite.Token(),
		invite.IsUsed(),
		invite.ExpiresAt(),
		invite.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save doctor invitation: %w", err)
	}

	return nil
}

func (r *doctorInvitationRepository) FindByToken(ctx context.Context, token string) (*domain.DoctorInvitation, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, email, specialization_id, token, is_used, expires_at, created_at
		FROM doctor_invitations
		WHERE token = $1
		LIMIT 1
	`

	row := exec.QueryRow(ctx, query, token)
	return mapRowToDoctorInvitation(row)
}

func (r *doctorInvitationRepository) FindByEmail(ctx context.Context, email domain.Email) (*domain.DoctorInvitation, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, email, specialization_id, token, is_used, expires_at, created_at
		FROM doctor_invitations
		WHERE email = $1
		LIMIT 1
	`

	row := exec.QueryRow(ctx, query, email.Value())
	return mapRowToDoctorInvitation(row)
}

func mapRowToDoctorInvitation(row pgx.Row) (*domain.DoctorInvitation, error) {
	var model DoctorInvitationModel

	err := row.Scan(
		&model.ID,
		&model.Email,
		&model.SpecializationID,
		&model.Token,
		&model.IsUsed,
		&model.ExpiresAt,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan doctor invitation row: %w", err)
	}

	email, err := domain.NewEmail(model.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email in database: %w", err)
	}

	return domain.ReconstituteDoctorInvitation(
		model.ID,
		email,
		model.SpecializationID,
		model.Token,
		model.IsUsed,
		model.ExpiresAt,
		model.CreatedAt,
	), nil
}
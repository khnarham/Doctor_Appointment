package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "my-appointment-app/Authentication_System/Internals/Domain"
)

var (
	ErrInvitationNotFound = errors.New("receptionist invitation not found")
)

type PostgresReceptionistInvitationRepository struct {
	db *pgxpool.Pool
}

func NewPostgresReceptionistInvitationRepository(db *pgxpool.Pool) *PostgresReceptionistInvitationRepository {
	return &PostgresReceptionistInvitationRepository{
		db: db,
	}
}

func (r *PostgresReceptionistInvitationRepository) Save(ctx context.Context, inv *domain.ReceptionistInvitation) error {
	query := `
		INSERT INTO receptionist_invitations (
			id, creater_id, email, token, status, expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	_, err := r.db.Exec(
		ctx,
		query,
		inv.ID(),
		inv.CreaterID(),
		inv.Email(),
		inv.Token(),
		string(inv.Status()),
		inv.ExpiresAt(),
		inv.CreatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert receptionist invitation: %w", err)
	}

	return nil
}

func (r *PostgresReceptionistInvitationRepository) GetByToken(ctx context.Context, token string) (*domain.ReceptionistInvitation, error) {
	query := `
		SELECT id, creater_id, email, token, status, expires_at, created_at
		FROM receptionist_invitations
		WHERE token = $1;
	`

	var (
		id        string
		createrID string
		email     string
		retToken  string
		statusStr string
		expiresAt time.Time
		createdAt time.Time
	)

	err := r.db.QueryRow(ctx, query, token).Scan(
		&id,
		&createrID,
		&email,
		&retToken,
		&statusStr,
		&expiresAt,
		&createdAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, fmt.Errorf("failed to scan receptionist invitation by token: %w", err)
	}

	status := domain.InvitationStatus(statusStr)

	invitation := domain.ReconstituteReceptionistInvitation(
		id,
		createrID,
		email,
		retToken,
		status,
		expiresAt,
		createdAt,
	)

	return invitation, nil
}

func (r *PostgresReceptionistInvitationRepository) Update(ctx context.Context, inv *domain.ReceptionistInvitation) error {
	query := `
		UPDATE receptionist_invitations
		SET status = $1, expires_at = $2
		WHERE id = $3;
	`

	cmdTag, err := r.db.Exec(
		ctx,
		query,
		string(inv.Status()),
		inv.ExpiresAt(),
		inv.ID(),
	)

	if err != nil {
		return fmt.Errorf("failed to update receptionist invitation: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrInvitationNotFound
	}

	return nil
}
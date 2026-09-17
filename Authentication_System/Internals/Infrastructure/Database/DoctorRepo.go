package database

import (
	"context"
	"errors"
	"fmt"
	domain "my-appointment-app/Authentication_System/Internals/Domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)




type GetDoctorRepository struct {
	db *pgxpool.Pool
}

func NewGetDoctorRepository(db *pgxpool.Pool) *GetDoctorRepository {
	return &GetDoctorRepository{db: db}
}

func (r *GetDoctorRepository) FindByUserID(ctx context.Context, userID string) (*domain.Doctor, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, user_id, specialization_id, license_number, consultation_fee, status, is_available, created_at, updated_at
		FROM doctors
		WHERE user_id = $1
		LIMIT 1
	`

	var docModel DoctorModel

	err := exec.QueryRow(ctx, query, userID).Scan(
		&docModel.ID,
		&docModel.UserID,
		&docModel.SpecializationID,
		&docModel.LicenseNumber,
		&docModel.ConsultationFee,
		&docModel.Status,
		&docModel.IsAvailable,
		&docModel.CreatedAt,
		&docModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to scan doctor row by user_id: %w", err)
	}

	doctorEntity := domain.ReconstituteDoctor(
		docModel.ID,
		docModel.UserID,
		docModel.SpecializationID,
		docModel.LicenseNumber,
		docModel.ConsultationFee,
		domain.DoctorStatus(docModel.Status),
		docModel.IsAvailable,
		docModel.CreatedAt,
		docModel.UpdatedAt,
	)

	return doctorEntity, nil
}
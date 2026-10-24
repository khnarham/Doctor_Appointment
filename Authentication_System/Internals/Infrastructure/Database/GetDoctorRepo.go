package database

import (
	"context"
	"fmt"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	usecase "my-appointment-app/Authentication_System/Internals/Usecase"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)




type doctorPgxRepository struct {
	db *pgxpool.Pool
}

func NewGetDoctorsRepository(db *pgxpool.Pool) *doctorPgxRepository {
	return &doctorPgxRepository{db: db}
}

func (r *doctorPgxRepository) GetDoctors(ctx context.Context, filter usecase.DoctorFilter) ([]*domain.Doctor, error) {
	var builder strings.Builder
	builder.WriteString(`
		SELECT id, user_id, specialization_id, license_number, consultation_fee, status, is_available, created_at, updated_at
		FROM doctors
		WHERE 1=1
	`)

	args := make([]interface{}, 0, 5)
	argID := 1

	if filter.SpecializationID != "" {
		builder.WriteString(" AND specialization_id = $")
		builder.WriteString(strconv.Itoa(argID))
		args = append(args, filter.SpecializationID)
		argID++
	}

	if filter.MaxFee > 0 {
		builder.WriteString(" AND consultation_fee <= $")
		builder.WriteString(strconv.Itoa(argID))
		args = append(args, filter.MaxFee)
		argID++
	}

	if filter.Search != "" {
		builder.WriteString(" AND license_number ILIKE $")
		builder.WriteString(strconv.Itoa(argID))
		args = append(args, "%"+filter.Search+"%")
		argID++
	}

	builder.WriteString(" ORDER BY created_at DESC")

	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	builder.WriteString(" LIMIT $")
	builder.WriteString(strconv.Itoa(argID))
	args = append(args, filter.Limit)
	argID++

	if filter.Offset > 0 {
		builder.WriteString(" OFFSET $")
		builder.WriteString(strconv.Itoa(argID))
		args = append(args, filter.Offset)
	}

	rows, err := r.db.Query(ctx, builder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query doctors: %w", err)
	}
	defer rows.Close()

	var doctors []*domain.Doctor
	for rows.Next() {
		var (
			id               string
			userID           string
			specializationID string
			licenseNumber    string
			consultationFee  float64
			status           domain.DoctorStatus
			isAvailable      bool
			createdAt        time.Time
			updatedAt        time.Time
		)

		err := rows.Scan(
			&id,
			&userID,
			&specializationID,
			&licenseNumber,
			&consultationFee,
			&status,
			&isAvailable,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan doctor row: %w", err)
		}

		doc := domain.ReconstituteDoctor(
			id,
			userID,
			specializationID,
			licenseNumber,
			consultationFee,
			status,
			isAvailable,
			createdAt,
			updatedAt,
		)
		doctors = append(doctors, doc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return doctors, nil
}
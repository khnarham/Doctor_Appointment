package database

import (
	"context"
	"fmt"
	domain "my-appointment-app/Authentication_System/Internals/Domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type patientUpdateRepository struct {
	db *pgxpool.Pool
}

func NewPatientUpdateRepository(db *pgxpool.Pool) *patientUpdateRepository {
	return &patientUpdateRepository{db: db}
}

func (r *patientUpdateRepository) Update(ctx context.Context, patient *domain.Patient) error {
	exec := getExecutor(ctx, r.db)

	query := `
		UPDATE patients
		SET 
			user_id = $1,
			walk_in_name = $2,
			walk_in_phone = $3,
			gender = $4,
			date_of_birth = $5,
			is_walk_in = $6
		WHERE id = $7
	`

	var walkInName *string
	if name := patient.WalkInName(); name != "" {
		walkInName = &name
	}

	_, err := exec.Exec(ctx, query,
		patient.UserID(),
		walkInName,
		patient.WalkInPhone(),
		string(patient.Gender()),
		patient.DateOfBirth(),
		patient.IsWalkIn(),
		patient.ID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update patient record: %w", err)
	}

	return nil
}
package database

import (
	"context"
	"fmt"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DoctorModel struct {
	ID               string    `db:"id"`
	UserID           string    `db:"user_id"`
	SpecializationID string    `db:"specialization_id"`
	LicenseNumber    string    `db:"license_number"`
	ConsultationFee  float64   `db:"consultation_fee"`
	Status           string    `db:"status"`
	IsAvailable      bool      `db:"is_available"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}


type doctorRepository struct {
	db *pgxpool.Pool
}

func NewDoctorRepository(db *pgxpool.Pool) *doctorRepository {
	return &doctorRepository{db: db}
}

func (r *doctorRepository) Save(ctx context.Context, doctor *domain.Doctor) error {
	exec := getExecutor(ctx, r.db)

	query := `
		INSERT INTO doctors (
			id, 
			user_id, 
			specialization_id, 
			license_number, 
			consultation_fee, 
			status, 
			is_available, 
			created_at, 
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			specialization_id = EXCLUDED.specialization_id,
			license_number = EXCLUDED.license_number,
			consultation_fee = EXCLUDED.consultation_fee,
			status = EXCLUDED.status,
			is_available = EXCLUDED.is_available,
			updated_at = EXCLUDED.updated_at
	`

	_, err := exec.Exec(ctx, query,
		doctor.ID(),
		doctor.UserID(),
		doctor.SpecializationID(),
		doctor.LicenseNumber(),
		doctor.ConsultationFee(),
		string(doctor.Status()),
		doctor.IsAvailable(),
		doctor.CreatedAt(),
		doctor.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save doctor record: %w", err)
	}

	return nil
}
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
)


type reviewRepository struct {
	db *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) *reviewRepository {
	return &reviewRepository{db: db}
}
func (r *reviewRepository) ExistsByAppointmentID(ctx context.Context, appointmentID string) (bool, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM reviews 
			WHERE appointment_id = $1
		);
	`

	var exists bool
	err := exec.QueryRow(ctx, query, appointmentID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check review existence by appointment_id: %w", err)
	}

	return exists, nil
}

func (r *reviewRepository) Save(ctx context.Context, review *domain.Review) error {
	if review == nil {
		return fmt.Errorf("cannot save nil review domain entity")
	}

	exec := getExecutor(ctx, r.db)

	query := `
		INSERT INTO reviews (
			id, 
			appointment_id, 
			patient_id, 
			doctor_id, 
			rating, 
			comment, 
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	_, err := exec.Exec(
		ctx,
		query,
		review.ID(),
		review.AppointmentID(),
		review.PatientID(),
		review.DoctorID(),
		review.Rating(),
		review.Comment(),
		review.CreatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert review record: %w", err)
	}

	return nil
}
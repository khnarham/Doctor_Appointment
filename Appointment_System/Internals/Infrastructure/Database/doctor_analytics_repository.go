package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
)

type DoctorAnalyticsRepository interface {
	GetDoctorMetrics(ctx context.Context, doctorID string) (*domain.DoctorAnalyticsSummary, error)
	GetDoctorPatients(ctx context.Context, doctorID string, limit int) ([]domain.DoctorPatientDetail, error)
}

type doctorAnalyticsRepo struct {
	pool *pgxpool.Pool
}

func NewDoctorAnalyticsRepository(pool *pgxpool.Pool) DoctorAnalyticsRepository {
	return &doctorAnalyticsRepo{pool: pool}
}

func (r *doctorAnalyticsRepo) GetDoctorMetrics(ctx context.Context, doctorID string) (*domain.DoctorAnalyticsSummary, error) {
	query := `
		SELECT 
			COALESCE(COUNT(a.id), 0) AS total_appointments,
			COALESCE(COUNT(CASE WHEN a.status = 'COMPLETED' THEN 1 END), 0) AS completed_appointments,
			COALESCE(COUNT(CASE WHEN a.status = 'PENDING' THEN 1 END), 0) AS pending_appointments,
			COALESCE(COUNT(CASE WHEN a.status = 'CANCELLED' THEN 1 END), 0) AS cancelled_appointments,
			COALESCE(COUNT(DISTINCT a.patient_id), 0) AS unique_patients_count,
			COALESCE(
				COUNT(CASE WHEN a.status = 'COMPLETED' THEN 1 END) * COALESCE(d.consultation_fee, 0), 
				0
			)::DOUBLE PRECISION AS total_earnings
		FROM doctors d
		LEFT JOIN appointments a ON d.id = a.doctor_id
		WHERE d.id = $1
		GROUP BY d.id, d.consultation_fee;
	`

	var summary domain.DoctorAnalyticsSummary
	err := r.pool.QueryRow(ctx, query, doctorID).Scan(
		&summary.TotalAppointments,
		&summary.CompletedAppointments,
		&summary.PendingAppointments,
		&summary.CancelledAppointments,
		&summary.UniquePatientsCount,
		&summary.TotalEarnings,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.DoctorAnalyticsSummary{}, nil
		}
		return nil, err
	}

	return &summary, nil
}

func (r *doctorAnalyticsRepo) GetDoctorPatients(ctx context.Context, doctorID string, limit int) ([]domain.DoctorPatientDetail, error) {
	// Query aligned with actual 'users' table schema (u.name):
	query := `
		SELECT 
			a.id AS appointment_id,
			p.id AS patient_id,
			CASE 
				WHEN p.is_walk_in = true THEN COALESCE(NULLIF(p.walk_in_name, ''), 'Walk-In Patient')
				ELSE COALESCE(u.name, 'Registered Patient')
			END AS patient_name,
			CASE 
				WHEN p.is_walk_in = true THEN p.walk_in_phone
				ELSE ''
			END AS patient_phone,
			p.is_walk_in,
			p.gender,
			p.date_of_birth,
			a.appointment_date,
			a.start_time,
			a.end_time,
			a.status
		FROM appointments a
		INNER JOIN patients p ON a.patient_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		WHERE a.doctor_id = $1
		ORDER BY a.appointment_date DESC, a.start_time DESC
		LIMIT $2;
	`

	rows, err := r.pool.Query(ctx, query, doctorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patients, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.DoctorPatientDetail])
	if err != nil {
		return nil, err
	}

	return patients, nil
}
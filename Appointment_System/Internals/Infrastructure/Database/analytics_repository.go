package database

import (
	"context"
	"fmt"
	domain "my-appointment-app/Appointment_System/Internals/Domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepository interface {
	GetTotalRevenueStats(ctx context.Context, filter domain.FinancialFilterDTO) (int64, float64, error)
	GetTopEarningDoctors(ctx context.Context, filter domain.FinancialFilterDTO) ([]domain.TopEarningDoctorDTO, error)
}

type postgresAnalyticsRepository struct {
	db *pgxpool.Pool
}

func NewPostgresAnalyticsRepository(db *pgxpool.Pool) AnalyticsRepository {
	return &postgresAnalyticsRepository{db: db}
}

// GetTotalRevenueStats Total COMPLETED appointments count aur Total Revenue calculate karta hai
func (r *postgresAnalyticsRepository) GetTotalRevenueStats(ctx context.Context, filter domain.FinancialFilterDTO) (int64, float64, error) {
	query := `
		SELECT 
			COUNT(a.id) AS total_completed,
			COALESCE(SUM(d.consultation_fee), 0.00) AS total_revenue
		FROM appointments a
		INNER JOIN doctors d ON a.doctor_id = d.id
		WHERE a.status = 'COMPLETED'
	`

	args := make([]interface{}, 0)
	argID := 1

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND a.appointment_date >= $%d", argID)
		args = append(args, *filter.StartDate)
		argID++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND a.appointment_date <= $%d", argID)
		args = append(args, *filter.EndDate)
		argID++
	}

	var totalCompleted int64
	var totalRevenue float64

	err := r.db.QueryRow(ctx, query, args...).Scan(&totalCompleted, &totalRevenue)
	if err != nil {
		return 0, 0, fmt.Errorf("analytics_repo.GetTotalRevenueStats: %w", err)
	}

	return totalCompleted, totalRevenue, nil
}

// GetTopEarningDoctors Sab se zyada earn karne waale doctors fetch karta hai (using pgx.CollectRows)
func (r *postgresAnalyticsRepository) GetTopEarningDoctors(ctx context.Context, filter domain.FinancialFilterDTO) ([]domain.TopEarningDoctorDTO, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 5 // Default top 5 doctors
	}

	query := `
		SELECT 
			d.id AS doctor_id,
			d.user_id,
			d.license_number,
			d.consultation_fee,
			COUNT(a.id) AS completed_appointments,
			COALESCE(SUM(d.consultation_fee), 0.00) AS total_earned
		FROM appointments a
		INNER JOIN doctors d ON a.doctor_id = d.id
		WHERE a.status = 'COMPLETED'
	`

	args := make([]interface{}, 0)
	argID := 1

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND a.appointment_date >= $%d", argID)
		args = append(args, *filter.StartDate)
		argID++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND a.appointment_date <= $%d", argID)
		args = append(args, *filter.EndDate)
		argID++
	}

	query += fmt.Sprintf(`
		GROUP BY d.id, d.user_id, d.license_number, d.consultation_fee
		ORDER BY total_earned DESC, completed_appointments DESC
		LIMIT $%d`, argID)

	args = append(args, limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("analytics_repo.GetTopEarningDoctors query failed: %w", err)
	}
	defer rows.Close()

	// pgx.RowToStructByPos se automatic mapping
	topDoctors, err := pgx.CollectRows(rows, pgx.RowToStructByPos[domain.TopEarningDoctorDTO])
	if err != nil {
		return nil, fmt.Errorf("analytics_repo.GetTopEarningDoctors scan failed: %w", err)
	}

	return topDoctors, nil
}
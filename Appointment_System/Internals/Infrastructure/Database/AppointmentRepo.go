package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxAppointmentRepository struct {
	db *pgxpool.Pool
}

func NewAppointmentRepository(db *pgxpool.Pool) *pgxAppointmentRepository {
	return &pgxAppointmentRepository{db: db}
}

func (r *pgxAppointmentRepository) GetBookedSlotTimesByDoctorAndDate(
	ctx context.Context,
	doctorID string,
	date time.Time,
) ([]string, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT start_time, end_time
		FROM appointments
		WHERE doctor_id = $1 
		  AND appointment_date = $2 
		  AND status NOT IN ('CANCELLED')
		ORDER BY start_time ASC
	`
	formattedDate := date.Format("2006-01-02")

	rows, err := exec.Query(ctx, query, doctorID, formattedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch booked slots: %w", err)
	}
	defer rows.Close()

	baseDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var bookedSlots []string

	for rows.Next() {
		var startTimeMin, endTimeMin int

		if err := rows.Scan(&startTimeMin, &endTimeMin); err != nil {
			return nil, fmt.Errorf("failed to scan slot row: %w", err)
		}

		startBoundary := baseDate.Add(time.Duration(startTimeMin) * time.Minute)
		endBoundary := baseDate.Add(time.Duration(endTimeMin) * time.Minute)

		slotRange := fmt.Sprintf("%s-%s", 
			startBoundary.Format("15:04"), 
			endBoundary.Format("15:04"),
		)

		bookedSlots = append(bookedSlots, slotRange)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return bookedSlots, nil
}
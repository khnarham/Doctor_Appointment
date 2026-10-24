package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	domain "my-appointment-app/Appointment_System/Internals/Domain"
)

type appointmentRepository struct {
	db *pgxpool.Pool
}

func NewDoctorAppointmentRepository(db *pgxpool.Pool) *appointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) GetByDoctorID(ctx context.Context, doctorID string) ([]*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, schedule_id, appointment_date, start_time, end_time, status
		FROM appointments
		WHERE doctor_id = $1
		ORDER BY appointment_date DESC, start_time DESC
	`

	rows, err := r.db.Query(ctx, query, doctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query doctor appointments: %w", err)
	}
	defer rows.Close()

	var appointments []*domain.Appointment

	for rows.Next() {
		var model AppointmentModel
		err := rows.Scan(
			&model.ID,
			&model.PatientID,
			&model.DoctorID,
			&model.ScheduleID,
			&model.Date,
			&model.StartTime,
			&model.EndTime,
			&model.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment row: %w", err)
		}

		app := domain.ReconstituteAppointment(
			model.ID,
			model.PatientID,
			model.DoctorID,
			model.ScheduleID,
			model.Date,
			model.StartTime,
			model.EndTime,
			domain.AppointmentStatus(model.Status),
		)
		appointments = append(appointments, app)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating appointment rows: %w", err)
	}

	return appointments, nil
}
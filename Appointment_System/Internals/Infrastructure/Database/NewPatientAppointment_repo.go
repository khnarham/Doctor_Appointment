package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
)

type appointmentPatientRepository struct {
	db *pgxpool.Pool
}

func NewPatientAppointmentRepository(db *pgxpool.Pool) *appointmentPatientRepository {
	return &appointmentPatientRepository{db: db}
}

type AppointmentModel struct {
	ID         string
	PatientID  string
	DoctorID   string
	ScheduleID string
	Date       time.Time
	StartTime  string
	EndTime    string
	Status     string
	CreatedAt  time.Time
}

func (r *appointmentPatientRepository) GetByPatientID(ctx context.Context, patientID string) ([]*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, schedule_id, appointment_date, start_time, end_time, status
		FROM appointments
		WHERE patient_id = $1
		ORDER BY appointment_date DESC, start_time DESC
	`

	rows, err := r.db.Query(ctx, query, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
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
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return appointments, nil
}
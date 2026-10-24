package database


import (
	"context"
	"errors"
	"fmt"
	"time"
	"strconv"
	"strings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	domain "my-appointment-app/Appointment_System/Internals/Domain"

)

type appointmentModel struct {
	ID              string    `db:"id"`
	PatientID       string    `db:"patient_id"`
	DoctorID        string    `db:"doctor_id"`
	ScheduleID      string    `db:"schedule_id"`
	AppointmentDate time.Time `db:"appointment_date"`
	StartTime       int       `db:"start_time"` 
	EndTime         int       `db:"end_time"`   
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
}

func timeStringToMinutes(tStr string) (int, error) {
	parts := strings.Split(tStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format: %s", tStr)
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	return (hours * 60) + minutes, nil
}

func minutesToTimeString(totalMinutes int) string {
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}






type pgxBookAppointmentRepository struct {
	db *pgxpool.Pool
}

func NewPgxBookAppointmentRepository(db *pgxpool.Pool) *pgxBookAppointmentRepository {
	return &pgxBookAppointmentRepository{db: db}
}

func (r *pgxBookAppointmentRepository) GetByID(ctx context.Context, id string) (*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, schedule_id, appointment_date, start_time, end_time, status, created_at
		FROM appointments
		WHERE id = $1
	`

	var model appointmentModel
	err := r.db.QueryRow(ctx, query, id).Scan(
		&model.ID,
		&model.PatientID,
		&model.DoctorID,
		&model.ScheduleID,
		&model.AppointmentDate,
		&model.StartTime,
		&model.EndTime,
		&model.Status,
		&model.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("appointment not found with id %s: %w", id, err)
		}
		return nil, fmt.Errorf("failed to fetch appointment: %w", err)
	}

	return domain.ReconstituteAppointment(
		model.ID,
		model.PatientID,
		model.DoctorID,
		model.ScheduleID,
		model.AppointmentDate,
		minutesToTimeString(model.StartTime),
		minutesToTimeString(model.EndTime),
		domain.AppointmentStatus(model.Status),
	), nil
}

func (r *pgxBookAppointmentRepository) IsSlotBooked(ctx context.Context, doctorID string, date time.Time, startTime string) (bool, error) {
	startMinutes, err := timeStringToMinutes(startTime)
	if err != nil {
		return false, fmt.Errorf("invalid start time format: %w", err)
	}

	formattedDate := date.Format("2006-01-02")

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM appointments 
			WHERE doctor_id = $1 
			  AND appointment_date = $2 
			  AND start_time = $3 
			  AND status NOT IN ('CANCELLED')
		)
	`

	var exists bool
	err = r.db.QueryRow(ctx, query, doctorID, formattedDate, startMinutes).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check slot booking status: %w", err)
	}

	return exists, nil
}

func (r *pgxBookAppointmentRepository) Save(ctx context.Context, appointment *domain.Appointment) error {
	startMinutes, err := timeStringToMinutes(appointment.StartTime())
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}

	endMinutes, err := timeStringToMinutes(appointment.EndTime())
	if err != nil {
		return fmt.Errorf("invalid end time: %w", err)
	}

	formattedDate := appointment.Date().Format("2006-01-02")

	query := `
		INSERT INTO appointments (
			id, patient_id, doctor_id, schedule_id, appointment_date, start_time, end_time, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.Exec(
		ctx,
		query,
		appointment.ID(),
		appointment.PatientID(),
		appointment.DoctorID(),
		appointment.ScheduleID(),
		formattedDate,
		startMinutes,
		endMinutes,
		string(appointment.Status()),
		appointment.CreatedAt(),
	)

	if err != nil {
		return fmt.Errorf("failed to save appointment: %w", err)
	}

	return nil
}

func (r *pgxBookAppointmentRepository) Update(ctx context.Context, appointment *domain.Appointment) error {
	startMinutes, err := timeStringToMinutes(appointment.StartTime())
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}

	endMinutes, err := timeStringToMinutes(appointment.EndTime())
	if err != nil {
		return fmt.Errorf("invalid end time: %w", err)
	}

	formattedDate := appointment.Date().Format("2006-01-02")

	query := `
		UPDATE appointments
		SET patient_id = $1,
		    doctor_id = $2,
		    schedule_id = $3,
		    appointment_date = $4,
		    start_time = $5,
		    end_time = $6,
		    status = $7
		WHERE id = $8
	`

	cmdTag, err := r.db.Exec(
		ctx,
		query,
		appointment.PatientID(),
		appointment.DoctorID(),
		appointment.ScheduleID(),
		formattedDate,
		startMinutes,
		endMinutes,
		string(appointment.Status()),
		appointment.ID(),
	)

	if err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("appointment with id %s not found for update", appointment.ID())
	}

	return nil
}
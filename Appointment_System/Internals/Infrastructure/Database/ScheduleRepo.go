package database

import (
	"context"
	"errors"
	"fmt"
	domain "my-appointment-app/Appointment_System/Internals/Domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type typeKey string

const txKey typeKey = "tx"

type DBExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) // Added Query method
}

func getExecutor(ctx context.Context, defaultDB *pgxpool.Pool) DBExecutor {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}
	return defaultDB
}

type DoctorScheduleModel struct {
	ID           string    `db:"id"`
	DoctorID     string    `db:"doctor_id"`
	DayOfWeek    string    `db:"day_of_week"`
	StartTime    int       `db:"start_time"`
	EndTime      int       `db:"end_time"`
	SlotDuration int64     `db:"slot_duration"` 
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type pgxScheduleRepository struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) *pgxScheduleRepository {
	return &pgxScheduleRepository{db: db}
}

func (r *pgxScheduleRepository) Save(ctx context.Context, schedule *domain.DoctorSchedule) error {
	exec := getExecutor(ctx, r.db)

	query := `
        INSERT INTO doctor_schedules (id, doctor_id, day_of_week, start_time, end_time, slot_duration, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

	_, err := exec.Exec(ctx, query,
		schedule.ID(),
		schedule.DoctorID(),
		schedule.Day().String(),                 
		schedule.TimeRange().Start().Minutes(),
		schedule.TimeRange().End().Minutes(),
		int64(schedule.SlotDuration().Minutes()), 
		schedule.IsActive(),
		schedule.CreatedAt(),
		schedule.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert schedule record: %w", err)
	}

	return nil
}

func (r *pgxScheduleRepository) Update(ctx context.Context, schedule *domain.DoctorSchedule) error {
	exec := getExecutor(ctx, r.db)

	query := `
        UPDATE doctor_schedules
        SET 
            doctor_id = $1,
            day_of_week = $2,
            start_time = $3,
            end_time = $4,
            slot_duration = $5,
            is_active = $6,
            updated_at = $7
        WHERE id = $8
    `

	_, err := exec.Exec(ctx, query,
		schedule.DoctorID(),
		schedule.Day().String(),                   // Fix: Pass ENUM string
		schedule.TimeRange().Start().Minutes(),
		schedule.TimeRange().End().Minutes(),
		int64(schedule.SlotDuration().Minutes()), // Fix: Minutes conversion
		schedule.IsActive(),
		schedule.UpdatedAt(),
		schedule.ID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update schedule record: %w", err)
	}

	return nil
}

func (r *pgxScheduleRepository) GetByID(ctx context.Context, scheduleID string) (*domain.DoctorSchedule, error) {
	exec := getExecutor(ctx, r.db)

	query := `
        SELECT id, doctor_id, day_of_week, start_time, end_time, slot_duration, is_active, created_at, updated_at
        FROM doctor_schedules
        WHERE id = $1
        LIMIT 1
    `

	row := exec.QueryRow(ctx, query, scheduleID)
	return mapRowToSchedule(row)
}

func (r *pgxScheduleRepository) GetActiveSchedulesByDoctorID(ctx context.Context, doctorID string) ([]domain.DoctorSchedule, error) {
	exec := getExecutor(ctx, r.db)

	query := `
        SELECT id, doctor_id, day_of_week, start_time, end_time, slot_duration, is_active, created_at, updated_at
        FROM doctor_schedules
        WHERE doctor_id = $1 AND is_active = true
    `

	rows, err := exec.Query(ctx, query, doctorID) 
	if err != nil {
		return nil, fmt.Errorf("failed to query active schedules: %w", err)
	}
	defer rows.Close()

	var schedules []domain.DoctorSchedule
	for rows.Next() {
		var model DoctorScheduleModel
		err := rows.Scan(
			&model.ID,
			&model.DoctorID,
			&model.DayOfWeek,
			&model.StartTime,
			&model.EndTime,
			&model.SlotDuration,
			&model.IsActive,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule row: %w", err)
		}

		scheduleEntity, err := reconstructScheduleEntity(model)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, *scheduleEntity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return schedules, nil
}

func mapRowToSchedule(row pgx.Row) (*domain.DoctorSchedule, error) {
	var model DoctorScheduleModel

	err := row.Scan(
		&model.ID,
		&model.DoctorID,
		&model.DayOfWeek,
		&model.StartTime,
		&model.EndTime,
		&model.SlotDuration,
		&model.IsActive,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan schedule row: %w", err)
	}

	return reconstructScheduleEntity(model)
}

func reconstructScheduleEntity(model DoctorScheduleModel) (*domain.DoctorSchedule, error) {
	startTime, err := domain.NewClockTimeFromMinutes(model.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time in db: %w", err)
	}

	endTime, err := domain.NewClockTimeFromMinutes(model.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time in db: %w", err)
	}

	timeRange, err := domain.NewTimeRange(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("invalid time range in db: %w", err)
	}


	dayOfWeek, err := domain.ParseDayOfWeek(model.DayOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to map day_of_week: %w", err)
	}

	slotDuration := time.Duration(model.SlotDuration) * time.Minute

	return domain.ReconstituteSchedule(
		model.ID,
		model.DoctorID,
		dayOfWeek,
		timeRange,
		slotDuration,
		model.IsActive,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}
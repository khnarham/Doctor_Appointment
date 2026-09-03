package database

import (
	"context"
	"errors"
	"fmt"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PatientModel struct {
	ID          string    `db:"id"`
	UserID      *string   `db:"user_id"`
	WalkInName  *string   `db:"walk_in_name"`
	WalkInPhone string    `db:"walk_in_phone"`
	Gender      string    `db:"gender"`
	DateOfBirth time.Time `db:"date_of_birth"`
	IsWalkIn    bool      `db:"is_walk_in"`
	CreatedAt   time.Time `db:"created_at"`
}

type DBExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type patientRegisterRepository struct {
	db *pgxpool.Pool
}

func NewPatientRegisterRepository(db *pgxpool.Pool) *patientRegisterRepository {
	return &patientRegisterRepository{db: db}
}

func getExecutor(ctx context.Context, defaultDB *pgxpool.Pool) DBExecutor {
	if tx, ok := ctx.Value("tx").(pgx.Tx); ok {
		return tx
	}
	return defaultDB
}

func (r *patientRegisterRepository) Save(ctx context.Context, patient *domain.Patient) error {
	exec := getExecutor(ctx, r.db)

	query := `
		INSERT INTO patients (id, user_id, walk_in_name, walk_in_phone, gender, date_of_birth, is_walk_in, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	var walkInName *string
	if name := patient.WalkInName(); name != "" {
		walkInName = &name
	}

	_, err := exec.Exec(ctx, query,
		patient.ID(),
		patient.UserID(),
		walkInName,
		patient.WalkInPhone(),
		string(patient.Gender()),
		patient.DateOfBirth(),
		patient.IsWalkIn(),
		patient.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert patient record: %w", err)
	}

	return nil
}

func (r *patientRegisterRepository) FindByUserId(ctx context.Context, userID string) (*domain.Patient, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, user_id, walk_in_name, walk_in_phone, gender, date_of_birth, is_walk_in, created_at
		FROM patients
		WHERE user_id = $1
		LIMIT 1
	`

	row := exec.QueryRow(ctx, query, userID)
	return mapRowToPatient(row)
}

func (r *patientRegisterRepository) FindByPhone(ctx context.Context, walkInPhone string) (*domain.Patient, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, user_id, walk_in_name, walk_in_phone, gender, date_of_birth, is_walk_in, created_at
		FROM patients
		WHERE walk_in_phone = $1
		LIMIT 1
	`

	row := exec.QueryRow(ctx, query, walkInPhone)
	return mapRowToPatient(row)
}

func mapRowToPatient(row pgx.Row) (*domain.Patient, error) {
	var model PatientModel

	err := row.Scan(
		&model.ID,
		&model.UserID,
		&model.WalkInName,
		&model.WalkInPhone,
		&model.Gender,
		&model.DateOfBirth,
		&model.IsWalkIn,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan patient row: %w", err)
	}

	var nameStr string
	if model.WalkInName != nil {
		nameStr = *model.WalkInName
	}

	return domain.ReconstructPatient(
		model.ID,
		model.UserID,
		nameStr,
		model.WalkInPhone,
		domain.Gender(model.Gender),
		model.DateOfBirth,
		model.IsWalkIn,
		model.CreatedAt,
	), nil
}
package database

import (
	"context"
	"errors"
	"fmt"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SpecializationModel struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type specializationRepository struct {
	db *pgxpool.Pool
}

func NewSpecializationRepository(db *pgxpool.Pool) *specializationRepository {
	return &specializationRepository{db: db}
}

func (r *specializationRepository) Save(ctx context.Context, spec *domain.Specialization) error {
	exec := getExecutor(ctx, r.db)

	query := `
		INSERT INTO specializations (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			updated_at = EXCLUDED.updated_at
	`

	_, err := exec.Exec(ctx, query,
		spec.ID(),
		spec.Name(),
		spec.Description(),
		spec.CreatedAt(),
		spec.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to save specialization: %w", err)
	}

	return nil
}

func (r *specializationRepository) FindByName(ctx context.Context, name string) (*domain.Specialization, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, name, description, created_at, updated_at
		FROM specializations
		WHERE name = $1
		LIMIT 1
	`

	row := exec.QueryRow(ctx, query, name)
	return mapRowToSpecialization(row)
}

func (r *specializationRepository) FindAll(ctx context.Context) ([]*domain.Specialization, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, name, description, created_at, updated_at
		FROM specializations
		ORDER BY name ASC
	`

	rows, err := exec.(interface {
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	}).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query specializations: %w", err)
	}
	defer rows.Close()

	var result []*domain.Specialization
	for rows.Next() {
		var model SpecializationModel
		err := rows.Scan(
			&model.ID,
			&model.Name,
			&model.Description,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan specialization row: %w", err)
		}

		result = append(result, domain.ReconstituteSpecialization(
			model.ID,
			model.Name,
			model.Description,
			model.CreatedAt,
			model.UpdatedAt,
		))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return result, nil
}


func (r *specializationRepository) FindByID(ctx context.Context, id string) (*domain.Specialization, error) {
	exec := getExecutor(ctx, r.db)

	query := `
		SELECT id, name, description, created_at, updated_at
		FROM specializations
		WHERE id = $1
		LIMIT 1
	`

	row := exec.QueryRow(ctx, query, id)
	return mapRowToSpecialization(row)
}


// 3. Helper Mapper
func mapRowToSpecialization(row pgx.Row) (*domain.Specialization, error) {
	var model SpecializationModel

	err := row.Scan(
		&model.ID,
		&model.Name,
		&model.Description,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan specialization row: %w", err)
	}

	return domain.ReconstituteSpecialization(
		model.ID,
		model.Name,
		model.Description,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}
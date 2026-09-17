package database

import (
	"context"
	"errors"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SpecializationRepository struct {
	db *pgxpool.Pool
}

func NewSpecializationOperationsRepository(db *pgxpool.Pool) *SpecializationRepository {
	return &SpecializationRepository{db: db}
}

func (r *SpecializationRepository) GetByID(ctx context.Context, id string) (*domain.Specialization, error) {
	query := `
		SELECT id, name, description, created_at, updated_at 
		FROM specializations 
		WHERE id = $1`

	var (
		specID, name, description string
		createdAt, updatedAt       time.Time
	)

	err := r.db.QueryRow(ctx, query, id).Scan(&specID, &name, &description, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSpecializationNotFound
		}
		return nil, err
	}

	return domain.ReconstituteSpecialization(specID, name, description, createdAt, updatedAt), nil
}

func (r *SpecializationRepository) GetByName(ctx context.Context, name string) (*domain.Specialization, error) {
	query := `
		SELECT id, name, description, created_at, updated_at 
		FROM specializations 
		WHERE LOWER(name) = LOWER($1)`

	var (
		specID, specName, description string
		createdAt, updatedAt           time.Time
	)

	err := r.db.QueryRow(ctx, query, name).Scan(&specID, &specName, &description, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSpecializationNotFound
		}
		return nil, err
	}

	return domain.ReconstituteSpecialization(specID, specName, description, createdAt, updatedAt), nil
}

func (r *SpecializationRepository) Update(ctx context.Context, spec *domain.Specialization) error {
	query := `
		UPDATE specializations 
		SET name = $1, description = $2, updated_at = $3 
		WHERE id = $4`

	cmdTag, err := r.db.Exec(ctx, query, spec.Name(), spec.Description(), spec.UpdatedAt(), spec.ID())
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrSpecializationNotFound
	}

	return nil
}

func (r *SpecializationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM specializations WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrSpecializationNotFound
	}

	return nil
}

// Integrity Check Rule: Check if any doctor is using this specialization before deletion
func (r *SpecializationRepository) HasAssociatedDoctors(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM doctor_specializations WHERE specialization_id = $1)`
	
	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
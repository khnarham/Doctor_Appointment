package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidSpecializationName = errors.New("specialization name cannot be empty")

type Specialization struct {
	id          string
	name        string
	description string 
	createdAt   time.Time
	updatedAt   time.Time
}

func NewSpecialization(name string, description string) (*Specialization, error) {
	if name == "" {
		return nil, ErrInvalidSpecializationName
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Specialization{
		id:          id.String(),
		name:        name,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func ReconstituteSpecialization(
	id string,
	name string,
	description string,
	createdAt time.Time,
	updatedAt time.Time,
) *Specialization {
	return &Specialization{
		id:          id,
		name:        name,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
func (s *Specialization) ID() string          { return s.id }
func (s *Specialization) Name() string        { return s.name }
func (s *Specialization) Description() string { return s.description }
func (s *Specialization) CreatedAt() time.Time { return s.createdAt }
func (s *Specialization) UpdatedAt() time.Time { return s.updatedAt }
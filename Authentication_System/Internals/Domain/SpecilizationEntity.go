package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidSpecializationName = errors.New("specialization name cannot be empty")
	ErrSpecializationNotFound    = errors.New("specialization not found")
	ErrSpecializationNameExists  = errors.New("specialization name already exists")
	ErrSpecializationInUse       = errors.New("cannot delete specialization assigned to active doctors")
)

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

func (s *Specialization) Update(name string, description string) error {
	if name == "" {
		return ErrInvalidSpecializationName
	}

	s.name = name
	s.description = description
	s.updatedAt = time.Now().UTC()
	return nil
}

func (s *Specialization) ID() string          { return s.id }
func (s *Specialization) Name() string        { return s.name }
func (s *Specialization) Description() string { return s.description }
func (s *Specialization) CreatedAt() time.Time { return s.createdAt }
func (s *Specialization) UpdatedAt() time.Time { return s.updatedAt }


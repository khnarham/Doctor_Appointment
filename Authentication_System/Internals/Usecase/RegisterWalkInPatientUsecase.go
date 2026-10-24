package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "my-appointment-app/Authentication_System/Internals/Domain"

	"github.com/google/uuid"
)

var (
	ErrPatientAlreadyExists = errors.New("patient with this phone number already exists")
	ErrInvalidDateFormat    = errors.New("date_of_birth must be in YYYY-MM-DD format")
)


type IDGenerator interface {
	Generate() string
}

type RegisterWalkInPatientInput struct {
	Name        string
	Phone       string
	Gender      string
	DateOfBirth string // YYYY-MM-DD format
}

type RegisterWalkInPatientOutput struct {
	ID          string    `json:"id"`
	WalkInName  string    `json:"walk_in_name"`
	WalkInPhone string    `json:"walk_in_phone"`
	Gender      string    `json:"gender"`
	DateOfBirth string    `json:"date_of_birth"`
	IsWalkIn    bool      `json:"is_walk_in"`
	CreatedAt   time.Time `json:"created_at"`
}

type RegisterWalkInPatientUsecase struct {
	repo  PatientRegisterReposaitory
	idGen IDGenerator
}

func NewRegisterWalkInPatientUsecase(
	repo PatientRegisterReposaitory,
) *RegisterWalkInPatientUsecase {
	return &RegisterWalkInPatientUsecase{
		repo:  repo,
	}
}

func (uc *RegisterWalkInPatientUsecase) Execute(
	ctx context.Context,
	input RegisterWalkInPatientInput,
) (*RegisterWalkInPatientOutput, error) {
	// 1. DOB string to time.Time parse karo
	parsedDOB, err := time.Parse("2006-01-02", input.DateOfBirth)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	// 2. Uniqueness check by Phone
	existingPatient, err := uc.repo.FindByPhone(ctx, input.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing patient: %w", err)
	}
	if existingPatient != nil {
		return nil, ErrPatientAlreadyExists
	}

		patientID, err := uuid.NewV7()
	     if err != nil {
		return nil ,err
	     }

	patient, err := domain.NewWalkInPatient(
		patientID.String(),
		input.Name,
		input.Phone,
		domain.Gender(input.Gender),
		parsedDOB,
	)
	if err != nil {
		return nil, fmt.Errorf("domain creation failed: %w", err)
	}

	// 4. Persistence
	if err := uc.repo.Save(ctx, patient); err != nil {
		return nil, fmt.Errorf("failed to save walk-in patient: %w", err)
	}

	return &RegisterWalkInPatientOutput{
		ID:          patient.ID(),
		WalkInName:  patient.WalkInName(),
		WalkInPhone: patient.WalkInPhone(),
		Gender:      string(patient.Gender()),
		DateOfBirth: patient.DateOfBirth().Format("2006-01-02"),
		IsWalkIn:    patient.IsWalkIn(),
		CreatedAt:   patient.CreatedAt(),
	}, nil
}
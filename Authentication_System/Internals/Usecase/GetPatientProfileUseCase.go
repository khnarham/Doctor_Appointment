package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Repository 



type GetPatientProfileInput struct {
	UserID string 
}

type PatientProfileOutput struct {
	ID          string    `json:"id"`
	UserID      *string   `json:"user_id,omitempty"`
	WalkInName  string    `json:"walk_in_name,omitempty"`
	WalkInPhone string    `json:"walk_in_phone,omitempty"`
	Gender      string    `json:"gender"`
	DateOfBirth string    `json:"date_of_birth"`
	IsWalkIn    bool      `json:"is_walk_in"`
	CreatedAt   time.Time `json:"created_at"`
}



var ErrPatientNotFound = errors.New("patient profile not found")

type GetPatientProfileUseCase struct {
	patientRepo PatientRegisterReposaitory
}

func NewGetPatientProfileUseCase(pRepo PatientRegisterReposaitory) *GetPatientProfileUseCase {
	return &GetPatientProfileUseCase{
		patientRepo: pRepo,
	}
}

func (uc *GetPatientProfileUseCase) Execute(ctx context.Context, input GetPatientProfileInput) (PatientProfileOutput, error) {
	if input.UserID == "" {
		return PatientProfileOutput{}, errors.New("user ID is required")
	}

	patient, err := uc.patientRepo.FindByUserId(ctx, input.UserID)
	if err != nil {
		return PatientProfileOutput{}, fmt.Errorf("%w: %v", ErrPatientNotFound, err)
	}

	return PatientProfileOutput{
		ID:          patient.ID(),
		UserID:      patient.UserID(),
		WalkInName:  patient.WalkInName(),
		WalkInPhone: patient.WalkInPhone(),
		Gender:      string(patient.Gender()),
		DateOfBirth: patient.DateOfBirth().Format("2006-01-02"),
		IsWalkIn:    patient.IsWalkIn(),
		CreatedAt:   patient.CreatedAt(),
	}, nil
}
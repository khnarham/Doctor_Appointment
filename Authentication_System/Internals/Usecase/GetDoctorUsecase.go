package usecase

import (
	"context"
	"errors"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
)

var (
	ErrInvalidUserID      = errors.New("invalid or empty user id")
	ErrDoctorNotFound     = errors.New("doctor profile not found")
	ErrDoctorNotAvailable = errors.New("doctor is currently not available")
)

type DoctorUsecase struct {
	doctorRepo GetDoctorRepository
}

func NewDoctorUsecase(repo GetDoctorRepository) *DoctorUsecase {
	return &DoctorUsecase{
		doctorRepo: repo,
	}
}

func (u *DoctorUsecase) Execute(ctx context.Context, userID string) (*domain.Doctor, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	doctor, err := u.doctorRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, ErrDoctorNotFound
	}

	if !doctor.IsAvailable() {
		return nil, ErrDoctorNotAvailable
	}

	return doctor, nil
}
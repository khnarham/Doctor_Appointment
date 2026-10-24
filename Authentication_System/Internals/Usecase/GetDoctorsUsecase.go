package usecase

import (
	"context"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"
)

type DoctorDTO struct {
	ID               string              `json:"id"`
	UserID           string              `json:"userID"`
	SpecializationID string              `json:"specializationID"`
	LicenseNumber    string              `json:"licenseNumber"`
	ConsultationFee  float64             `json:"consultationFee"`
	Status           domain.DoctorStatus `json:"status"`
	IsAvailable      bool                `json:"isAvailable"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
}


type DoctorFilter struct {
	SpecializationID string
	MaxFee           float64
	Search           string
	Limit            int
	Offset           int
}


type GetdoctorUsecase struct {
	repo GetDoctorsRepository
}

func NewGetDoctorUsecase(repo GetDoctorsRepository) *GetdoctorUsecase {
	return &GetdoctorUsecase{repo: repo}
}

func (u *GetdoctorUsecase) ListDoctors(ctx context.Context, filter DoctorFilter) ([]DoctorDTO, error) {
	doctors, err := u.repo.GetDoctors(ctx, filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]DoctorDTO, 0, len(doctors))
	for _, doc := range doctors {
		dtos = append(dtos, DoctorDTO{
			ID:               doc.ID(),
			UserID:           doc.UserID(),
			SpecializationID: doc.SpecializationID(),
			LicenseNumber:    doc.LicenseNumber(),
			ConsultationFee:  doc.ConsultationFee(),
			Status:           doc.Status(),
			IsAvailable:      doc.IsAvailable(),
			CreatedAt:        doc.CreatedAt(),
			UpdatedAt:        doc.UpdatedAt(),
		})
	}

	return dtos, nil
}
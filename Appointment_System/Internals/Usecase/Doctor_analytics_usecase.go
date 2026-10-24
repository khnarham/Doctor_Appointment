package usecase

import (
	"context"
	domain "my-appointment-app/Appointment_System/Internals/Domain"
	"time"
)


type DoctorAnalyticsRepository interface {
	GetDoctorMetrics(ctx context.Context, doctorID string) (*domain.DoctorAnalyticsSummary, error)
	GetDoctorPatients(ctx context.Context, doctorID string, limit int) ([]domain.DoctorPatientDetail, error)
}

type DoctorAnalyticsUsecase struct {
	repo           DoctorAnalyticsRepository
	contextTimeout time.Duration
}

func NewDoctorAnalyticsUsecase(repo DoctorAnalyticsRepository, timeout time.Duration) *DoctorAnalyticsUsecase {
	return &DoctorAnalyticsUsecase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *DoctorAnalyticsUsecase) GetDashboardData(c context.Context, doctorID string) (*domain.DoctorDashboardResponse, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// 1. Fetch aggregated stats for the doctor
	summary, err := u.repo.GetDoctorMetrics(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch latest patient appointments specific to this doctor (e.g., top 10)
	patients, err := u.repo.GetDoctorPatients(ctx, doctorID, 10)
	if err != nil {
		return nil, err
	}

	return &domain.DoctorDashboardResponse{
		Summary:        *summary,
		RecentPatients: patients,
	}, nil
}
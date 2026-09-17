package usecase

import (
	"context"
	domain "my-appointment-app/Appointment_System/Internals/Domain"
	database "my-appointment-app/Appointment_System/Internals/Infrastructure/Database"
)

type AnalyticsUsecase interface {
	GetFinancialOverview(ctx context.Context, filter domain.FinancialFilterDTO) (*domain.FinancialOverviewResponse, error)
}

type analyticsUsecase struct {
	analyticsRepo database.AnalyticsRepository
}

func NewAnalyticsUsecase(repo database.AnalyticsRepository) AnalyticsUsecase {
	return &analyticsUsecase{
		analyticsRepo: repo,
	}
}

func (u *analyticsUsecase) GetFinancialOverview(ctx context.Context, filter domain.FinancialFilterDTO) (*domain.FinancialOverviewResponse, error) {
	// 1. Compute total completed appointments and total revenue
	totalCompleted, totalRevenue, err := u.analyticsRepo.GetTotalRevenueStats(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 2. Fetch top earning doctors
	topDoctors, err := u.analyticsRepo.GetTopEarningDoctors(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Nil slice protection
	if topDoctors == nil {
		topDoctors = []domain.TopEarningDoctorDTO{}
	}

	return &domain.FinancialOverviewResponse{
		TotalCompletedAppointments: totalCompleted,
		TotalRevenue:               totalRevenue,
		TopEarningDoctors:          topDoctors,
	}, nil
}
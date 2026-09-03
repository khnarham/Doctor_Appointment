package usecase

import (
	"context"
	"fmt"
)

type CancelScheduleInput struct {
	ScheduleID string
	DoctorID   string
}

type CancelScheduleUsecase struct {
	scheduleRepo ScheduleRepository
}

func NewCancelScheduleUsecase(repo ScheduleRepository) *CancelScheduleUsecase {
	return &CancelScheduleUsecase{scheduleRepo: repo}
}

func (uc *CancelScheduleUsecase) Execute(ctx context.Context, input CancelScheduleInput) error {
	schedule, err := uc.scheduleRepo.GetByID(ctx, input.ScheduleID)
	if err != nil {
		return fmt.Errorf("schedule not found: %w", err)
	}

	if schedule.DoctorID() != input.DoctorID {
		return fmt.Errorf("unauthorized: schedule does not belong to this doctor")
	}

	if err := schedule.Deactivate(); err != nil {
		return fmt.Errorf("deactivation failed: %w", err)
	}

	if err := uc.scheduleRepo.Update(ctx, schedule); err != nil {
		return fmt.Errorf("failed to deactivate schedule: %w", err)
	}

	return nil
}
package usecase

import (
	"context"
	"fmt"
	"time"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
)

type UpdateScheduleInput struct {
	ScheduleID   string
	DoctorID     string
	Day          domain.DayOfWeek
	StartHour    int
	StartMinute  int
	EndHour      int
	EndMinute    int
	SlotDuration time.Duration
}

type UpdateScheduleOutput struct {
	ID           string           `json:"id"`
	DoctorID     string           `json:"doctor_id"`
	Day          domain.DayOfWeek `json:"day"`
	StartTime    string           `json:"start_time"`
	EndTime      string           `json:"end_time"`
	SlotDuration int              `json:"slot_duration_minutes"`
	IsActive     bool             `json:"is_active"`
}

type UpdateScheduleUsecase struct {
	scheduleRepo ScheduleRepository
}

func NewUpdateScheduleUsecase(repo ScheduleRepository) *UpdateScheduleUsecase {
	return &UpdateScheduleUsecase{scheduleRepo: repo}
}

func (uc *UpdateScheduleUsecase) Execute(ctx context.Context, input UpdateScheduleInput) (UpdateScheduleOutput, error) {
	targetSchedule, err := uc.scheduleRepo.GetByID(ctx, input.ScheduleID)
	if err != nil {
		return UpdateScheduleOutput{}, fmt.Errorf("schedule not found: %w", err)
	}

	if targetSchedule.DoctorID() != input.DoctorID {
		return UpdateScheduleOutput{}, fmt.Errorf("unauthorized: schedule does not belong to doctor")
	}

	startClock, err := domain.NewClockTime(input.StartHour, input.StartMinute)
	if err != nil {
		return UpdateScheduleOutput{}, fmt.Errorf("invalid start time: %w", err)
	}

	endClock, err := domain.NewClockTime(input.EndHour, input.EndMinute)
	if err != nil {
		return UpdateScheduleOutput{}, fmt.Errorf("invalid end time: %w", err)
	}

	timeRange, err := domain.NewTimeRange(startClock, endClock)
	if err != nil {
		return UpdateScheduleOutput{}, fmt.Errorf("invalid time range: %w", err)
	}

	if err := targetSchedule.Update(input.Day, timeRange, input.SlotDuration); err != nil {
		return UpdateScheduleOutput{}, fmt.Errorf("domain update failed: %w", err)
	}

	activeSchedules, err := uc.scheduleRepo.GetActiveSchedulesByDoctorID(ctx, input.DoctorID)
	if err != nil {
		return UpdateScheduleOutput{}, fmt.Errorf("failed to fetch active schedules: %w", err)
	}

	otherSchedules := make([]domain.DoctorSchedule, 0, len(activeSchedules))
	for _, s := range activeSchedules {
		if s.ID() != input.ScheduleID {
			otherSchedules = append(otherSchedules, s)
		}
	}

	if err := targetSchedule.CheckOverlapWith(otherSchedules); err != nil {
		return UpdateScheduleOutput{}, fmt.Errorf("schedule overlap error: %w", err)
	}

	if err := uc.scheduleRepo.Update(ctx, targetSchedule); err != nil {
		return UpdateScheduleOutput{}, fmt.Errorf("failed to save schedule update: %w", err)
	}

	return UpdateScheduleOutput{
		ID:           targetSchedule.ID(),
		DoctorID:     targetSchedule.DoctorID(),
		Day:          targetSchedule.Day(),
		StartTime:    targetSchedule.TimeRange().Start().String(),
		EndTime:      targetSchedule.TimeRange().End().String(),
		SlotDuration: int(targetSchedule.SlotDuration().Minutes()),
		IsActive:     targetSchedule.IsActive(),
	}, nil
}
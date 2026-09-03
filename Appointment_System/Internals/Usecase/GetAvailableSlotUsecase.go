package usecase

import (
	"context"
	"fmt"
	"time"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
)


type GetAvailableSlotsInput struct {
	DoctorID   string
	TargetDate string 
}

type TimeSlotOutput struct {
	StartTime string 
	EndTime   string 
	IsBooked  bool  
}

type GetAvailableSlotsOutput struct {
	DoctorID   string           
	TargetDate string          
	Slots      []TimeSlotOutput
}

type GetAvailableSlotsUseCase struct {
	scheduleRepo    ScheduleRepository
	appointmentRepo AppointmentReader
}

func NewGetAvailableSlotsUseCase(sRepo ScheduleRepository, aRepo AppointmentReader) *GetAvailableSlotsUseCase {
	return &GetAvailableSlotsUseCase{
		scheduleRepo:    sRepo,
		appointmentRepo: aRepo,
	}
}

func (uc *GetAvailableSlotsUseCase) Execute(ctx context.Context, input GetAvailableSlotsInput) (GetAvailableSlotsOutput, error) {
	parsedDate, err := time.Parse("2006-01-02", input.TargetDate)
	if err != nil {
		return GetAvailableSlotsOutput{}, fmt.Errorf("invalid date format: %w", err)
	}

	parsedDate = parsedDate.UTC()
	now := time.Now().UTC()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	if parsedDate.Before(today) {
		return GetAvailableSlotsOutput{
			DoctorID:   input.DoctorID,
			TargetDate: input.TargetDate,
			Slots:      []TimeSlotOutput{},
		}, fmt.Errorf("failed to fetch schedules: %w", err)
	}
   
	targetDay , err := parseDayOfWeek(parsedDate.Weekday().String())
	if err != nil {
		return GetAvailableSlotsOutput{} , fmt.Errorf("Nikal chutiya: %w", err)
	}

	schedules, err := uc.scheduleRepo.GetActiveSchedulesByDoctorID(ctx, input.DoctorID)
	if err != nil {
		return GetAvailableSlotsOutput{}, fmt.Errorf("failed to fetch schedules: %w", err)
	}

	var matchedSchedules []domain.DoctorSchedule
	for _, s := range schedules {
		if s.IsAvailableOn(targetDay) {
			matchedSchedules = append(matchedSchedules, s)
		}
	}

	if len(matchedSchedules) == 0 {
		return GetAvailableSlotsOutput{
			DoctorID:   input.DoctorID,
			TargetDate: input.TargetDate,
			Slots:      []TimeSlotOutput{},
		}, fmt.Errorf("failed to fetch schedules: %w", err)
	}

	bookedTimes, err := uc.appointmentRepo.GetBookedSlotTimesByDoctorAndDate(ctx, input.DoctorID, parsedDate)
	if err != nil {
		return GetAvailableSlotsOutput{}, fmt.Errorf("failed to fetch bookings: %w", err)
	}

	bookedMap := make(map[string]bool, len(bookedTimes))
	for _, t := range bookedTimes {
		bookedMap[t] = true
	}

	outputSlots := make([]TimeSlotOutput, 0)
	isToday := parsedDate.Equal(today)

	for _, schedule := range matchedSchedules {
		rawSlots, err := schedule.GenerateSlotsForDate(parsedDate)
		if err != nil {
			return GetAvailableSlotsOutput{}, fmt.Errorf("failed to generate slots for schedule %s: %w", schedule.ID(), err)
		}

		for _, slot := range rawSlots {
			if isToday && slot.Start.Before(now) {
				continue
			}

			startStr := slot.Start.Format("15:04")
			endStr := slot.End.Format("15:04")

			outputSlots = append(outputSlots, TimeSlotOutput{
				StartTime: startStr,
				EndTime:   endStr,
				IsBooked:  bookedMap[startStr],
			})
		}
	}

	return GetAvailableSlotsOutput{
		DoctorID:   input.DoctorID,
		TargetDate: input.TargetDate,
		Slots:      outputSlots,
	}, nil
}



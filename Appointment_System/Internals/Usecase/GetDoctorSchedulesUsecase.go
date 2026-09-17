package usecase

import (
	"context"
	"fmt"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
)

type ScheduleSlotOutput struct {
	StartTime string `json:"start_time"` 
	EndTime   string `json:"end_time"`  
	IsBooked  bool   `json:"is_booked"` 
}

type DoctorScheduleOutput struct {
	ID           string               `json:"id"`
	Day          string               `json:"day"`
	StartTime    string               `json:"start_time"`
	EndTime      string               `json:"end_time"`
	SlotDuration int                  `json:"slot_duration_minutes"`
	Slots        []ScheduleSlotOutput `json:"slots"`
}

type GetDoctorSchedulesOutput struct {
	DoctorID  string                 `json:"doctor_id"`
	Schedules []DoctorScheduleOutput `json:"schedules"`
}

type GetDoctorSchedulesUseCase struct {
	scheduleRepo    ScheduleRepository
	appointmentRepo AppointmentReader
}

func NewGetDoctorSchedulesUseCase(sRepo ScheduleRepository, aRepo AppointmentReader) *GetDoctorSchedulesUseCase {
	return &GetDoctorSchedulesUseCase{
		scheduleRepo:    sRepo,
		appointmentRepo: aRepo,
	}
}

func (uc *GetDoctorSchedulesUseCase) Execute(ctx context.Context, doctorID string) (GetDoctorSchedulesOutput, error) {
	if doctorID == "" {
		return GetDoctorSchedulesOutput{}, domain.ErrEmptyDoctorID
	}

	schedules, err := uc.scheduleRepo.GetActiveSchedulesByDoctorID(ctx, doctorID)
	if err != nil {
		return GetDoctorSchedulesOutput{}, fmt.Errorf("failed to fetch doctor schedules: %w", err)
	}

	if len(schedules) == 0 {
		return GetDoctorSchedulesOutput{
			DoctorID:  doctorID,
			Schedules: []DoctorScheduleOutput{},
		}, nil
	}

	appointments, err := uc.appointmentRepo.GetBookedAppointmentsByDoctorID(ctx, doctorID)
	if err != nil {
		return GetDoctorSchedulesOutput{}, fmt.Errorf("failed to fetch doctor appointments: %w", err)
	}

	bookedSlotsMap := make(map[string]bool)
	for _, appt := range appointments {
		if appt.Status() != domain.StatusCancelled {
			key := fmt.Sprintf("%s_%s", appt.ScheduleID(), appt.StartTime())
			bookedSlotsMap[key] = true
		}
	}

	scheduleOutputs := make([]DoctorScheduleOutput, 0, len(schedules))

	for _, s := range schedules {
		rawSlots := s.GenerateBaseSlots()

		slotOutputs := make([]ScheduleSlotOutput, 0, len(rawSlots))
		for _, slot := range rawSlots {
			lookupKey := fmt.Sprintf("%s_%s", s.ID(), slot.Start)
			
			slotOutputs = append(slotOutputs, ScheduleSlotOutput{
				StartTime: slot.Start,
				EndTime:   slot.End,
				IsBooked:  bookedSlotsMap[lookupKey],
			})
		}

		scheduleOutputs = append(scheduleOutputs, DoctorScheduleOutput{
			ID:           s.ID(),
			Day:          s.Day().String(),
			StartTime:    s.TimeRange().Start().String(),
			EndTime:      s.TimeRange().End().String(),
			SlotDuration: int(s.SlotDuration().Minutes()),
			Slots:        slotOutputs,
		})
	}

	return GetDoctorSchedulesOutput{
		DoctorID:  doctorID,
		Schedules: scheduleOutputs,
	}, nil
}
package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
)


type IDGenerator interface {
	Generate() string
}

type BookAppointmentInput struct {
	PatientID string
	DoctorID  string
	Date      string
	StartTime string 
}

type BookAppointmentOutput struct {
	AppointmentID string `json:"appointment_id"`
	DoctorID      string `json:"doctor_id"`
	PatientID     string `json:"patient_id"`
	Date          string `json:"date"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	Status        string `json:"status"`
}

type BookAppointmentUseCase struct {
	scheduleRepo    ScheduleRepository
	appointmentRepo AppointmentRepository
}

func NewBookAppointmentUseCase(
	sRepo ScheduleRepository,
	aRepo AppointmentRepository,
) *BookAppointmentUseCase {
	return &BookAppointmentUseCase{
		scheduleRepo:    sRepo,
		appointmentRepo: aRepo,
	}
}

func (uc *BookAppointmentUseCase) Execute(ctx context.Context, input BookAppointmentInput) (BookAppointmentOutput, error) {
	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return BookAppointmentOutput{}, fmt.Errorf("invalid date format: %w", err)
	}
	parsedDate = parsedDate.UTC()

	slotStartDateTime, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", input.Date, input.StartTime))
	if err != nil {
		return BookAppointmentOutput{}, fmt.Errorf("invalid time format: %w", err)
	}
	slotStartDateTime = slotStartDateTime.UTC()

	now := time.Now().UTC()
	if slotStartDateTime.Before(now) {
		return BookAppointmentOutput{}, fmt.Errorf("cannot book appointment for past date or time")
	}
     
	targetDay,  err := parseDayOfWeek(parsedDate.Weekday().String())
	if err != nil {
		return BookAppointmentOutput{} , err
	}

	schedules, err := uc.scheduleRepo.GetActiveSchedulesByDoctorID(ctx, input.DoctorID)
	if err != nil {
		return BookAppointmentOutput{}, fmt.Errorf("failed to fetch schedules: %w", err)
	}

	var matchedSlot *domain.SlotInterval
	var matchedScheduleID string
	for _, s := range schedules {
		if !s.IsAvailableOn(targetDay) {
			continue
		}

		slots, err := s.GenerateSlotsForDate(parsedDate)
		if err != nil {
			continue
		}

		if slots == nil {
			return  BookAppointmentOutput{} , fmt.Errorf("No Slotes their")
		}

		for _, slot := range slots {
			if slot.Start.Format("15:04") == input.StartTime {
				matchedSlot = &slot
				matchedScheduleID = s.ID()
				break
			}
		}

		if matchedSlot != nil {
			break
		}
	}
   
	if matchedSlot == nil {
       return BookAppointmentOutput{}, fmt.Errorf("schedule match error" )
	}

	isBooked, err := uc.appointmentRepo.IsSlotBooked(ctx, input.DoctorID, parsedDate, input.StartTime)
	if err != nil {
		return BookAppointmentOutput{}, fmt.Errorf("failed to check slot booking status: %w", err)
	}
	if isBooked {
		return BookAppointmentOutput{}, fmt.Errorf("slot already booked")
	}

	endTimeStr := matchedSlot.End.Format("15:04")

	appointment, err := domain.NewAppointment(
		input.PatientID,
		input.DoctorID,
		matchedScheduleID,
		parsedDate,
		input.StartTime,
		endTimeStr,
	)
	if err != nil {
		return BookAppointmentOutput{}, fmt.Errorf("failed to create appointment entity: %w", err)
	}

	if err := uc.appointmentRepo.Save(ctx, appointment); err != nil {
		return BookAppointmentOutput{}, fmt.Errorf("failed to save appointment: %w", err)
	}

	return BookAppointmentOutput{
		AppointmentID: appointment.ID(),
		DoctorID:      appointment.DoctorID(),
		PatientID:     appointment.PatientID(),
		Date:          appointment.Date().Format("2006-01-02"),
		StartTime:     appointment.StartTime(),
		EndTime:       appointment.EndTime(),
		Status:        string(appointment.Status()),
	}, nil
}



func parseDayOfWeek(dayStr string) (domain.DayOfWeek, error) {
	switch strings.ToLower(strings.TrimSpace(dayStr)) {
	case "sunday":
		return domain.Sunday, nil
	case "monday":
		return domain.Monday, nil
	case "tuesday":
		return domain.Tuesday, nil
	case "wednesday":
		return domain.Wednesday, nil
	case "thursday":
		return domain.Thursday, nil
	case "friday":
		return domain.Friday, nil
	case "saturday":
		return domain.Saturday, nil
	default:
		return 0, fmt.Errorf("invalid day of week: %s", dayStr)
	}
}
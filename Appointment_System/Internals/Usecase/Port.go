package usecase

import (
	"context"
	domain "my-appointment-app/Appointment_System/Internals/Domain"
	"time"
)

type ScheduleRepository interface {
    Update(ctx context.Context, schedule *domain.DoctorSchedule) error
	GetActiveSchedulesByDoctorID(ctx context.Context, doctorID string) ([]domain.DoctorSchedule, error)
    GetByID(ctx context.Context, scheduleID string) (*domain.DoctorSchedule, error) 
	Save(ctx context.Context, schedule *domain.DoctorSchedule) error
}

type AppointmentReader interface {
	GetBookedSlotTimesByDoctorAndDate(ctx context.Context, doctorID string, date time.Time) ([]string, error)
}

type AppointmentRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Appointment, error)
	IsSlotBooked(ctx context.Context, doctorID string, date time.Time, startTime string) (bool, error)
	Save(ctx context.Context, appointment *domain.Appointment) error
	Update(ctx context.Context, appointment *domain.Appointment) error
}
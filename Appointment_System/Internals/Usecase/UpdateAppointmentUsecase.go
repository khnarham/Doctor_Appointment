package usecase

import (
	"context"
	"fmt"
	domain "my-appointment-app/Appointment_System/Internals/Domain"
	"time"
)

type UpdateAppointmentStatusInput struct {
	AppointmentID string                  `json:"appointment_id"`
	DoctorID      string                  `json:"doctor_id"`
	NewStatus     domain.AppointmentStatus `json:"new_status"`
}

type UpdateAppointmentStatusUseCase struct {
	appointmentRepo AppointmentRepository
}

func NewUpdateAppointmentStatusUseCase(repo AppointmentRepository) *UpdateAppointmentStatusUseCase {
	return &UpdateAppointmentStatusUseCase{appointmentRepo: repo}
}

func (uc *UpdateAppointmentStatusUseCase) Execute(ctx context.Context, input UpdateAppointmentStatusInput) error {
	appointment, err := uc.appointmentRepo.GetByID(ctx, input.AppointmentID)
	if err != nil {
		return fmt.Errorf("appointment not found: %w", err)
	}

	if appointment.DoctorID() != input.DoctorID {
		return fmt.Errorf("unauthorized: appointment belongs to another doctor")
	}

	if err := appointment.UpdateStatus(input.NewStatus , time.Now()); err != nil {
		return err
	}

	if err := uc.appointmentRepo.Update(ctx, appointment); err != nil {
		return fmt.Errorf("failed to update appointment status: %w", err)
	}

	return nil
}
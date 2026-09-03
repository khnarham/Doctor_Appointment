package usecase

import (
	"context"
	"fmt"
	"time"
)

type CancelAppointmentInput struct {
	AppointmentID string 
	PatientID     string 
}

type CancelAppointmentUseCase struct {
	appointmentRepo AppointmentRepository
}

func NewCancelAppointmentUseCase(repo AppointmentRepository) *CancelAppointmentUseCase {
	return &CancelAppointmentUseCase{appointmentRepo: repo}
}

func (uc *CancelAppointmentUseCase) Execute(ctx context.Context, input CancelAppointmentInput) error {

	appointment, err := uc.appointmentRepo.GetByID(ctx, input.AppointmentID)
	if err != nil {
		return fmt.Errorf("appointment not found: %w", err)
	}

	if appointment.PatientID() != input.PatientID {
		return fmt.Errorf("unauthorized access to appointment")
	}

	if err := appointment.CancelByPatient(time.Now()); err != nil {
		return err
	}

	if err := uc.appointmentRepo.Update(ctx, appointment); err != nil {
		return fmt.Errorf("failed to cancel appointment: %w", err)
	}

	return nil
}
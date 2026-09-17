// Usecase/get_patient_appointments.go
package usecase

import (
	"context"
	"errors"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
)

type GetPatientAppointmentsInput struct {
	PatientID string
}

type PatientAppointmentReader interface {
	GetByPatientID(ctx context.Context, patientID string) ([]*domain.Appointment, error)
}

type AppointmentOutput struct {
	ID         string `json:"id"`
	PatientID  string `json:"patient_id"`
	DoctorID   string `json:"doctor_id"`
	ScheduleID string `json:"schedule_id"`
	Date       string `json:"date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Status     string `json:"status"`
}

type GetPatientAppointmentsUseCase struct {
	repo PatientAppointmentReader
}

func NewGetPatientAppointmentsUseCase(repo PatientAppointmentReader) *GetPatientAppointmentsUseCase {
	return &GetPatientAppointmentsUseCase{repo: repo}
}

func (uc *GetPatientAppointmentsUseCase) Execute(ctx context.Context, input GetPatientAppointmentsInput) ([]AppointmentOutput, error) {
	if input.PatientID == "" {
		return nil, errors.New("patient ID is required")
	}

	appointments, err := uc.repo.GetByPatientID(ctx, input.PatientID)
	if err != nil {
		return nil, err
	}

	output := make([]AppointmentOutput, 0, len(appointments))
	for _, app := range appointments {
		if app == nil {
			continue // Safety check to prevent nil pointer panics
		}

		output = append(output, AppointmentOutput{
			ID:         app.ID(),
			PatientID:  app.PatientID(),
			DoctorID:   app.DoctorID(),
			ScheduleID: app.ScheduleID(),
			Date:       app.Date().Format("2006-01-02"),
			StartTime:  app.StartTime(),
			EndTime:    app.EndTime(),
			Status:     string(app.Status()),
		})
	}

	return output, nil
}
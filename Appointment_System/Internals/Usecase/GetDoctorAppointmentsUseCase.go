package usecase

import (
	"context"
	"errors"

)

var ErrDoctorIDRequired = errors.New("doctor ID is required")

type GetDoctorAppointmentsInput struct {
	DoctorID string
}

type DoctorAppointmentOutput struct {
	ID         string `json:"id"`
	PatientID  string `json:"patient_id"`
	DoctorID   string `json:"doctor_id"`
	ScheduleID string `json:"schedule_id"`
	Date       string `json:"date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Status     string `json:"status"`
}

type GetDoctorAppointmentsUseCase struct {
	repo DoctorAppointmentReader
}

func NewGetDoctorAppointmentsUseCase(repo DoctorAppointmentReader) *GetDoctorAppointmentsUseCase {
	return &GetDoctorAppointmentsUseCase{repo: repo}
}

func (uc *GetDoctorAppointmentsUseCase) Execute(ctx context.Context, input GetDoctorAppointmentsInput) ([]DoctorAppointmentOutput, error) {
	if input.DoctorID == "" {
		return nil, ErrDoctorIDRequired
	}

	appointments, err := uc.repo.GetByDoctorID(ctx, input.DoctorID)
	if err != nil {
		return nil, err
	}

	output := make([]DoctorAppointmentOutput, 0, len(appointments))
	for _, app := range appointments {
		if app == nil {
			continue 
		}

		output = append(output, DoctorAppointmentOutput{
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
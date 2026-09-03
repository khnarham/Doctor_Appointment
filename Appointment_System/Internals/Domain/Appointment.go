package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "PENDING"
	StatusConfirmed AppointmentStatus = "CONFIRMED"
	StatusCancelled AppointmentStatus = "CANCELLED"
	StatusCompleted AppointmentStatus = "COMPLETED"
)

var (
	ErrAppointmentAlreadyCancelled = errors.New("appointment is already cancelled")
	ErrAppointmentAlreadyCompleted = errors.New("appointment is already completed")
	ErrInvalidStatusTransition    = errors.New("invalid appointment status transition")
	ErrCancellationWindowPassed   = errors.New("cannot cancel appointment less than 2 hours before start time")
	ErrInvalidAppointmentDate = errors.New("appointment date must be in the future")
	ErrSlotAlreadyBooked      = errors.New("the requested slot is already booked")
	ErrSlotNotInSchedule      = errors.New("requested time slot is not within doctor schedule")
)




type Appointment struct {
	id         string
	patientID  string
	doctorID   string
	scheduleID string
	date       time.Time
	startTime  string
	endTime    string
	status     AppointmentStatus
	createdAt  time.Time
}


func ReconstituteAppointment(
	id string,
	patientID string,
	doctorID string,
	scheduleID string,
	date time.Time,
	startTime string,
	endTime string,
	status AppointmentStatus,
) *Appointment {
	return &Appointment{
		id:         id,
		patientID:  patientID,
		doctorID:   doctorID,
		scheduleID: scheduleID,
		date:       date,
		startTime:  startTime,
		endTime:    endTime,
		status:     status,
		createdAt:  time.Now(),
	}
}


func NewAppointment(
	patientID, doctorID, scheduleID string,
	date time.Time,
	startTime, endTime string,
) (*Appointment, error) {

	if patientID == "" || doctorID == "" {
		return nil, errors.New("patientID and doctorID are required")
	}

	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, ErrInvalidAppointmentDate
	}
		id, err := uuid.NewV7()
	     if err != nil {
		return nil ,err
	    }



	return &Appointment{
		id:         id.String(),
		patientID:  patientID,
		doctorID:   doctorID,
		scheduleID: scheduleID,
		date:       date,
		startTime:  startTime,
		endTime:    endTime,
		status:     StatusConfirmed,
		createdAt:  time.Now(),
	}, nil
}



func (a *Appointment) CancelByPatient(now time.Time) error {
	if a.status == StatusCancelled {
		return ErrAppointmentAlreadyCancelled
	}
	if a.status == StatusCompleted {
		return ErrAppointmentAlreadyCompleted
	}

	startDateTime, err := time.Parse("2006-01-02 15:04", a.date.Format("2006-01-02")+" "+a.startTime)
	if err == nil {
		if startDateTime.Sub(now) < 2*time.Hour {
			return ErrCancellationWindowPassed
		}
	}

	a.status = StatusCancelled
	return nil
}

func (a *Appointment) UpdateStatus(newStatus AppointmentStatus, now time.Time) error {
	if a.status == StatusCancelled {
		return ErrAppointmentAlreadyCancelled
	}
	if a.status == StatusCompleted {
		return ErrAppointmentAlreadyCompleted
	}

	switch newStatus {
	case StatusConfirmed:
		if a.status != StatusPending {
			return ErrInvalidStatusTransition
		}
		a.status = StatusConfirmed

	case StatusCompleted:
		if a.status != StatusConfirmed {
			return ErrInvalidStatusTransition
		}
		a.status = StatusCompleted

	case StatusCancelled:

		if a.status == StatusConfirmed {
			startDateTime, err := time.Parse("2006-01-02 15:04", a.date.Format("2006-01-02")+" "+a.startTime)
			if err == nil && startDateTime.Sub(now) < 2*time.Hour {
				return ErrCancellationWindowPassed
			}
		}
		a.status = StatusCancelled

	default:
		return ErrInvalidStatusTransition
	}

	return nil
}

func (a *Appointment) ID() string             { return a.id }
func (a *Appointment) PatientID() string      { return a.patientID }
func (a *Appointment) DoctorID() string       { return a.doctorID }
func (a *Appointment) ScheduleID() string     { return a.scheduleID }
func (a *Appointment) Date() time.Time        { return a.date }
func (a *Appointment) StartTime() string      { return a.startTime }
func (a *Appointment) EndTime() string        { return a.endTime }
func (a *Appointment) Status() AppointmentStatus { return a.status }
func (a *Appointment) CreatedAt() time.Time {return  a.createdAt}
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidReviewID     = errors.New("review ID cannot be empty")
	ErrInvalidAppointment  = errors.New("appointment ID cannot be empty")
	ErrInvalidRating       = errors.New("rating must be between 1 and 5")
	ErrAppointmentNotDone  = errors.New("cannot review an incomplete appointment")
)

type Review struct {
	id            string
	appointmentId string
	patientId     string
	doctorId      string
	rating        int   
	comment       string
	createdAt     time.Time
}

func NewReview(
	appointmentId string,
	patientId string,
	doctorId string,
	rating int,
	comment string,
	now time.Time,
) (*Review, error) {


	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRating
	}


	id, err := uuid.NewV7()
	     if err != nil {
		return nil ,err
	     }
	return &Review{
		id:            id.String(),
		appointmentId: appointmentId,
		patientId:     patientId,
		doctorId:      doctorId,
		rating:        rating,
		comment:       comment,
		createdAt:     now.UTC(),
	}, nil
}

func (r *Review) ID() string            { return r.id }
func (r *Review) AppointmentID() string { return r.appointmentId }
func (r *Review) PatientID() string     { return r.patientId }
func (r *Review) DoctorID() string      { return r.doctorId }
func (r *Review) Rating() int           { return r.rating }
func (r *Review) Comment() string       { return r.comment }
func (r *Review) CreatedAt() time.Time  { return r.createdAt }
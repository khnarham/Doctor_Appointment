package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidReceptionistID = errors.New("receptionist ID cannot be empty")
	ErrInvalidShiftTime      = errors.New("shift start time must be before end time")
	ErrReceptionistInactive  = errors.New("receptionist is currently inactive")
)

type ReceptionistID string

func NewReceptionistID(id string) (ReceptionistID, error) {
	if id == "" {
		return "", ErrInvalidReceptionistID
	}
	return ReceptionistID(id), nil
}

type Shift struct {
	StartTime time.Time
	EndTime   time.Time
}

func NewShift(start, end time.Time) (Shift, error) {
	if !start.Before(end) {
		return Shift{}, ErrInvalidShiftTime
	}
	return Shift{StartTime: start, EndTime: end}, nil
}

type ReceptionistStatus string

const (
	StatusActive   ReceptionistStatus = "ACTIVE"
	StatusInactive ReceptionistStatus = "INACTIVE"
	StatusOnLeave  ReceptionistStatus = "ON_LEAVE"
)

type Receptionist struct {
	id        ReceptionistID
	userID    string
	shift     Shift
	status    ReceptionistStatus
	createdAt time.Time
	updatedAt time.Time
}

func NewReceptionist(userID string, shift Shift) *Receptionist {
	now := time.Now()

		id, err := uuid.NewV7()
	     if err != nil {
		return nil
	     }

	return &Receptionist{
		id:        ReceptionistID(id.String()),
		userID:    userID,
		shift:     shift,
		status:    StatusActive,
		createdAt: now,
		updatedAt: now,
	}
}

func (r *Receptionist) ID() ReceptionistID         { return r.id }
func (r *Receptionist) UserID() string            { return r.userID }
func (r *Receptionist) Shift() Shift              { return r.shift }
func (r *Receptionist) Status() ReceptionistStatus { return r.status }

func (r *Receptionist) ChangeShift(newShift Shift) error {
	if r.status != StatusActive {
		return ErrReceptionistInactive
	}
	r.shift = newShift
	r.updatedAt = time.Now()
	return nil
}

func (r *Receptionist) Deactivate() {
	r.status = StatusInactive
	r.updatedAt = time.Now()
}

func (r *Receptionist) MarkOnLeave() {
	r.status = StatusOnLeave
	r.updatedAt = time.Now()
}
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInviteExpired = errors.New("doctor invitation has expired")
	ErrInviteUsed    = errors.New("doctor invitation has already been used")
)

type DoctorInvitation struct {
	id               string
	email            Email
	specializationID string
	token            string
	isUsed           bool
	expiresAt        time.Time
	createdAt        time.Time
}

func NewDoctorInvitation(email Email, specializationID string, token string, expirationDuration time.Duration) (*DoctorInvitation, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &DoctorInvitation{
		id:               id.String(),
		email:            email,
		specializationID: specializationID,
		token:            token,
		isUsed:           false,
		expiresAt:        now.Add(expirationDuration),
		createdAt:        now,
	}, nil
}


func ReconstituteDoctorInvitation(
	id string,
	email Email,
	specializationID string,
	token string,
	isUsed bool,
	expiresAt time.Time,
	createdAt time.Time,
) *DoctorInvitation {
	return &DoctorInvitation{
		id:               id,
		email:            email,
		specializationID: specializationID,
		token:            token,
		isUsed:           isUsed,
		expiresAt:        expiresAt,
		createdAt:        createdAt,
	}
}


func (i *DoctorInvitation) Validate() error {
	if i.isUsed {
		return ErrInviteUsed
	}
	if time.Now().UTC().After(i.expiresAt) {
		return ErrInviteExpired
	}
	return nil
}

func (i *DoctorInvitation) MarkAsUsed() error {
	if err := i.Validate(); err != nil {
		return err
	}
	i.isUsed = true
	return nil
}

func (i *DoctorInvitation) ID() string               { return i.id }
func (i *DoctorInvitation) Email() Email             { return i.email }
func (i *DoctorInvitation) SpecializationID() string { return i.specializationID }
func (i *DoctorInvitation) Token() string            { return i.token }
func (i *DoctorInvitation) IsUsed() bool             { return i.isUsed }
func (i *DoctorInvitation) ExpiresAt() time.Time     { return i.expiresAt }
func (i *DoctorInvitation) CreatedAt() time.Time {return i.createdAt}
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type InvitationStatus string

const (
	InvitationStatusPending InvitationStatus = "PENDING"
	InvitationStatusUsed    InvitationStatus = "USED"
	InvitationStatusExpired InvitationStatus = "EXPIRED"
)

var (
	ErrInvalidInvitationID     = errors.New("invitation ID cannot be empty")
	ErrInvalidEmail            = errors.New("invalid or empty receptionist email")
	ErrInvalidToken            = errors.New("invitation token cannot be empty")
	ErrInvitationAlreadyUsed   = errors.New("invitation token has already been used")
	ErrInvitationExpired       = errors.New("invitation token has expired")
	ErrInvitationNotPending    = errors.New("invitation is not in pending status")
)



type ReceptionistInvitation struct {
	id        string
	createrID  string
	email     string
	token     string
	status    InvitationStatus
	expiresAt time.Time
	createdAt time.Time
}

func NewReceptionistInvitation(
	doctorId string,
	email string,
	token string,
	now time.Time,
) (*ReceptionistInvitation, error) {

	if doctorId == "" {
		return nil, ErrInvalidDoctorID
	}
	if email == "" { 
		return nil, ErrInvalidEmail
	}
	if token == "" {
		return nil, ErrInvalidToken
	}
	id, err := uuid.NewV7()
	     if err != nil {
		return nil ,err
	     }

	utcNow := now.UTC()
	return &ReceptionistInvitation{
		id:        id.String(),
		createrID:  doctorId,
		email:     email,
		token:     token,
		status:    InvitationStatusPending,
		expiresAt: utcNow.Add(48 * time.Hour),
		createdAt: utcNow,
	}, nil
}




func ReconstituteReceptionistInvitation(
	id string,
	doctorId string,
	email string,
	token string,
	status InvitationStatus,
	expiresAt time.Time,
	createdAt time.Time,
) *ReceptionistInvitation {
	return &ReceptionistInvitation{
		id:        id,
		createrID:  doctorId,
		email:     email,
		token:     token,
		status:    status,
		expiresAt: expiresAt.UTC(),
		createdAt: createdAt.UTC(),
	}
}
func (inv *ReceptionistInvitation) IsValid(now time.Time) error {
	if inv.status == InvitationStatusUsed {
		return ErrInvitationAlreadyUsed
	}
	if inv.status == InvitationStatusExpired || now.UTC().After(inv.expiresAt) {
		return ErrInvitationExpired
	}
	if inv.status != InvitationStatusPending {
		return ErrInvitationNotPending
	}
	return nil
}

func (inv *ReceptionistInvitation) MarkAsUsed(now time.Time) error {
	if err := inv.IsValid(now); err != nil {
		return err
	}
	inv.status = InvitationStatusUsed
	return nil
}

func (inv *ReceptionistInvitation) MarkAsExpired(now time.Time) error {
	if inv.status == InvitationStatusUsed {
		return ErrInvitationAlreadyUsed
	}
	if now.UTC().Before(inv.expiresAt) && inv.status == InvitationStatusPending {
		return errors.New("cannot expire active invitation before expiration timestamp")
	}
	inv.status = InvitationStatusExpired
	return nil
}


func (inv *ReceptionistInvitation) ID() string               { return inv.id }
func (inv *ReceptionistInvitation) CreaterID() string         { return inv.createrID }
func (inv *ReceptionistInvitation) Email() string            { return inv.email }
func (inv *ReceptionistInvitation) Token() string            { return inv.token }
func (inv *ReceptionistInvitation) Status() InvitationStatus { return inv.status }
func (inv *ReceptionistInvitation) ExpiresAt() time.Time     { return inv.expiresAt }
func (inv *ReceptionistInvitation) CreatedAt() time.Time     { return inv.createdAt }
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Business Domain Errors
var (
	ErrInvalidDoctorID       = errors.New("doctor ID cannot be empty")
	ErrInvalidUserID         = errors.New("user ID cannot be empty")
	ErrInvalidSpecialization = errors.New("specialization ID cannot be empty")
	ErrInvalidLicense        = errors.New("license number is required")
	ErrInvalidFee            = errors.New("consultation fee cannot be negative")
	ErrDoctorInactive        = errors.New("doctor account is not active")
)

type DoctorStatus string

const (
	DoctorStatusActive   DoctorStatus = "ACTIVE"
	DoctorStatusOnLeave  DoctorStatus = "ON_LEAVE"
	DoctorStatusInactive DoctorStatus = "INACTIVE"
)

type Doctor struct {
	id               string
	userID           string
	specializationID string
	licenseNumber    string
	consultationFee  float64
	status           DoctorStatus
	isAvailable      bool
	createdAt        time.Time
	updatedAt        time.Time
}

func NewDoctor(userID, specializationID, licenseNumber string, fee float64) (*Doctor, error) {

	if userID == "" {
		return nil, ErrInvalidUserID
	}
	if specializationID == "" {
		return nil, ErrInvalidSpecialization
	}
	if licenseNumber == "" {
		return nil, ErrInvalidLicense
	}
	if fee < 0 {
		return nil, ErrInvalidFee
	}
		id, err := uuid.NewV7()
	     if err != nil {
		return nil ,err
	     }

	now := time.Now().UTC()
	return &Doctor{
		id:               id.String(),
		userID:           userID,
		specializationID: specializationID,
		licenseNumber:    licenseNumber,
		consultationFee:  fee,
		status:           DoctorStatusActive,
		isAvailable:      true,
		createdAt:        now,
		updatedAt:        now,
	}, nil
}


func ReconstituteDoctor(
	id string,
	userID string,
	specializationID string,
	licenseNumber string,
	consultationFee float64,
	status DoctorStatus,
	isAvailable bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Doctor {
	return &Doctor{
		id:               id,
		userID:           userID,
		specializationID: specializationID,
		licenseNumber:    licenseNumber,
		consultationFee:  consultationFee,
		status:           status,
		isAvailable:      isAvailable,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}


func (d *Doctor) UpdateConsultationFee(newFee float64) error {
	if newFee < 0 {
		return ErrInvalidFee
	}
	d.consultationFee = newFee
	d.updatedAt = time.Now().UTC()
	return nil
}

func (d *Doctor) ToggleAvailability(available bool) error {
	if d.status != DoctorStatusActive {
		return ErrDoctorInactive
	}
	d.isAvailable = available
	d.updatedAt = time.Now().UTC()
	return nil
}

func (d *Doctor) MarkOnLeave() {
	d.status = DoctorStatusOnLeave
	d.isAvailable = false
	d.updatedAt = time.Now().UTC()
}

func (d *Doctor) Activate() {
	d.status = DoctorStatusActive
	d.updatedAt = time.Now().UTC()
}

func (d *Doctor) Deactivate() {
	d.status = DoctorStatusInactive
	d.isAvailable = false
	d.updatedAt = time.Now().UTC()
}

// Getters
func (d *Doctor) ID() string               { return d.id }
func (d *Doctor) UserID() string           { return d.userID }
func (d *Doctor) SpecializationID() string { return d.specializationID }
func (d *Doctor) LicenseNumber() string    { return d.licenseNumber }
func (d *Doctor) ConsultationFee() float64 { return d.consultationFee }
func (d *Doctor) Status() DoctorStatus     { return d.status }
func (d *Doctor) IsAvailable() bool        { return d.isAvailable }
func (d *Doctor) CreatedAt() time.Time     { return d.createdAt }
func (d *Doctor) UpdatedAt() time.Time     { return d.updatedAt }
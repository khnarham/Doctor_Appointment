package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCannotUpdateWalkInDirectly = errors.New("walk-in patient credentials cannot be updated without linking an account")
	ErrInvalidContactDetails      = errors.New("name and contact details are required")
	ErrInvalidPatientName = errors.New("patient name cannot be empty or less than 2 characters")
	ErrInvalidAge         = errors.New("patient date of birth cannot be in the future or exceed 150 years")
	ErrWalkInIdentity     = errors.New("walk-in patient must have explicit contact name and phone")
	ErrRegisteredIdentity = errors.New("registered patient must be linked to a valid user id")
)



type Patient struct {
	id          string
	userID      *string 
	walkInName  string  
	walkInPhone string  
	gender      Gender
	dateOfBirth time.Time
	isWalkIn    bool
	createdAt   time.Time
}

func ReconstructPatient(
	id string,
	userID *string,
	walkInName string,
	walkInPhone string,
	gender Gender,
	dateOfBirth time.Time,
	isWalkIn bool,
	createdAt time.Time,
) *Patient {
	return &Patient{
		id:          id,
		userID:      userID,
		walkInName:  walkInName,
		walkInPhone: walkInPhone,
		gender:      gender,
		dateOfBirth: dateOfBirth,
		isWalkIn:    isWalkIn,
		createdAt:   createdAt,
	}
}


func (p *Patient) UpdateCredentials(name string , phone string , gender Gender , dataOfBirth time.Time) error {
	trimmedName := strings.TrimSpace(name)
	trimmedPhone := strings.TrimSpace(phone)

	if trimmedName == "" || len(trimmedName) < 2 {
		return ErrInvalidPatientName
	}

	if p.isWalkIn {
		if trimmedPhone == "" {
			return ErrInvalidContactDetails
		}
		p.walkInName = trimmedName
		p.walkInPhone = trimmedPhone
	} else {
		if p.userID == nil || *p.userID == "" {
			return ErrRegisteredIdentity
		}
	}
   
	if !gender.IsValid(){
		return InvalidGender
	}
	p.gender = gender
	p.dateOfBirth = dataOfBirth

	if err := p.Validate(); err != nil {
		return err
	}

	return nil
}


func NewWalkInPatient(id, name, phone string, gender Gender, dob time.Time) (*Patient, error) {
	p := &Patient{
		id:          id,
		walkInName:  strings.TrimSpace(name),
		walkInPhone: strings.TrimSpace(phone),
		gender:      gender,
		dateOfBirth: dob,
		isWalkIn:    true,
		createdAt:   time.Now().UTC(),
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func NewRegisteredPatient(userID string, gender Gender, dob time.Time , phone string) (*Patient, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrRegisteredIdentity
	}

     
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	p := &Patient{
		id:          id.String(),
		userID:      &userID,
		gender:      gender,
		dateOfBirth: dob,
		walkInPhone: phone,
		isWalkIn:    false,
		createdAt:   time.Now().UTC(),
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Patient) Validate() error {
	if p.isWalkIn {
		if len(p.walkInName) < 2 {
			return ErrInvalidPatientName
		}
		if p.walkInPhone == "" {
			return ErrWalkInIdentity
		}
	} else {
		if p.userID == nil || *p.userID == "" {
			return ErrRegisteredIdentity
		}
	}

	return p.validateAgeConstraint()
}

func (p *Patient) validateAgeConstraint() error {
	now := time.Now().UTC()
	if p.dateOfBirth.After(now) {
		return ErrInvalidAge
	}
	
	age := now.Year() - p.dateOfBirth.Year()
	if age > 150 {
		return ErrInvalidAge
	}
	return nil
}

func (p *Patient) CalculateAge() int {
	now := time.Now().UTC()
	age := now.Year() - p.dateOfBirth.Year()

	if now.YearDay() < p.dateOfBirth.YearDay() {
		age--
	}
	return age
}

func (p *Patient) LinkToUserAccount(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrRegisteredIdentity
	}
	
	p.userID = &userID
	p.isWalkIn = false
	p.walkInName = "" 
	p.walkInPhone = ""

	return p.Validate()
}


func (p *Patient) ID() string           { return p.id }
func (p *Patient) UserID() *string      { return p.userID }
func (p *Patient) IsWalkIn() bool       { return p.isWalkIn }
func (p *Patient) Gender() Gender       { return p.gender }
func (p *Patient) DateOfBirth() time.Time { return p.dateOfBirth }
func (p *Patient) WalkInName() string {return  p.walkInName}
func (p *Patient) WalkInPhone() string {return p.walkInPhone}
func (p *Patient) CreatedAt() time.Time {return p.createdAt}
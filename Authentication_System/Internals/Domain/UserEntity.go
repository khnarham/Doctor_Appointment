package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var AlreadyVerified = errors.New("The User is already verified")
var InvalidRole = errors.New("Invalid user role")
var InvalidGender = errors.New("Invalid gender provided")

type Role string

const (
	RolePatient Role = "PATIENT"
	RoleDoctor  Role = "DOCTOR"
	RoleAdmin   Role = "ADMIN"
	RoleReceptionst Role = "RECEPTIONIST"
)

type Gender string

const (
	MaleGender Gender = "Male"
	FemaleGender Gender = "Female"
)

type User struct {
	id string
	name string
    email Email
	gender Gender
	passwordHash string 
	isVerified bool
	role         Role
	createdAt time.Time
	updatedAt time.Time
}

func (r Role) IsValid() bool {
	switch r {
	case RolePatient, RoleDoctor, RoleAdmin:
		return true
	}
	return false
}

func (g Gender) IsValid() bool {
	switch g {
	case MaleGender, FemaleGender:
		return true
	}
	return false
}

func NewUser(name string , email Email , passwordHash string , role Role, gender Gender) (*User , error){
	if !role.IsValid() {
		return nil, InvalidRole
	}
	if !gender.IsValid() {
		return  nil , InvalidGender
	}
       id, err := uuid.NewV7()
	   if err != nil {
		    return nil , err
	   }
	   return &User{
		  id: id.String(),
		  gender: gender,
		  role: role,
		  name: name,
		  email: email,
		  passwordHash: passwordHash,
		  createdAt: time.Now().UTC(),
		  updatedAt: time.Now().UTC(),
	   }, nil
}


func ReconstituteUser(
	id string,
	name string,
	email Email,
	passwordHash string,
	isVerified bool,
	role Role,
	gender Gender,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {
	return &User{
		id:           id,
		name:         name,
		email:        email,
		passwordHash: passwordHash,
		isVerified:   isVerified,
		role:         role,
		gender:       gender,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}, nil
}


func (u *User) UpdateCradentials(name string , passwordHash string , gender Gender) error {
	if name == "" {
		return nil
	}
	if passwordHash == "" {
		return  nil
	}
	u.name = name;
	u.gender = gender
	u.passwordHash = passwordHash
	u.createdAt = time.Now().UTC()
	u.updatedAt = time.Now().UTC()
	return  nil
}

func (u *User) MarkAsVerified() error{
	if u.isVerified{
		return AlreadyVerified
	}
	u.isVerified = true;
	u.createdAt = time.Now().UTC();
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) ID() string { return u.id};
func (u *User) Name() string { return u.name};
func (u *User) Email() Email { return u.email};
func (u *User) PasswordHash() string { return u.passwordHash};
func (u *User) Role() Role {return u.role}
func (u *User) Gender() Gender {return u.gender}
func (u *User) IsVerified() bool { return u.isVerified};
func (u *User) CreatedAt() time.Time { return u.createdAt};
func (u *User) UpdatedAt() time.Time { return u.updatedAt};

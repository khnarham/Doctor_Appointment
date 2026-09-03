package usecase

import (
	"context"
	"errors"
	"my-appointment-app/Authentication_System/Internals/Domain"
	"time"
)

var ErrTokenExpired = errors.New("verification token has expired")
var ErrTokenMismatch = errors.New("verification token does not match")

type UserRegisterReposaitory interface {
	FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) 
	Save(ctx context.Context, user *domain.User) error
	
}


type PatientRegisterReposaitory interface {
	FindByUserId(ctx context.Context, id string) (*domain.Patient, error)
	FindByPhone(ctx context.Context , walkInPhone string) (*domain.Patient, error)
	Save(ctx context.Context, patient *domain.Patient) error
}



type EmailVerificationTokenRepoitory interface {
	FindByToken(ctx context.Context , token string) (*EmailVerificationData , error)
	FindByID(ctx context.Context, UserID string) (*domain.User , error)
	Save(ctx context.Context , user *domain.User) error
	Delete(ctx context.Context , TokenID string) error
}

type ReceptionistInvitationRepository interface {
	Save(ctx context.Context, inv *domain.ReceptionistInvitation) error
	GetByToken(ctx context.Context, token string) (*domain.ReceptionistInvitation, error)
	Update(ctx context.Context, inv *domain.ReceptionistInvitation) error
}

type UserReceptionistRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateReceptionistUser(ctx context.Context, email, hashedPassword, name string) (string, error) // Returns new receptionistID
}

type SendVerificationEmail interface{
	SendVerificationToken(ctx context.Context, to string, token string) error
}

type EmailVerificationRepository interface {
    SaveToken(ctx context.Context, userID string, token string, expiresAt time.Time) error
}
type PasswordHasher interface{
    Hash(password string) (string, error)
	Compare(password, hash string) error
}

type IToken interface {
	GenerateToken() string
}

type DoctorInvitationRepository interface {
	Save(ctx context.Context, invite *domain.DoctorInvitation) error
	FindByToken(ctx context.Context, token string) (*domain.DoctorInvitation, error)
	FindByEmail(ctx context.Context, email domain.Email) (*domain.DoctorInvitation, error)
}

type DoctorRepository interface {
	Save(ctx context.Context, doctor *domain.Doctor) error
}

type EmailSender interface {
	SendDoctorInviteEmail(ctx context.Context, toEmail string, token string) error
}

type EmailReceptionistSender interface {
	SendDoctorInviteEmail(ctx context.Context, toEmail string, token string) error
}



type UserRepository interface {
	// Single method: Directly creates receptionist user in `users` table with doctor_id foreign key
	CreateReceptionistUser(ctx context.Context, email, hashedPassword, fullName, doctorID string) (string, error)
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type TokenProvider interface {
	GenerateTokenPair(userID string , role domain.Role) (accessToken string, refreshToken string, err error)
	GenerateAccessToken(userID string , role domain.Role) (string, error)
	GenerateRefreshToken(userID string , role domain.Role) (string, error)
}

type ExistingUserRegistrationRepo interface{
	Update(ctx context.Context , user *domain.User) error
}
type ExistingPatientRegistrationRepo interface{
	Update(ctx context.Context , user *domain.Patient) error
}


type SpecializationRepository interface {
	Save(ctx context.Context, spec *domain.Specialization) error
	FindByName(ctx context.Context, name string) (*domain.Specialization, error)
	FindAll(ctx context.Context) ([]*domain.Specialization, error)
	FindByID(ctx context.Context , id string) (*domain.Specialization , error)
}


type RegisterUserInput struct {
	Name string 
	Email string 
	Password string
	Gender domain.Gender
	Phone string
	DateOfBirth time.Time
}

type RegisterUserOutput struct {
	Id string
	Name string
	Email string
	CreatedAt time.Time
	UpdatedAt time.Time
}


type EmailVerificationData struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
}


func (e *EmailVerificationData) IsValidToken(token string) (error){
	 if e.Token != token {
		return ErrTokenMismatch
	 }
	 return nil
}


func (e *EmailVerificationData) ValidateExpiration() error {
	if time.Now().After(e.ExpiresAt) {
		return ErrTokenExpired
	}
	return nil
}


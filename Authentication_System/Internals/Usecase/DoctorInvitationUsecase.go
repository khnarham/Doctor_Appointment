package usecase

import (
	"context"
	"errors"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"
)



type InviteDoctorInput struct {
	Email            string
	SpecializationID string
}

type InviteDoctorOutput struct {
	InviteID  string    `json:"invite_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RegisterDoctorWithInviteInput struct {
	Token           string
	Name            string
	Password        string
	Gender          domain.Gender
	LicenseNumber   string
	ConsultationFee float64
}



var (
	ErrUserAlreadyExists   = errors.New("a user with this email already exists")
	ErrInviteAlreadyActive = errors.New("an active invitation already exists for this email")
	ErrSpecializationNotFound = errors.New("specialization not found")
)

const InviteExpirationDuration = 48 * time.Hour

type InviteDoctorUsecase struct {
	inviteRepo DoctorInvitationRepository
	userRepo   UserRegisterReposaitory
	specRepo   SpecializationRepository
	emailServ  SendVerificationEmail
	generator  IToken
}

func NewInviteDoctorUsecase(
	inviteRepo DoctorInvitationRepository,
	userRepo UserRegisterReposaitory,
	specRepo SpecializationRepository,
	emailServ SendVerificationEmail,
	generator IToken,
) *InviteDoctorUsecase {
	return &InviteDoctorUsecase{
		inviteRepo: inviteRepo,
		userRepo:   userRepo,
		specRepo:   specRepo,
		emailServ:  emailServ,
		generator:  generator,
	}
}

func (uc *InviteDoctorUsecase) Execute(ctx context.Context, input InviteDoctorInput) (*InviteDoctorOutput, error) {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	specialization, err := uc.specRepo.FindByID(ctx, input.SpecializationID)
    if err != nil {
    return nil, err
    }
    if specialization == nil {
    return nil,  ErrSpecializationNotFound
    }
     
	existingUser, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	existingInvite, err := uc.inviteRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingInvite != nil && !existingInvite.IsUsed() && time.Now().UTC().Before(existingInvite.ExpiresAt()) {
		return nil, ErrInviteAlreadyActive
	}


	token := uc.generator.GenerateToken()
	invite, err := domain.NewDoctorInvitation(email, input.SpecializationID, token, InviteExpirationDuration)
	if err != nil {
		return nil, err
	}

	if err := uc.inviteRepo.Save(ctx, invite); err != nil {
		return nil, err
	}

	if err := uc.emailServ.SendVerificationToken(ctx, email.Value(), token); err != nil {
		return nil, err
	}

	return &InviteDoctorOutput{
		InviteID:  invite.ID(),
		Email:     invite.Email().Value(),
		Token:     invite.Token(),
		ExpiresAt: invite.ExpiresAt(),
	}, nil
}
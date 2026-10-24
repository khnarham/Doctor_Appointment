package usecase

import (
	"context"
	"errors"
	"fmt"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"
)

type InviteReceptionistInput struct {
	DoctorID string
	Email    string
}

type InviteReceptionistOutput struct {
	InvitationID string
	Status       string
}

type InviteReceptionistUsecase struct {
	invRepo     ReceptionistInvitationRepository
	userRepo    UserReceptionistRepository
	emailSender EmailReceptionistSender
    generater  IToken

}

func NewInviteReceptionistUsecase(
	invRepo ReceptionistInvitationRepository,
	userRepo UserReceptionistRepository,
	emailSender EmailReceptionistSender,
	generater  IToken,
) *InviteReceptionistUsecase {
	return &InviteReceptionistUsecase{
		invRepo:     invRepo,
		userRepo:    userRepo,
		emailSender: emailSender,
		generater: generater,
	}
}

func (uc *InviteReceptionistUsecase) Execute(ctx context.Context, input InviteReceptionistInput) (*InviteReceptionistOutput, error) {
	emailExists, err := uc.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if emailExists {
		return nil, errors.New("user with this email already exists")
	}

	
    token := uc.generater.GenerateToken()
    now := time.Now()
	invitation, err := domain.NewReceptionistInvitation(
		input.DoctorID,
		input.Email,
		token,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	if err := uc.invRepo.Save(ctx, invitation); err != nil {
		return nil, fmt.Errorf("failed to save invitation: %w", err)
	}

	if err := uc.emailSender.SendDoctorInviteEmail(ctx , input.Email , token); err != nil {
		return nil, fmt.Errorf("failed to dispatch invite email: %w", err)
	}

	return &InviteReceptionistOutput{
		InvitationID: invitation.ID(),
		Status:       string(invitation.Status()),
	}, nil
}
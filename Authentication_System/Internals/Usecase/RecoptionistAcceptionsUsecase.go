package usecase

import (
	"context"
	"fmt"
	"time"
)

type AcceptInvitationInput struct {
	Token    string
	FullName string
	Password string
}

type AcceptInvitationOutput struct {
	ReceptionistID string
	Email          string
	DoctorID       string
}

type AcceptInvitationUsecase struct {
	invRepo  ReceptionistInvitationRepository
	userRepo UserRepository
	hasher   PasswordHasher
}

func NewAcceptInvitationUsecase(
	invRepo ReceptionistInvitationRepository,
	userRepo UserRepository,
	hasher PasswordHasher,
) *AcceptInvitationUsecase {
	return &AcceptInvitationUsecase{
		invRepo:  invRepo,
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (uc *AcceptInvitationUsecase) Execute(ctx context.Context, input AcceptInvitationInput) (*AcceptInvitationOutput, error) {
	now := time.Now().UTC()

	inv, err := uc.invRepo.GetByToken(ctx, input.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if err := inv.IsValid(now); err != nil {
		return nil, fmt.Errorf("invitation invalid: %w", err)
	}

	hashedPassword, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	receptionistID, err := uc.userRepo.CreateReceptionistUser(
		ctx,
		inv.Email(),
		hashedPassword,
		input.FullName,
		inv.DoctorID(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create receptionist user: %w", err)
	}

	if err := inv.MarkAsUsed(now); err != nil {
		return nil, err
	}

	if err := uc.invRepo.Update(ctx, inv); err != nil {
		return nil, fmt.Errorf("failed to update token status: %w", err)
	}

	return &AcceptInvitationOutput{
		ReceptionistID: receptionistID,
		Email:          inv.Email(),
		DoctorID:       inv.DoctorID(),
	}, nil
}
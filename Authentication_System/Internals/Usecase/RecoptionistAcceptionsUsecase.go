package usecase

import (
	"context"
	"errors"
	"time"
)

// Custom Sentinel Errors
var (
	ErrInvitationNotFound  = errors.New("invitation token not found or invalid")
	ErrInvitationInvalid   = errors.New("invitation token is expired or already used")
	ErrPasswordHashing     = errors.New("failed to hash receptionist password")
	ErrReceptionistCreationFailed = errors.New("failed to create receptionist user")
	ErrTokenUpdateFailed   = errors.New("failed to update invitation token status")
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
		return nil,  ErrInvitationNotFound
	}

	if err := inv.IsValid(now); err != nil {
		return nil,  ErrInvitationInvalid
	}

	hashedPassword, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, ErrPasswordHashing
	}

	receptionistID, err := uc.userRepo.CreateReceptionistUser(
		ctx,
		inv.Email(),
		hashedPassword,
		input.FullName,
		inv.CreaterID(),
	)
	if err != nil {
		return nil, ErrReceptionistCreationFailed
	}

	if err := inv.MarkAsUsed(now); err != nil {
		return nil, err
	}

	if err := uc.invRepo.Update(ctx, inv); err != nil {
		return nil, ErrTokenUpdateFailed
	}

	return &AcceptInvitationOutput{
		ReceptionistID: receptionistID,
		Email:          inv.Email(),
		DoctorID:       inv.CreaterID(),
	}, nil
}
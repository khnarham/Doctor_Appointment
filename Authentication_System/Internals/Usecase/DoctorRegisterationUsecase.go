package usecase

import (
	"context"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
)

type RegisterDoctorWithInviteUsecase struct {
	inviteRepo DoctorInvitationRepository
	userRepo   UserRegisterReposaitory
	doctorRepo DoctorRepository
	hasher     PasswordHasher
	txManager  TransactionManager
}

func NewRegisterDoctorWithInviteUsecase(
	inviteRepo DoctorInvitationRepository,
	userRepo UserRegisterReposaitory,
	doctorRepo DoctorRepository,
	hasher PasswordHasher,
	txManager TransactionManager,
) *RegisterDoctorWithInviteUsecase {
	return &RegisterDoctorWithInviteUsecase{
		inviteRepo: inviteRepo,
		userRepo:   userRepo,
		doctorRepo: doctorRepo,
		hasher:     hasher,
		txManager:  txManager,
	}
}

func (uc *RegisterDoctorWithInviteUsecase) Execute(ctx context.Context, input RegisterDoctorWithInviteInput) (*RegisterUserOutput, error) {
	invite, err := uc.inviteRepo.FindByToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, domain.ErrInviteExpired 
	}

	if err := invite.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(input.Name, invite.Email(),  hashedPassword, domain.RoleDoctor , input.Gender)
	if err != nil {
		return nil, err
	}

	doctor, err := domain.NewDoctor(
		user.ID(),
		invite.SpecializationID(),
		input.LicenseNumber,
		input.ConsultationFee,
	)
	if err != nil {
		return nil, err
	}

	err = uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {

		if err := user.MarkAsVerified(); err != nil{
			return  err
		} 
		if err := uc.userRepo.Save(txCtx, user); err != nil {
			return err
		}

		if err := uc.doctorRepo.Save(txCtx, doctor); err != nil {
			return err
		}

		if err := invite.MarkAsUsed(); err != nil {
			return err
		}

		return uc.inviteRepo.Save(txCtx, invite)
	})

	if err != nil {
		return nil, err
	}

	return &RegisterUserOutput{
		Id:        user.ID(),
		Name:      user.Name(),
		Email:     user.Email().Value(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}, nil
}
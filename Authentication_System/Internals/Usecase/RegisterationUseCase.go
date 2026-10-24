package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "my-appointment-app/Authentication_System/Internals/Domain"
)

	var ErrEmailAlreadyExists = errors.New("email is already registered")



type RegisterUserUsecase struct {
    userRepo UserRegisterReposaitory
    hasher PasswordHasher
    patientRepo PatientRegisterReposaitory
    txManager TransactionManager
    verificationRepo EmailVerificationRepository
    generator IToken
    sendEmail SendVerificationEmail
	ExistedUserRepo ExistingUserRegistrationRepo
	ExitedPatientRepo ExistingPatientRegistrationRepo
}

const TokenExpirationDuration = 24 * time.Hour

func NewRegisterUserUsecase(userRepo UserRegisterReposaitory , hash PasswordHasher , patientRepo PatientRegisterReposaitory, VerRepo EmailVerificationRepository , generater IToken , txManager TransactionManager , sendEmail SendVerificationEmail , ExistedUserRepo ExistingUserRegistrationRepo , ExitedPatientRepo ExistingPatientRegistrationRepo) *RegisterUserUsecase {
    return &RegisterUserUsecase{
    userRepo: userRepo,
    hasher: hash,
    patientRepo: patientRepo,
    txManager: txManager,
    verificationRepo: VerRepo,
    generator: generater,
    sendEmail: sendEmail,
	ExistedUserRepo: ExistedUserRepo,
	ExitedPatientRepo: ExitedPatientRepo,
};

}

func (r *RegisterUserUsecase) Execute(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	existingUser, err := r.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrEmailAlreadyExists
	}

	if existingUser != nil && existingUser.IsVerified() {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := r.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now().UTC()
	token := r.generator.GenerateToken()
	expiresAt := now.Add(TokenExpirationDuration)

	var targetUser *domain.User

	if existingUser != nil {
		user, err := r.handleExistingUnverifiedUser(ctx, existingUser, input, token, expiresAt)
		if err != nil {
			return nil, err
		}
		targetUser = user
	} else {
		user, err := r.handleNewUserRegistration(ctx, email, hashedPassword, input, token, expiresAt)
		if err != nil {
			return nil, err
		}
		targetUser = user
	}

	

	return &RegisterUserOutput{
		Id:        targetUser.ID(),
		Name:      targetUser.Name(),
		Email:     targetUser.Email().Value(),
		CreatedAt: targetUser.CreatedAt(),
		UpdatedAt: targetUser.UpdatedAt(),
	}, nil
}

func (r *RegisterUserUsecase) handleExistingUnverifiedUser(
	ctx context.Context,
	user *domain.User,
	input RegisterUserInput,
	token string,
	expiresAt time.Time,
) (*domain.User, error) {
	patient, err := r.patientRepo.FindByUserId(ctx, user.ID())
	if err != nil {
		return nil, err
	}

	if err := user.UpdateCradentials(input.Name, user.PasswordHash(), input.Gender); err != nil {
		return nil, err
	}
	if err := patient.UpdateCredentials(input.Name, input.Phone, input.Gender, input.DateOfBirth,); err != nil {
		return nil, err
	}

	err = r.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := r.ExistedUserRepo.Update(txCtx, user); err != nil {
			return err
		}
		if err := r.ExitedPatientRepo.Update(txCtx, patient); err != nil {
			return err
		}
		if err := r.verificationRepo.SaveToken(txCtx, user.ID(), token, expiresAt); err != nil {
			return err
		}
		if err := r.sendEmail.SendVerificationToken(ctx, input.Email, token); err != nil {
		return err
	}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("tx failed for existing user: %w", err)
	}

	return user, nil
}

func (r *RegisterUserUsecase) handleNewUserRegistration(
	ctx context.Context,
	email domain.Email,
	hashedPassword string,
	input RegisterUserInput,
	token string,
	expiresAt time.Time,
) (*domain.User, error) {
	user, err := domain.NewUser(input.Name, email, hashedPassword, domain.RolePatient, input.Gender)
	if err != nil {
		return nil, err
	}

	var patient *domain.Patient
	existingWalkIn, fetchErr := r.patientRepo.FindByPhone(ctx, input.Phone)

	if fetchErr == nil && existingWalkIn != nil && existingWalkIn.IsWalkIn() {
		if err := existingWalkIn.LinkToUserAccount(user.ID()); err != nil {
			return nil, err
		}
		patient = existingWalkIn
	} else {
		patient, err = domain.NewRegisteredPatient(user.ID(), input.Gender, input.DateOfBirth, input.Phone)
		if err != nil {
			return nil, err
		}
	}

	err = r.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := r.userRepo.Save(txCtx, user); err != nil {
			return err
		}
		if err := r.patientRepo.Save(txCtx, patient); err != nil {
			return err
		}
		if err := r.verificationRepo.SaveToken(txCtx, user.ID(), token, expiresAt); err != nil {
			return err
		}
		if err := r.sendEmail.SendVerificationToken(ctx, input.Email, token); err != nil {
		return err
	}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("tx failed for new user: %w", err)
	}

	return user, nil
}
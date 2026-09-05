package usecase

import (
	"context"
	"errors"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailNotVerified   = errors.New("email is not verified")
)

type LoginUserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}



type LoginUsecase struct {
	userRepo      UserRegisterReposaitory
	passwordHasher PasswordHasher
	tokenGenerator TokenProvider
}

func NewLoginUsecase(
	ur UserRegisterReposaitory,
	ph PasswordHasher,
	tp TokenProvider,
) *LoginUsecase {
	return &LoginUsecase{
		userRepo:       ur,
		passwordHasher: ph,
		tokenGenerator: tp,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

type UserOutput struct {
	ID string   `json:"user_id"`
	Name string   `json:"name"`
    Email string `json:"email"`
	Gender string   `json:"gender"`
	IsVerified bool        `json:"is_verified"`
	Role         string   `json:"role"`
}



type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	User        UserOutput
	LoggedInAt   time.Time
}

func (uc *LoginUsecase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	email , err := domain.NewEmail(input.Email)
	if err != nil {
		return nil , err
	}
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := uc.passwordHasher.Compare(input.Password , user.PasswordHash()); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsVerified() {
		return nil, ErrEmailNotVerified
	}

	accessToken, refreshToken, err := uc.tokenGenerator.GenerateTokenPair(user.ID() , user.Role())
	if err != nil {
		return nil, err
	}
   
   

	
	UserOutput := UserOutput{
		ID: user.ID(),
		Name: user.Name(),
		Email: string(user.Email().Value()),
		Gender: string(user.Gender()),
		Role: string(user.Role()),
	}
	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:  UserOutput,
		LoggedInAt:   time.Now().UTC(),
	}, nil
}
package usecase

import (
	"context"
	"time"

)


type AdminLoginInput struct {
	Email    string
	Password string
}


type AuthUseCase struct {
	userRepo  AdminRepository
	passwordHasher PasswordHasher
	tokenGenerator TokenProvider

}

func NewAuthUseCase(userRepo AdminRepository , 	tokenGenerator TokenProvider,passwordHasher PasswordHasher) *AuthUseCase  {
	return &AuthUseCase{
		userRepo:  userRepo, 
		tokenGenerator: tokenGenerator,
		passwordHasher: passwordHasher,
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, input AdminLoginInput ) (*VerificationOutput, error) {
	user, err := uc.userRepo.FindByEmail(input.Email)
	if err != nil {
		return &VerificationOutput{}, err
	}

	// 2. Verify password hash using bcrypt
	if err := uc.passwordHasher.Compare(input.Password , user.PasswordHash()); err != nil {
		return &VerificationOutput{}, ErrInvalidCredentials
	}

	// 3. Generate JWT Token
 
	accessToken , refreshToken , err := uc.tokenGenerator.GenerateTokenPair(user.ID() , user.Role());
	 if err != nil {
		return &VerificationOutput{},err
	 }

	if err != nil {
		return  nil , err
	}

	     return &VerificationOutput{
          AccessToken: accessToken,
		  RefreshToken: refreshToken,
		  Role: string(user.Role()),
		  CreatedAt: time.Now().UTC(),
		  UpdatedAt: time.Now().UTC(),
	 }, nil

}


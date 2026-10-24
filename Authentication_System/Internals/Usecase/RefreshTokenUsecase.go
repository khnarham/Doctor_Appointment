package usecase

import (
	"context"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	jwt "my-appointment-app/Authentication_System/Internals/Infrastructure/JWT"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenServiceContract interface {
	ValidateRefreshToken(tokenStr string) (*jwt.TokenClaims, error)
	GenerateTokenPair(userID string, role domain.Role) (string, string, error)
}

type RefreshTokenUsecase struct {
	jwtService TokenServiceContract
}

func NewRefreshTokenUsecase(jwtService TokenServiceContract) *RefreshTokenUsecase {
	return &RefreshTokenUsecase{jwtService: jwtService}
}

func (u *RefreshTokenUsecase) Execute(ctx context.Context, refreshTokenStr string) (*TokenPair, error) {
	claims, err := u.jwtService.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, err
	}

	access, refresh, err := u.jwtService.GenerateTokenPair(claims.UserID, claims.Role)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
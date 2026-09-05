package usecase

import (
	"context"
	"errors"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
)

var ErrUserNotFound = errors.New("user not found")

type GetMeUserRepository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
}

type GetMeOutput struct {
	ID         string `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Gender     string `json:"gender"`
	IsVerified bool   `json:"is_verified"`
	Role       string `json:"role"`
}



type GetMeUsecase struct {
	userRepo GetMeUserRepository
}

func NewGetMeUsecase(ur GetMeUserRepository) *GetMeUsecase {
	return &GetMeUsecase{
		userRepo: ur,
	}
}

func (uc *GetMeUsecase) Execute(ctx context.Context, userID string) (*GetMeOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &GetMeOutput{
		ID:         user.ID(),
		Name:       user.Name(),
		Email:      string(user.Email().Value()),
		Gender:     string(user.Gender()),
		IsVerified: user.IsVerified(),
		Role:       string(user.Role()),
	}, nil
}
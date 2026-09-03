package usecase

import (
	"context"
	"errors"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"
)

type CreateSpecializationInput struct {
	Name        string
	Description string
}

type SpecializationOutput struct {
	ID          string    
	Name        string   
	Description string    
	CreatedAt   time.Time 
	UpdatedAt   time.Time 
}


var ErrSpecializationAlreadyExists = errors.New("specialization with this name already exists")

type CreateSpecializationUsecase struct {
	specRepo SpecializationRepository
}

func NewCreateSpecializationUsecase(repo SpecializationRepository) *CreateSpecializationUsecase {
	return &CreateSpecializationUsecase{
		specRepo: repo,
	}
}

func (uc *CreateSpecializationUsecase) Execute(ctx context.Context, input CreateSpecializationInput) (*SpecializationOutput, error) {

	existing, err := uc.specRepo.FindByName(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrSpecializationAlreadyExists
	}

	spec, err := domain.NewSpecialization(input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	if err := uc.specRepo.Save(ctx, spec); err != nil {
		return nil, err
	}

	return &SpecializationOutput{
		ID:          spec.ID(),
		Name:        spec.Name(),
		Description: spec.Description(),
		CreatedAt:   spec.CreatedAt(),
		UpdatedAt:   spec.UpdatedAt(),
	}, nil
}
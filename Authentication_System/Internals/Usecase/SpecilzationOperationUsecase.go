package usecase

import (
	"context"
	"errors"
	domain "my-appointment-app/Authentication_System/Internals/Domain"
)

type SpecializationUsecase struct {
	repo SpecializationOperationRepository
}

func NewSpecializationUsecase(repo SpecializationOperationRepository) *SpecializationUsecase {
	return &SpecializationUsecase{repo: repo}
}

type UpdateSpecializationInput struct {
	ID          string
	Name        string
	Description string
}

func (u *SpecializationUsecase) Update(ctx context.Context, input UpdateSpecializationInput) (*domain.Specialization, error) {
	spec, err := u.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Name != spec.Name() {
		existing, err := u.repo.GetByName(ctx, input.Name)
		if err == nil && existing.ID() != spec.ID() {
			return nil, domain.ErrSpecializationNameExists
		}
		if err != nil && !errors.Is(err, domain.ErrSpecializationNotFound) {
			return nil, err
		}
	}

	if err := spec.Update(input.Name, input.Description); err != nil {
		return nil, err
	}

	if err := u.repo.Update(ctx, spec); err != nil {
		return nil, err
	}

	return spec, nil
}

func (u *SpecializationUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}
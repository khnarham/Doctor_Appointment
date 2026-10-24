package usecase

import "context"

type ListSpecializationsUsecase struct {
	specRepo SpecializationRepository
}

func NewListSpecializationsUsecase(repo SpecializationRepository) *ListSpecializationsUsecase {
	return &ListSpecializationsUsecase{
		specRepo: repo,
	}
}

func (uc *ListSpecializationsUsecase) Execute(ctx context.Context) ([]*SpecializationOutput, error) {
	specs, err := uc.specRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var output []*SpecializationOutput
	for _, spec := range specs {
		output = append(output, &SpecializationOutput{
			ID:          spec.ID(),
			Name:        spec.Name(),
			Description: spec.Description(),
			CreatedAt:   spec.CreatedAt(),
			UpdatedAt:   spec.UpdatedAt(),
		})
	}

	return output, nil
}
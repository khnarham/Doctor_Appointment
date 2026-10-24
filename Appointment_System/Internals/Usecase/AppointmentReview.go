package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "my-appointment-app/Appointment_System/Internals/Domain"
)

var (
	ErrAppointmentNotFound    = errors.New("appointment record not found")
	ErrUnauthorizedReview     = errors.New("you are not authorized to review this appointment")
	ErrAppointmentNotInCompletedState = errors.New("cannot review an appointment that is not completed")
	ErrReviewAlreadySubmitted = errors.New("a review for this appointment has already been submitted")
	ErrInvalidRatingRange     = errors.New("rating must be an integer between 1 and 5")
	ErrReviewPeriodExpired    = errors.New("review submission window (30 days post completion) has expired")
)



type ReviewRepository interface {
	ExistsByAppointmentID(ctx context.Context, appointmentID string) (bool, error)
	Save(ctx context.Context, review *domain.Review) error
}


type SubmitReviewInput struct {
	AuthenticatedPatientID string // Extract from JWT Auth Middleware Context
	AppointmentID          string
	Rating                 int
	Comment                string
}

type SubmitReviewUsecase struct {
	appointmentRepo AppointmentRepository
	reviewRepo      ReviewRepository
	idGen           IDGenerator
}

func NewSubmitReviewUsecase(
	appRepo AppointmentRepository,
	revRepo ReviewRepository,
) *SubmitReviewUsecase {
	return &SubmitReviewUsecase{
		appointmentRepo: appRepo,
		reviewRepo:      revRepo,
	}
}

func (uc *SubmitReviewUsecase) Execute(ctx context.Context, input SubmitReviewInput) error {
	now := time.Now().UTC()

	if input.Rating < 1 || input.Rating > 5 {
		return ErrInvalidRatingRange
	}

	app, err := uc.appointmentRepo.GetByID(ctx, input.AppointmentID)
	if err != nil {
		return fmt.Errorf("failed to fetch appointment: %w", err)
	}
	if app == nil {
		return ErrAppointmentNotFound
	}
	if app.PatientID() != input.AuthenticatedPatientID {
		return ErrUnauthorizedReview
	}

	if app.Status() != domain.StatusCompleted {
		return ErrAppointmentNotInCompletedState
	}



	exists, err := uc.reviewRepo.ExistsByAppointmentID(ctx, input.AppointmentID)
	if err != nil {
		return fmt.Errorf("database check for existing review failed: %w", err)
	}
	if exists {
		return ErrReviewAlreadySubmitted
	}


	review, err := domain.NewReview(
		app.ID(),
		app.PatientID(),
		app.DoctorID(),
		input.Rating,
		input.Comment,
		now,
	)
	if err != nil {
		return fmt.Errorf("domain instantiation failed: %w", err)
	}


	if err := uc.reviewRepo.Save(ctx, review); err != nil {
		return fmt.Errorf("failed to persist review: %w", err)
	}

	return nil
}
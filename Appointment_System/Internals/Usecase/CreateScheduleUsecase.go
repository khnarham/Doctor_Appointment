package usecase

import (
	"context"
	domain "my-appointment-app/Appointment_System/Internals/Domain"
	"time"
)


type CreateScheduleInput struct {
	DoctorID     string
	Days         domain.DayOfWeek
	TimeRange    domain.TimeRange
	SlotDuration int  
}

type CreateScheduleOutput struct {
	ID           string   
	DoctorID     string   
	Days         string 
	StartTime    string
	EndTime      string
	SlotDuration int      
	IsActive     bool     
}


type CreateScheduleUsecase struct {
	scheduleRepo ScheduleRepository
}

func NewCreateScheduleUsecase(scheduleRepo ScheduleRepository) *CreateScheduleUsecase {
	return &CreateScheduleUsecase{
		scheduleRepo: scheduleRepo,
	}
}

func (cs *CreateScheduleUsecase) Execute(ctx context.Context , cmd CreateScheduleInput) (CreateScheduleOutput ,error){
	existingSchedules , err := cs.scheduleRepo.GetActiveSchedulesByDoctorID(ctx , cmd.DoctorID); if err != nil{
		return CreateScheduleOutput{} , err
	}

	newSchedule, err := domain.NewDoctorSchedule(
		cmd.DoctorID,
		cmd.Days,
		cmd.TimeRange,
		time.Duration(cmd.SlotDuration)*time.Minute,
	)

	if err != nil {
		return CreateScheduleOutput{} , err
	}
	
	if err := newSchedule.CheckOverlapWith(existingSchedules); err != nil{
		return CreateScheduleOutput{} ,err
	}

	if err := cs.scheduleRepo.Save(ctx, newSchedule); err != nil {
		return CreateScheduleOutput{} ,err
	}

   return CreateScheduleOutput{
		ID:           newSchedule.ID(),
		DoctorID:     newSchedule.DoctorID(),
		Days:         newSchedule.Day().String(),
		StartTime:    newSchedule.TimeRange().Start().String(),
		EndTime:      newSchedule.TimeRange().End().String(),
		SlotDuration: int(newSchedule.SlotDuration().Minutes()),
		IsActive:     newSchedule.IsActive(),
	}, nil
}
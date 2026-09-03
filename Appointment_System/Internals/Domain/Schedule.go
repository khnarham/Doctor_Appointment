package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyScheduleID  = errors.New("schedule ID cannot be empty")
	ErrEmptyDoctorID    = errors.New("doctor ID cannot be empty")
	ErrInvalidClockTime = errors.New("minutes must be between 0 and 1439 (00:00 - 23:59)")
	ErrInvalidTimeRange = errors.New("start clock time must be strictly before end time")
	ErrInvalidDuration  = errors.New("slot duration must be greater than zero")
	ErrScheduleOverlap  = errors.New("new schedule overlaps with an existing active schedule")
	ErrScheduleInactive = errors.New("cannot perform action on an inactive schedule")
)

type DayOfWeek time.Weekday

const (
	Sunday    DayOfWeek = DayOfWeek(time.Sunday)
	Monday    DayOfWeek = DayOfWeek(time.Monday)
	Tuesday   DayOfWeek = DayOfWeek(time.Tuesday)
	Wednesday DayOfWeek = DayOfWeek(time.Wednesday)
	Thursday  DayOfWeek = DayOfWeek(time.Thursday)
	Friday    DayOfWeek = DayOfWeek(time.Friday)
	Saturday  DayOfWeek = DayOfWeek(time.Saturday)
)

type ClockTime struct {
	minutes int
}

func NewClockTime(hour, minute int) (ClockTime, error) {
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return ClockTime{}, ErrInvalidClockTime
	}
	return ClockTime{minutes: hour*60 + minute}, nil
}

func NewClockTimeFromMinutes(mins int) (ClockTime, error) {
	if mins < 0 || mins >= 1440 {
		return ClockTime{}, ErrInvalidClockTime
	}
	return ClockTime{minutes: mins}, nil
}

func (ct ClockTime) Minutes() int { return ct.minutes }
func (ct ClockTime) String() string {
	return fmt.Sprintf("%02d:%02d", ct.minutes/60, ct.minutes%60)
}

type TimeRange struct {
	start ClockTime
	end   ClockTime
}

func NewTimeRange(start, end ClockTime) (TimeRange, error) {
	if start.Minutes() >= end.Minutes() {
		return TimeRange{}, ErrInvalidTimeRange
	}
	return TimeRange{start: start, end: end}, nil
}

func (tr TimeRange) Start() ClockTime { return tr.start }
func (tr TimeRange) End() ClockTime   { return tr.end }

func (tr TimeRange) Overlaps(other TimeRange) bool {
	return tr.start.Minutes() < other.end.Minutes() && tr.end.Minutes() > other.start.Minutes()
}

type SlotInterval struct {
	Start time.Time
	End   time.Time
}

type DoctorSchedule struct {
	id           string
	doctorID     string
	day          DayOfWeek
	timeRange    TimeRange
	slotDuration time.Duration
	isActive     bool
	createdAt    time.Time
	updatedAt    time.Time
}

func NewDoctorSchedule(
	doctorID string,
	day DayOfWeek,
	timeRange TimeRange,
	slotDuration time.Duration,
) (*DoctorSchedule, error) {

	if doctorID == "" {
		return nil, ErrEmptyDoctorID
	}
	if slotDuration <= 0 {
		return nil, ErrInvalidDuration
	}
	id, err := uuid.NewV7()
	     if err != nil {
		return nil ,err
	     }

	now := time.Now().UTC()
	return &DoctorSchedule{
		id:           id.String(),
		doctorID:     doctorID,
		day:          day,
		timeRange:    timeRange,
		slotDuration: slotDuration,
		isActive:     true,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func ReconstituteSchedule(
	id string,
	doctorID string,
	day DayOfWeek,
	timeRange TimeRange,
	slotDuration time.Duration,
	isActive bool,
	createdAt, updatedAt time.Time,
) *DoctorSchedule {
	return &DoctorSchedule{
		id:           id,
		doctorID:     doctorID,
		day:          day,
		timeRange:    timeRange,
		slotDuration: slotDuration,
		isActive:     isActive,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (s *DoctorSchedule) Update(newDay DayOfWeek, newTimeRange TimeRange, newDuration time.Duration) error {
	if !s.isActive {
		return ErrScheduleInactive
	}
	if newDuration <= 0 {
		return ErrInvalidDuration
	}

	s.day = newDay
	s.timeRange = newTimeRange
	s.slotDuration = newDuration
	s.updatedAt = time.Now().UTC()
	return nil
}

func (s *DoctorSchedule) Deactivate() error {
	if !s.isActive {
		return ErrScheduleInactive
	}
	s.isActive = false
	s.updatedAt = time.Now().UTC()
	return nil
}

func (s *DoctorSchedule) CheckOverlapWith(existingSchedules []DoctorSchedule) error {
	for _, existing := range existingSchedules {
		if !existing.isActive || existing.id == s.id {
			continue
		}
		if s.day == existing.day && s.timeRange.Overlaps(existing.timeRange) {
			return ErrScheduleOverlap
		}
	}
	return nil
}

func (s *DoctorSchedule) IsAvailableOn(day DayOfWeek) bool {
	if !s.isActive {
		return false
	}
	return s.day == day
}

func (s *DoctorSchedule) GenerateSlotsForDate(targetDate time.Time) ([]SlotInterval, error) {
	if !s.isActive {
		return nil, ErrScheduleInactive
	}


	year, month, day := targetDate.Date()
	baseDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	startBoundary := baseDate.Add(time.Duration(s.timeRange.start.Minutes()) * time.Minute)
	endBoundary := baseDate.Add(time.Duration(s.timeRange.end.Minutes()) * time.Minute)

	var slots []SlotInterval
	current := startBoundary

	for {
		next := current.Add(s.slotDuration)
		if next.After(endBoundary) {
			break
		}
		slots = append(slots, SlotInterval{
			Start: current,
			End:   next,
		})
		current = next
	}

	return slots, nil
}


func ParseDayOfWeek(s string) (DayOfWeek, error) {
	switch s {
	case "SUNDAY":
		return Sunday, nil
	case "MONDAY":
		return Monday, nil
	case "TUESDAY":
		return Tuesday, nil
	case "WEDNESDAY":
		return Wednesday, nil
	case "THURSDAY":
		return Thursday, nil
	case "FRIDAY":
		return Friday, nil
	case "SATURDAY":
		return Saturday, nil
	default:
		return 0, fmt.Errorf("invalid day of week: %s", s)
	}
}

func (d DayOfWeek) String() string {
	days := [...]string{"SUNDAY", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY"}
	if int(d) < 0 || int(d) >= len(days) {
		return "UNKNOWN"
	}
	return days[d]
}

func (s *DoctorSchedule) ID() string                   { return s.id }
func (s *DoctorSchedule) DoctorID() string             { return s.doctorID }
func (s *DoctorSchedule) Day() DayOfWeek               { return s.day }
func (s *DoctorSchedule) IsActive() bool               { return s.isActive }
func (s *DoctorSchedule) TimeRange() TimeRange         { return s.timeRange }
func (s *DoctorSchedule) SlotDuration() time.Duration { return s.slotDuration }
func (s *DoctorSchedule) CreatedAt() time.Time         { return s.createdAt }
func (s *DoctorSchedule) UpdatedAt() time.Time         { return s.updatedAt }
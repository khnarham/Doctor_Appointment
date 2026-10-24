package domain

import (
	domain "my-appointment-app/Authentication_System/Internals/Domain"
	"time"
)

type DoctorAnalyticsSummary struct {
	TotalAppointments     int64   `json:"total_appointments"`
	CompletedAppointments int64   `json:"completed_appointments"`
	PendingAppointments   int64   `json:"pending_appointments"`
	CancelledAppointments int64   `json:"cancelled_appointments"`
	TotalEarnings         float64 `json:"total_earnings"`
	UniquePatientsCount   int64   `json:"unique_patients_count"`
}

type DoctorPatientDetail struct {
	AppointmentID   string            `json:"appointment_id" db:"appointment_id"`
	PatientID       string            `json:"patient_id" db:"patient_id"`
	PatientName     string            `json:"patient_name" db:"patient_name"`
	PatientPhone    string            `json:"patient_phone" db:"patient_phone"`
	IsWalkIn        bool              `json:"is_walk_in" db:"is_walk_in"`
	Gender          domain.Gender            `json:"gender" db:"gender"`
	DateOfBirth     time.Time         `json:"date_of_birth" db:"date_of_birth"`
	AppointmentDate time.Time         `json:"appointment_date" db:"appointment_date"`
	StartTime       int               `json:"start_time" db:"start_time"`
	EndTime         int               `json:"end_time" db:"end_time"`
	Status          AppointmentStatus `json:"status" db:"status"`
}

type DoctorDashboardResponse struct {
	Summary        DoctorAnalyticsSummary `json:"summary"`
	RecentPatients []DoctorPatientDetail  `json:"recent_patients"`
}
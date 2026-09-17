package domain

import "time"

// FinancialOverviewResponse Admin Dashboard ke top metrics ke liye
type FinancialOverviewResponse struct {
	TotalCompletedAppointments int64                 `json:"total_completed_appointments"`
	TotalRevenue               float64               `json:"total_revenue"`
	TopEarningDoctors          []TopEarningDoctorDTO `json:"top_earning_doctors"`
}

// TopEarningDoctorDTO Top Earning Doctor ki summary
type TopEarningDoctorDTO struct {
	DoctorID              string  `json:"doctor_id"`
	UserID                string  `json:"user_id"`
	LicenseNumber         string  `json:"license_number"`
	ConsultationFee       float64 `json:"consultation_fee"`
	CompletedAppointments int64   `json:"completed_appointments"`
	TotalEarned           float64 `json:"total_earned"`
}

// FinancialFilterDTO Date range filtering ke liye
type FinancialFilterDTO struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Limit     int        `json:"limit"`
}
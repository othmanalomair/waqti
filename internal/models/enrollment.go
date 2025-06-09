package models

import (
	"time"

	"github.com/google/uuid"
)

type Enrollment struct {
	ID             uuid.UUID  `json:"id"`
	WorkshopID     uuid.UUID  `json:"workshop_id"`
	SessionID      *uuid.UUID `json:"session_id"`      // Links to specific session
	OrderID        *uuid.UUID `json:"order_id"`        // Links to order
	WorkshopName   string     `json:"workshop_name"`
	WorkshopNameAr string     `json:"workshop_name_ar"`
	StudentName    string     `json:"student_name"`
	StudentEmail   string     `json:"student_email"`
	StudentPhone   string     `json:"student_phone"`
	TotalPrice     float64    `json:"total_price"`
	Status         string     `json:"status"` // "successful", "pending", "rejected"
	StatusAr       string     `json:"status_ar"`
	EnrollmentDate time.Time  `json:"enrollment_date"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type EnrollmentStats struct {
	SuccessfulSales int     `json:"successful_sales"`
	TotalSales      float64 `json:"total_sales"`
	RejectedSales   int     `json:"rejected_sales"`
	PendingSales    int     `json:"pending_sales"`
}

type EnrollmentFilter struct {
	TimeRange string `json:"time_range"` // "days", "months", "year"
	OrderBy   string `json:"order_by"`   // "date", "price", "name"
	OrderDir  string `json:"order_dir"`  // "asc", "desc"
}

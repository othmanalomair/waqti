package models

import "time"

type Creator struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	NameAr    string    `json:"name_ar"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Plan      string    `json:"plan"`
	PlanAr    string    `json:"plan_ar"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Workshop struct {
	ID            int       `json:"id"`
	CreatorID     int       `json:"creator_id"`
	Title         string    `json:"title"`
	TitleAr       string    `json:"title_ar"`
	Description   string    `json:"description"`
	DescriptionAr string    `json:"description_ar"`
	Price         float64   `json:"price"`
	Duration      int       `json:"duration"` // in minutes
	MaxStudents   int       `json:"max_students"`
	Category      string    `json:"category"`
	CategoryAr    string    `json:"category_ar"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DashboardStats struct {
	TotalWorkshops   int     `json:"total_workshops"`
	ActiveWorkshops  int     `json:"active_workshops"`
	TotalEnrollments int     `json:"total_enrollments"`
	MonthlyRevenue   float64 `json:"monthly_revenue"`
	ProjectedSales   float64 `json:"projected_sales"`
	RemainingSeats   int     `json:"remaining_seats"`
}

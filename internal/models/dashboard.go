package models

type DashboardStats struct {
	TotalWorkshops   int     `json:"total_workshops"`
	ActiveWorkshops  int     `json:"active_workshops"`
	TotalEnrollments int     `json:"total_enrollments"`
	MonthlyRevenue   float64 `json:"monthly_revenue"`
	ProjectedSales   float64 `json:"projected_sales"`
	RemainingSeats   int     `json:"remaining_seats"`
}

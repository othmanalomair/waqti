package models

import (
	"time"
	"waqti/internal/database"
)

// AdminDashboardData contains data for the admin dashboard
type AdminDashboardData struct {
	Admin *database.Creator    `json:"admin"`
	Stats *AdminDashboardStats `json:"stats"`
}

// AdminDashboardStats contains dashboard statistics
type AdminDashboardStats struct {
	TotalUsers        int                 `json:"total_users"`
	TotalCreators     int                 `json:"total_creators"`
	TotalAdmins       int                 `json:"total_admins"`
	TotalWorkshops    int                 `json:"total_workshops"`
	ActiveWorkshops   int                 `json:"active_workshops"`
	TotalEnrollments  int                 `json:"total_enrollments"`
	RecentUsers       []database.Creator  `json:"recent_users"`
	TrafficToday      int                 `json:"traffic_today"`
	ConversionRate    float64             `json:"conversion_rate"`
	PopularPages      map[string]int      `json:"popular_pages"`
}

// AdminUserManagementData contains data for user management page
type AdminUserManagementData struct {
	Admin *database.Creator    `json:"admin"`
	Users []database.Creator   `json:"users"`
	Total int                  `json:"total"`
}

// AdminAnalyticsData contains data for analytics page
type AdminAnalyticsData struct {
	Admin     *database.Creator `json:"admin"`
	Analytics *SystemAnalytics  `json:"analytics"`
}

// SystemAnalytics contains system-wide analytics data
type SystemAnalytics struct {
	TotalPageViews    int                     `json:"total_page_views"`
	LandingPageViews  int                     `json:"landing_page_views"`
	SignInPageViews   int                     `json:"sign_in_page_views"`
	SignUpPageViews   int                     `json:"sign_up_page_views"`
	StorePageViews    int                     `json:"store_page_views"`
	RecentActivity    []AdminAnalyticsEntry  `json:"recent_activity"`
	StartDate         string                  `json:"start_date"`
	EndDate           string                  `json:"end_date"`
	SelectedStore     string                  `json:"selected_store"`
	SelectedPageType  string                  `json:"selected_page_type"`
	AvailableStores   []StoreInfo             `json:"available_stores"`
}

// AdminAnalyticsEntry represents a single analytics entry
type AdminAnalyticsEntry struct {
	ID        int       `json:"id"`
	PageType  string    `json:"page_type"`
	Country   string    `json:"country"`
	Device    string    `json:"device"`
	Browser   string    `json:"browser"`
	StoreName string    `json:"store_name"`
	StoreURL  string    `json:"store_url"`
	CreatedAt time.Time `json:"created_at"`
}

// StoreInfo represents store information for filtering
type StoreInfo struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	NameAr   string `json:"name_ar"`
}

// AdminCreateUserData contains data for create user page
type AdminCreateUserData struct {
	Admin *database.Creator `json:"admin"`
}
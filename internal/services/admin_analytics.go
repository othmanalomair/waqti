package services

import (
	"fmt"
	"time"
	"waqti/internal/database"
	"waqti/internal/models"

	"github.com/google/uuid"
)

// AdminAnalyticsService handles admin-level analytics
type AdminAnalyticsService struct {
	db *database.DB
}

// NewAdminAnalyticsService creates a new admin analytics service
func NewAdminAnalyticsService() *AdminAnalyticsService {
	return &AdminAnalyticsService{
		db: database.Instance,
	}
}

// AdminAnalytic represents an admin analytics record
type AdminAnalytic struct {
	ID        uuid.UUID `json:"id"`
	PageType  string    `json:"page_type"`
	IPAddress *string   `json:"ip_address"`
	UserAgent *string   `json:"user_agent"`
	Referrer  *string   `json:"referrer"`
	Country   *string   `json:"country"`
	Device    *string   `json:"device"`
	Browser   *string   `json:"browser"`
	Timestamp time.Time `json:"timestamp"`
}

// AnalyticsSummary represents daily analytics summary
type AnalyticsSummary struct {
	PageType       string    `json:"page_type"`
	Date           time.Time `json:"date"`
	TotalVisits    int       `json:"total_visits"`
	UniqueVisitors int       `json:"unique_visitors"`
	MobileVisits   int       `json:"mobile_visits"`
	DesktopVisits  int       `json:"desktop_visits"`
}

// DashboardStats represents admin dashboard statistics
type DashboardStats struct {
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

// TrackPageVisit tracks a page visit for admin analytics
func (as *AdminAnalyticsService) TrackPageVisit(pageType, ipAddress, userAgent, referrer, country, device, browser, storeUsername string) error {
	query := `
		INSERT INTO admin_analytics (page_type, ip_address, user_agent, referrer, country, device, browser, store_username)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	// Only set store_username for store_visit page type
	var storeUsernameParam interface{}
	if pageType == "store_visit" && storeUsername != "" {
		storeUsernameParam = storeUsername
	} else {
		storeUsernameParam = nil
	}

	_, err := as.db.Exec(query, pageType, ipAddress, userAgent, referrer, country, device, browser, storeUsernameParam)
	if err != nil {
		return fmt.Errorf("failed to track page visit: %w", err)
	}

	return nil
}

// GetDashboardStats retrieves comprehensive admin dashboard statistics
func (as *AdminAnalyticsService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// Get user counts (all creators are the same role now)
	userCountsQuery := `
		SELECT COUNT(*) as total
		FROM creators 
		WHERE is_active = true
	`

	err := as.db.QueryRow(userCountsQuery).Scan(&stats.TotalUsers)
	if err == nil {
		stats.TotalCreators = stats.TotalUsers // All users are creators
		stats.TotalAdmins = 0 // Admins are in separate table
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user counts: %w", err)
	}

	// Get workshop counts
	workshopCountsQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN is_active = true THEN 1 END) as active
		FROM workshops
	`

	err = as.db.QueryRow(workshopCountsQuery).Scan(&stats.TotalWorkshops, &stats.ActiveWorkshops)
	if err != nil {
		return nil, fmt.Errorf("failed to get workshop counts: %w", err)
	}

	// Get enrollment count
	enrollmentQuery := `SELECT COUNT(*) FROM enrollments WHERE status = 'successful'`
	err = as.db.QueryRow(enrollmentQuery).Scan(&stats.TotalEnrollments)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment count: %w", err)
	}

	// Get recent users (last 10)
	recentUsersQuery := `
		SELECT id, name, name_ar, username, email, plan, is_active, created_at
		FROM creators
		ORDER BY created_at DESC
		LIMIT 10
	`

	rows, err := as.db.Query(recentUsersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent users: %w", err)
	}
	defer rows.Close()

	var recentUsers []database.Creator
	for rows.Next() {
		var user database.Creator
		err := rows.Scan(
			&user.ID, &user.Name, &user.NameAr, &user.Username,
			&user.Email, &user.Plan, &user.IsActive, &user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recent user: %w", err)
		}
		recentUsers = append(recentUsers, user)
	}
	stats.RecentUsers = recentUsers

	// Get traffic today from admin analytics
	trafficTodayQuery := `
		SELECT COUNT(*) 
		FROM admin_analytics 
		WHERE DATE(timestamp) = CURRENT_DATE
	`
	err = as.db.QueryRow(trafficTodayQuery).Scan(&stats.TrafficToday)
	if err != nil {
		// If table doesn't exist yet, set to 0
		stats.TrafficToday = 0
	}

	// Get popular pages from admin analytics
	popularPagesQuery := `
		SELECT page_type, COUNT(*) as visits
		FROM admin_analytics 
		WHERE timestamp >= CURRENT_DATE - INTERVAL '7 days'
		GROUP BY page_type
		ORDER BY visits DESC
		LIMIT 5
	`

	rows, err = as.db.Query(popularPagesQuery)
	if err != nil {
		// If table doesn't exist yet, initialize empty map
		stats.PopularPages = make(map[string]int)
	} else {
		defer rows.Close()
		stats.PopularPages = make(map[string]int)
		for rows.Next() {
			var pageType string
			var visits int
			err := rows.Scan(&pageType, &visits)
			if err != nil {
				continue
			}
			stats.PopularPages[pageType] = visits
		}
	}

	// Calculate conversion rate (signups vs landing page visits)
	conversionQuery := `
		SELECT 
			COALESCE((
				SELECT COUNT(*) FROM admin_analytics 
				WHERE page_type = 'signup' AND timestamp >= CURRENT_DATE - INTERVAL '7 days'
			), 0) as signups,
			COALESCE((
				SELECT COUNT(*) FROM admin_analytics 
				WHERE page_type = 'landing' AND timestamp >= CURRENT_DATE - INTERVAL '7 days'
			), 1) as landing_visits
	`

	var signups, landingVisits int
	err = as.db.QueryRow(conversionQuery).Scan(&signups, &landingVisits)
	if err != nil {
		stats.ConversionRate = 0.0
	} else {
		if landingVisits > 0 {
			stats.ConversionRate = (float64(signups) / float64(landingVisits)) * 100
		} else {
			stats.ConversionRate = 0.0
		}
	}

	return stats, nil
}

// GetAnalyticsSummary retrieves analytics summary for a date range
func (as *AdminAnalyticsService) GetAnalyticsSummary(startDate, endDate time.Time) ([]AnalyticsSummary, error) {
	query := `
		SELECT 
			page_type,
			DATE(timestamp) as date,
			COUNT(*) as total_visits,
			COUNT(DISTINCT ip_address) as unique_visitors,
			COUNT(CASE WHEN device = 'Mobile' THEN 1 END) as mobile_visits,
			COUNT(CASE WHEN device = 'Desktop' THEN 1 END) as desktop_visits
		FROM admin_analytics 
		WHERE DATE(timestamp) BETWEEN $1 AND $2
		GROUP BY page_type, DATE(timestamp)
		ORDER BY date DESC, page_type
	`

	rows, err := as.db.Query(query, startDate, endDate)
	if err != nil {
		// If table doesn't exist, return empty slice
		return []AnalyticsSummary{}, nil
	}
	defer rows.Close()

	var summaries []AnalyticsSummary
	for rows.Next() {
		var summary AnalyticsSummary
		err := rows.Scan(
			&summary.PageType, &summary.Date, &summary.TotalVisits,
			&summary.UniqueVisitors, &summary.MobileVisits, &summary.DesktopVisits,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan analytics summary: %w", err)
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// GetPageAnalytics retrieves analytics for a specific page type
func (as *AdminAnalyticsService) GetPageAnalytics(pageType string, days int) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Get total visits for the period
	totalVisitsQuery := `
		SELECT COUNT(*) 
		FROM admin_analytics 
		WHERE page_type = $1 AND timestamp >= $2
	`

	var totalVisits int
	startDate := time.Now().AddDate(0, 0, -days)
	err := as.db.QueryRow(totalVisitsQuery, pageType, startDate).Scan(&totalVisits)
	if err != nil {
		totalVisits = 0
	}
	result["total_visits"] = totalVisits

	// Get unique visitors
	uniqueVisitorsQuery := `
		SELECT COUNT(DISTINCT ip_address) 
		FROM admin_analytics 
		WHERE page_type = $1 AND timestamp >= $2
	`

	var uniqueVisitors int
	err = as.db.QueryRow(uniqueVisitorsQuery, pageType, startDate).Scan(&uniqueVisitors)
	if err != nil {
		uniqueVisitors = 0
	}
	result["unique_visitors"] = uniqueVisitors

	// Get device breakdown
	deviceQuery := `
		SELECT 
			device,
			COUNT(*) as count
		FROM admin_analytics 
		WHERE page_type = $1 AND timestamp >= $2
		GROUP BY device
		ORDER BY count DESC
	`

	rows, err := as.db.Query(deviceQuery, pageType, startDate)
	if err != nil {
		result["devices"] = make(map[string]int)
	} else {
		defer rows.Close()
		devices := make(map[string]int)
		for rows.Next() {
			var device string
			var count int
			err := rows.Scan(&device, &count)
			if err != nil {
				continue
			}
			if device != "" {
				devices[device] = count
			}
		}
		result["devices"] = devices
	}

	// Get daily visits
	dailyQuery := `
		SELECT 
			DATE(timestamp) as date,
			COUNT(*) as visits
		FROM admin_analytics 
		WHERE page_type = $1 AND timestamp >= $2
		GROUP BY DATE(timestamp)
		ORDER BY date DESC
	`

	rows, err = as.db.Query(dailyQuery, pageType, startDate)
	if err != nil {
		result["daily_visits"] = []map[string]interface{}{}
	} else {
		defer rows.Close()
		var dailyVisits []map[string]interface{}
		for rows.Next() {
			var date time.Time
			var visits int
			err := rows.Scan(&date, &visits)
			if err != nil {
				continue
			}
			dailyVisits = append(dailyVisits, map[string]interface{}{
				"date":   date.Format("2006-01-02"),
				"visits": visits,
			})
		}
		result["daily_visits"] = dailyVisits
	}

	return result, nil
}

// CreateAdminAnalyticsTable creates the admin analytics table if it doesn't exist
func (as *AdminAnalyticsService) CreateAdminAnalyticsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS admin_analytics (
			id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
			page_type VARCHAR(50) NOT NULL,
			ip_address INET,
			user_agent TEXT,
			referrer TEXT,
			country VARCHAR(100),
			device VARCHAR(50),
			browser VARCHAR(50),
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			
			CONSTRAINT admin_analytics_page_type_check 
				CHECK (page_type IN ('landing', 'signin', 'signup', 'store_visit'))
		);

		CREATE INDEX IF NOT EXISTS idx_admin_analytics_page_type ON admin_analytics(page_type);
		CREATE INDEX IF NOT EXISTS idx_admin_analytics_timestamp ON admin_analytics(timestamp);
		CREATE INDEX IF NOT EXISTS idx_admin_analytics_page_date ON admin_analytics(page_type, date(timestamp));
	`

	_, err := as.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create admin analytics table: %w", err)
	}

	return nil
}

// GetTrafficTrends retrieves traffic trends for the last N days
func (as *AdminAnalyticsService) GetTrafficTrends(days int) (map[string][]map[string]interface{}, error) {
	query := `
		SELECT 
			page_type,
			DATE(timestamp) as date,
			COUNT(*) as visits,
			COUNT(DISTINCT ip_address) as unique_visitors
		FROM admin_analytics 
		WHERE timestamp >= $1
		GROUP BY page_type, DATE(timestamp)
		ORDER BY date DESC
	`

	startDate := time.Now().AddDate(0, 0, -days)
	rows, err := as.db.Query(query, startDate)
	if err != nil {
		return make(map[string][]map[string]interface{}), nil
	}
	defer rows.Close()

	trends := make(map[string][]map[string]interface{})
	for rows.Next() {
		var pageType string
		var date time.Time
		var visits, uniqueVisitors int

		err := rows.Scan(&pageType, &date, &visits, &uniqueVisitors)
		if err != nil {
			continue
		}

		if trends[pageType] == nil {
			trends[pageType] = []map[string]interface{}{}
		}

		trends[pageType] = append(trends[pageType], map[string]interface{}{
			"date":             date.Format("2006-01-02"),
			"visits":           visits,
			"unique_visitors":  uniqueVisitors,
		})
	}

	return trends, nil
}

// GetRecentActivity retrieves recent analytics entries for display with optional filtering
func (as *AdminAnalyticsService) GetRecentActivity(limit int, storeFilter string, startDate, endDate time.Time) ([]models.AdminAnalyticsEntry, error) {
	query := `
		SELECT 
			aa.page_type, 
			COALESCE(aa.country, 'N/A') as country, 
			aa.device, 
			aa.browser, 
			aa.timestamp,
			CASE 
				WHEN aa.page_type = 'store_visit' AND aa.store_username IS NOT NULL 
				THEN COALESCE(c.name, aa.store_username)
				ELSE ''
			END as store_name,
			COALESCE(aa.store_username, '') as store_url
		FROM admin_analytics aa
		LEFT JOIN creators c ON aa.store_username = c.username
		WHERE aa.timestamp >= $1 AND aa.timestamp <= $2
	`
	
	args := []interface{}{startDate, endDate}
	
	// Add store filter if specified
	if storeFilter != "" {
		query += ` AND (aa.page_type != 'store_visit' OR aa.store_username = $3)`
		args = append(args, storeFilter)
	}
	
	query += ` ORDER BY aa.timestamp DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit)

	rows, err := as.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}
	defer rows.Close()

	var activities []models.AdminAnalyticsEntry
	for rows.Next() {
		var activity models.AdminAnalyticsEntry
		err := rows.Scan(
			&activity.PageType, &activity.Country, &activity.Device, 
			&activity.Browser, &activity.CreatedAt, &activity.StoreName, &activity.StoreURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recent activity: %w", err)
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// GetRecentActivityWithPageFilter retrieves recent analytics entries with page type filtering
func (as *AdminAnalyticsService) GetRecentActivityWithPageFilter(limit int, storeFilter, pageTypeFilter string, startDate, endDate time.Time) ([]models.AdminAnalyticsEntry, error) {
	var activities []models.AdminAnalyticsEntry
	
	// If filtering for store visits or no page type filter, get store visit data from analytics_clicks
	if pageTypeFilter == "" || pageTypeFilter == "store_visit" {
		storeQuery := `
			SELECT 
				'store_visit' as page_type,
				COALESCE(ac.country, 'N/A') as country,
				ac.device,
				ac.browser,
				ac.clicked_at as timestamp,
				c.name as store_name,
				c.username as store_url
			FROM analytics_clicks ac
			JOIN creators c ON ac.creator_id = c.id
			WHERE ac.clicked_at >= $1 AND ac.clicked_at <= $2
		`
		
		storeArgs := []interface{}{startDate, endDate}
		storeArgCount := 2
		
		// Add store filter if specified
		if storeFilter != "" {
			storeArgCount++
			storeQuery += fmt.Sprintf(` AND c.username = $%d`, storeArgCount)
			storeArgs = append(storeArgs, storeFilter)
		}
		
		storeArgCount++
		storeQuery += fmt.Sprintf(` ORDER BY ac.clicked_at DESC LIMIT $%d`, storeArgCount)
		storeArgs = append(storeArgs, limit)

		rows, err := as.db.Query(storeQuery, storeArgs...)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var activity models.AdminAnalyticsEntry
				err := rows.Scan(
					&activity.PageType, &activity.Country, &activity.Device, 
					&activity.Browser, &activity.CreatedAt, &activity.StoreName, &activity.StoreURL,
				)
				if err == nil {
					activities = append(activities, activity)
				}
			}
		}
	}
	
	// If not filtering for store visits only, get other page types from admin_analytics
	if pageTypeFilter == "" || pageTypeFilter != "store_visit" {
		adminQuery := `
			SELECT 
				aa.page_type, 
				COALESCE(aa.country, 'N/A') as country, 
				aa.device, 
				aa.browser, 
				aa.timestamp,
				'' as store_name,
				'' as store_url
			FROM admin_analytics aa
			WHERE aa.timestamp >= $1 AND aa.timestamp <= $2
		`
		
		adminArgs := []interface{}{startDate, endDate}
		adminArgCount := 2
		
		// Add page type filter if specified (exclude store_visit)
		if pageTypeFilter != "" && pageTypeFilter != "store_visit" {
			adminArgCount++
			adminQuery += fmt.Sprintf(` AND aa.page_type = $%d`, adminArgCount)
			adminArgs = append(adminArgs, pageTypeFilter)
		} else if pageTypeFilter == "" {
			// Exclude store_visit from admin_analytics since we get it from analytics_clicks
			adminQuery += ` AND aa.page_type != 'store_visit'`
		}
		
		adminArgCount++
		adminQuery += fmt.Sprintf(` ORDER BY aa.timestamp DESC LIMIT $%d`, adminArgCount)
		adminArgs = append(adminArgs, limit)

		rows, err := as.db.Query(adminQuery, adminArgs...)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var activity models.AdminAnalyticsEntry
				err := rows.Scan(
					&activity.PageType, &activity.Country, &activity.Device, 
					&activity.Browser, &activity.CreatedAt, &activity.StoreName, &activity.StoreURL,
				)
				if err == nil {
					activities = append(activities, activity)
				}
			}
		}
	}

	return activities, nil
}

// GetPageAnalyticsWithFilters retrieves analytics for a specific page type with date and store filtering
func (as *AdminAnalyticsService) GetPageAnalyticsWithFilters(pageType string, storeFilter string, startDate, endDate time.Time) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var totalVisits int
	
	// Use analytics_clicks table for store visits, admin_analytics for other page types
	if pageType == "store_visit" {
		query := `
			SELECT COUNT(*) 
			FROM analytics_clicks ac
			JOIN creators c ON ac.creator_id = c.id
			WHERE ac.clicked_at >= $1 AND ac.clicked_at <= $2
		`
		
		args := []interface{}{startDate, endDate}
		
		// Add store filter if specified
		if storeFilter != "" {
			query += ` AND c.username = $3`
			args = append(args, storeFilter)
		}

		err := as.db.QueryRow(query, args...).Scan(&totalVisits)
		if err != nil {
			totalVisits = 0
		}
	} else {
		// Use admin_analytics for other page types
		query := `
			SELECT COUNT(*) 
			FROM admin_analytics aa
			WHERE aa.page_type = $1 AND aa.timestamp >= $2 AND aa.timestamp <= $3
		`
		
		args := []interface{}{pageType, startDate, endDate}

		err := as.db.QueryRow(query, args...).Scan(&totalVisits)
		if err != nil {
			totalVisits = 0
		}
	}
	
	result["total_visits"] = totalVisits

	return result, nil
}

// GetAvailableStores retrieves all creators who have stores for filtering
func (as *AdminAnalyticsService) GetAvailableStores() ([]models.StoreInfo, error) {
	query := `
		SELECT DISTINCT c.username, c.name, c.name_ar
		FROM creators c
		INNER JOIN analytics_clicks ac ON ac.creator_id = c.id
		ORDER BY c.name
	`

	rows, err := as.db.Query(query)
	if err != nil {
		return []models.StoreInfo{}, nil // Return empty slice if error
	}
	defer rows.Close()

	var stores []models.StoreInfo
	for rows.Next() {
		var store models.StoreInfo
		err := rows.Scan(&store.Username, &store.Name, &store.NameAr)
		if err != nil {
			continue
		}
		stores = append(stores, store)
	}

	return stores, nil
}

// Helper function to check if admin analytics table exists
func (as *AdminAnalyticsService) TableExists() bool {
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'admin_analytics'
		);
	`

	var exists bool
	err := as.db.QueryRow(query).Scan(&exists)
	return err == nil && exists
}
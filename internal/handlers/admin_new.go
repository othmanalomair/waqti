package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"waqti/internal/database"
	"waqti/internal/middleware"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AdminNewHandler handles admin-related requests with separate authentication
type AdminNewHandler struct {
	adminAuthService *services.AdminAuthService
	analyticsService *services.AdminAnalyticsService
	creatorService   *middleware.AuthService // For managing creators
}

// NewAdminNewHandler creates a new admin handler with separate auth
func NewAdminNewHandler(adminAuthService *services.AdminAuthService, analyticsService *services.AdminAnalyticsService, creatorService *middleware.AuthService) *AdminNewHandler {
	return &AdminNewHandler{
		adminAuthService: adminAuthService,
		analyticsService: analyticsService,
		creatorService:   creatorService,
	}
}

// ShowAdminLogin displays the admin login page
func (ah *AdminNewHandler) ShowAdminLogin(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	errorMessage := c.QueryParam("error")

	return templates.AdminLoginPage(errorMessage, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

// HandleAdminLogin processes admin login form
func (ah *AdminNewHandler) HandleAdminLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Redirect(http.StatusFound, "/admin?error=Please+provide+username+and+password")
	}

	// Authenticate admin
	adminUser, sessionToken, err := ah.adminAuthService.Login(username, password)
	if err != nil {
		log.Printf("Admin login failed for %s: %v", username, err)
		return c.Redirect(http.StatusFound, "/admin?error=Invalid+credentials")
	}

	// Set session cookie
	middleware.SetAdminSessionCookie(c, sessionToken)

	log.Printf("Admin login successful: %s (%s)", adminUser.Email, adminUser.Role)
	return c.Redirect(http.StatusFound, "/admin/dashboard")
}

// HandleAdminLogout processes admin logout
func (ah *AdminNewHandler) HandleAdminLogout(c echo.Context) error {
	// Get session token from cookie
	if cookie, err := c.Cookie("admin_session"); err == nil {
		// Remove session from database
		ah.adminAuthService.Logout(cookie.Value)
	}

	// Clear session cookie
	cookie := &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/admin")
}

// GetCurrentAdmin helper to get current admin user from context
func (ah *AdminNewHandler) GetCurrentAdmin(c echo.Context) *services.AdminUser {
	adminUser, ok := c.Get("admin_user").(*services.AdminUser)
	if !ok {
		return nil
	}
	return adminUser
}

// ShowAdminDashboard displays the main admin dashboard
func (ah *AdminNewHandler) ShowAdminDashboard(c echo.Context) error {
	adminUser := ah.GetCurrentAdmin(c)
	if adminUser == nil {
		return c.Redirect(http.StatusFound, "/admin")
	}

	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	// Convert AdminUser to Creator for template compatibility
	admin := &database.Creator{
		ID:       adminUser.ID,
		Name:     adminUser.Name,
		Email:    adminUser.Email,
		Username: adminUser.Username,
		IsActive: adminUser.IsActive,
	}
	if adminUser.NameAr != nil {
		admin.NameAr = *adminUser.NameAr
	}

	// Get dashboard statistics
	stats, err := ah.analyticsService.GetDashboardStats()
	if err != nil {
		log.Printf("Error getting dashboard stats: %v", err)
		stats = &services.DashboardStats{}
	}

	adminStats := &models.AdminDashboardStats{
		TotalUsers:       stats.TotalUsers,
		TotalCreators:    stats.TotalCreators,
		TotalAdmins:      stats.TotalAdmins,
		TotalWorkshops:   stats.TotalWorkshops,
		ActiveWorkshops:  stats.ActiveWorkshops,
		TotalEnrollments: stats.TotalEnrollments,
		RecentUsers:      stats.RecentUsers,
		TrafficToday:     stats.TrafficToday,
		ConversionRate:   stats.ConversionRate,
		PopularPages:     stats.PopularPages,
	}

	data := &models.AdminDashboardData{
		Admin: admin,
		Stats: adminStats,
	}

	return templates.AdminDashboardPage(data, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

// ShowUserManagement displays the user management interface
func (ah *AdminNewHandler) ShowUserManagement(c echo.Context) error {
	adminUser := ah.GetCurrentAdmin(c)
	if adminUser == nil {
		return c.Redirect(http.StatusFound, "/admin")
	}

	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	// Convert AdminUser to Creator for template compatibility
	admin := &database.Creator{
		ID:       adminUser.ID,
		Name:     adminUser.Name,
		Email:    adminUser.Email,
		Username: adminUser.Username,
		IsActive: adminUser.IsActive,
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit := 20
	offset := (page - 1) * limit

	// Get users with pagination
	users, err := ah.creatorService.GetAllUsers(limit+1, offset)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":    "Failed to load users",
			"error_ar": "فشل في تحميل المستخدمين",
		})
	}

	hasMore := len(users) > limit
	if hasMore {
		users = users[:limit]
	}

	stats, _ := ah.creatorService.GetUserStats()
	totalUsers := 0
	for _, count := range stats {
		totalUsers += count
	}

	data := &models.AdminUserManagementData{
		Admin: admin,
		Users: users,
		Total: totalUsers,
	}

	return templates.AdminUserManagementPage(data, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

// ShowAnalyticsDashboard displays the analytics dashboard with filtering
func (ah *AdminNewHandler) ShowAnalyticsDashboard(c echo.Context) error {
	adminUser := ah.GetCurrentAdmin(c)
	if adminUser == nil {
		return c.Redirect(http.StatusFound, "/admin")
	}

	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	// Convert AdminUser to Creator for template compatibility
	admin := &database.Creator{
		ID:       adminUser.ID,
		Name:     adminUser.Name,
		Email:    adminUser.Email,
		Username: adminUser.Username,
		IsActive: adminUser.IsActive,
	}

	// Get filter parameters from query string
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	storeFilter := c.QueryParam("store_filter")
	pageTypeFilter := c.QueryParam("page_type_filter")

	// Parse date filters
	if startParam := c.QueryParam("start_date"); startParam != "" {
		if parsed, err := time.Parse("2006-01-02", startParam); err == nil {
			startDate = parsed
		}
	}

	if endParam := c.QueryParam("end_date"); endParam != "" {
		if parsed, err := time.Parse("2006-01-02", endParam); err == nil {
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	// Get available stores for the filter dropdown
	availableStores, err := ah.analyticsService.GetAvailableStores()
	if err != nil {
		log.Printf("Error getting available stores: %v", err)
		availableStores = []models.StoreInfo{}
	}

	// Get page analytics for each page type with filters
	pageAnalytics := make(map[string]map[string]interface{})
	var pageTypes []string

	if pageTypeFilter != "" {
		pageTypes = []string{pageTypeFilter}
	} else {
		pageTypes = []string{"landing", "signin", "signup", "store_visit"}
	}

	for _, pageType := range pageTypes {
		analytics, err := ah.analyticsService.GetPageAnalyticsWithFilters(pageType, storeFilter, startDate, endDate)
		if err != nil {
			log.Printf("Error getting %s analytics: %v", pageType, err)
			analytics = make(map[string]interface{})
		}
		pageAnalytics[pageType] = analytics
	}

	// Get recent activity with filters
	recentActivity, err := ah.analyticsService.GetRecentActivityWithPageFilter(20, storeFilter, pageTypeFilter, startDate, endDate)
	if err != nil {
		log.Printf("Error getting recent activity: %v", err)
		recentActivity = []models.AdminAnalyticsEntry{}
	}

	analytics := &models.SystemAnalytics{
		TotalPageViews:   0,
		LandingPageViews: 0,
		SignInPageViews:  0,
		SignUpPageViews:  0,
		StorePageViews:   0,
		RecentActivity:   recentActivity,
		StartDate:        startDate.Format("2006-01-02"),
		EndDate:          endDate.Format("2006-01-02"),
		SelectedStore:    storeFilter,
		SelectedPageType: pageTypeFilter,
		AvailableStores:  availableStores,
	}

	// Calculate page views from pageAnalytics
	if landing, ok := pageAnalytics["landing"]; ok {
		if views, exists := landing["total_visits"].(int); exists {
			analytics.LandingPageViews = views
			analytics.TotalPageViews += views
		}
	}
	if signin, ok := pageAnalytics["signin"]; ok {
		if views, exists := signin["total_visits"].(int); exists {
			analytics.SignInPageViews = views
			analytics.TotalPageViews += views
		}
	}
	if store, ok := pageAnalytics["store_visit"]; ok {
		if views, exists := store["total_visits"].(int); exists {
			analytics.StorePageViews = views
			analytics.TotalPageViews += views
		}
	}
	if signup, ok := pageAnalytics["signup"]; ok {
		if views, exists := signup["total_visits"].(int); exists {
			analytics.SignUpPageViews = views
			analytics.TotalPageViews += views
		}
	}

	data := &models.AdminAnalyticsData{
		Admin:     admin,
		Analytics: analytics,
	}

	return templates.AdminAnalyticsPage(data, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

// ShowCreateAdminForm displays the form to create a new admin user
func (ah *AdminNewHandler) ShowCreateAdminForm(c echo.Context) error {
	adminUser := ah.GetCurrentAdmin(c)
	if adminUser == nil {
		return c.Redirect(http.StatusFound, "/admin")
	}

	// Only super admin can create admin users
	if adminUser.Role != "super_admin" {
		return echo.NewHTTPError(http.StatusForbidden, "Super admin access required")
	}

	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	// Convert AdminUser to Creator for template compatibility
	admin := &database.Creator{
		ID:       adminUser.ID,
		Name:     adminUser.Name,
		Email:    adminUser.Email,
		Username: adminUser.Username,
		IsActive: adminUser.IsActive,
	}

	data := &models.AdminCreateUserData{
		Admin: admin,
	}

	return templates.AdminCreateUserPage(data, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

// CreateAdminUser creates a new admin user (super admin only)
func (ah *AdminNewHandler) CreateAdminUser(c echo.Context) error {
	adminUser := ah.GetCurrentAdmin(c)
	if adminUser == nil {
		return c.Redirect(http.StatusFound, "/admin")
	}

	// Only super admin can create admin users
	if adminUser.Role != "super_admin" {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error":    "Super admin access required",
			"error_ar": "مطلوب صلاحية المدير الأعلى",
		})
	}

	// Get form data
	name := c.FormValue("name")
	nameAr := c.FormValue("name_ar")
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	role := c.FormValue("role")

	// Validate required fields
	if name == "" || username == "" || email == "" || password == "" || role == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":    "All fields are required",
			"error_ar": "جميع الحقول مطلوبة",
		})
	}

	// Validate role
	if role != "admin" && role != "super_admin" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":    "Invalid role selected",
			"error_ar": "دور غير صحيح",
		})
	}

	// Create the admin user
	newAdmin, err := ah.adminAuthService.CreateAdminUser(username, email, name, nameAr, password, role)
	if err != nil {
		log.Printf("Error creating admin user: %v", err)

		if strings.Contains(err.Error(), "email already exists") {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error":    "Email already exists",
				"error_ar": "البريد الإلكتروني موجود بالفعل",
			})
		}

		if strings.Contains(err.Error(), "username already exists") {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error":    "Username already exists",
				"error_ar": "اسم المستخدم موجود بالفعل",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":    "Failed to create admin user",
			"error_ar": "فشل في إنشاء مستخدم إداري",
		})
	}

	log.Printf("Admin user created: %s (role: %s) by %s", newAdmin.Email, newAdmin.Role, adminUser.Email)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"message":    "Admin user created successfully",
		"message_ar": "تم إنشاء المستخدم الإداري بنجاح",
		"user": map[string]interface{}{
			"id":       newAdmin.ID,
			"name":     newAdmin.Name,
			"username": newAdmin.Username,
			"email":    newAdmin.Email,
			"role":     newAdmin.Role,
		},
	})
}

// ToggleUserStatus toggles a user's active status
func (ah *AdminNewHandler) ToggleUserStatus(c echo.Context) error {
	adminUser := ah.GetCurrentAdmin(c)
	if adminUser == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Admin authentication required",
		})
	}

	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":    "Invalid user ID",
			"error_ar": "معرف المستخدم غير صحيح",
		})
	}

	// Get the new status
	isActive := c.FormValue("is_active") == "true"

	// Update user status
	query := `UPDATE creators SET is_active = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err = database.Instance.Exec(query, isActive, userID)
	if err != nil {
		log.Printf("Error toggling user status: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":    "Failed to update user status",
			"error_ar": "فشل في تحديث حالة المستخدم",
		})
	}

	status := "activated"
	statusAr := "مفعل"
	if !isActive {
		status = "deactivated"
		statusAr = "معطل"
	}

	log.Printf("User %s %s by admin %s", userIDStr, status, adminUser.Email)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":    true,
		"message":    fmt.Sprintf("User %s successfully", status),
		"message_ar": fmt.Sprintf("تم %s المستخدم بنجاح", statusAr),
	})
}

// GetUserDetails retrieves detailed information about a user
func (ah *AdminNewHandler) GetUserDetails(c echo.Context) error {
	adminUser := ah.GetCurrentAdmin(c)
	if adminUser == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Admin authentication required",
		})
	}

	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":    "Invalid user ID",
			"error_ar": "معرف المستخدم غير صحيح",
		})
	}

	// Get user by ID
	query := `
		SELECT id, name, name_ar, username, email, avatar,
			   plan, plan_ar, is_active, email_verified, created_at, updated_at
		FROM creators
		WHERE id = $1
	`

	var user database.Creator
	err = database.Instance.QueryRow(query, userID).Scan(
		&user.ID, &user.Name, &user.NameAr, &user.Username,
		&user.Email, &user.Avatar, &user.Plan, &user.PlanAr,
		&user.IsActive, &user.EmailVerified,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error":    "User not found",
				"error_ar": "المستخدم غير موجود",
			})
		}
		log.Printf("Error getting user details: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":    "Failed to get user details",
			"error_ar": "فشل في الحصول على تفاصيل المستخدم",
		})
	}

	// Get additional statistics for this user
	workshopCountQuery := `SELECT COUNT(*) FROM workshops WHERE creator_id = $1`
	var workshopCount int
	database.Instance.QueryRow(workshopCountQuery, userID).Scan(&workshopCount)

	enrollmentCountQuery := `
		SELECT COUNT(*) FROM enrollments e 
		JOIN workshops w ON e.workshop_id = w.id 
		WHERE w.creator_id = $1 AND e.status = 'successful'
	`
	var enrollmentCount int
	database.Instance.QueryRow(enrollmentCountQuery, userID).Scan(&enrollmentCount)

	revenueQuery := `
		SELECT COALESCE(SUM(e.total_price), 0) FROM enrollments e 
		JOIN workshops w ON e.workshop_id = w.id 
		WHERE w.creator_id = $1 AND e.status = 'successful'
	`
	var totalRevenue float64
	database.Instance.QueryRow(revenueQuery, userID).Scan(&totalRevenue)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
		"stats": map[string]interface{}{
			"workshop_count":   workshopCount,
			"enrollment_count": enrollmentCount,
			"total_revenue":    totalRevenue,
		},
	})
}
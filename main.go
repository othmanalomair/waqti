package main

import (
	"log"
	"net/http"
	"os"
	"waqti/internal/database"
	"waqti/internal/handlers"
	"waqti/internal/middleware"
	"waqti/internal/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		log.Println("Continuing with system environment variables...")
	}

	// Initialize database connection
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize Echo
	e := echo.New()

	// Basic middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Language middleware (should be applied before auth)
	e.Use(middleware.LanguageMiddleware())

	// Static files
	e.Static("/static", "web/static")

	// Initialize services
	creatorService := services.NewCreatorService()
	workshopService := services.NewWorkshopService()
	enrollmentService := services.NewEnrollmentService()
	analyticsService := services.NewAnalyticsService()
	settingsService := services.NewSettingsService()
	urlService := services.NewURLService()
	orderService := services.NewOrderService()

	// Initialize auth service
	authService := middleware.NewAuthService(database.Instance)

	// Apply conditional authentication middleware
	e.Use(middleware.ConditionalAuthMiddleware(authService))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, workshopService)
	dashboardHandler := handlers.NewDashboardHandler(workshopService, orderService)
	workshopHandler := handlers.NewWorkshopHandler(workshopService)
	enrollmentHandler := handlers.NewEnrollmentHandler(creatorService, workshopService, enrollmentService)
	analyticsHandler := handlers.NewAnalyticsHandler(creatorService, analyticsService)
	settingsHandler := handlers.NewSettingsHandler(creatorService, settingsService)
	qrHandler := handlers.NewQRHandler(creatorService, settingsService)
	urlHandler := handlers.NewURLHandler(creatorService, urlService)
	orderHandler := handlers.NewOrderHandler(creatorService, orderService)

	// Public routes (no authentication required)
	e.GET("/", authHandler.ShowLandingPage)
	e.POST("/toggle-language", authHandler.ToggleLanguage)
	e.GET("/signin", authHandler.ShowSignIn)
	e.POST("/signin", authHandler.ProcessSignIn)
	e.GET("/signup", authHandler.ShowSignUp)
	e.POST("/signup", authHandler.ProcessSignUp)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Store page route (public, must be after other routes to avoid conflicts)
	e.GET("/:username", authHandler.ShowStorePage)

	// API routes (some public, some protected)
	api := e.Group("/api")
	api.POST("/orders", orderHandler.CreateOrder) // Public for WhatsApp integration

	// Protected routes (authentication required)
	protected := e.Group("")
	// Note: Authentication is handled by ConditionalAuthMiddleware above

	// Auth routes
	protected.POST("/signout", dashboardHandler.ProcessSignOut)

	// Dashboard routes
	protected.GET("/dashboard", dashboardHandler.ShowDashboard)
	protected.POST("/dashboard/toggle-language", dashboardHandler.ToggleLanguage)

	// Workshop routes
	protected.GET("/workshops/add", workshopHandler.ShowAddWorkshop)
	protected.POST("/workshops/create", workshopHandler.CreateWorkshop)
	protected.GET("/workshops/edit/:id", workshopHandler.ShowEditWorkshop)
	protected.POST("/workshops/update/:id", workshopHandler.UpdateWorkshop)
	protected.DELETE("/workshops/delete/:id", workshopHandler.DeleteWorkshop)
	protected.DELETE("/workshops/sessions/:session_id", workshopHandler.DeleteWorkshopSession)
	protected.GET("/workshops/reorder", workshopHandler.ShowReorderWorkshops)
	protected.POST("/workshops/reorder", workshopHandler.ReorderWorkshop)
	protected.POST("/workshops/toggle-status", workshopHandler.ToggleWorkshopStatus)

	// Enrollment routes
	protected.GET("/enrollments/tracking", enrollmentHandler.ShowEnrollmentTracking)
	protected.POST("/enrollments/filter", enrollmentHandler.FilterEnrollments)
	protected.POST("/enrollments/delete", enrollmentHandler.DeleteEnrollment)

	// Order routes
	protected.GET("/orders/tracking", orderHandler.ShowOrderTracking)
	protected.POST("/orders/filter", orderHandler.FilterOrders)
	protected.POST("/orders/update-status", orderHandler.UpdateOrderStatus)
	protected.POST("/orders/delete", orderHandler.DeleteOrder)
	protected.POST("/orders/bulk-action", orderHandler.BulkAction)
	protected.GET("/orders/:id", orderHandler.GetOrderDetails)
	protected.GET("/orders/stats", orderHandler.GetOrderStats)
	protected.GET("/orders/export", orderHandler.ExportOrders)
	protected.POST("/orders/mark-viewed", orderHandler.MarkOrderAsViewed)

	// Analytics routes
	protected.GET("/analytics", analyticsHandler.ShowAnalytics)
	protected.POST("/analytics/filter", analyticsHandler.FilterAnalytics)

	// Settings routes
	protected.GET("/settings/shop", settingsHandler.ShowShopSettings)
	protected.POST("/settings/shop", settingsHandler.UpdateShopSettings)
	protected.POST("/settings/upload-logo", settingsHandler.UploadLogo)

	// QR routes
	protected.GET("/qr/modal", qrHandler.ShowQRModal)

	// URL routes
	protected.GET("/url/edit", urlHandler.ShowEditURLModal)
	protected.POST("/url/validate", urlHandler.ValidateUsername)
	protected.POST("/url/update", urlHandler.UpdateURL)

	// Get port from environment or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting server on 0.0.0.0:%s", port)
	log.Println("Database connected successfully")
	log.Fatal(e.Start("0.0.0.0:" + port))
}

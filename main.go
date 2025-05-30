package main

import (
	"log"
	"net/http"
	"waqti/internal/handlers"
	"waqti/internal/middleware"
	"waqti/internal/services"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
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

	// Initialize handlers
	dashboardHandler := handlers.NewDashboardHandler(creatorService, workshopService)
	workshopHandler := handlers.NewWorkshopHandler(creatorService, workshopService)
	enrollmentHandler := handlers.NewEnrollmentHandler(creatorService, workshopService, enrollmentService)
	analyticsHandler := handlers.NewAnalyticsHandler(creatorService, analyticsService)
	settingsHandler := handlers.NewSettingsHandler(creatorService, settingsService)
	qrHandler := handlers.NewQRHandler(creatorService, settingsService)
	urlHandler := handlers.NewURLHandler(creatorService, urlService)
	authHandler := handlers.NewAuthHandler(creatorService)

	// Root route - redirect to landing page
	e.GET("/", authHandler.ShowLandingPage)

	// General language toggle route
	e.POST("/toggle-language", authHandler.ToggleLanguage)

	// Store page route (must be after other routes to avoid conflicts)
	e.GET("/:username", authHandler.ShowStorePage)

	// Auth routes
	e.GET("/signin", authHandler.ShowSignIn)
	e.POST("/signin", authHandler.ProcessSignIn)
	e.GET("/signup", authHandler.ShowSignUp)
	e.POST("/signup", authHandler.ProcessSignUp)
	e.POST("/signout", authHandler.ProcessSignOut)

	// Health check route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Protected routes (in real app, these would have auth middleware)
	// Dashboard routes
	e.GET("/dashboard", dashboardHandler.ShowDashboard)
	e.POST("/dashboard/toggle-language", dashboardHandler.ToggleLanguage)

	// Workshop routes
	e.GET("/workshops/add", workshopHandler.ShowAddWorkshop)
	e.POST("/workshops/create", workshopHandler.CreateWorkshop)
	e.GET("/workshops/reorder", workshopHandler.ShowReorderWorkshops)
	e.POST("/workshops/reorder", workshopHandler.ReorderWorkshop)
	e.POST("/workshops/toggle-status", workshopHandler.ToggleWorkshopStatus)

	// Enrollment routes
	e.GET("/enrollments/tracking", enrollmentHandler.ShowEnrollmentTracking)
	e.POST("/enrollments/filter", enrollmentHandler.FilterEnrollments)
	e.POST("/enrollments/delete", enrollmentHandler.DeleteEnrollment)

	// Analytics routes
	e.GET("/analytics", analyticsHandler.ShowAnalytics)
	e.POST("/analytics/filter", analyticsHandler.FilterAnalytics)

	// Settings routes
	e.GET("/settings/shop", settingsHandler.ShowShopSettings)
	e.POST("/settings/shop", settingsHandler.UpdateShopSettings)
	e.POST("/settings/upload-logo", settingsHandler.UploadLogo)

	// QR routes
	e.GET("/qr/modal", qrHandler.ShowQRModal)

	// URL routes
	e.GET("/url/edit", urlHandler.ShowEditURLModal)
	e.POST("/url/validate", urlHandler.ValidateUsername)
	e.POST("/url/update", urlHandler.UpdateURL)

	log.Println("Starting server on 0.0.0.0:8080")
	log.Fatal(e.Start("0.0.0.0:8080"))
}

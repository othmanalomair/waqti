package handlers

import (
	"net/http"
	"waqti/internal/database"
	"waqti/internal/middleware"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

type EnrollmentHandler struct {
	creatorService    *services.CreatorService
	workshopService   *services.WorkshopService
	enrollmentService *services.EnrollmentService
}

// Helper function to load shop settings
func (h *EnrollmentHandler) loadShopSettings(creatorID uuid.UUID) *models.ShopSettings {
	dbShopSettings, err := database.Instance.GetShopSettingsByCreatorID(creatorID)
	if err != nil || dbShopSettings == nil {
		return nil
	}

	// Helper function to safely dereference string pointers
	getStringValue := func(s *string) string {
		if s != nil {
			return *s
		}
		return ""
	}
	
	return &models.ShopSettings{
		ID:                 dbShopSettings.ID,
		CreatorID:          dbShopSettings.CreatorID,
		LogoURL:            getStringValue(dbShopSettings.LogoURL),
		CreatorName:        getStringValue(dbShopSettings.CreatorName),
		CreatorNameAr:      getStringValue(dbShopSettings.CreatorNameAr),
		SubHeader:          getStringValue(dbShopSettings.SubHeader),
		SubHeaderAr:        getStringValue(dbShopSettings.SubHeaderAr),
		EnrollmentWhatsApp: getStringValue(dbShopSettings.EnrollmentWhatsApp),
		ContactWhatsApp:    getStringValue(dbShopSettings.ContactWhatsApp),
		CheckoutLanguage:   dbShopSettings.CheckoutLanguage,
		GreetingMessage:    getStringValue(dbShopSettings.GreetingMessage),
		GreetingMessageAr:  getStringValue(dbShopSettings.GreetingMessageAr),
		CurrencySymbol:     dbShopSettings.CurrencySymbol,
		CurrencySymbolAr:   dbShopSettings.CurrencySymbolAr,
		CreatedAt:          dbShopSettings.CreatedAt,
		UpdatedAt:          dbShopSettings.UpdatedAt,
	}
}

func NewEnrollmentHandler(creatorService *services.CreatorService, workshopService *services.WorkshopService, enrollmentService *services.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{
		creatorService:    creatorService,
		workshopService:   workshopService,
		enrollmentService: enrollmentService,
	}
}

func (h *EnrollmentHandler) ShowEnrollmentTracking(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	timeRange := c.QueryParam("time_range")
	if timeRange == "" {
		timeRange = "days"
	}

	orderBy := c.QueryParam("order_by")
	if orderBy == "" {
		orderBy = "date"
	}

	orderDir := c.QueryParam("order_dir")
	if orderDir == "" {
		orderDir = "desc"
	}

	filter := models.EnrollmentFilter{
		TimeRange: timeRange,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
	}

	creatorID := dbCreator.ID
	enrollments := h.enrollmentService.GetEnrollmentsByCreatorID(creatorID, filter)
	stats := h.enrollmentService.GetEnrollmentStats(creatorID, timeRange)

	// Get shop settings
	dbShopSettings, err := database.Instance.GetShopSettingsByCreatorID(dbCreator.ID)
	var shopSettings *models.ShopSettings
	if err == nil && dbShopSettings != nil {
		// Helper function to safely dereference string pointers
		getStringValue := func(s *string) string {
			if s != nil {
				return *s
			}
			return ""
		}
		
		shopSettings = &models.ShopSettings{
			ID:                 dbShopSettings.ID,
			CreatorID:          dbShopSettings.CreatorID,
			LogoURL:            getStringValue(dbShopSettings.LogoURL),
			CreatorName:        getStringValue(dbShopSettings.CreatorName),
			CreatorNameAr:      getStringValue(dbShopSettings.CreatorNameAr),
			SubHeader:          getStringValue(dbShopSettings.SubHeader),
			SubHeaderAr:        getStringValue(dbShopSettings.SubHeaderAr),
			EnrollmentWhatsApp: getStringValue(dbShopSettings.EnrollmentWhatsApp),
			ContactWhatsApp:    getStringValue(dbShopSettings.ContactWhatsApp),
			CheckoutLanguage:   dbShopSettings.CheckoutLanguage,
			GreetingMessage:    getStringValue(dbShopSettings.GreetingMessage),
			GreetingMessageAr:  getStringValue(dbShopSettings.GreetingMessageAr),
			CurrencySymbol:     dbShopSettings.CurrencySymbol,
			CurrencySymbolAr:   dbShopSettings.CurrencySymbolAr,
			CreatedAt:          dbShopSettings.CreatedAt,
			UpdatedAt:          dbShopSettings.UpdatedAt,
		}
	}

	component := templates.EnrollmentTrackingPage(enrollments, stats, filter, shopSettings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *EnrollmentHandler) FilterEnrollments(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	timeRange := c.FormValue("time_range")
	orderBy := c.FormValue("order_by")
	orderDir := c.FormValue("order_dir")

	filter := models.EnrollmentFilter{
		TimeRange: timeRange,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
	}

	creatorID := dbCreator.ID
	enrollments := h.enrollmentService.GetEnrollmentsByCreatorID(creatorID, filter)
	stats := h.enrollmentService.GetEnrollmentStats(creatorID, timeRange)

	shopSettings := h.loadShopSettings(dbCreator.ID)
	component := templates.EnrollmentContent(enrollments, stats, filter, shopSettings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *EnrollmentHandler) DeleteEnrollment(c echo.Context) error {
	enrollmentIDStr := c.FormValue("enrollment_id")

	enrollmentID, err := uuid.Parse(enrollmentIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid enrollment ID")
	}

	err = h.enrollmentService.DeleteEnrollment(enrollmentID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error deleting enrollment")
	}

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	timeRange := c.FormValue("current_time_range")
	orderBy := c.FormValue("current_order_by")
	orderDir := c.FormValue("current_order_dir")

	filter := models.EnrollmentFilter{
		TimeRange: timeRange,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
	}

	creatorID := dbCreator.ID
	enrollments := h.enrollmentService.GetEnrollmentsByCreatorID(creatorID, filter)
	stats := h.enrollmentService.GetEnrollmentStats(creatorID, timeRange)

	shopSettings := h.loadShopSettings(dbCreator.ID)
	component := templates.EnrollmentContent(enrollments, stats, filter, shopSettings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

package handlers

import (
	"net/http"
	"waqti/internal/middleware"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type AnalyticsHandler struct {
	creatorService   *services.CreatorService
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(creatorService *services.CreatorService, analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		creatorService:   creatorService,
		analyticsService: analyticsService,
	}
}

func (h *AnalyticsHandler) ShowAnalytics(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current creator from session
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	// Convert database.Creator to models.Creator for template compatibility
	creator := &models.Creator{
		ID:       dbCreator.ID,
		Name:     dbCreator.Name,
		NameAr:   dbCreator.NameAr,
		Username: dbCreator.Username,
		Email:    dbCreator.Email,
		Plan:     dbCreator.Plan,
		PlanAr:   dbCreator.PlanAr,
		IsActive: dbCreator.IsActive,
	}

	dateRange := c.QueryParam("date_range")
	if dateRange == "" {
		dateRange = "30days"
	}

	filterType := c.QueryParam("filter_type")
	if filterType == "" {
		filterType = "all"
	}

	filter := models.AnalyticsFilter{
		FilterType: filterType,
		DateRange:  dateRange,
	}

	clicks := h.analyticsService.GetClicksByCreatorID(dbCreator.ID, filter)
	stats := h.analyticsService.GetAnalyticsStats(dbCreator.ID, filter)
	scriptTag := h.analyticsService.GenerateCompleteScript(clicks, stats)

	component := templates.AnalyticsPage(creator, clicks, stats, filter, lang, isRTL, scriptTag)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AnalyticsHandler) FilterAnalytics(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current creator from session
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	// Convert database.Creator to models.Creator for template compatibility
	creator := &models.Creator{
		ID:       dbCreator.ID,
		Name:     dbCreator.Name,
		NameAr:   dbCreator.NameAr,
		Username: dbCreator.Username,
		Email:    dbCreator.Email,
		Plan:     dbCreator.Plan,
		PlanAr:   dbCreator.PlanAr,
		IsActive: dbCreator.IsActive,
	}

	dateRange := c.FormValue("date_range")
	filterType := c.FormValue("filter_type")

	filter := models.AnalyticsFilter{
		FilterType: filterType,
		DateRange:  dateRange,
	}

	clicks := h.analyticsService.GetClicksByCreatorID(dbCreator.ID, filter)
	stats := h.analyticsService.GetAnalyticsStats(dbCreator.ID, filter)
	scriptTag := h.analyticsService.GenerateCompleteScript(clicks, stats)

	component := templates.AnalyticsContent(creator, clicks, stats, filter, lang, isRTL, scriptTag)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

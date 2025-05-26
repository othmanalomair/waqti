package handlers

import (
	"net/http"
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
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get creator data
	creator, err := h.creatorService.GetCreatorByID(1)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
	}

	// Get query parameters for filtering
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

	// Get analytics data
	clicks := h.analyticsService.GetClicksByCreatorID(1, filter)
	stats := h.analyticsService.GetAnalyticsStats(1, filter)

	// Render template
	component := templates.AnalyticsPage(creator, clicks, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AnalyticsHandler) FilterAnalytics(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get creator data
	creator, err := h.creatorService.GetCreatorByID(1)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
	}

	// Get form parameters
	dateRange := c.FormValue("date_range")
	filterType := c.FormValue("filter_type")

	filter := models.AnalyticsFilter{
		FilterType: filterType,
		DateRange:  dateRange,
	}

	// Get filtered data
	clicks := h.analyticsService.GetClicksByCreatorID(1, filter)
	stats := h.analyticsService.GetAnalyticsStats(1, filter)

	// Return updated content via HTMX
	component := templates.AnalyticsContent(creator, clicks, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

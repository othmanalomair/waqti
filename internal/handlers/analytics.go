package handlers

import (
	"net/http"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"

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

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	creator, err := h.creatorService.GetCreatorByID(creatorID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
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

	clicks := h.analyticsService.GetClicksByCreatorID(creatorID, filter)
	stats := h.analyticsService.GetAnalyticsStats(creatorID, filter)

	component := templates.AnalyticsPage(creator, clicks, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AnalyticsHandler) FilterAnalytics(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	creator, err := h.creatorService.GetCreatorByID(creatorID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
	}

	dateRange := c.FormValue("date_range")
	filterType := c.FormValue("filter_type")

	filter := models.AnalyticsFilter{
		FilterType: filterType,
		DateRange:  dateRange,
	}

	clicks := h.analyticsService.GetClicksByCreatorID(creatorID, filter)
	stats := h.analyticsService.GetAnalyticsStats(creatorID, filter)

	component := templates.AnalyticsContent(creator, clicks, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

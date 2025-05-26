package handlers

import (
	"net/http"
	"strconv"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type EnrollmentHandler struct {
	creatorService    *services.CreatorService
	workshopService   *services.WorkshopService
	enrollmentService *services.EnrollmentService
}

func NewEnrollmentHandler(creatorService *services.CreatorService, workshopService *services.WorkshopService, enrollmentService *services.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{
		creatorService:    creatorService,
		workshopService:   workshopService,
		enrollmentService: enrollmentService,
	}
}

func (h *EnrollmentHandler) ShowEnrollmentTracking(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get query parameters for filtering
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

	// Get enrollments and stats
	enrollments := h.enrollmentService.GetEnrollmentsByCreatorID(1, filter)
	stats := h.enrollmentService.GetEnrollmentStats(1, timeRange)

	// Render template
	component := templates.EnrollmentTrackingPage(enrollments, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *EnrollmentHandler) FilterEnrollments(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get form parameters
	timeRange := c.FormValue("time_range")
	orderBy := c.FormValue("order_by")
	orderDir := c.FormValue("order_dir")

	filter := models.EnrollmentFilter{
		TimeRange: timeRange,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
	}

	// Get filtered data
	enrollments := h.enrollmentService.GetEnrollmentsByCreatorID(1, filter)
	stats := h.enrollmentService.GetEnrollmentStats(1, timeRange)

	// Return updated content via HTMX
	component := templates.EnrollmentContent(enrollments, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *EnrollmentHandler) DeleteEnrollment(c echo.Context) error {
	enrollmentID := c.FormValue("enrollment_id")

	// Convert to int
	id, err := strconv.Atoi(enrollmentID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid enrollment ID")
	}

	// Delete enrollment
	err = h.enrollmentService.DeleteEnrollment(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error deleting enrollment")
	}

	// Return updated list via HTMX
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current filter from form
	timeRange := c.FormValue("current_time_range")
	orderBy := c.FormValue("current_order_by")
	orderDir := c.FormValue("current_order_dir")

	filter := models.EnrollmentFilter{
		TimeRange: timeRange,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
	}

	enrollments := h.enrollmentService.GetEnrollmentsByCreatorID(1, filter)
	stats := h.enrollmentService.GetEnrollmentStats(1, timeRange)

	component := templates.EnrollmentContent(enrollments, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

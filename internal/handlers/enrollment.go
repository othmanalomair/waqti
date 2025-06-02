package handlers

import (
	"net/http"
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

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	enrollments := h.enrollmentService.GetEnrollmentsByCreatorID(creatorID, filter)
	stats := h.enrollmentService.GetEnrollmentStats(creatorID, timeRange)

	component := templates.EnrollmentTrackingPage(enrollments, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *EnrollmentHandler) FilterEnrollments(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	timeRange := c.FormValue("time_range")
	orderBy := c.FormValue("order_by")
	orderDir := c.FormValue("order_dir")

	filter := models.EnrollmentFilter{
		TimeRange: timeRange,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
	}

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	enrollments := h.enrollmentService.GetEnrollmentsByCreatorID(creatorID, filter)
	stats := h.enrollmentService.GetEnrollmentStats(creatorID, timeRange)

	component := templates.EnrollmentContent(enrollments, stats, filter, lang, isRTL)
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

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	enrollments := h.enrollmentService.GetEnrollmentsByCreatorID(creatorID, filter)
	stats := h.enrollmentService.GetEnrollmentStats(creatorID, timeRange)

	component := templates.EnrollmentContent(enrollments, stats, filter, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

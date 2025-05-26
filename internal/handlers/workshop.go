package handlers

import (
	"net/http"
	"strconv"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type WorkshopHandler struct {
	creatorService  *services.CreatorService
	workshopService *services.WorkshopService
}

func NewWorkshopHandler(creatorService *services.CreatorService, workshopService *services.WorkshopService) *WorkshopHandler {
	return &WorkshopHandler{
		creatorService:  creatorService,
		workshopService: workshopService,
	}
}

func (h *WorkshopHandler) ShowReorderWorkshops(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get creator data (hardcoded ID for demo)
	creator, err := h.creatorService.GetCreatorByID(1)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
	}

	// Get workshops and stats
	workshops := h.workshopService.GetWorkshopsByCreatorID(1)
	stats := h.workshopService.GetDashboardStats(1)

	// Render template
	component := templates.ReorderWorkshopsPage(creator, workshops, stats, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) ReorderWorkshop(c echo.Context) error {
	workshopID := c.FormValue("workshop_id")
	direction := c.FormValue("direction") // "up" or "down"

	// Convert to int
	id, err := strconv.Atoi(workshopID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid workshop ID")
	}

	// Reorder logic (dummy implementation)
	err = h.workshopService.ReorderWorkshop(id, direction)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reordering workshop")
	}

	// Return updated workshops list via HTMX
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)
	workshops := h.workshopService.GetWorkshopsByCreatorID(1)

	component := templates.WorkshopsList(workshops, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) ToggleWorkshopStatus(c echo.Context) error {
	workshopID := c.FormValue("workshop_id")

	// Convert to int
	id, err := strconv.Atoi(workshopID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid workshop ID")
	}

	// Toggle status logic
	err = h.workshopService.ToggleWorkshopStatus(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error toggling workshop status")
	}

	// Return updated workshops list via HTMX
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)
	workshops := h.workshopService.GetWorkshopsByCreatorID(1)

	component := templates.WorkshopsList(workshops, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

package handlers

import (
	"net/http"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"

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

func (h *WorkshopHandler) ShowAddWorkshop(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	creator := h.creatorService.GetCreator(creatorID)

	return templates.AddWorkshopPage(creator, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) CreateWorkshop(c echo.Context) error {
	workshop := &models.Workshop{
		ID:          uuid.New(),
		CreatorID:   uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Status:      "draft",
		Currency:    "KWD",
	}

	if price := c.FormValue("price"); price != "" {
		// Parse price logic here
	}

	if currency := c.FormValue("currency"); currency != "" {
		workshop.Currency = currency
	}

	workshop.IsFree = c.FormValue("is_free") == "true" || c.FormValue("is_free") == "on"
	workshop.IsRecurring = c.FormValue("is_recurring") == "true" || c.FormValue("is_recurring") == "on"
	workshop.RecurrenceType = c.FormValue("recurrence_type")

	status := c.FormValue("status")
	if status == "published" {
		workshop.Status = "published"
	} else {
		workshop.Status = "draft"
	}

	return c.Redirect(http.StatusSeeOther, "/workshops/reorder")
}

func (h *WorkshopHandler) ShowReorderWorkshops(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	creator := h.creatorService.GetCreator(creatorID)
	workshops := h.workshopService.GetWorkshops(creatorID)
	stats := h.workshopService.GetDashboardStats(creatorID)

	return templates.ReorderWorkshopsPage(creator, workshops, stats, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) ReorderWorkshop(c echo.Context) error {
	workshopIDStr := c.FormValue("workshop_id")
	direction := c.FormValue("direction")

	workshopID, err := uuid.Parse(workshopIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid workshop ID")
	}

	// Process reorder logic here
	_ = workshopID
	_ = direction

	lang := c.Get("lang").(string)
	isRTL := lang == "ar"
	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	workshops := h.workshopService.GetWorkshops(creatorID)

	return templates.WorkshopsList(workshops, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) ToggleWorkshopStatus(c echo.Context) error {
	workshopIDStr := c.FormValue("workshop_id")

	workshopID, err := uuid.Parse(workshopIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid workshop ID")
	}

	// Process toggle logic here
	_ = workshopID

	lang := c.Get("lang").(string)
	isRTL := lang == "ar"
	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	workshops := h.workshopService.GetWorkshops(creatorID)

	return templates.WorkshopsList(workshops, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

package handlers

import (
	"net/http"
	"waqti/internal/middleware"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

type WorkshopHandler struct {
	workshopService *services.WorkshopService
}

func NewWorkshopHandler(workshopService *services.WorkshopService) *WorkshopHandler {
	return &WorkshopHandler{
		workshopService: workshopService,
	}
}

func (h *WorkshopHandler) ShowAddWorkshop(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	// Convert to models.Creator for template compatibility
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

	return templates.AddWorkshopPage(creator, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) CreateWorkshop(c echo.Context) error {
	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	workshop := &models.Workshop{
		ID:          uuid.New(),
		CreatorID:   dbCreator.ID,
		Name:        c.FormValue("name"),
		Title:       c.FormValue("name"), // Use name as title for now
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

	// TODO: Implement actual database storage for workshops
	// For now, redirect to reorder page
	return c.Redirect(http.StatusSeeOther, "/workshops/reorder")
}

func (h *WorkshopHandler) ShowReorderWorkshops(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	// Convert to models.Creator for template compatibility
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

	workshops := h.workshopService.GetWorkshopsByCreatorID(dbCreator.ID)
	stats := h.workshopService.GetDashboardStats(dbCreator.ID)

	return templates.ReorderWorkshopsPage(creator, workshops, stats, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) ReorderWorkshop(c echo.Context) error {
	// Check authentication
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

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

	workshops := h.workshopService.GetWorkshopsByCreatorID(dbCreator.ID)

	return templates.WorkshopsList(workshops, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) ToggleWorkshopStatus(c echo.Context) error {
	// Check authentication
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	workshopIDStr := c.FormValue("workshop_id")

	workshopID, err := uuid.Parse(workshopIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid workshop ID")
	}

	// Process toggle logic here
	_ = workshopID

	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	workshops := h.workshopService.GetWorkshopsByCreatorID(dbCreator.ID)

	return templates.WorkshopsList(workshops, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

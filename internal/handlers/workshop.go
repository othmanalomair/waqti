package handlers

import (
	"net/http"
	"strconv"
	"waqti/internal/models"
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

// ShowAddWorkshop displays the add workshop form
func (h *WorkshopHandler) ShowAddWorkshop(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	creator := h.creatorService.GetCreator(1) // Using dummy data

	return templates.AddWorkshopPage(creator, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

// CreateWorkshop handles workshop creation
func (h *WorkshopHandler) CreateWorkshop(c echo.Context) error {
	// Parse form data
	workshop := &models.Workshop{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Status:      "draft", // Default status
		Currency:    "KWD",   // Default currency
	}

	// Process price
	if price := c.FormValue("price"); price != "" {
		if p, err := strconv.ParseFloat(price, 64); err == nil {
			workshop.Price = p
		}
	}

	// Set currency (default to KWD if empty)
	if currency := c.FormValue("currency"); currency != "" {
		workshop.Currency = currency
	}

	// Parse boolean fields
	workshop.IsFree = c.FormValue("is_free") == "true" || c.FormValue("is_free") == "on"
	workshop.IsRecurring = c.FormValue("is_recurring") == "true" || c.FormValue("is_recurring") == "on"
	workshop.RecurrenceType = c.FormValue("recurrence_type")

	// Set status based on submit button
	status := c.FormValue("status")
	if status == "published" {
		workshop.Status = "published"
	} else {
		workshop.Status = "draft"
	}

	// In a real app, you'd save to database and handle file uploads
	// For now, just log the workshop data
	// fmt.Printf("Workshop created: %+v\n", workshop)

	// Redirect back to reorder page
	return c.Redirect(http.StatusSeeOther, "/workshops/reorder")
}

// ShowReorderWorkshops displays the workshop reorder page
func (h *WorkshopHandler) ShowReorderWorkshops(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	creator := h.creatorService.GetCreator(1)
	workshops := h.workshopService.GetWorkshops(1)
	stats := h.workshopService.GetDashboardStats(1)

	return templates.ReorderWorkshopsPage(creator, workshops, stats, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

// ReorderWorkshop handles workshop reordering
func (h *WorkshopHandler) ReorderWorkshop(c echo.Context) error {
	// Parse workshop ID and direction
	workshopID, _ := strconv.Atoi(c.FormValue("workshop_id"))
	direction := c.FormValue("direction")

	// In a real app, you'd update the order in the database
	_ = workshopID
	_ = direction

	// Return updated workshops list
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"
	workshops := h.workshopService.GetWorkshops(1)

	return templates.WorkshopsList(workshops, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

// ToggleWorkshopStatus handles workshop status toggling
func (h *WorkshopHandler) ToggleWorkshopStatus(c echo.Context) error {
	// Parse workshop ID
	workshopID, _ := strconv.Atoi(c.FormValue("workshop_id"))

	// In a real app, you'd toggle the status in the database
	_ = workshopID

	// Return updated workshops list
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"
	workshops := h.workshopService.GetWorkshops(1)

	return templates.WorkshopsList(workshops, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

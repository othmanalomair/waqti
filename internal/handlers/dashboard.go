package handlers

import (
	"net/http"
	"time"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	creatorService  *services.CreatorService
	workshopService *services.WorkshopService
	orderService    *services.OrderService // Added order service
}

func NewDashboardHandler(creatorService *services.CreatorService, workshopService *services.WorkshopService, orderService *services.OrderService) *DashboardHandler {
	return &DashboardHandler{
		creatorService:  creatorService,
		workshopService: workshopService,
		orderService:    orderService, // Initialize order service
	}
}

func (h *DashboardHandler) ShowDashboard(c echo.Context) error {
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

	// Get pending orders count for the notification badge
	pendingOrdersCount := h.orderService.GetPendingOrdersCount(1)

	// Render template with pending orders count
	component := templates.DashboardPage(creator, workshops, stats, pendingOrdersCount, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *DashboardHandler) ToggleLanguage(c echo.Context) error {
	currentLang := c.Get("lang").(string)
	newLang := "en"
	if currentLang == "en" {
		newLang = "ar"
	}

	// Set language cookie
	cookie := &http.Cookie{
		Name:     "lang",
		Value:    newLang,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	// Redirect back to dashboard
	return c.Redirect(http.StatusSeeOther, "/dashboard")
}

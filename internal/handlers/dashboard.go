package handlers

import (
	"net/http"
	"time"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	creatorService  *services.CreatorService
	workshopService *services.WorkshopService
	orderService    *services.OrderService
}

func NewDashboardHandler(creatorService *services.CreatorService, workshopService *services.WorkshopService, orderService *services.OrderService) *DashboardHandler {
	return &DashboardHandler{
		creatorService:  creatorService,
		workshopService: workshopService,
		orderService:    orderService,
	}
}

func (h *DashboardHandler) ShowDashboard(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Use the fixed demo creator ID
	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	creator, err := h.creatorService.GetCreatorByID(creatorID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
	}

	workshops := h.workshopService.GetWorkshopsByCreatorID(creatorID)
	stats := h.workshopService.GetDashboardStats(creatorID)
	pendingOrdersCount := h.orderService.GetPendingOrdersCount(creatorID)

	component := templates.DashboardPage(creator, workshops, stats, pendingOrdersCount, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *DashboardHandler) ToggleLanguage(c echo.Context) error {
	currentLang := c.Get("lang").(string)
	newLang := "en"
	if currentLang == "en" {
		newLang = "ar"
	}

	cookie := &http.Cookie{
		Name:     "lang",
		Value:    newLang,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/dashboard")
}

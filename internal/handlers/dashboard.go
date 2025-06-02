package handlers

import (
	"net/http"
	"time"
	"waqti/internal/middleware"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	workshopService *services.WorkshopService
	orderService    *services.OrderService
}

func NewDashboardHandler(workshopService *services.WorkshopService, orderService *services.OrderService) *DashboardHandler {
	return &DashboardHandler{
		workshopService: workshopService,
		orderService:    orderService,
	}
}

func (h *DashboardHandler) ShowDashboard(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		// This should be handled by middleware, but just in case
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
	pendingOrdersCount := h.orderService.GetPendingOrdersCount(dbCreator.ID)

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

func (h *DashboardHandler) ProcessSignOut(c echo.Context) error {
	// Clear the session cookie
	cookie := &http.Cookie{
		Name:     "waqti_session",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   c.Scheme() == "https" || c.Request().Header.Get("X-Forwarded-Proto") == "https",
	}
	c.SetCookie(cookie)

	// Redirect to the landing page with a success message
	return c.Redirect(http.StatusSeeOther, "/?signed_out=1")
}

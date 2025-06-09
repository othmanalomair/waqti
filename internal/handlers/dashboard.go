package handlers

import (
	"fmt"
	"net/http"
	"time"
	"waqti/internal/database"
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

	fmt.Printf("ShowDashboard: Creator ID: %s, Username: %s, Email: %s\n",
		dbCreator.ID, dbCreator.Username, dbCreator.Email)

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

	// Get shop settings
	dbShopSettings, err := database.Instance.GetShopSettingsByCreatorID(dbCreator.ID)
	var shopSettings *models.ShopSettings
	if err == nil && dbShopSettings != nil {
		// Helper function to safely dereference string pointers
		getStringValue := func(s *string) string {
			if s != nil {
				return *s
			}
			return ""
		}
		
		shopSettings = &models.ShopSettings{
			ID:                 dbShopSettings.ID,
			CreatorID:          dbShopSettings.CreatorID,
			LogoURL:            getStringValue(dbShopSettings.LogoURL),
			CreatorName:        getStringValue(dbShopSettings.CreatorName),
			CreatorNameAr:      getStringValue(dbShopSettings.CreatorNameAr),
			SubHeader:          getStringValue(dbShopSettings.SubHeader),
			SubHeaderAr:        getStringValue(dbShopSettings.SubHeaderAr),
			EnrollmentWhatsApp: getStringValue(dbShopSettings.EnrollmentWhatsApp),
			ContactWhatsApp:    getStringValue(dbShopSettings.ContactWhatsApp),
			CheckoutLanguage:   dbShopSettings.CheckoutLanguage,
			GreetingMessage:    getStringValue(dbShopSettings.GreetingMessage),
			GreetingMessageAr:  getStringValue(dbShopSettings.GreetingMessageAr),
			CurrencySymbol:     dbShopSettings.CurrencySymbol,
			CurrencySymbolAr:   dbShopSettings.CurrencySymbolAr,
			CreatedAt:          dbShopSettings.CreatedAt,
			UpdatedAt:          dbShopSettings.UpdatedAt,
		}
	}

	// Get URL settings to show remaining changes
	dbURLSettings, err := database.Instance.GetURLSettingsByCreatorID(dbCreator.ID)
	if err != nil {
		fmt.Printf("ShowDashboard: Error getting URL settings: %v\n", err)
		// If URL settings don't exist, create them
		tx, txErr := database.Instance.Begin()
		if txErr == nil {
			_, execErr := tx.Exec(`INSERT INTO url_settings (creator_id, username, changes_used, max_changes) VALUES ($1, $2, 0, 5)`,
				dbCreator.ID, dbCreator.Username)
			if execErr == nil {
				tx.Commit()
				fmt.Printf("ShowDashboard: Created new URL settings for creator %s\n", dbCreator.ID)
				// Try to get them again
				dbURLSettings, _ = database.Instance.GetURLSettingsByCreatorID(dbCreator.ID)
			} else {
				fmt.Printf("ShowDashboard: Error creating URL settings: %v\n", execErr)
				tx.Rollback()
			}
		}
	}

	// Convert database.URLSettings to models.URLSettings
	var urlSettings *models.URLSettings
	if dbURLSettings != nil {
		urlSettings = &models.URLSettings{
			ID:          dbURLSettings.ID,
			CreatorID:   dbURLSettings.CreatorID,
			Username:    dbURLSettings.Username,
			ChangesUsed: dbURLSettings.ChangesUsed,
			MaxChanges:  dbURLSettings.MaxChanges,
			LastChanged: dbURLSettings.LastChanged,
			CreatedAt:   dbURLSettings.CreatedAt,
			UpdatedAt:   dbURLSettings.UpdatedAt,
		}
		fmt.Printf("ShowDashboard: URL Settings found - Username: %s, Changes: %d/%d\n",
			urlSettings.Username, urlSettings.ChangesUsed, urlSettings.MaxChanges)
	} else {
		// Fallback if still no URL settings - create a default one
		urlSettings = &models.URLSettings{
			CreatorID:   dbCreator.ID,
			Username:    dbCreator.Username,
			ChangesUsed: 0,
			MaxChanges:  5,
		}
		fmt.Printf("ShowDashboard: Using fallback URL settings for creator %s\n", dbCreator.ID)
	}

	component := templates.DashboardPageWithURLSettings(creator, workshops, stats, pendingOrdersCount, urlSettings, shopSettings, lang, isRTL)
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

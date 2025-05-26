package handlers

import (
	"net/http"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type SettingsHandler struct {
	creatorService  *services.CreatorService
	settingsService *services.SettingsService
}

func NewSettingsHandler(creatorService *services.CreatorService, settingsService *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		creatorService:  creatorService,
		settingsService: settingsService,
	}
}

func (h *SettingsHandler) ShowShopSettings(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get creator data
	creator, err := h.creatorService.GetCreatorByID(1)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
	}

	// Get settings
	settings, err := h.settingsService.GetSettingsByCreatorID(1)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading settings")
	}

	// Check for success message
	successMsg := c.QueryParam("success")

	// Render template
	component := templates.ShopSettingsPage(creator, settings, successMsg, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *SettingsHandler) UpdateShopSettings(c echo.Context) error {
	// Parse form data
	var request models.SettingsUpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid form data")
	}

	// Update settings
	err := h.settingsService.UpdateSettings(1, request)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating settings")
	}

	// Redirect with success message
	return c.Redirect(http.StatusSeeOther, "/settings/shop?success=1")
}

func (h *SettingsHandler) UploadLogo(c echo.Context) error {
	// Handle file upload (simplified implementation)
	// In real implementation, you would:
	// 1. Validate file type and size
	// 2. Save to storage (local/cloud)
	// 3. Update database with new URL

	// file, err := c.FormFile("logo")
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, "No file uploaded")
	// }

	// For demo, just return success with dummy URL
	logoURL := "/static/images/uploaded-logo.png"

	err := h.settingsService.UpdateLogo(1, logoURL)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating logo")
	}

	// Return updated logo section via HTMX
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	settings, _ := h.settingsService.GetSettingsByCreatorID(1)
	component := templates.LogoSection(settings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

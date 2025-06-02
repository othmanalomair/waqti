package handlers

import (
	"net/http"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"

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
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	creator, err := h.creatorService.GetCreatorByID(creatorID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
	}

	settings, err := h.settingsService.GetSettingsByCreatorID(creatorID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading settings")
	}

	successMsg := c.QueryParam("success")

	component := templates.ShopSettingsPage(creator, settings, successMsg, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *SettingsHandler) UpdateShopSettings(c echo.Context) error {
	var request models.SettingsUpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid form data")
	}

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	err := h.settingsService.UpdateSettings(creatorID, request)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating settings")
	}

	return c.Redirect(http.StatusSeeOther, "/settings/shop?success=1")
}

func (h *SettingsHandler) UploadLogo(c echo.Context) error {
	logoURL := "/static/images/uploaded-logo.png"

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	err := h.settingsService.UpdateLogo(creatorID, logoURL)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error updating logo")
	}

	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	settings, _ := h.settingsService.GetSettingsByCreatorID(creatorID)
	component := templates.LogoSection(settings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

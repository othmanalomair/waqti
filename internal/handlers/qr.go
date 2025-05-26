package handlers

import (
	"net/http"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type QRHandler struct {
	creatorService  *services.CreatorService
	settingsService *services.SettingsService
}

func NewQRHandler(creatorService *services.CreatorService, settingsService *services.SettingsService) *QRHandler {
	return &QRHandler{
		creatorService:  creatorService,
		settingsService: settingsService,
	}
}

func (h *QRHandler) ShowQRModal(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get creator data
	creator, err := h.creatorService.GetCreatorByID(1)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
	}

	// Get settings for description
	settings, err := h.settingsService.GetSettingsByCreatorID(1)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading settings")
	}

	// Render QR modal component
	component := templates.QRModal(creator, settings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

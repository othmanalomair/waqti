package handlers

import (
	"net/http"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"

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

	component := templates.QRModal(creator, settings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

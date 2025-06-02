package handlers

import (
	"net/http"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	creatorService *services.CreatorService
	urlService     *services.URLService
}

func NewURLHandler(creatorService *services.CreatorService, urlService *services.URLService) *URLHandler {
	return &URLHandler{
		creatorService: creatorService,
		urlService:     urlService,
	}
}

func (h *URLHandler) ShowEditURLModal(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	creator, err := h.creatorService.GetCreatorByID(creatorID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
	}

	urlSettings, err := h.urlService.GetURLSettingsByCreatorID(creatorID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading URL settings")
	}

	component := templates.EditURLModal(creator, urlSettings, "", lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *URLHandler) ValidateUsername(c echo.Context) error {
	var request models.URLUpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.URLValidationResult{
			IsValid:      false,
			ErrorMessage: "Invalid request",
		})
	}

	validation := h.urlService.ValidateUsername(request.Username)
	return c.JSON(http.StatusOK, validation)
}

func (h *URLHandler) UpdateURL(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	var request models.URLUpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid form data")
	}

	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	err := h.urlService.UpdateUsername(creatorID, request.Username)
	if err != nil {
		creator, _ := h.creatorService.GetCreatorByID(creatorID)
		urlSettings, _ := h.urlService.GetURLSettingsByCreatorID(creatorID)

		component := templates.EditURLModal(creator, urlSettings, err.Error(), lang, isRTL)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	creator, _ := h.creatorService.GetCreatorByID(creatorID)
	urlSettings, _ := h.urlService.GetURLSettingsByCreatorID(creatorID)

	component := templates.EditURLModal(creator, urlSettings, "success", lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

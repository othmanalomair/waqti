package handlers

import (
	"net/http"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

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
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get creator data
	creator, err := h.creatorService.GetCreatorByID(1)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading creator")
	}

	// Get URL settings
	urlSettings, err := h.urlService.GetURLSettingsByCreatorID(1)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading URL settings")
	}

	// Render modal component
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
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	var request models.URLUpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid form data")
	}

	// Update URL
	err := h.urlService.UpdateUsername(1, request.Username)
	if err != nil {
		// Return modal with error
		creator, _ := h.creatorService.GetCreatorByID(1)
		urlSettings, _ := h.urlService.GetURLSettingsByCreatorID(1)

		component := templates.EditURLModal(creator, urlSettings, err.Error(), lang, isRTL)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// Return modal with success
	creator, _ := h.creatorService.GetCreatorByID(1)
	urlSettings, _ := h.urlService.GetURLSettingsByCreatorID(1)

	component := templates.EditURLModal(creator, urlSettings, "success", lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

package handlers

import (
	"net/http"
	"waqti/internal/middleware"
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
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
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

	urlSettings, err := h.urlService.GetURLSettingsByCreatorID(dbCreator.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error loading URL settings")
	}

	component := templates.EditURLModal(creator, urlSettings, "", lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *URLHandler) ValidateUsername(c echo.Context) error {
	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.JSON(http.StatusUnauthorized, models.URLValidationResult{
			IsValid:      false,
			ErrorMessage: "Not authenticated",
		})
	}

	var request models.URLUpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.URLValidationResult{
			IsValid:      false,
			ErrorMessage: "Invalid request",
		})
	}

	// Use the new validation method that excludes current creator
	validation := h.urlService.ValidateUsernameForCreator(request.Username, dbCreator.ID)
	return c.JSON(http.StatusOK, validation)
}

func (h *URLHandler) UpdateURL(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
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

	var request models.URLUpdateRequest
	if err := c.Bind(&request); err != nil {
		urlSettings, _ := h.urlService.GetURLSettingsByCreatorID(dbCreator.ID)
		component := templates.EditURLModal(creator, urlSettings, "Invalid form data", lang, isRTL)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	err := h.urlService.UpdateUsername(dbCreator.ID, request.Username)
	if err != nil {
		urlSettings, _ := h.urlService.GetURLSettingsByCreatorID(dbCreator.ID)
		component := templates.EditURLModal(creator, urlSettings, err.Error(), lang, isRTL)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// Successfully updated - get fresh data and show success
	urlSettings, _ := h.urlService.GetURLSettingsByCreatorID(dbCreator.ID)

	// Update the creator's username in the creator object to reflect the change
	creator.Username = urlSettings.Username

	component := templates.EditURLModal(creator, urlSettings, "success", lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

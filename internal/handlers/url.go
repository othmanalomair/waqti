package handlers

import (
	"fmt"
	"net/http"
	"waqti/internal/database"
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

	fmt.Printf("ShowEditURLModal: Creator ID: %s, Username: %s\n", dbCreator.ID, dbCreator.Username)

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
		fmt.Printf("ShowEditURLModal: Error getting URL settings: %v\n", err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading URL settings: %v", err))
	}

	if urlSettings == nil {
		fmt.Printf("ShowEditURLModal: URL settings is nil for creator %s\n", dbCreator.ID)
		return c.String(http.StatusInternalServerError, "URL settings not found")
	}

	fmt.Printf("ShowEditURLModal: URL Settings - Username: %s, Changes: %d/%d\n",
		urlSettings.Username, urlSettings.ChangesUsed, urlSettings.MaxChanges)

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

	var request models.URLUpdateRequest
	if err := c.Bind(&request); err != nil {
		// Get current URL settings for error display
		urlSettings, _ := h.urlService.GetURLSettingsByCreatorID(dbCreator.ID)
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
		component := templates.EditURLModal(creator, urlSettings, "Invalid form data", lang, isRTL)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// Log the update attempt
	fmt.Printf("Attempting to update username from %s to %s for creator %s\n",
		dbCreator.Username, request.Username, dbCreator.ID)

	err := h.urlService.UpdateUsername(dbCreator.ID, request.Username)
	if err != nil {
		fmt.Printf("Error updating username: %v\n", err)
		// Get current URL settings for error display
		urlSettings, _ := h.urlService.GetURLSettingsByCreatorID(dbCreator.ID)
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
		component := templates.EditURLModal(creator, urlSettings, err.Error(), lang, isRTL)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	fmt.Printf("Successfully updated username to %s\n", request.Username)

	// Successfully updated - get fresh data from database
	updatedCreator, err := database.Instance.GetCreatorByID(dbCreator.ID)
	if err != nil {
		fmt.Printf("Error getting updated creator: %v\n", err)
		updatedCreator = dbCreator // fallback to current creator
	}

	urlSettings, _ := h.urlService.GetURLSettingsByCreatorID(dbCreator.ID)

	// Convert to models.Creator with updated username
	creator := &models.Creator{
		ID:       updatedCreator.ID,
		Name:     updatedCreator.Name,
		NameAr:   updatedCreator.NameAr,
		Username: updatedCreator.Username, // This should now be the updated username
		Email:    updatedCreator.Email,
		Plan:     updatedCreator.Plan,
		PlanAr:   updatedCreator.PlanAr,
		IsActive: updatedCreator.IsActive,
	}

	fmt.Printf("Final creator username: %s\n", creator.Username)

	component := templates.EditURLModal(creator, urlSettings, "success", lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

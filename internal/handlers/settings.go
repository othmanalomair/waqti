package handlers

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"waqti/internal/database"
	"waqti/internal/middleware"
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

	// Get shop settings from database
	dbSettings, err := database.Instance.GetShopSettingsByCreatorID(dbCreator.ID)
	if err != nil {
		c.Logger().Error("Error loading shop settings:", err)
		return c.String(http.StatusInternalServerError, "Error loading settings")
	}

	// If no settings exist, create default ones
	if dbSettings == nil {
		dbSettings = &database.ShopSettings{
			CreatorID:        dbCreator.ID,
			LogoURL:          "",
			CreatorName:      dbCreator.Name,
			CreatorNameAr:    dbCreator.NameAr,
			SubHeader:        "",
			SubHeaderAr:      "",
			ContactWhatsApp:  "",
			CheckoutLanguage: "both",
			GreetingMessage:  "Welcome to my workshop!",
			GreetingMessageAr: "مرحباً بكم في ورشتي!",
			CurrencySymbol:   "KWD",
			CurrencySymbolAr: "د.ك",
		}
		
		err = database.Instance.CreateShopSettings(dbSettings)
		if err != nil {
			c.Logger().Error("Error creating default shop settings:", err)
			// Continue with defaults even if database save fails
		}
	}

	// Convert to models.ShopSettings for template compatibility
	settings := &models.ShopSettings{
		ID:                dbSettings.ID,
		CreatorID:         dbSettings.CreatorID,
		LogoURL:           dbSettings.LogoURL,
		CreatorName:       dbSettings.CreatorName,
		CreatorNameAr:     dbSettings.CreatorNameAr,
		SubHeader:         dbSettings.SubHeader,
		SubHeaderAr:       dbSettings.SubHeaderAr,
		ContactWhatsApp:   dbSettings.ContactWhatsApp,
		CheckoutLanguage:  dbSettings.CheckoutLanguage,
		GreetingMessage:   dbSettings.GreetingMessage,
		GreetingMessageAr: dbSettings.GreetingMessageAr,
		CurrencySymbol:    dbSettings.CurrencySymbol,
		CurrencySymbolAr:  dbSettings.CurrencySymbolAr,
		CreatedAt:         dbSettings.CreatedAt,
		UpdatedAt:         dbSettings.UpdatedAt,
	}

	successMsg := c.QueryParam("success")

	component := templates.ShopSettingsPage(creator, settings, successMsg, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *SettingsHandler) UpdateShopSettings(c echo.Context) error {
	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	// Parse form data
	creatorName := strings.TrimSpace(c.FormValue("creator_name"))
	creatorNameAr := strings.TrimSpace(c.FormValue("creator_name_ar"))
	subHeader := strings.TrimSpace(c.FormValue("sub_header"))
	subHeaderAr := strings.TrimSpace(c.FormValue("sub_header_ar"))
	contactWhatsApp := strings.TrimSpace(c.FormValue("contact_whatsapp"))
	checkoutLanguage := c.FormValue("checkout_language")
	greetingMessage := strings.TrimSpace(c.FormValue("greeting_message"))
	greetingMessageAr := strings.TrimSpace(c.FormValue("greeting_message_ar"))
	currencySymbol := strings.TrimSpace(c.FormValue("currency_symbol"))
	currencySymbolAr := strings.TrimSpace(c.FormValue("currency_symbol_ar"))

	// Validate checkout language
	if checkoutLanguage != "ar" && checkoutLanguage != "en" && checkoutLanguage != "both" {
		checkoutLanguage = "both"
	}

	// Get existing settings or create defaults
	dbSettings, err := database.Instance.GetShopSettingsByCreatorID(dbCreator.ID)
	if err != nil {
		c.Logger().Error("Error loading shop settings:", err)
		return c.String(http.StatusInternalServerError, "Error loading settings")
	}

	logoURL := ""
	if dbSettings != nil {
		logoURL = dbSettings.LogoURL
	}

	// Handle profile picture upload if provided
	profileFile, err := c.FormFile("profile_picture")
	if err == nil && profileFile != nil {
		// Validate file size (5MB limit)
		if profileFile.Size > 5*1024*1024 {
			return c.String(http.StatusBadRequest, "File size exceeds 5MB limit")
		}

		// Validate file type
		fileExt := strings.ToLower(filepath.Ext(profileFile.Filename))
		if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
			return c.String(http.StatusBadRequest, "Invalid file type. Only JPG, JPEG, and PNG are allowed")
		}

		// Generate unique filename
		src, err := profileFile.Open()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error opening uploaded file")
		}
		defer src.Close()

		// Create hash for unique filename
		hash := md5.New()
		if _, err := io.Copy(hash, src); err != nil {
			return c.String(http.StatusInternalServerError, "Error processing file")
		}
		src.Seek(0, 0) // Reset file pointer

		filename := fmt.Sprintf("%x%s", hash.Sum(nil), fileExt)
		
		// Ensure upload directory exists
		uploadDir := "web/static/images/upload"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.Logger().Error("Error creating upload directory:", err)
			return c.String(http.StatusInternalServerError, "Error creating upload directory")
		}

		uploadPath := filepath.Join(uploadDir, filename)

		// Create destination file
		dst, err := os.Create(uploadPath)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error saving file")
		}
		defer dst.Close()

		// Copy file content
		if _, err := io.Copy(dst, src); err != nil {
			return c.String(http.StatusInternalServerError, "Error saving file")
		}

		// Update logo URL
		logoURL = fmt.Sprintf("/static/images/upload/%s", filename)
	}

	// Create or update settings
	if dbSettings == nil {
		// Create new settings
		newSettings := &database.ShopSettings{
			CreatorID:         dbCreator.ID,
			LogoURL:           logoURL,
			CreatorName:       creatorName,
			CreatorNameAr:     creatorNameAr,
			SubHeader:         subHeader,
			SubHeaderAr:       subHeaderAr,
			ContactWhatsApp:   contactWhatsApp,
			CheckoutLanguage:  checkoutLanguage,
			GreetingMessage:   greetingMessage,
			GreetingMessageAr: greetingMessageAr,
			CurrencySymbol:    currencySymbol,
			CurrencySymbolAr:  currencySymbolAr,
		}

		err = database.Instance.CreateShopSettings(newSettings)
		if err != nil {
			c.Logger().Error("Error creating shop settings:", err)
			return c.String(http.StatusInternalServerError, "Error creating settings")
		}
	} else {
		// Update existing settings
		dbSettings.LogoURL = logoURL
		dbSettings.CreatorName = creatorName
		dbSettings.CreatorNameAr = creatorNameAr
		dbSettings.SubHeader = subHeader
		dbSettings.SubHeaderAr = subHeaderAr
		dbSettings.ContactWhatsApp = contactWhatsApp
		dbSettings.CheckoutLanguage = checkoutLanguage
		dbSettings.GreetingMessage = greetingMessage
		dbSettings.GreetingMessageAr = greetingMessageAr
		dbSettings.CurrencySymbol = currencySymbol
		dbSettings.CurrencySymbolAr = currencySymbolAr

		err = database.Instance.UpdateShopSettings(dbSettings)
		if err != nil {
			c.Logger().Error("Error updating shop settings:", err)
			return c.String(http.StatusInternalServerError, "Error updating settings")
		}
	}

	c.Logger().Infof("Shop settings updated successfully for creator: %s", dbCreator.Username)

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
	component := templates.ProfilePictureSection(settings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"waqti/internal/database"
	"waqti/internal/middleware"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type QRHandler struct {
	creatorService *services.CreatorService
	qrService      *services.QRService
}

func NewQRHandler(creatorService *services.CreatorService, qrService *services.QRService) *QRHandler {
	return &QRHandler{
		creatorService: creatorService,
		qrService:      qrService,
	}
}

func (h *QRHandler) ShowQRModal(c echo.Context) error {
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

	// Get shop settings directly from database for the authenticated creator
	dbShopSettings, err := database.Instance.GetShopSettingsByCreatorID(dbCreator.ID)
	var settings *models.ShopSettings
	if err == nil && dbShopSettings != nil {
		// Helper function to safely dereference string pointers
		getStringValue := func(s *string) string {
			if s != nil {
				return *s
			}
			return ""
		}
		
		settings = &models.ShopSettings{
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

	// Generate QR code for the modal
	storeURL := fmt.Sprintf("https://waqti.me/%s", dbCreator.Username)
	qrCodeDataURL, err := h.qrService.GenerateQRCode(storeURL, 192)
	if err != nil {
		fmt.Printf("Error generating QR code: %v\n", err)
		// Pass empty string if QR generation fails
		qrCodeDataURL = ""
	}

	component := templates.QRModalWithQR(creator, settings, qrCodeDataURL, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// GenerateQRCode generates QR code for the current user's store
func (h *QRHandler) GenerateQRCode(c echo.Context) error {
	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Get size parameter (default to 192)
	size := 192
	if sizeParam := c.QueryParam("size"); sizeParam != "" {
		if parsedSize, err := strconv.Atoi(sizeParam); err == nil && parsedSize > 0 && parsedSize <= 1000 {
			size = parsedSize
		}
	}

	// Generate store URL
	storeURL := fmt.Sprintf("https://waqti.me/%s", dbCreator.Username)

	// Generate QR code as base64 data URL
	dataURL, err := h.qrService.GenerateQRCode(storeURL, size)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate QR code"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"qrCode": dataURL,
		"url":    storeURL,
	})
}

// DownloadQRCode generates and downloads QR code as PNG file
func (h *QRHandler) DownloadQRCode(c echo.Context) error {
	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Get size parameter (default to 300 for downloads)
	size := 300
	if sizeParam := c.QueryParam("size"); sizeParam != "" {
		if parsedSize, err := strconv.Atoi(sizeParam); err == nil && parsedSize > 0 && parsedSize <= 1000 {
			size = parsedSize
		}
	}

	// Generate store URL
	storeURL := fmt.Sprintf("https://waqti.me/%s", dbCreator.Username)

	// Generate QR code bytes
	qrBytes, err := h.qrService.GenerateQRCodeBytes(storeURL, size)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate QR code"})
	}

	// Set headers for file download
	filename := fmt.Sprintf("qr-code-%s.png", dbCreator.Username)
	c.Response().Header().Set("Content-Type", "image/png")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", len(qrBytes)))

	return c.Blob(http.StatusOK, "image/png", qrBytes)
}

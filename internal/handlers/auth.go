package handlers

import (
	"net/http"
	"strings"
	"time"
	"waqti/internal/database"
	"waqti/internal/middleware"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService     *middleware.AuthService
	workshopService *services.WorkshopService
	settingsService *services.SettingsService
}

func NewAuthHandler(authService *middleware.AuthService, workshopService *services.WorkshopService, settingsService *services.SettingsService) *AuthHandler {
	return &AuthHandler{
		authService:     authService,
		workshopService: workshopService,
		settingsService: settingsService,
	}
}

func (h *AuthHandler) ShowLandingPage(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Check if user was just signed out
	signedOut := c.QueryParam("signed_out")

	component := templates.LandingPageWithSignout(lang, isRTL, signedOut)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ToggleLanguage(c echo.Context) error {
	currentLang := c.Get("lang").(string)
	newLang := "en"
	if currentLang == "en" {
		newLang = "ar"
	}

	cookie := &http.Cookie{
		Name:     "lang",
		Value:    newLang,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	redirectTo := c.FormValue("redirect_to")
	if redirectTo == "" {
		redirectTo = "/"
	}

	return c.Redirect(http.StatusSeeOther, redirectTo)
}

func (h *AuthHandler) ShowSignIn(c echo.Context) error {
	// If user is already logged in, redirect to dashboard
	if creator := middleware.GetCurrentCreator(c); creator != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}

	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)
	errorMsg := c.QueryParam("error")

	component := templates.SignInPage(errorMsg, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ProcessSignIn(c echo.Context) error {
	email := strings.TrimSpace(c.FormValue("email"))
	password := c.FormValue("password")

	// Validate input
	if email == "" || password == "" {
		return c.Redirect(http.StatusSeeOther, "/signin?error=empty_fields")
	}

	// Attempt login
	creator, err := h.authService.LoginCreator(c, email, password)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			return c.Redirect(http.StatusSeeOther, "/signin?error=invalid_credentials")
		}
		// Log the actual error but show generic message to user
		c.Logger().Error("Login error:", err)
		return c.Redirect(http.StatusSeeOther, "/signin?error=server_error")
	}

	// Successful login
	c.Logger().Infof("User %s (%s) logged in successfully", creator.Email, creator.ID)
	return c.Redirect(http.StatusSeeOther, "/dashboard")
}

func (h *AuthHandler) ShowSignUp(c echo.Context) error {
	// If user is already logged in, redirect to dashboard
	if creator := middleware.GetCurrentCreator(c); creator != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}

	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)
	errorMsg := c.QueryParam("error")
	successMsg := c.QueryParam("success")

	component := templates.SignUpPage(errorMsg, successMsg, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ProcessSignUp(c echo.Context) error {
	name := strings.TrimSpace(c.FormValue("name"))
	email := strings.TrimSpace(c.FormValue("email"))
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	// Validate input
	if name == "" || email == "" || password == "" {
		return c.Redirect(http.StatusSeeOther, "/signup?error=empty_fields")
	}

	if password != confirmPassword {
		return c.Redirect(http.StatusSeeOther, "/signup?error=password_mismatch")
	}

	if len(password) < 6 {
		return c.Redirect(http.StatusSeeOther, "/signup?error=password_too_short")
	}

	// Generate username from name (basic implementation)
	username := generateUsername(name)
	nameAr := name // For now, use the same name for Arabic

	// Attempt registration
	creator, err := h.authService.RegisterCreator(name, nameAr, username, email, password)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			return c.Redirect(http.StatusSeeOther, "/signup?error=email_exists")
		}
		if strings.Contains(err.Error(), "username already exists") {
			return c.Redirect(http.StatusSeeOther, "/signup?error=username_exists")
		}
		// Log the actual error but show generic message to user
		c.Logger().Error("Registration error:", err)
		return c.Redirect(http.StatusSeeOther, "/signup?error=server_error")
	}

	// Successful registration
	c.Logger().Infof("User %s (%s) registered successfully", creator.Email, creator.ID)
	return c.Redirect(http.StatusSeeOther, "/signup?success=account_created")
}

func (h *AuthHandler) ShowStorePage(c echo.Context) error {
	username := c.Param("username")

	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get creator by username from database
	creator, err := database.Instance.GetCreatorByUsername(username)
	if err != nil {
		c.Logger().Error("Error getting creator:", err)
		return c.String(http.StatusInternalServerError, "Server error")
	}
	if creator == nil {
		return c.String(http.StatusNotFound, "Creator not found")
	}

	// Get shop settings from database
	dbSettings, err := database.Instance.GetShopSettingsByCreatorID(creator.ID)
	if err != nil {
		c.Logger().Error("Error getting shop settings:", err)
	}

	// Convert database settings to models.ShopSettings or use defaults
	var settings *models.ShopSettings
	if dbSettings != nil {
		logoURL := ""
		if dbSettings.LogoURL != nil {
			logoURL = *dbSettings.LogoURL
		}
		
		creatorName := ""
		if dbSettings.CreatorName != nil {
			creatorName = *dbSettings.CreatorName
		}
		
		creatorNameAr := ""
		if dbSettings.CreatorNameAr != nil {
			creatorNameAr = *dbSettings.CreatorNameAr
		}
		
		subHeader := ""
		if dbSettings.SubHeader != nil {
			subHeader = *dbSettings.SubHeader
		}
		
		subHeaderAr := ""
		if dbSettings.SubHeaderAr != nil {
			subHeaderAr = *dbSettings.SubHeaderAr
		}
		
		contactWhatsApp := ""
		if dbSettings.ContactWhatsApp != nil {
			contactWhatsApp = *dbSettings.ContactWhatsApp
		}
		
		greetingMessage := ""
		if dbSettings.GreetingMessage != nil {
			greetingMessage = *dbSettings.GreetingMessage
		}
		
		greetingMessageAr := ""
		if dbSettings.GreetingMessageAr != nil {
			greetingMessageAr = *dbSettings.GreetingMessageAr
		}
		
		settings = &models.ShopSettings{
			ID:                dbSettings.ID,
			CreatorID:         dbSettings.CreatorID,
			LogoURL:           logoURL,
			CreatorName:       creatorName,
			CreatorNameAr:     creatorNameAr,
			SubHeader:         subHeader,
			SubHeaderAr:       subHeaderAr,
			ContactWhatsApp:   contactWhatsApp,
			CheckoutLanguage:  dbSettings.CheckoutLanguage,
			GreetingMessage:   greetingMessage,
			GreetingMessageAr: greetingMessageAr,
			CurrencySymbol:    dbSettings.CurrencySymbol,
			CurrencySymbolAr:  dbSettings.CurrencySymbolAr,
			CreatedAt:         dbSettings.CreatedAt,
			UpdatedAt:         dbSettings.UpdatedAt,
		}
	} else {
		// Use default settings if not found
		settings = &models.ShopSettings{
			CreatorName:     creator.Name,
			CreatorNameAr:   creator.NameAr,
			SubHeader:       "Certified Design Trainer",
			SubHeaderAr:     "مدرب معتمد في التصميم",
			ContactWhatsApp: "+965-9999-7777",
			LogoURL:         "/static/images/creator-avatar.jpg",
		}
	}

	// Get workshops with upcoming sessions for this creator
	workshops := h.workshopService.GetActiveWorkshopsWithUpcomingSessions(creator.ID)

	// Enhance workshops with images
	for i := range workshops {
		// Get workshop images
		images, err := h.workshopService.GetWorkshopImagesByWorkshopID(workshops[i].ID)
		if err != nil {
			c.Logger().Error("Error getting workshop images:", err)
		} else {
			workshops[i].Images = images
		}
	}

	// Convert database.Creator to models.Creator for template compatibility
	templateCreator := &models.Creator{
		ID:       creator.ID,
		Name:     creator.Name,
		NameAr:   creator.NameAr,
		Username: creator.Username,
		Email:    creator.Email,
		Plan:     creator.Plan,
		PlanAr:   creator.PlanAr,
		IsActive: creator.IsActive,
	}

	component := templates.StorePage(templateCreator, workshops, settings, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ProcessSignOut(c echo.Context) error {
	// Logout the user
	if err := h.authService.LogoutCreator(c); err != nil {
		c.Logger().Error("Logout error:", err)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

// generateUsername creates a username from a name
func generateUsername(name string) string {
	// Simple implementation: lowercase, remove spaces, limit length
	username := strings.ToLower(name)
	username = strings.ReplaceAll(username, " ", "")
	username = strings.ReplaceAll(username, ".", "")

	// Remove non-alphanumeric characters except underscores and hyphens
	var result strings.Builder
	for _, r := range username {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			result.WriteRune(r)
		}
	}

	username = result.String()

	// Ensure minimum length and maximum length
	if len(username) < 3 {
		username = username + "123"
	}
	if len(username) > 20 {
		username = username[:20]
	}

	return username
}

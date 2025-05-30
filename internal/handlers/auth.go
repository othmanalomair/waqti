package handlers

import (
	"net/http"
	"time"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	creatorService  *services.CreatorService
	workshopService *services.WorkshopService
}

func NewAuthHandler(creatorService *services.CreatorService) *AuthHandler {
	return &AuthHandler{
		creatorService:  creatorService,
		workshopService: services.NewWorkshopService(),
	}
}

func (h *AuthHandler) ShowLandingPage(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Render landing page
	component := templates.LandingPage(lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ToggleLanguage(c echo.Context) error {
	currentLang := c.Get("lang").(string)
	newLang := "en"
	if currentLang == "en" {
		newLang = "ar"
	}

	// Set language cookie
	cookie := &http.Cookie{
		Name:     "lang",
		Value:    newLang,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	// Get redirect URL from form data, default to home page
	redirectTo := c.FormValue("redirect_to")
	if redirectTo == "" {
		redirectTo = "/"
	}

	// Redirect back to the same page
	return c.Redirect(http.StatusSeeOther, redirectTo)
}

func (h *AuthHandler) ShowSignIn(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Check for error message
	errorMsg := c.QueryParam("error")

	// Render sign in page
	component := templates.SignInPage(errorMsg, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ProcessSignIn(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Simple validation for demo
	if email == "" || password == "" {
		return c.Redirect(http.StatusSeeOther, "/signin?error=empty_fields")
	}

	// Demo authentication - accept any email/password for now
	if email == "demo@waqti.me" && password == "password" {
		// Set auth cookie (in real app, use proper JWT or session)
		cookie := &http.Cookie{
			Name:     "auth",
			Value:    "authenticated",
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		c.SetCookie(cookie)

		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}

	return c.Redirect(http.StatusSeeOther, "/signin?error=invalid_credentials")
}

func (h *AuthHandler) ShowSignUp(c echo.Context) error {
	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Check for error message
	errorMsg := c.QueryParam("error")
	successMsg := c.QueryParam("success")

	// Render sign up page
	component := templates.SignUpPage(errorMsg, successMsg, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ProcessSignUp(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	// Simple validation for demo
	if name == "" || email == "" || password == "" {
		return c.Redirect(http.StatusSeeOther, "/signup?error=empty_fields")
	}

	if password != confirmPassword {
		return c.Redirect(http.StatusSeeOther, "/signup?error=password_mismatch")
	}

	if len(password) < 6 {
		return c.Redirect(http.StatusSeeOther, "/signup?error=password_too_short")
	}

	// Demo - just show success and redirect to sign in
	return c.Redirect(http.StatusSeeOther, "/signup?success=account_created")
}

func (h *AuthHandler) ShowStorePage(c echo.Context) error {
	// Get username from URL parameter
	username := c.Param("username")

	// Get language from context
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	// Get creator by username (dummy data for now)
	creator, err := h.creatorService.GetCreatorByUsername(username)
	if err != nil || creator == nil {
		// Return 404 if creator not found
		return c.String(http.StatusNotFound, "Creator not found")
	}

	// Get creator's workshops/courses (dummy data)
	workshops := h.workshopService.GetWorkshopsByCreatorID(creator.ID)

	// Render store page
	component := templates.StorePage(creator, workshops, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ProcessSignOut(c echo.Context) error {
	// Clear auth cookie
	cookie := &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/")
}

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
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	component := templates.LandingPage(lang, isRTL)
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
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	errorMsg := c.QueryParam("error")

	component := templates.SignInPage(errorMsg, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ProcessSignIn(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.Redirect(http.StatusSeeOther, "/signin?error=empty_fields")
	}

	if email == "demo@waqti.me" && password == "password" {
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
	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	errorMsg := c.QueryParam("error")
	successMsg := c.QueryParam("success")

	component := templates.SignUpPage(errorMsg, successMsg, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ProcessSignUp(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	if name == "" || email == "" || password == "" {
		return c.Redirect(http.StatusSeeOther, "/signup?error=empty_fields")
	}

	if password != confirmPassword {
		return c.Redirect(http.StatusSeeOther, "/signup?error=password_mismatch")
	}

	if len(password) < 6 {
		return c.Redirect(http.StatusSeeOther, "/signup?error=password_too_short")
	}

	return c.Redirect(http.StatusSeeOther, "/signup?success=account_created")
}

func (h *AuthHandler) ShowStorePage(c echo.Context) error {
	username := c.Param("username")

	lang := c.Get("lang").(string)
	isRTL := c.Get("isRTL").(bool)

	creator, err := h.creatorService.GetCreatorByUsername(username)
	if err != nil || creator == nil {
		return c.String(http.StatusNotFound, "Creator not found")
	}

	workshops := h.workshopService.GetWorkshopsByCreatorID(creator.ID)

	component := templates.StorePage(creator, workshops, lang, isRTL)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) ProcessSignOut(c echo.Context) error {
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

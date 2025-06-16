package middleware

import (
	"net/http"
	"strings"
	"waqti/internal/services"

	"github.com/labstack/echo/v4"
)

// AdminAuth middleware for admin authentication
func AdminAuth(adminAuthService *services.AdminAuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get session token from cookie
			cookie, err := c.Cookie("admin_session")
			if err != nil {
				return redirectToAdminLogin(c)
			}

			// Validate session
			adminUser, err := adminAuthService.ValidateSession(cookie.Value)
			if err != nil {
				// Clear invalid cookie
				clearAdminSessionCookie(c)
				return redirectToAdminLogin(c)
			}

			// Store admin user in context
			c.Set("admin_user", adminUser)
			return next(c)
		}
	}
}

// AdminAuthOptional middleware for optional admin authentication
func AdminAuthOptional(adminAuthService *services.AdminAuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get session token from cookie
			cookie, err := c.Cookie("admin_session")
			if err == nil {
				// Validate session
				adminUser, err := adminAuthService.ValidateSession(cookie.Value)
				if err == nil {
					// Store admin user in context if valid
					c.Set("admin_user", adminUser)
				} else {
					// Clear invalid cookie
					clearAdminSessionCookie(c)
				}
			}

			return next(c)
		}
	}
}

// RequireAdminRole middleware to require specific admin role
func RequireAdminRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			adminUser, ok := c.Get("admin_user").(*services.AdminUser)
			if !ok {
				return redirectToAdminLogin(c)
			}

			// Check if user has required role
			if role == "super_admin" && adminUser.Role != "super_admin" {
				return echo.NewHTTPError(http.StatusForbidden, "Super admin access required")
			}

			return next(c)
		}
	}
}

// NoAdminAuth middleware to prevent access if already logged in as admin
func NoAdminAuth(adminAuthService *services.AdminAuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if already logged in as admin
			cookie, err := c.Cookie("admin_session")
			if err == nil {
				// Validate session
				_, err := adminAuthService.ValidateSession(cookie.Value)
				if err == nil {
					// Already logged in, redirect to admin dashboard
					return c.Redirect(http.StatusFound, "/admin/dashboard")
				}
			}

			return next(c)
		}
	}
}

// Helper functions
func redirectToAdminLogin(c echo.Context) error {
	// If this is an AJAX request, return JSON error
	if strings.Contains(c.Request().Header.Get("Accept"), "application/json") ||
		strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		return echo.NewHTTPError(http.StatusUnauthorized, "Admin authentication required")
	}

	// For regular requests, redirect to admin login
	return c.Redirect(http.StatusFound, "/admin")
}

func clearAdminSessionCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
}

// SetAdminSessionCookie helper to set admin session cookie
func SetAdminSessionCookie(c echo.Context, sessionToken string) {
	cookie := &http.Cookie{
		Name:     "admin_session",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60, // 30 days in seconds
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
}
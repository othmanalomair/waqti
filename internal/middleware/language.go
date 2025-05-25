package middleware

import (
	"github.com/labstack/echo/v4"
)

func LanguageMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check for language preference in cookie
			lang := "ar" // default to Arabic
			if cookie, err := c.Cookie("lang"); err == nil {
				if cookie.Value == "en" || cookie.Value == "ar" {
					lang = cookie.Value
				}
			}

			// Set language in context
			c.Set("lang", lang)
			c.Set("isRTL", lang == "ar")

			return next(c)
		}
	}
}

package middleware

import (
	"log"
	"net"
	"strings"
	"waqti/internal/services"

	"github.com/labstack/echo/v4"
)

// AnalyticsTrackingMiddleware tracks page visits for admin analytics
func AnalyticsTrackingMiddleware(adminAnalyticsService *services.AdminAnalyticsService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get request details
			path := c.Request().URL.Path
			userAgent := c.Request().UserAgent()
			referrer := c.Request().Referer()
			
			// Get IP address
			ipAddress := getClientIP(c)
			
			// Determine device type
			device := getDeviceType(userAgent)
			
			// Determine browser
			browser := getBrowserType(userAgent)
			
			// Map paths to page types for tracking
			pageType := getPageType(path)
			
			// Track the visit if it's a page we're interested in
			if pageType != "" {
				// Extract store username for store visits
				storeUsername := ""
				if pageType == "store_visit" {
					storeUsername = extractUsernameFromPath(path)
				}
				
				// Track asynchronously to not slow down the request
				go func() {
					err := adminAnalyticsService.TrackPageVisit(
						pageType,
						ipAddress,
						userAgent,
						referrer,
						"", // country - would need geolocation service
						device,
						browser,
						storeUsername,
					)
					if err != nil {
						log.Printf("Error tracking page visit: %v", err)
					}
				}()
			}

			return next(c)
		}
	}
}

// getClientIP extracts the client IP address from the request
func getClientIP(c echo.Context) string {
	// Check X-Forwarded-For header first (for reverse proxies)
	xForwardedFor := c.Request().Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, we want the first one
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	xRealIP := c.Request().Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if err != nil {
		return c.Request().RemoteAddr
	}
	return ip
}

// getDeviceType determines device type from user agent
func getDeviceType(userAgent string) string {
	userAgent = strings.ToLower(userAgent)
	
	// Mobile indicators
	mobileIndicators := []string{
		"mobile", "android", "iphone", "ipad", "ipod", 
		"blackberry", "windows phone", "nokia", "samsung",
		"htc", "lg", "motorola", "sony", "xiaomi", "huawei",
	}
	
	for _, indicator := range mobileIndicators {
		if strings.Contains(userAgent, indicator) {
			return "Mobile"
		}
	}
	
	// Tablet indicators
	tabletIndicators := []string{"tablet", "ipad"}
	for _, indicator := range tabletIndicators {
		if strings.Contains(userAgent, indicator) {
			return "Tablet"
		}
	}
	
	return "Desktop"
}

// getBrowserType determines browser type from user agent
func getBrowserType(userAgent string) string {
	userAgent = strings.ToLower(userAgent)
	
	// Browser detection logic
	if strings.Contains(userAgent, "chrome") && !strings.Contains(userAgent, "edge") {
		return "Chrome"
	}
	if strings.Contains(userAgent, "firefox") {
		return "Firefox"
	}
	if strings.Contains(userAgent, "safari") && !strings.Contains(userAgent, "chrome") {
		return "Safari"
	}
	if strings.Contains(userAgent, "edge") {
		return "Edge"
	}
	if strings.Contains(userAgent, "opera") {
		return "Opera"
	}
	if strings.Contains(userAgent, "msie") || strings.Contains(userAgent, "trident") {
		return "Internet Explorer"
	}
	
	return "Other"
}

// getPageType maps URL paths to trackable page types
func getPageType(path string) string {
	switch path {
	case "/":
		return "landing"
	case "/signin":
		return "signin"
	case "/signup":
		return "signup"
	default:
		// Check if it's a store visit (/:username pattern)
		if isUsernameRoute(path) {
			return "store_visit"
		}
		return "" // Don't track other pages
	}
}

// extractUsernameFromPath extracts the username from a /:username path
func extractUsernameFromPath(path string) string {
	// Remove leading slash
	if path == "/" {
		return ""
	}
	
	pathWithoutSlash := strings.TrimPrefix(path, "/")
	
	// Should not contain additional slashes (simple username check)
	if strings.Contains(pathWithoutSlash, "/") {
		return ""
	}
	
	// Return the username if it's a valid format
	if len(pathWithoutSlash) >= 3 && len(pathWithoutSlash) <= 50 {
		return pathWithoutSlash
	}
	
	return ""
}

// TrackPageVisitManual allows manual tracking of page visits
func TrackPageVisitManual(c echo.Context, adminAnalyticsService *services.AdminAnalyticsService, pageType string) {
	// Get request details
	userAgent := c.Request().UserAgent()
	referrer := c.Request().Referer()
	ipAddress := getClientIP(c)
	device := getDeviceType(userAgent)
	browser := getBrowserType(userAgent)
	
	// Extract store username if it's a store visit
	storeUsername := ""
	if pageType == "store_visit" {
		storeUsername = extractUsernameFromPath(c.Request().URL.Path)
	}
	
	// Track asynchronously
	go func() {
		err := adminAnalyticsService.TrackPageVisit(
			pageType,
			ipAddress,
			userAgent,
			referrer,
			"", // country
			device,
			browser,
			storeUsername,
		)
		if err != nil {
			log.Printf("Error tracking manual page visit: %v", err)
		}
	}()
}

// Enhanced analytics middleware for specific tracking needs
func EnhancedAnalyticsMiddleware(adminAnalyticsService *services.AdminAnalyticsService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Store the analytics service in context for manual tracking
			c.Set("adminAnalyticsService", adminAnalyticsService)
			
			return next(c)
		}
	}
}

// GetAnalyticsService retrieves the admin analytics service from context
func GetAnalyticsService(c echo.Context) *services.AdminAnalyticsService {
	if service, ok := c.Get("adminAnalyticsService").(*services.AdminAnalyticsService); ok {
		return service
	}
	return nil
}
package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
	"waqti/internal/database"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	SessionCookieName = "waqti_session"
	SessionDuration   = 30 * 24 * time.Hour // 30 days
	ContextKeyCreator = "creator"
	ContextKeySession = "session"
)

// AuthService handles authentication logic
type AuthService struct {
	db *database.DB
}

func NewAuthService(db *database.DB) *AuthService {
	return &AuthService{db: db}
}

// GenerateSessionToken generates a secure random session token
func (as *AuthService) GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashPassword hashes a password using bcrypt
func (as *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword checks if a password matches the hash
func (as *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateSession creates a new session for a creator
func (as *AuthService) CreateSession(creatorID uuid.UUID, userAgent, ipAddress string) (*database.Session, error) {
	token, err := as.GenerateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	session := &database.Session{
		CreatorID:    creatorID,
		SessionToken: token,
		DeviceInfo:   &userAgent,
		IPAddress:    &ipAddress,
		ExpiresAt:    time.Now().Add(SessionDuration),
	}

	if err := as.db.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// ValidateSession validates a session token and returns the associated creator
func (as *AuthService) ValidateSession(token string) (*database.Creator, *database.Session, error) {
	session, err := as.db.GetSessionByToken(token)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return nil, nil, nil // Session not found or expired
	}

	creator, err := as.db.GetCreatorByID(session.CreatorID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get creator: %w", err)
	}
	if creator == nil {
		// Creator not found, delete the invalid session
		as.db.DeleteSession(token)
		return nil, nil, nil
	}

	// Update session activity
	if err := as.db.UpdateSessionActivity(session.ID); err != nil {
		// Log but don't fail - this is not critical
		fmt.Printf("Warning: failed to update session activity: %v\n", err)
	}

	return creator, session, nil
}

// AuthMiddleware provides authentication middleware
func AuthMiddleware(authService *AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get session token from cookie
			cookie, err := c.Cookie(SessionCookieName)
			if err != nil {
				// No session cookie found
				return redirectToSignin(c)
			}

			// Validate session
			creator, session, err := authService.ValidateSession(cookie.Value)
			if err != nil {
				c.Logger().Errorf("Session validation error: %v", err)
				return redirectToSignin(c)
			}
			if creator == nil || session == nil {
				// Invalid or expired session
				clearSessionCookie(c)
				return redirectToSignin(c)
			}

			// Set creator and session in context
			c.Set(ContextKeyCreator, creator)
			c.Set(ContextKeySession, session)

			return next(c)
		}
	}
}

// OptionalAuthMiddleware provides optional authentication (doesn't redirect)
func OptionalAuthMiddleware(authService *AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get session token from cookie
			cookie, err := c.Cookie(SessionCookieName)
			if err != nil {
				// No session cookie found, continue without auth
				return next(c)
			}

			// Validate session
			creator, session, err := authService.ValidateSession(cookie.Value)
			if err != nil {
				c.Logger().Errorf("Session validation error: %v", err)
				clearSessionCookie(c)
				return next(c)
			}
			if creator == nil || session == nil {
				// Invalid or expired session
				clearSessionCookie(c)
				return next(c)
			}

			// Set creator and session in context
			c.Set(ContextKeyCreator, creator)
			c.Set(ContextKeySession, session)

			return next(c)
		}
	}
}

// GetCurrentCreator gets the current authenticated creator from context
func GetCurrentCreator(c echo.Context) *database.Creator {
	if creator, ok := c.Get(ContextKeyCreator).(*database.Creator); ok {
		return creator
	}
	return nil
}

// GetCurrentSession gets the current session from context
func GetCurrentSession(c echo.Context) *database.Session {
	if session, ok := c.Get(ContextKeySession).(*database.Session); ok {
		return session
	}
	return nil
}

// RequireAuth ensures a user is authenticated and returns creator or redirects
func RequireAuth(c echo.Context) *database.Creator {
	creator := GetCurrentCreator(c)
	if creator == nil {
		// This should trigger a redirect to signin
		panic("RequireAuth called on unauthenticated user - this should be caught by middleware")
	}
	return creator
}

// SetSessionCookie sets the session cookie
func SetSessionCookie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(SessionDuration),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   isHTTPS(c),
	}
	c.SetCookie(cookie)
}

// ClearSessionCookie clears the session cookie
func clearSessionCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   isHTTPS(c),
	}
	c.SetCookie(cookie)
}

// redirectToSignin redirects to the signin page
func redirectToSignin(c echo.Context) error {
	// For HTMX requests, return 401 to trigger client-side redirect
	if c.Request().Header.Get("HX-Request") == "true" {
		c.Response().Header().Set("HX-Redirect", "/signin")
		return c.NoContent(http.StatusUnauthorized)
	}

	// For regular requests, do a server-side redirect
	return c.Redirect(http.StatusSeeOther, "/signin")
}

// isHTTPS checks if the request is HTTPS
func isHTTPS(c echo.Context) bool {
	return c.Scheme() == "https" ||
		c.Request().Header.Get("X-Forwarded-Proto") == "https"
}

// IsPublicRoute checks if a route should be publicly accessible
func IsPublicRoute(path string) bool {
	publicRoutes := map[string]bool{
		"/":                true,
		"/signin":          true,
		"/signup":          true,
		"/toggle-language": true,
		"/api/orders":      true,
		"/health":          true,
		"/static":          true, // This will be handled by prefix
	}

	// Check exact matches
	if publicRoutes[path] {
		return true
	}

	// Check prefixes
	publicPrefixes := []string{
		"/static/",
	}

	for _, prefix := range publicPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	// Check if it's a username route (/:username)
	// This should be the last check as it's a catch-all
	if isUsernameRoute(path) {
		return true
	}

	return false
}

// isUsernameRoute checks if the path matches the /:username pattern
func isUsernameRoute(path string) bool {
	// Remove leading slash
	if path == "/" {
		return false
	}

	pathWithoutSlash := strings.TrimPrefix(path, "/")

	// Should not contain additional slashes (simple username check)
	if strings.Contains(pathWithoutSlash, "/") {
		return false
	}

	// Should be a valid username format (basic check)
	if len(pathWithoutSlash) < 3 || len(pathWithoutSlash) > 50 {
		return false
	}

	return true
}

// ConditionalAuthMiddleware applies auth middleware only to protected routes
func ConditionalAuthMiddleware(authService *AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			// If it's a public route, use optional auth
			if IsPublicRoute(path) {
				return OptionalAuthMiddleware(authService)(next)(c)
			}

			// Otherwise, require authentication
			return AuthMiddleware(authService)(next)(c)
		}
	}
}

// LoginCreator handles the login process
func (as *AuthService) LoginCreator(c echo.Context, email, password string) (*database.Creator, error) {
	// Get creator by email
	creator, err := as.db.GetCreatorByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if creator == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if !as.CheckPassword(password, creator.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Create session
	userAgent := c.Request().UserAgent()
	ipAddress := c.RealIP()

	session, err := as.CreateSession(creator.ID, userAgent, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Set session cookie
	SetSessionCookie(c, session.SessionToken)

	return creator, nil
}

// LogoutCreator handles the logout process
func (as *AuthService) LogoutCreator(c echo.Context) error {
	// Get session token from cookie
	cookie, err := c.Cookie(SessionCookieName)
	if err == nil {
		// Delete session from database
		as.db.DeleteSession(cookie.Value)
	}

	// Clear session cookie
	clearSessionCookie(c)

	return nil
}

// RegisterCreator handles the registration process and creates URL settings
func (as *AuthService) RegisterCreator(name, nameAr, username, email, password string) (*database.Creator, error) {
	// Check if email already exists
	exists, err := as.db.CheckEmailExists(email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Check if username already exists
	exists, err = as.db.CheckUsernameExists(username)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// Hash password
	hashedPassword, err := as.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Start transaction
	tx, err := as.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create creator
	creator := &database.Creator{
		Name:         name,
		NameAr:       nameAr,
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Plan:         "free",
		PlanAr:       "مجاني",
		IsActive:     true,
	}

	// Insert creator and get the ID
	creatorQuery := `
		INSERT INTO creators (name, name_ar, username, email, password_hash, plan, plan_ar, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(creatorQuery,
		creator.Name, creator.NameAr, creator.Username,
		creator.Email, creator.PasswordHash, creator.Plan, creator.PlanAr, creator.IsActive,
	).Scan(&creator.ID, &creator.CreatedAt, &creator.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create creator: %w", err)
	}

	// Create URL settings for the new creator
	urlSettingsQuery := `
		INSERT INTO url_settings (creator_id, username, changes_used, max_changes)
		VALUES ($1, $2, 0, 5)
	`
	_, err = tx.Exec(urlSettingsQuery, creator.ID, creator.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL settings: %w", err)
	}

	// Create shop settings for the new creator
	shopSettingsQuery := `
		INSERT INTO shop_settings (creator_id, creator_name, creator_name_ar)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(shopSettingsQuery, creator.ID, creator.Name, creator.NameAr)
	if err != nil {
		return nil, fmt.Errorf("failed to create shop settings: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return creator, nil
}

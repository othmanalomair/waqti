package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
	"waqti/internal/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AdminUser represents an admin user
type AdminUser struct {
	ID           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	NameAr       *string    `json:"name_ar"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// AdminSession represents an admin session
type AdminSession struct {
	ID           uuid.UUID `json:"id"`
	AdminUserID  uuid.UUID `json:"admin_user_id"`
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	LastAccessed time.Time `json:"last_accessed"`
}

// AdminAuthService handles admin authentication
type AdminAuthService struct {
	db *database.DB
}

// NewAdminAuthService creates a new admin auth service
func NewAdminAuthService() *AdminAuthService {
	return &AdminAuthService{
		db: database.Instance,
	}
}

// Login authenticates an admin user and creates a session
func (as *AdminAuthService) Login(username, password string) (*AdminUser, string, error) {
	// Get admin user by username or email
	user, err := as.GetAdminUserByUsernameOrEmail(username)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, "", fmt.Errorf("account is deactivated")
	}

	// Create session
	sessionToken, err := as.CreateSession(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	err = as.UpdateLastLogin(user.ID)
	if err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	return user, sessionToken, nil
}

// GetAdminUserByUsernameOrEmail retrieves admin user by username or email
func (as *AdminAuthService) GetAdminUserByUsernameOrEmail(identifier string) (*AdminUser, error) {
	query := `
		SELECT id, username, email, name, name_ar, password_hash, role, is_active, last_login, created_at, updated_at
		FROM admin_users 
		WHERE (username = $1 OR email = $1) AND is_active = true
	`

	var user AdminUser
	err := as.db.QueryRow(query, identifier).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name, &user.NameAr,
		&user.PasswordHash, &user.Role, &user.IsActive, &user.LastLogin,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin user not found")
		}
		return nil, fmt.Errorf("failed to get admin user: %w", err)
	}

	return &user, nil
}

// CreateSession creates a new admin session
func (as *AdminAuthService) CreateSession(adminUserID uuid.UUID) (string, error) {
	// Clean up expired sessions first
	as.CleanupExpiredSessions()

	// Generate session token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}
	sessionToken := hex.EncodeToString(tokenBytes)

	// Session expires in 30 days
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	query := `
		INSERT INTO admin_sessions (admin_user_id, session_token, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := as.db.Exec(query, adminUserID, sessionToken, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return sessionToken, nil
}

// ValidateSession validates an admin session token and returns the admin user
func (as *AdminAuthService) ValidateSession(sessionToken string) (*AdminUser, error) {
	// Get session info with user data
	query := `
		SELECT 
			au.id, au.username, au.email, au.name, au.name_ar, au.role, au.is_active,
			au.last_login, au.created_at, au.updated_at,
			s.expires_at
		FROM admin_sessions s
		JOIN admin_users au ON s.admin_user_id = au.id
		WHERE s.session_token = $1 AND s.expires_at > NOW() AND au.is_active = true
	`

	var user AdminUser
	var expiresAt time.Time

	err := as.db.QueryRow(query, sessionToken).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name, &user.NameAr,
		&user.Role, &user.IsActive, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
		&expiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid or expired session")
		}
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}

	// Update last accessed time
	updateQuery := `UPDATE admin_sessions SET last_accessed = NOW() WHERE session_token = $1`
	as.db.Exec(updateQuery, sessionToken)

	return &user, nil
}

// Logout removes an admin session
func (as *AdminAuthService) Logout(sessionToken string) error {
	query := `DELETE FROM admin_sessions WHERE session_token = $1`
	_, err := as.db.Exec(query, sessionToken)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

// UpdateLastLogin updates the last login timestamp
func (as *AdminAuthService) UpdateLastLogin(adminUserID uuid.UUID) error {
	query := `UPDATE admin_users SET last_login = NOW() WHERE id = $1`
	_, err := as.db.Exec(query, adminUserID)
	return err
}

// CleanupExpiredSessions removes expired admin sessions
func (as *AdminAuthService) CleanupExpiredSessions() error {
	query := `DELETE FROM admin_sessions WHERE expires_at < NOW()`
	_, err := as.db.Exec(query)
	return err
}

// CreateAdminUser creates a new admin user
func (as *AdminAuthService) CreateAdminUser(username, email, name, nameAr, password, role string) (*AdminUser, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO admin_users (username, email, name, name_ar, password_hash, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, email, name, name_ar, role, is_active, created_at, updated_at
	`

	var user AdminUser
	err = as.db.QueryRow(query, username, email, name, nameAr, string(hashedPassword), role).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name, &user.NameAr,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	return &user, nil
}

// GetAdminUserByID retrieves an admin user by ID
func (as *AdminAuthService) GetAdminUserByID(id uuid.UUID) (*AdminUser, error) {
	query := `
		SELECT id, username, email, name, name_ar, role, is_active, last_login, created_at, updated_at
		FROM admin_users 
		WHERE id = $1
	`

	var user AdminUser
	err := as.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name, &user.NameAr,
		&user.Role, &user.IsActive, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin user not found")
		}
		return nil, fmt.Errorf("failed to get admin user: %w", err)
	}

	return &user, nil
}

// GetAllAdminUsers retrieves all admin users
func (as *AdminAuthService) GetAllAdminUsers() ([]AdminUser, error) {
	query := `
		SELECT id, username, email, name, name_ar, role, is_active, last_login, created_at, updated_at
		FROM admin_users 
		ORDER BY created_at DESC
	`

	rows, err := as.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin users: %w", err)
	}
	defer rows.Close()

	var users []AdminUser
	for rows.Next() {
		var user AdminUser
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Name, &user.NameAr,
			&user.Role, &user.IsActive, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan admin user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// ToggleAdminUserStatus toggles the active status of an admin user
func (as *AdminAuthService) ToggleAdminUserStatus(id uuid.UUID) error {
	query := `UPDATE admin_users SET is_active = NOT is_active WHERE id = $1`
	_, err := as.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to toggle admin user status: %w", err)
	}
	return nil
}
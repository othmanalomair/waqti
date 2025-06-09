package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

var Instance *DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfig() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "waqti"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func Connect() error {
	config := NewConfig()

	db, err := sql.Open("postgres", config.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	Instance = &DB{db}
	log.Println("Successfully connected to PostgreSQL database")
	return nil
}

func Close() error {
	if Instance != nil {
		return Instance.Close()
	}
	return nil
}

// Creator represents a creator in the database
type Creator struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	NameAr        string    `json:"name_ar"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	Avatar        *string   `json:"avatar"`
	Plan          string    `json:"plan"`
	PlanAr        string    `json:"plan_ar"`
	IsActive      bool      `json:"is_active"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Session represents a user session
type Session struct {
	ID           uuid.UUID `json:"id"`
	CreatorID    uuid.UUID `json:"creator_id"`
	SessionToken string    `json:"session_token"`
	DeviceInfo   *string   `json:"device_info"`
	IPAddress    *string   `json:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at"`
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetCreatorByEmail retrieves a creator by email
func (db *DB) GetCreatorByEmail(email string) (*Creator, error) {
	query := `
		SELECT id, name, name_ar, username, email, password_hash, avatar,
			   plan, plan_ar, is_active, email_verified, created_at, updated_at
		FROM creators
		WHERE email = $1 AND is_active = true
	`

	creator := &Creator{}
	err := db.QueryRow(query, email).Scan(
		&creator.ID, &creator.Name, &creator.NameAr, &creator.Username,
		&creator.Email, &creator.PasswordHash, &creator.Avatar,
		&creator.Plan, &creator.PlanAr, &creator.IsActive, &creator.EmailVerified,
		&creator.CreatedAt, &creator.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get creator by email: %w", err)
	}

	return creator, nil
}

// GetCreatorByUsername retrieves a creator by username
func (db *DB) GetCreatorByUsername(username string) (*Creator, error) {
	query := `
		SELECT id, name, name_ar, username, email, password_hash, avatar,
			   plan, plan_ar, is_active, email_verified, created_at, updated_at
		FROM creators
		WHERE username = $1 AND is_active = true
	`

	creator := &Creator{}
	err := db.QueryRow(query, username).Scan(
		&creator.ID, &creator.Name, &creator.NameAr, &creator.Username,
		&creator.Email, &creator.PasswordHash, &creator.Avatar,
		&creator.Plan, &creator.PlanAr, &creator.IsActive, &creator.EmailVerified,
		&creator.CreatedAt, &creator.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get creator by username: %w", err)
	}

	return creator, nil
}

// GetCreatorByID retrieves a creator by ID
func (db *DB) GetCreatorByID(id uuid.UUID) (*Creator, error) {
	query := `
		SELECT id, name, name_ar, username, email, password_hash, avatar,
			   plan, plan_ar, is_active, email_verified, created_at, updated_at
		FROM creators
		WHERE id = $1 AND is_active = true
	`

	creator := &Creator{}
	err := db.QueryRow(query, id).Scan(
		&creator.ID, &creator.Name, &creator.NameAr, &creator.Username,
		&creator.Email, &creator.PasswordHash, &creator.Avatar,
		&creator.Plan, &creator.PlanAr, &creator.IsActive, &creator.EmailVerified,
		&creator.CreatedAt, &creator.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get creator by ID: %w", err)
	}

	return creator, nil
}

// CreateCreator creates a new creator
func (db *DB) CreateCreator(creator *Creator) error {
	query := `
		INSERT INTO creators (name, name_ar, username, email, password_hash, plan, plan_ar)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(query,
		creator.Name, creator.NameAr, creator.Username,
		creator.Email, creator.PasswordHash, creator.Plan, creator.PlanAr,
	).Scan(&creator.ID, &creator.CreatedAt, &creator.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create creator: %w", err)
	}

	return nil
}

// CheckEmailExists checks if an email already exists
func (db *DB) CheckEmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM creators WHERE email = $1)`
	err := db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// CheckUsernameExists checks if a username already exists
func (db *DB) CheckUsernameExists(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM creators WHERE username = $1)`
	err := db.QueryRow(query, username).Scan(&exists)
	return exists, err
}

// CreateSession creates a new user session
func (db *DB) CreateSession(session *Session) error {
	query := `
		INSERT INTO creator_sessions (creator_id, session_token, device_info, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, last_activity
	`

	err := db.QueryRow(query,
		session.CreatorID, session.SessionToken, session.DeviceInfo,
		session.IPAddress, session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt, &session.LastActivity)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSessionByToken retrieves a session by token
func (db *DB) GetSessionByToken(token string) (*Session, error) {
	query := `
		SELECT id, creator_id, session_token, device_info, ip_address,
			   expires_at, last_activity, created_at
		FROM creator_sessions
		WHERE session_token = $1 AND expires_at > NOW()
	`

	session := &Session{}
	err := db.QueryRow(query, token).Scan(
		&session.ID, &session.CreatorID, &session.SessionToken,
		&session.DeviceInfo, &session.IPAddress, &session.ExpiresAt,
		&session.LastActivity, &session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// UpdateSessionActivity updates the last activity of a session
func (db *DB) UpdateSessionActivity(sessionID uuid.UUID) error {
	query := `UPDATE creator_sessions SET last_activity = NOW() WHERE id = $1`
	_, err := db.Exec(query, sessionID)
	return err
}

// DeleteSession deletes a session
func (db *DB) DeleteSession(token string) error {
	query := `DELETE FROM creator_sessions WHERE session_token = $1`
	_, err := db.Exec(query, token)
	return err
}

// CleanupExpiredSessions removes expired sessions
func (db *DB) CleanupExpiredSessions() error {
	query := `DELETE FROM creator_sessions WHERE expires_at <= NOW()`
	_, err := db.Exec(query)
	return err
}

// GetURLSettingsByCreatorID retrieves URL settings for a creator - FIXED
func (db *DB) GetURLSettingsByCreatorID(creatorID uuid.UUID) (*URLSettings, error) {
	query := `
		SELECT id, creator_id, username, changes_used, max_changes, last_changed, created_at, updated_at
		FROM url_settings
		WHERE creator_id = $1
	`

	urlSettings := &URLSettings{}
	err := db.QueryRow(query, creatorID).Scan(
		&urlSettings.ID,
		&urlSettings.CreatorID,
		&urlSettings.Username,
		&urlSettings.ChangesUsed,
		&urlSettings.MaxChanges,
		&urlSettings.LastChanged, // This will now handle NULL properly
		&urlSettings.CreatedAt,
		&urlSettings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Return nil instead of error when not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get URL settings: %w", err)
	}

	return urlSettings, nil
}

// UpdateCreatorUsername updates both creator and URL settings tables
func (db *DB) UpdateCreatorUsername(creatorID uuid.UUID, newUsername string) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Update creator's username
	creatorQuery := `UPDATE creators SET username = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err = tx.Exec(creatorQuery, newUsername, creatorID)
	if err != nil {
		return fmt.Errorf("failed to update creator username: %w", err)
	}

	// Update URL settings
	urlQuery := `
		UPDATE url_settings
		SET username = $1, changes_used = changes_used + 1, last_changed = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE creator_id = $2
	`
	_, err = tx.Exec(urlQuery, newUsername, creatorID)
	if err != nil {
		return fmt.Errorf("failed to update URL settings: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CheckUsernameExists checks if a username is already taken (excluding current user)
func (db *DB) CheckUsernameExistsExcluding(username string, excludeCreatorID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM creators WHERE LOWER(username) = LOWER($1) AND id != $2)`
	err := db.QueryRow(query, username, excludeCreatorID).Scan(&exists)
	return exists, err
}

// URLSettings represents URL settings in the database
type URLSettings struct {
	ID          uuid.UUID  `json:"id"`
	CreatorID   uuid.UUID  `json:"creator_id"`
	Username    string     `json:"username"`
	ChangesUsed int        `json:"changes_used"`
	MaxChanges  int        `json:"max_changes"`
	LastChanged *time.Time `json:"last_changed"` // Changed to pointer to handle NULL
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ShopSettings represents shop settings in the database
type ShopSettings struct {
	ID                 uuid.UUID `json:"id"`
	CreatorID          uuid.UUID `json:"creator_id"`
	LogoURL            *string   `json:"logo_url"`
	CreatorName        *string   `json:"creator_name"`
	CreatorNameAr      *string   `json:"creator_name_ar"`
	SubHeader          *string   `json:"sub_header"`
	SubHeaderAr        *string   `json:"sub_header_ar"`
	EnrollmentWhatsApp *string   `json:"enrollment_whatsapp"`
	ContactWhatsApp    *string   `json:"contact_whatsapp"`
	CheckoutLanguage   string    `json:"checkout_language"`
	GreetingMessage    *string   `json:"greeting_message"`
	GreetingMessageAr  *string   `json:"greeting_message_ar"`
	CurrencySymbol     string    `json:"currency_symbol"`
	CurrencySymbolAr   string    `json:"currency_symbol_ar"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// GetShopSettingsByCreatorID retrieves shop settings for a creator
func (db *DB) GetShopSettingsByCreatorID(creatorID uuid.UUID) (*ShopSettings, error) {
	query := `
		SELECT id, creator_id, logo_url, creator_name, creator_name_ar, sub_header, sub_header_ar,
		       enrollment_whatsapp, contact_whatsapp, checkout_language, greeting_message, greeting_message_ar,
		       currency_symbol, currency_symbol_ar, created_at, updated_at
		FROM shop_settings
		WHERE creator_id = $1
	`

	settings := &ShopSettings{}
	err := db.QueryRow(query, creatorID).Scan(
		&settings.ID,
		&settings.CreatorID,
		&settings.LogoURL,
		&settings.CreatorName,
		&settings.CreatorNameAr,
		&settings.SubHeader,
		&settings.SubHeaderAr,
		&settings.EnrollmentWhatsApp,
		&settings.ContactWhatsApp,
		&settings.CheckoutLanguage,
		&settings.GreetingMessage,
		&settings.GreetingMessageAr,
		&settings.CurrencySymbol,
		&settings.CurrencySymbolAr,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get shop settings: %w", err)
	}

	return settings, nil
}

// CreateShopSettings creates initial shop settings for a creator
func (db *DB) CreateShopSettings(settings *ShopSettings) error {
	query := `
		INSERT INTO shop_settings (creator_id, logo_url, creator_name, creator_name_ar, sub_header, sub_header_ar,
		                          enrollment_whatsapp, contact_whatsapp, checkout_language, greeting_message, greeting_message_ar,
		                          currency_symbol, currency_symbol_ar)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(query,
		settings.CreatorID,
		settings.LogoURL,
		settings.CreatorName,
		settings.CreatorNameAr,
		settings.SubHeader,
		settings.SubHeaderAr,
		settings.EnrollmentWhatsApp,
		settings.ContactWhatsApp,
		settings.CheckoutLanguage,
		settings.GreetingMessage,
		settings.GreetingMessageAr,
		settings.CurrencySymbol,
		settings.CurrencySymbolAr,
	).Scan(&settings.ID, &settings.CreatedAt, &settings.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create shop settings: %w", err)
	}

	return nil
}

// UpdateShopSettings updates existing shop settings
func (db *DB) UpdateShopSettings(settings *ShopSettings) error {
	query := `
		UPDATE shop_settings 
		SET logo_url = $2, creator_name = $3, creator_name_ar = $4, sub_header = $5, sub_header_ar = $6,
		    enrollment_whatsapp = $7, contact_whatsapp = $8, checkout_language = $9, greeting_message = $10, 
		    greeting_message_ar = $11, currency_symbol = $12, currency_symbol_ar = $13, updated_at = CURRENT_TIMESTAMP
		WHERE creator_id = $1
		RETURNING updated_at
	`

	err := db.QueryRow(query,
		settings.CreatorID,
		settings.LogoURL,
		settings.CreatorName,
		settings.CreatorNameAr,
		settings.SubHeader,
		settings.SubHeaderAr,
		settings.EnrollmentWhatsApp,
		settings.ContactWhatsApp,
		settings.CheckoutLanguage,
		settings.GreetingMessage,
		settings.GreetingMessageAr,
		settings.CurrencySymbol,
		settings.CurrencySymbolAr,
	).Scan(&settings.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update shop settings: %w", err)
	}

	return nil
}

// AnalyticsClick represents an analytics click in the database
type AnalyticsClick struct {
	ID          uuid.UUID `json:"id"`
	CreatorID   uuid.UUID `json:"creator_id"`
	IPAddress   *string   `json:"ip_address"`
	UserAgent   *string   `json:"user_agent"`
	Referrer    *string   `json:"referrer"`
	Country     string    `json:"country"`
	CountryAr   string    `json:"country_ar"`
	City        *string   `json:"city"`
	CityAr      *string   `json:"city_ar"`
	Device      string    `json:"device"`
	DeviceAr    string    `json:"device_ar"`
	OS          string    `json:"os"`
	OSAr        string    `json:"os_ar"`
	Browser     *string   `json:"browser"`
	BrowserAr   *string   `json:"browser_ar"`
	Platform    string    `json:"platform"`
	PlatformAr  string    `json:"platform_ar"`
	ClickedAt   time.Time `json:"clicked_at"`
}

// CreateAnalyticsClick creates a new analytics click record
func (db *DB) CreateAnalyticsClick(click *AnalyticsClick) error {
	query := `
		INSERT INTO analytics_clicks (
			creator_id, ip_address, user_agent, referrer, country, country_ar, 
			city, city_ar, device, device_ar, os, os_ar, browser, browser_ar, 
			platform, platform_ar, clicked_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id
	`

	err := db.QueryRow(query,
		click.CreatorID, click.IPAddress, click.UserAgent, click.Referrer,
		click.Country, click.CountryAr, click.City, click.CityAr,
		click.Device, click.DeviceAr, click.OS, click.OSAr,
		click.Browser, click.BrowserAr, click.Platform, click.PlatformAr,
		click.ClickedAt,
	).Scan(&click.ID)

	if err != nil {
		return fmt.Errorf("failed to create analytics click: %w", err)
	}

	return nil
}

// GetAnalyticsClicksByCreatorID retrieves analytics clicks for a creator
func (db *DB) GetAnalyticsClicksByCreatorID(creatorID uuid.UUID, limit int) ([]AnalyticsClick, error) {
	query := `
		SELECT id, creator_id, ip_address, user_agent, referrer, country, country_ar,
		       city, city_ar, device, device_ar, os, os_ar, browser, browser_ar,
		       platform, platform_ar, clicked_at
		FROM analytics_clicks
		WHERE creator_id = $1
		ORDER BY clicked_at DESC
		LIMIT $2
	`

	rows, err := db.Query(query, creatorID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics clicks: %w", err)
	}
	defer rows.Close()

	var clicks []AnalyticsClick
	for rows.Next() {
		var click AnalyticsClick
		err := rows.Scan(
			&click.ID, &click.CreatorID, &click.IPAddress, &click.UserAgent,
			&click.Referrer, &click.Country, &click.CountryAr, &click.City,
			&click.CityAr, &click.Device, &click.DeviceAr, &click.OS, &click.OSAr,
			&click.Browser, &click.BrowserAr, &click.Platform, &click.PlatformAr,
			&click.ClickedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan analytics click: %w", err)
		}
		clicks = append(clicks, click)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate analytics clicks: %w", err)
	}

	return clicks, nil
}

// GetAnalyticsStats gets analytics statistics for a creator
func (db *DB) GetAnalyticsStats(creatorID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM analytics_clicks WHERE creator_id = $1`
	
	var totalClicks int
	err := db.QueryRow(query, creatorID).Scan(&totalClicks)
	if err != nil {
		return 0, fmt.Errorf("failed to get analytics stats: %w", err)
	}

	return totalClicks, nil
}

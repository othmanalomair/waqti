package services

import (
	"fmt"
	"regexp"
	"strings"
	"waqti/internal/database"
	"waqti/internal/models"

	"github.com/google/uuid"
)

type URLService struct {
	// In real app, this would store reserved/taken usernames
	reservedUsernames []string
}

func NewURLService() *URLService {
	reservedUsernames := []string{
		"admin", "api", "www", "mail", "support", "help", "blog", "news",
		"shop", "store", "app", "mobile", "web", "dashboard", "settings",
		"profile", "user", "account", "login", "register", "signup",
		"waqti", "team", "about", "contact", "privacy", "terms",
	}

	return &URLService{
		reservedUsernames: reservedUsernames,
	}
}

func (s *URLService) GetURLSettingsByCreatorID(creatorID uuid.UUID) (*models.URLSettings, error) {
	// Get URL settings from database
	fmt.Printf("URLService.GetURLSettingsByCreatorID: Looking for creator %s\n", creatorID)

	dbURLSettings, err := database.Instance.GetURLSettingsByCreatorID(creatorID)
	if err != nil {
		fmt.Printf("URLService.GetURLSettingsByCreatorID: Database error: %v\n", err)
		return nil, fmt.Errorf("failed to get URL settings: %w", err)
	}

	if dbURLSettings == nil {
		fmt.Printf("URLService.GetURLSettingsByCreatorID: No URL settings found for creator %s\n", creatorID)
		return nil, fmt.Errorf("no URL settings found for creator")
	}

	fmt.Printf("URLService.GetURLSettingsByCreatorID: Found settings - Username: %s, Changes: %d/%d\n",
		dbURLSettings.Username, dbURLSettings.ChangesUsed, dbURLSettings.MaxChanges)

	// Convert database type to models type
	urlSettings := &models.URLSettings{
		ID:          dbURLSettings.ID,
		CreatorID:   dbURLSettings.CreatorID,
		Username:    dbURLSettings.Username,
		ChangesUsed: dbURLSettings.ChangesUsed,
		MaxChanges:  dbURLSettings.MaxChanges,
		LastChanged: dbURLSettings.LastChanged, // Both are now *time.Time
		CreatedAt:   dbURLSettings.CreatedAt,
		UpdatedAt:   dbURLSettings.UpdatedAt,
	}

	return urlSettings, nil
}

func (s *URLService) ValidateUsername(username string) models.URLValidationResult {
	username = strings.ToLower(strings.TrimSpace(username))

	// Check if empty
	if username == "" {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Username cannot be empty",
			ErrorMessageAr: "اسم المستخدم لا يمكن أن يكون فارغ",
		}
	}

	// Check length (3-20 characters)
	if len(username) < 3 {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Username must be at least 3 characters",
			ErrorMessageAr: "اسم المستخدم يجب أن يكون 3 أحرف على الأقل",
		}
	}

	if len(username) > 20 {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Username must be 20 characters or less",
			ErrorMessageAr: "اسم المستخدم يجب أن يكون 20 حرف أو أقل",
		}
	}

	// Check format (alphanumeric and underscores/hyphens only)
	validFormat := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validFormat.MatchString(username) {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Username can only contain letters, numbers, underscores, and hyphens",
			ErrorMessageAr: "اسم المستخدم يمكن أن يحتوي على حروف وأرقام وشرطات فقط",
		}
	}

	// Check if starts with letter or number
	if !regexp.MustCompile(`^[a-zA-Z0-9]`).MatchString(username) {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Username must start with a letter or number",
			ErrorMessageAr: "اسم المستخدم يجب أن يبدأ بحرف أو رقم",
		}
	}

	// Check if reserved
	for _, reserved := range s.reservedUsernames {
		if username == reserved {
			return models.URLValidationResult{
				IsValid:        false,
				ErrorMessage:   "This username is not available",
				ErrorMessageAr: "اسم المستخدم هذا غير متاح",
			}
		}
	}

	// Check if already taken by using database
	exists, err := database.Instance.CheckUsernameExists(username)
	if err != nil {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Error checking username availability",
			ErrorMessageAr: "خطأ في التحقق من توفر اسم المستخدم",
		}
	}

	if exists {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "This username is already taken",
			ErrorMessageAr: "اسم المستخدم مأخوذ بالفعل",
		}
	}

	return models.URLValidationResult{
		IsValid: true,
	}
}

// ValidateUsernameForCreator validates username excluding current creator
func (s *URLService) ValidateUsernameForCreator(username string, creatorID uuid.UUID) models.URLValidationResult {
	username = strings.ToLower(strings.TrimSpace(username))

	// Check if empty
	if username == "" {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Username cannot be empty",
			ErrorMessageAr: "اسم المستخدم لا يمكن أن يكون فارغ",
		}
	}

	// Check length (3-20 characters)
	if len(username) < 3 {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Username must be at least 3 characters",
			ErrorMessageAr: "اسم المستخدم يجب أن يكون 3 أحرف على الأقل",
		}
	}

	if len(username) > 20 {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Username must be 20 characters or less",
			ErrorMessageAr: "اسم المستخدم يجب أن يكون 20 حرف أو أقل",
		}
	}

	// Check format (alphanumeric and underscores/hyphens only)
	validFormat := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validFormat.MatchString(username) {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Username can only contain letters, numbers, underscores, and hyphens",
			ErrorMessageAr: "اسم المستخدم يمكن أن يحتوي على حروف وأرقام وشرطات فقط",
		}
	}

	// Check if starts with letter or number
	if !regexp.MustCompile(`^[a-zA-Z0-9]`).MatchString(username) {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Username must start with a letter or number",
			ErrorMessageAr: "اسم المستخدم يجب أن يبدأ بحرف أو رقم",
		}
	}

	// Check if reserved
	for _, reserved := range s.reservedUsernames {
		if username == reserved {
			return models.URLValidationResult{
				IsValid:        false,
				ErrorMessage:   "This username is not available",
				ErrorMessageAr: "اسم المستخدم هذا غير متاح",
			}
		}
	}

	// Check if already taken by another creator (excluding current creator)
	exists, err := database.Instance.CheckUsernameExistsExcluding(username, creatorID)
	if err != nil {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "Error checking username availability",
			ErrorMessageAr: "خطأ في التحقق من توفر اسم المستخدم",
		}
	}

	if exists {
		return models.URLValidationResult{
			IsValid:        false,
			ErrorMessage:   "This username is already taken",
			ErrorMessageAr: "اسم المستخدم مأخوذ بالفعل",
		}
	}

	return models.URLValidationResult{
		IsValid: true,
	}
}

func (s *URLService) UpdateUsername(creatorID uuid.UUID, newUsername string) error {
	// Get current URL settings
	urlSettings, err := s.GetURLSettingsByCreatorID(creatorID)
	if err != nil {
		return fmt.Errorf("failed to get URL settings: %w", err)
	}

	// Check if user has remaining changes
	if urlSettings.ChangesUsed >= urlSettings.MaxChanges {
		return fmt.Errorf("maximum number of changes reached")
	}

	// Validate new username (excluding current creator)
	validation := s.ValidateUsernameForCreator(newUsername, creatorID)
	if !validation.IsValid {
		return fmt.Errorf(validation.ErrorMessage)
	}

	newUsername = strings.ToLower(strings.TrimSpace(newUsername))

	// Use database method to update both tables in a transaction
	err = database.Instance.UpdateCreatorUsername(creatorID, newUsername)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}

func (s *URLService) GetRemainingChanges(creatorID uuid.UUID) int {
	urlSettings, err := s.GetURLSettingsByCreatorID(creatorID)
	if err != nil {
		return 0
	}
	return urlSettings.MaxChanges - urlSettings.ChangesUsed
}

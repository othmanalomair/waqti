package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"waqti/internal/models"
)

type URLService struct {
	urlSettings models.URLSettings
	// In real app, this would store reserved/taken usernames
	reservedUsernames []string
}

func NewURLService() *URLService {
	// Dummy URL settings
	urlSettings := models.URLSettings{
		ID:          1,
		CreatorID:   1,
		Username:    "ahmed",
		ChangesUsed: 2,
		MaxChanges:  5,
		LastChanged: time.Now().AddDate(0, -1, 0),
		CreatedAt:   time.Now().AddDate(0, -6, 0),
		UpdatedAt:   time.Now().AddDate(0, -1, 0),
	}

	reservedUsernames := []string{
		"admin", "api", "www", "mail", "support", "help", "blog", "news",
		"shop", "store", "app", "mobile", "web", "dashboard", "settings",
		"profile", "user", "account", "login", "register", "signup",
		"waqti", "team", "about", "contact", "privacy", "terms",
	}

	return &URLService{
		urlSettings:       urlSettings,
		reservedUsernames: reservedUsernames,
	}
}

func (s *URLService) GetURLSettingsByCreatorID(creatorID int) (*models.URLSettings, error) {
	return &s.urlSettings, nil
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

	// Check if already taken (simulate database check)
	if username == "john" || username == "sara" || username == "mohammed" {
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

func (s *URLService) UpdateUsername(creatorID int, newUsername string) error {
	// Check if user has remaining changes
	if s.urlSettings.ChangesUsed >= s.urlSettings.MaxChanges {
		return fmt.Errorf("maximum number of changes reached")
	}

	// Validate new username
	validation := s.ValidateUsername(newUsername)
	if !validation.IsValid {
		return fmt.Errorf(validation.ErrorMessage)
	}

	// Update settings
	s.urlSettings.Username = strings.ToLower(strings.TrimSpace(newUsername))
	s.urlSettings.ChangesUsed++
	s.urlSettings.LastChanged = time.Now()
	s.urlSettings.UpdatedAt = time.Now()

	return nil
}

func (s *URLService) GetRemainingChanges(creatorID int) int {
	return s.urlSettings.MaxChanges - s.urlSettings.ChangesUsed
}

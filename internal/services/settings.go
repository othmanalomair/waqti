package services

import (
	"time"
	"waqti/internal/models"
)

type SettingsService struct {
	settings models.ShopSettings
}

func NewSettingsService() *SettingsService {
	// Dummy settings data
	settings := models.ShopSettings{
		ID:                 1,
		CreatorID:          1,
		LogoURL:            "/static/images/default.jpg",
		CreatorName:        "Ahmed Al-Kuwaiti",
		CreatorNameAr:      "أحمد الكويتي",
		SubHeader:          "Certified Design Trainer",
		SubHeaderAr:        "مدرب معتمد في التصميم",
		EnrollmentWhatsApp: "+965-9999-8888",
		ContactWhatsApp:    "+965-9999-7777",
		CheckoutLanguage:   "both",
		GreetingMessage:    "Welcome to my workshop! Ready to learn?",
		GreetingMessageAr:  "مرحباً بك في ورشتي! هل أنت مستعد للتعلم؟",
		CurrencySymbol:     "KD",
		CurrencySymbolAr:   "د.ك",
		CreatedAt:          time.Now().AddDate(0, -1, 0),
		UpdatedAt:          time.Now(),
	}

	return &SettingsService{
		settings: settings,
	}
}

func (s *SettingsService) GetSettingsByCreatorID(creatorID int) (*models.ShopSettings, error) {
	// In real implementation, query database by creator ID
	return &s.settings, nil
}

func (s *SettingsService) UpdateSettings(creatorID int, request models.SettingsUpdateRequest) error {
	// Update settings (in real implementation, update database)
	s.settings.CreatorName = request.CreatorName
	s.settings.CreatorNameAr = request.CreatorNameAr
	s.settings.SubHeader = request.SubHeader
	s.settings.SubHeaderAr = request.SubHeaderAr
	s.settings.EnrollmentWhatsApp = request.EnrollmentWhatsApp
	s.settings.ContactWhatsApp = request.ContactWhatsApp
	s.settings.CheckoutLanguage = request.CheckoutLanguage
	s.settings.GreetingMessage = request.GreetingMessage
	s.settings.GreetingMessageAr = request.GreetingMessageAr
	s.settings.CurrencySymbol = request.CurrencySymbol
	s.settings.CurrencySymbolAr = request.CurrencySymbolAr
	s.settings.UpdatedAt = time.Now()

	return nil
}

func (s *SettingsService) UpdateLogo(creatorID int, logoURL string) error {
	s.settings.LogoURL = logoURL
	s.settings.UpdatedAt = time.Now()
	return nil
}

package models

import "time"

type ShopSettings struct {
	ID        int `json:"id"`
	CreatorID int `json:"creator_id"`
	// Branding Preferences
	LogoURL            string `json:"logo_url"`
	CreatorName        string `json:"creator_name"`
	CreatorNameAr      string `json:"creator_name_ar"`
	SubHeader          string `json:"sub_header"`
	SubHeaderAr        string `json:"sub_header_ar"`
	EnrollmentWhatsApp string `json:"enrollment_whatsapp"`
	ContactWhatsApp    string `json:"contact_whatsapp"`
	// Checkout Preferences
	CheckoutLanguage  string `json:"checkout_language"` // "ar", "en", "both"
	GreetingMessage   string `json:"greeting_message"`
	GreetingMessageAr string `json:"greeting_message_ar"`
	CurrencySymbol    string `json:"currency_symbol"`
	CurrencySymbolAr  string `json:"currency_symbol_ar"`
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SettingsUpdateRequest struct {
	// Branding
	CreatorName        string `form:"creator_name"`
	CreatorNameAr      string `form:"creator_name_ar"`
	SubHeader          string `form:"sub_header"`
	SubHeaderAr        string `form:"sub_header_ar"`
	EnrollmentWhatsApp string `form:"enrollment_whatsapp"`
	ContactWhatsApp    string `form:"contact_whatsapp"`
	// Checkout
	CheckoutLanguage  string `form:"checkout_language"`
	GreetingMessage   string `form:"greeting_message"`
	GreetingMessageAr string `form:"greeting_message_ar"`
	CurrencySymbol    string `form:"currency_symbol"`
	CurrencySymbolAr  string `form:"currency_symbol_ar"`
}

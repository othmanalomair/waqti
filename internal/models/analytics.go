package models

import (
	"time"

	"github.com/google/uuid"
)

type AnalyticsClick struct {
	ID         uuid.UUID `json:"id"`
	CreatorID  uuid.UUID `json:"creator_id"`
	Country    string    `json:"country"`
	CountryAr  string    `json:"country_ar"`
	Device     string    `json:"device"`
	DeviceAr   string    `json:"device_ar"`
	OS         string    `json:"os"`
	OSAr       string    `json:"os_ar"`
	Browser    string    `json:"browser"`
	BrowserAr  string    `json:"browser_ar"`
	Platform   string    `json:"platform"`
	PlatformAr string    `json:"platform_ar"`
	ClickedAt  time.Time `json:"clicked_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type AnalyticsStats struct {
	TotalClicks       int            `json:"total_clicks"`
	DateRange         string         `json:"date_range"`
	CountryBreakdown  map[string]int `json:"country_breakdown"`
	DeviceBreakdown   map[string]int `json:"device_breakdown"`
	OSBreakdown       map[string]int `json:"os_breakdown"`
	PlatformBreakdown map[string]int `json:"platform_breakdown"`
}

type AnalyticsFilter struct {
	FilterType string `json:"filter_type"` // "all", "country", "device", "os", "platform"
	DateRange  string `json:"date_range"`  // "7days", "30days", "90days", "all"
}

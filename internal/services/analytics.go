package services

import (
	"time"
	"waqti/internal/models"

	"github.com/google/uuid"
)

type AnalyticsService struct {
	clicks []models.AnalyticsClick
}

func NewAnalyticsService() *AnalyticsService {
	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	clicks := []models.AnalyticsClick{
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440050"),
			CreatorID:  creatorID,
			Country:    "Kuwait",
			CountryAr:  "الكويت",
			Device:     "Mobile",
			DeviceAr:   "جوال",
			OS:         "iOS",
			OSAr:       "آي أو إس",
			Platform:   "Instagram",
			PlatformAr: "إنستغرام",
			ClickedAt:  time.Now().AddDate(0, 0, -1),
			CreatedAt:  time.Now().AddDate(0, 0, -1),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440051"),
			CreatorID:  creatorID,
			Country:    "Saudi Arabia",
			CountryAr:  "السعودية",
			Device:     "Desktop",
			DeviceAr:   "سطح المكتب",
			OS:         "Windows",
			OSAr:       "ويندوز",
			Platform:   "WhatsApp",
			PlatformAr: "واتساب",
			ClickedAt:  time.Now().AddDate(0, 0, -2),
			CreatedAt:  time.Now().AddDate(0, 0, -2),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440052"),
			CreatorID:  creatorID,
			Country:    "UAE",
			CountryAr:  "الإمارات",
			Device:     "Mobile",
			DeviceAr:   "جوال",
			OS:         "Android",
			OSAr:       "أندرويد",
			Platform:   "Snapchat",
			PlatformAr: "سناب شات",
			ClickedAt:  time.Now().AddDate(0, 0, -3),
			CreatedAt:  time.Now().AddDate(0, 0, -3),
		},
	}

	return &AnalyticsService{
		clicks: clicks,
	}
}

func (s *AnalyticsService) GetClicksByCreatorID(creatorID uuid.UUID, filter models.AnalyticsFilter) []models.AnalyticsClick {
	var filtered []models.AnalyticsClick
	for _, click := range s.clicks {
		if click.CreatorID == creatorID {
			filtered = append(filtered, click)
		}
	}
	return filtered
}

func (s *AnalyticsService) GetAnalyticsStats(creatorID uuid.UUID, filter models.AnalyticsFilter) models.AnalyticsStats {
	clicks := s.GetClicksByCreatorID(creatorID, filter)

	stats := models.AnalyticsStats{
		TotalClicks:       len(clicks),
		DateRange:         filter.DateRange,
		CountryBreakdown:  make(map[string]int),
		DeviceBreakdown:   make(map[string]int),
		OSBreakdown:       make(map[string]int),
		PlatformBreakdown: make(map[string]int),
	}

	for _, click := range clicks {
		stats.CountryBreakdown[click.Country]++
		stats.DeviceBreakdown[click.Device]++
		stats.OSBreakdown[click.OS]++
		stats.PlatformBreakdown[click.Platform]++
	}

	return stats
}

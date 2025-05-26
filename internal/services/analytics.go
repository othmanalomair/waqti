package services

import (
	"sort"
	"time"
	"waqti/internal/models"
)

type AnalyticsService struct {
	clicks []models.AnalyticsClick
}

func NewAnalyticsService() *AnalyticsService {
	// Dummy analytics data
	clicks := []models.AnalyticsClick{
		{
			ID:         1,
			CreatorID:  1,
			Country:    "Kuwait",
			CountryAr:  "الكويت",
			Device:     "Mobile",
			DeviceAr:   "جوال",
			OS:         "iOS",
			OSAr:       "آي أو إس",
			Platform:   "Instagram",
			PlatformAr: "إنستغرام",
			ClickedAt:  time.Now().Add(-2 * time.Hour),
			CreatedAt:  time.Now().Add(-2 * time.Hour),
		},
		{
			ID:         2,
			CreatorID:  1,
			Country:    "Saudi Arabia",
			CountryAr:  "السعودية",
			Device:     "Desktop",
			DeviceAr:   "سطح المكتب",
			OS:         "Windows",
			OSAr:       "ويندوز",
			Platform:   "Facebook",
			PlatformAr: "فيسبوك",
			ClickedAt:  time.Now().Add(-5 * time.Hour),
			CreatedAt:  time.Now().Add(-5 * time.Hour),
		},
		{
			ID:         3,
			CreatorID:  1,
			Country:    "UAE",
			CountryAr:  "الإمارات",
			Device:     "Mobile",
			DeviceAr:   "جوال",
			OS:         "Android",
			OSAr:       "أندرويد",
			Platform:   "TikTok",
			PlatformAr: "تيك توك",
			ClickedAt:  time.Now().Add(-8 * time.Hour),
			CreatedAt:  time.Now().Add(-8 * time.Hour),
		},
		{
			ID:         4,
			CreatorID:  1,
			Country:    "Kuwait",
			CountryAr:  "الكويت",
			Device:     "Tablet",
			DeviceAr:   "تابلت",
			OS:         "iPadOS",
			OSAr:       "آي باد أو إس",
			Platform:   "Snapchat",
			PlatformAr: "سناب شات",
			ClickedAt:  time.Now().Add(-12 * time.Hour),
			CreatedAt:  time.Now().Add(-12 * time.Hour),
		},
		{
			ID:         5,
			CreatorID:  1,
			Country:    "Bahrain",
			CountryAr:  "البحرين",
			Device:     "Mobile",
			DeviceAr:   "جوال",
			OS:         "iOS",
			OSAr:       "آي أو إس",
			Platform:   "Instagram",
			PlatformAr: "إنستغرام",
			ClickedAt:  time.Now().Add(-24 * time.Hour),
			CreatedAt:  time.Now().Add(-24 * time.Hour),
		},
		{
			ID:         6,
			CreatorID:  1,
			Country:    "Qatar",
			CountryAr:  "قطر",
			Device:     "Desktop",
			DeviceAr:   "سطح المكتب",
			OS:         "macOS",
			OSAr:       "ماك أو إس",
			Platform:   "Facebook",
			PlatformAr: "فيسبوك",
			ClickedAt:  time.Now().Add(-36 * time.Hour),
			CreatedAt:  time.Now().Add(-36 * time.Hour),
		},
		{
			ID:         7,
			CreatorID:  1,
			Country:    "Kuwait",
			CountryAr:  "الكويت",
			Device:     "Mobile",
			DeviceAr:   "جوال",
			OS:         "Android",
			OSAr:       "أندرويد",
			Platform:   "Instagram",
			PlatformAr: "إنستغرام",
			ClickedAt:  time.Now().Add(-48 * time.Hour),
			CreatedAt:  time.Now().Add(-48 * time.Hour),
		},
		{
			ID:         8,
			CreatorID:  1,
			Country:    "Oman",
			CountryAr:  "عمان",
			Device:     "Mobile",
			DeviceAr:   "جوال",
			OS:         "iOS",
			OSAr:       "آي أو إس",
			Platform:   "TikTok",
			PlatformAr: "تيك توك",
			ClickedAt:  time.Now().Add(-72 * time.Hour),
			CreatedAt:  time.Now().Add(-72 * time.Hour),
		},
		{
			ID:         9,
			CreatorID:  1,
			Country:    "Saudi Arabia",
			CountryAr:  "السعودية",
			Device:     "Desktop",
			DeviceAr:   "سطح المكتب",
			OS:         "Windows",
			OSAr:       "ويندوز",
			Platform:   "Snapchat",
			PlatformAr: "سناب شات",
			ClickedAt:  time.Now().Add(-96 * time.Hour),
			CreatedAt:  time.Now().Add(-96 * time.Hour),
		},
		{
			ID:         10,
			CreatorID:  1,
			Country:    "Kuwait",
			CountryAr:  "الكويت",
			Device:     "Mobile",
			DeviceAr:   "جوال",
			OS:         "iOS",
			OSAr:       "آي أو إس",
			Platform:   "Instagram",
			PlatformAr: "إنستغرام",
			ClickedAt:  time.Now().Add(-120 * time.Hour),
			CreatedAt:  time.Now().Add(-120 * time.Hour),
		},
		{
			ID:         11,
			CreatorID:  1,
			Country:    "UAE",
			CountryAr:  "الإمارات",
			Device:     "Tablet",
			DeviceAr:   "تابلت",
			OS:         "Android",
			OSAr:       "أندرويد",
			Platform:   "Facebook",
			PlatformAr: "فيسبوك",
			ClickedAt:  time.Now().Add(-144 * time.Hour),
			CreatedAt:  time.Now().Add(-144 * time.Hour),
		},
		{
			ID:         12,
			CreatorID:  1,
			Country:    "Kuwait",
			CountryAr:  "الكويت",
			Device:     "Mobile",
			DeviceAr:   "جوال",
			OS:         "iOS",
			OSAr:       "آي أو إس",
			Platform:   "Instagram",
			PlatformAr: "إنستغرام",
			ClickedAt:  time.Now().Add(-168 * time.Hour),
			CreatedAt:  time.Now().Add(-168 * time.Hour),
		},
		{
			ID:         13,
			CreatorID:  1,
			Country:    "Saudi Arabia",
			CountryAr:  "السعودية",
			Device:     "Desktop",
			DeviceAr:   "سطح المكتب",
			OS:         "macOS",
			OSAr:       "ماك أو إس",
			Platform:   "TikTok",
			PlatformAr: "تيك توك",
			ClickedAt:  time.Now().Add(-192 * time.Hour),
			CreatedAt:  time.Now().Add(-192 * time.Hour),
		},
	}

	return &AnalyticsService{
		clicks: clicks,
	}
}

func (s *AnalyticsService) GetClicksByCreatorID(creatorID int, filter models.AnalyticsFilter) []models.AnalyticsClick {
	// Filter by date range
	filteredClicks := s.filterByDateRange(s.clicks, filter.DateRange)

	// Sort by most recent first
	sort.Slice(filteredClicks, func(i, j int) bool {
		return filteredClicks[i].ClickedAt.After(filteredClicks[j].ClickedAt)
	})

	return filteredClicks
}

func (s *AnalyticsService) GetAnalyticsStats(creatorID int, filter models.AnalyticsFilter) models.AnalyticsStats {
	clicks := s.filterByDateRange(s.clicks, filter.DateRange)

	stats := models.AnalyticsStats{
		TotalClicks:       len(clicks),
		DateRange:         s.getDateRangeText(filter.DateRange),
		CountryBreakdown:  make(map[string]int),
		DeviceBreakdown:   make(map[string]int),
		OSBreakdown:       make(map[string]int),
		PlatformBreakdown: make(map[string]int),
	}

	// Calculate breakdowns
	for _, click := range clicks {
		stats.CountryBreakdown[click.Country]++
		stats.DeviceBreakdown[click.Device]++
		stats.OSBreakdown[click.OS]++
		stats.PlatformBreakdown[click.Platform]++
	}

	return stats
}

func (s *AnalyticsService) filterByDateRange(clicks []models.AnalyticsClick, dateRange string) []models.AnalyticsClick {
	now := time.Now()
	var cutoff time.Time

	switch dateRange {
	case "7days":
		cutoff = now.AddDate(0, 0, -7)
	case "30days":
		cutoff = now.AddDate(0, 0, -30)
	case "90days":
		cutoff = now.AddDate(0, 0, -90)
	case "all":
		cutoff = time.Time{} // Beginning of time
	default:
		cutoff = now.AddDate(0, 0, -30) // Default to 30 days
	}

	var filtered []models.AnalyticsClick
	for _, click := range clicks {
		if click.ClickedAt.After(cutoff) {
			filtered = append(filtered, click)
		}
	}

	return filtered
}

func (s *AnalyticsService) getDateRangeText(dateRange string) string {
	switch dateRange {
	case "7days":
		return "Last 7 Days"
	case "30days":
		return "Last 30 Days"
	case "90days":
		return "Last 90 Days"
	case "all":
		return "All Time"
	default:
		return "Last 30 Days"
	}
}

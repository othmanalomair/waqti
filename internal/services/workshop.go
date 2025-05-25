package services

import (
	"encoding/json"
	"time"
	"waqti/internal/models"
)

type WorkshopService struct {
	workshops []models.Workshop
}

func NewWorkshopService() *WorkshopService {
	// Dummy data - replace with database later
	workshops := []models.Workshop{
		{
			ID:            1,
			CreatorID:     1,
			Title:         "Photography Basics",
			TitleAr:       "أساسيات التصوير",
			Description:   "Learn the fundamentals of photography",
			DescriptionAr: "تعلم أساسيات التصوير الفوتوغرافي",
			Price:         25.0,
			Duration:      120,
			MaxStudents:   15,
			Category:      "Photography",
			CategoryAr:    "التصوير",
			IsActive:      true,
			CreatedAt:     time.Now().AddDate(0, 0, -10),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            2,
			CreatorID:     1,
			Title:         "Digital Marketing",
			TitleAr:       "التسويق الرقمي",
			Description:   "Master social media marketing strategies",
			DescriptionAr: "إتقن استراتيجيات التسويق عبر وسائل التواصل",
			Price:         35.0,
			Duration:      90,
			MaxStudents:   20,
			Category:      "Marketing",
			CategoryAr:    "التسويق",
			IsActive:      true,
			CreatedAt:     time.Now().AddDate(0, 0, -5),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            3,
			CreatorID:     1,
			Title:         "Arabic Calligraphy",
			TitleAr:       "الخط العربي",
			Description:   "Traditional Arabic calligraphy techniques",
			DescriptionAr: "تقنيات الخط العربي التقليدية",
			Price:         20.0,
			Duration:      150,
			MaxStudents:   10,
			Category:      "Art",
			CategoryAr:    "الفنون",
			IsActive:      false,
			CreatedAt:     time.Now().AddDate(0, 0, -15),
			UpdatedAt:     time.Now(),
		},
	}

	return &WorkshopService{
		workshops: workshops,
	}
}

func (s *WorkshopService) GetWorkshopsByCreatorID(creatorID int) []models.Workshop {
	var result []models.Workshop
	for _, workshop := range s.workshops {
		if workshop.CreatorID == creatorID {
			result = append(result, workshop)
		}
	}
	return result
}

func (s *WorkshopService) GetDashboardStats(creatorID int) models.DashboardStats {
	workshops := s.GetWorkshopsByCreatorID(creatorID)

	stats := models.DashboardStats{
		TotalWorkshops:   len(workshops),
		ActiveWorkshops:  0,
		TotalEnrollments: 45,     // Dummy data
		MonthlyRevenue:   1250.0, // Dummy data
	}

	for _, workshop := range workshops {
		if workshop.IsActive {
			stats.ActiveWorkshops++
		}
	}

	return stats
}

func (s *WorkshopService) ToJSON(workshops []models.Workshop) string {
	data, _ := json.Marshal(workshops)
	return string(data)
}

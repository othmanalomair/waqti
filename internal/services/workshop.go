package services

import (
	"encoding/json"
	"fmt"
	"time"
	"waqti/internal/models"

	"github.com/google/uuid"
)

type WorkshopService struct {
	workshops []models.Workshop
}

func NewWorkshopService() *WorkshopService {
	// Generate fixed UUIDs for demo data consistency
	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	workshop1ID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	workshop2ID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")
	workshop3ID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")
	workshop4ID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440004")

	workshops := []models.Workshop{
		{
			ID:            workshop1ID,
			CreatorID:     creatorID,
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
			ID:            workshop2ID,
			CreatorID:     creatorID,
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
			ID:            workshop3ID,
			CreatorID:     creatorID,
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
		{
			ID:            workshop4ID,
			CreatorID:     creatorID,
			Title:         "Business English",
			TitleAr:       "الإنجليزية التجارية",
			Description:   "Professional English for business communication",
			DescriptionAr: "الإنجليزية المهنية للتواصل التجاري",
			Price:         30.0,
			Duration:      60,
			MaxStudents:   12,
			Category:      "Language",
			CategoryAr:    "اللغات",
			IsActive:      true,
			CreatedAt:     time.Now().AddDate(0, 0, -3),
			UpdatedAt:     time.Now(),
		},
	}

	return &WorkshopService{
		workshops: workshops,
	}
}

func (s *WorkshopService) GetWorkshopsByCreatorID(creatorID uuid.UUID) []models.Workshop {
	var result []models.Workshop
	for _, workshop := range s.workshops {
		if workshop.CreatorID == creatorID {
			result = append(result, workshop)
		}
	}
	return result
}

func (s *WorkshopService) GetDashboardStats(creatorID uuid.UUID) models.DashboardStats {
	workshops := s.GetWorkshopsByCreatorID(creatorID)

	totalSeats := 0
	projectedSales := 0.0

	for _, workshop := range workshops {
		if workshop.IsActive {
			totalSeats += workshop.MaxStudents
			projectedSales += workshop.Price * float64(workshop.MaxStudents) * 0.7
		}
	}

	stats := models.DashboardStats{
		TotalWorkshops:   len(workshops),
		ActiveWorkshops:  0,
		TotalEnrollments: 45,
		MonthlyRevenue:   1250.0,
		ProjectedSales:   projectedSales,
		RemainingSeats:   totalSeats - 23,
	}

	for _, workshop := range workshops {
		if workshop.IsActive {
			stats.ActiveWorkshops++
		}
	}

	return stats
}

func (s *WorkshopService) ReorderWorkshop(workshopID uuid.UUID, direction string) error {
	var workshopIndex = -1
	for i, workshop := range s.workshops {
		if workshop.ID == workshopID {
			workshopIndex = i
			break
		}
	}

	if workshopIndex == -1 {
		return fmt.Errorf("workshop not found")
	}

	if direction == "up" && workshopIndex > 0 {
		s.workshops[workshopIndex], s.workshops[workshopIndex-1] = s.workshops[workshopIndex-1], s.workshops[workshopIndex]
	} else if direction == "down" && workshopIndex < len(s.workshops)-1 {
		s.workshops[workshopIndex], s.workshops[workshopIndex+1] = s.workshops[workshopIndex+1], s.workshops[workshopIndex]
	}

	return nil
}

func (s *WorkshopService) ToggleWorkshopStatus(workshopID uuid.UUID) error {
	for i, workshop := range s.workshops {
		if workshop.ID == workshopID {
			s.workshops[i].IsActive = !workshop.IsActive
			s.workshops[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("workshop not found")
}

func (s *WorkshopService) ToJSON(workshops []models.Workshop) string {
	data, _ := json.Marshal(workshops)
	return string(data)
}

func (s *WorkshopService) GetWorkshops(creatorID uuid.UUID) []models.Workshop {
	return []models.Workshop{
		{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440010"),
			CreatorID:   creatorID,
			Name:        "Web Development Basics",
			Title:       "Web Development Basics",
			TitleAr:     "أساسيات تطوير الويب",
			Description: "Learn the fundamentals of web development",
			Price:       50.0,
			Currency:    "KWD",
			MaxStudents: 20,
			IsActive:    true,
			Status:      "published",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440011"),
			CreatorID:   creatorID,
			Name:        "JavaScript Advanced",
			Title:       "JavaScript Advanced",
			TitleAr:     "جافاسكريبت متقدم",
			Description: "Advanced JavaScript concepts and patterns",
			Price:       75.0,
			Currency:    "KWD",
			MaxStudents: 15,
			IsActive:    false,
			Status:      "draft",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

func (s *WorkshopService) CreateWorkshop(workshop *models.Workshop) error {
	return nil
}

func (s *WorkshopService) UpdateWorkshop(workshop *models.Workshop) error {
	return nil
}

func (s *WorkshopService) DeleteWorkshop(id uuid.UUID) error {
	return nil
}

func (s *WorkshopService) ReorderWorkshops(workshopID uuid.UUID, direction string) error {
	return nil
}

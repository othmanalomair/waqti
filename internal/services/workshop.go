package services

import (
	"encoding/json"
	"fmt"
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
		{
			ID:            4,
			CreatorID:     1,
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

	totalSeats := 0
	projectedSales := 0.0

	for _, workshop := range workshops {
		if workshop.IsActive {
			totalSeats += workshop.MaxStudents
			projectedSales += workshop.Price * float64(workshop.MaxStudents) * 0.7 // Assuming 70% booking rate
		}
	}

	stats := models.DashboardStats{
		TotalWorkshops:   len(workshops),
		ActiveWorkshops:  0,
		TotalEnrollments: 45,     // Dummy data
		MonthlyRevenue:   1250.0, // Dummy data
		ProjectedSales:   projectedSales,
		RemainingSeats:   totalSeats - 23, // Dummy enrolled count
	}

	for _, workshop := range workshops {
		if workshop.IsActive {
			stats.ActiveWorkshops++
		}
	}

	return stats
}

func (s *WorkshopService) ReorderWorkshop(workshopID int, direction string) error {
	// Find workshop index
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

	// Reorder logic
	if direction == "up" && workshopIndex > 0 {
		s.workshops[workshopIndex], s.workshops[workshopIndex-1] = s.workshops[workshopIndex-1], s.workshops[workshopIndex]
	} else if direction == "down" && workshopIndex < len(s.workshops)-1 {
		s.workshops[workshopIndex], s.workshops[workshopIndex+1] = s.workshops[workshopIndex+1], s.workshops[workshopIndex]
	}

	return nil
}

func (s *WorkshopService) ToggleWorkshopStatus(workshopID int) error {
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

// GetWorkshops returns workshops for a creator
func (s *WorkshopService) GetWorkshops(creatorID int) []models.Workshop {
	// Return dummy data for now
	return []models.Workshop{
		{
			ID:          1,
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
			ID:          2,
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

// CreateWorkshop creates a new workshop
func (s *WorkshopService) CreateWorkshop(workshop *models.Workshop) error {
	// In a real app, you'd save to database
	return nil
}

// UpdateWorkshop updates an existing workshop
func (s *WorkshopService) UpdateWorkshop(workshop *models.Workshop) error {
	// In a real app, you'd update in database
	return nil
}

// DeleteWorkshop deletes a workshop
func (s *WorkshopService) DeleteWorkshop(id int) error {
	// In a real app, you'd delete from database
	return nil
}

// ReorderWorkshops changes workshop order
func (s *WorkshopService) ReorderWorkshops(workshopID int, direction string) error {
	// In a real app, you'd update order in database
	return nil
}

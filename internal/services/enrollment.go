package services

import (
	"sort"
	"time"
	"waqti/internal/models"

	"github.com/google/uuid"
)

type EnrollmentService struct {
	enrollments []models.Enrollment
}

func NewEnrollmentService() *EnrollmentService {
	// Generate fixed UUIDs for demo data consistency
	workshop1ID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	workshop2ID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")
	workshop3ID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")
	workshop4ID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440004")

	enrollments := []models.Enrollment{
		{
			ID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440020"),
			WorkshopID:     workshop1ID,
			WorkshopName:   "Photography Basics",
			WorkshopNameAr: "أساسيات التصوير",
			StudentName:    "سارة أحمد",
			StudentEmail:   "sara@example.com",
			TotalPrice:     25.0,
			Status:         "successful",
			StatusAr:       "مكتمل",
			EnrollmentDate: time.Now().AddDate(0, 0, -2),
			CreatedAt:      time.Now().AddDate(0, 0, -2),
			UpdatedAt:      time.Now().AddDate(0, 0, -2),
		},
		{
			ID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440021"),
			WorkshopID:     workshop2ID,
			WorkshopName:   "Digital Marketing",
			WorkshopNameAr: "التسويق الرقمي",
			StudentName:    "محمد الكويتي",
			StudentEmail:   "mohammed@example.com",
			TotalPrice:     35.0,
			Status:         "successful",
			StatusAr:       "مكتمل",
			EnrollmentDate: time.Now().AddDate(0, 0, -5),
			CreatedAt:      time.Now().AddDate(0, 0, -5),
			UpdatedAt:      time.Now().AddDate(0, 0, -5),
		},
		{
			ID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440022"),
			WorkshopID:     workshop4ID,
			WorkshopName:   "Business English",
			WorkshopNameAr: "الإنجليزية التجارية",
			StudentName:    "فاطمة الزهراء",
			StudentEmail:   "fatima@example.com",
			TotalPrice:     30.0,
			Status:         "rejected",
			StatusAr:       "مرفوض",
			EnrollmentDate: time.Now().AddDate(0, 0, -1),
			CreatedAt:      time.Now().AddDate(0, 0, -1),
			UpdatedAt:      time.Now().AddDate(0, 0, -1),
		},
		{
			ID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440023"),
			WorkshopID:     workshop1ID,
			WorkshopName:   "Photography Basics",
			WorkshopNameAr: "أساسيات التصوير",
			StudentName:    "أحمد عبدالله",
			StudentEmail:   "ahmed@example.com",
			TotalPrice:     25.0,
			Status:         "successful",
			StatusAr:       "مكتمل",
			EnrollmentDate: time.Now().AddDate(0, 0, -7),
			CreatedAt:      time.Now().AddDate(0, 0, -7),
			UpdatedAt:      time.Now().AddDate(0, 0, -7),
		},
		{
			ID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440024"),
			WorkshopID:     workshop2ID,
			WorkshopName:   "Digital Marketing",
			WorkshopNameAr: "التسويق الرقمي",
			StudentName:    "نورا السالم",
			StudentEmail:   "nora@example.com",
			TotalPrice:     35.0,
			Status:         "pending",
			StatusAr:       "قيد المراجعة",
			EnrollmentDate: time.Now(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440025"),
			WorkshopID:     workshop3ID,
			WorkshopName:   "Arabic Calligraphy",
			WorkshopNameAr: "الخط العربي",
			StudentName:    "يوسف المطيري",
			StudentEmail:   "youssef@example.com",
			TotalPrice:     20.0,
			Status:         "successful",
			StatusAr:       "مكتمل",
			EnrollmentDate: time.Now().AddDate(0, 0, -10),
			CreatedAt:      time.Now().AddDate(0, 0, -10),
			UpdatedAt:      time.Now().AddDate(0, 0, -10),
		},
	}

	return &EnrollmentService{
		enrollments: enrollments,
	}
}

func (s *EnrollmentService) GetEnrollmentsByCreatorID(creatorID uuid.UUID, filter models.EnrollmentFilter) []models.Enrollment {
	enrollments := s.enrollments
	filteredEnrollments := s.filterByTimeRange(enrollments, filter.TimeRange)
	s.sortEnrollments(filteredEnrollments, filter.OrderBy, filter.OrderDir)
	return filteredEnrollments
}

func (s *EnrollmentService) GetEnrollmentStats(creatorID uuid.UUID, timeRange string) models.EnrollmentStats {
	enrollments := s.filterByTimeRange(s.enrollments, timeRange)

	stats := models.EnrollmentStats{}

	for _, enrollment := range enrollments {
		switch enrollment.Status {
		case "successful":
			stats.SuccessfulSales++
			stats.TotalSales += enrollment.TotalPrice
		case "rejected":
			stats.RejectedSales++
		case "pending":
			stats.PendingSales++
		}
	}

	return stats
}

func (s *EnrollmentService) filterByTimeRange(enrollments []models.Enrollment, timeRange string) []models.Enrollment {
	now := time.Now()
	var cutoff time.Time

	switch timeRange {
	case "days":
		cutoff = now.AddDate(0, 0, -30)
	case "months":
		cutoff = now.AddDate(0, -12, 0)
	case "year":
		cutoff = now.AddDate(-5, 0, 0)
	default:
		cutoff = now.AddDate(0, 0, -30)
	}

	var filtered []models.Enrollment
	for _, enrollment := range enrollments {
		if enrollment.EnrollmentDate.After(cutoff) {
			filtered = append(filtered, enrollment)
		}
	}

	return filtered
}

func (s *EnrollmentService) sortEnrollments(enrollments []models.Enrollment, orderBy, orderDir string) {
	sort.Slice(enrollments, func(i, j int) bool {
		var less bool

		switch orderBy {
		case "date":
			less = enrollments[i].EnrollmentDate.Before(enrollments[j].EnrollmentDate)
		case "price":
			less = enrollments[i].TotalPrice < enrollments[j].TotalPrice
		case "name":
			less = enrollments[i].WorkshopName < enrollments[j].WorkshopName
		default:
			less = enrollments[i].EnrollmentDate.Before(enrollments[j].EnrollmentDate)
		}

		if orderDir == "desc" {
			return !less
		}
		return less
	})
}

func (s *EnrollmentService) CreateEnrollment(enrollment *models.Enrollment) error {
	if enrollment.ID == uuid.Nil {
		enrollment.ID = uuid.New()
	}
	enrollment.CreatedAt = time.Now()
	enrollment.UpdatedAt = time.Now()
	
	// In a real implementation, this would insert into database
	// For now, add to in-memory slice for demo
	s.enrollments = append(s.enrollments, *enrollment)
	
	return nil
}

func (s *EnrollmentService) DeleteEnrollment(enrollmentID uuid.UUID) error {
	for i, enrollment := range s.enrollments {
		if enrollment.ID == enrollmentID {
			s.enrollments = append(s.enrollments[:i], s.enrollments[i+1:]...)
			return nil
		}
	}
	return nil
}

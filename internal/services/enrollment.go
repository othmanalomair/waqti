package services

import (
	"sort"
	"time"
	"waqti/internal/models"
)

type EnrollmentService struct {
	enrollments []models.Enrollment
}

func NewEnrollmentService() *EnrollmentService {
	// Dummy enrollment data
	enrollments := []models.Enrollment{
		{
			ID:             1,
			WorkshopID:     1,
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
			ID:             2,
			WorkshopID:     2,
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
			ID:             3,
			WorkshopID:     4,
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
			ID:             4,
			WorkshopID:     1,
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
			ID:             5,
			WorkshopID:     2,
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
			ID:             6,
			WorkshopID:     3,
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

func (s *EnrollmentService) GetEnrollmentsByCreatorID(creatorID int, filter models.EnrollmentFilter) []models.Enrollment {
	// For demo, return all enrollments (in real app, filter by creator)
	enrollments := s.enrollments

	// Apply time range filter
	filteredEnrollments := s.filterByTimeRange(enrollments, filter.TimeRange)

	// Apply sorting
	s.sortEnrollments(filteredEnrollments, filter.OrderBy, filter.OrderDir)

	return filteredEnrollments
}

func (s *EnrollmentService) GetEnrollmentStats(creatorID int, timeRange string) models.EnrollmentStats {
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
		cutoff = now.AddDate(0, 0, -30) // Last 30 days
	case "months":
		cutoff = now.AddDate(0, -12, 0) // Last 12 months
	case "year":
		cutoff = now.AddDate(-5, 0, 0) // Last 5 years
	default:
		cutoff = now.AddDate(0, 0, -30) // Default to 30 days
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

func (s *EnrollmentService) DeleteEnrollment(enrollmentID int) error {
	for i, enrollment := range s.enrollments {
		if enrollment.ID == enrollmentID {
			s.enrollments = append(s.enrollments[:i], s.enrollments[i+1:]...)
			return nil
		}
	}
	return nil
}

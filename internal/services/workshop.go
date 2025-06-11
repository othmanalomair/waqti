package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"waqti/internal/database"
	"waqti/internal/models"

	"github.com/google/uuid"
)

type WorkshopService struct {
	workshops []models.Workshop // Keep for fallback/demo data
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

// CreateWorkshop creates a new workshop in the database
func (s *WorkshopService) CreateWorkshop(workshop *models.Workshop) error {
	query := `
		INSERT INTO workshops (
			id, creator_id, name, title, title_ar, description, description_ar,
			price, currency, duration, max_students, status, is_active,
			is_free, is_recurring, recurrence_type, workshop_type, location_name,
			location_link, sort_order, deleted
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)
		RETURNING created_at, updated_at
	`

	// Get the next sort order
	sortOrder, err := s.getNextSortOrder(workshop.CreatorID)
	if err != nil {
		return fmt.Errorf("failed to get sort order: %w", err)
	}
	workshop.SortOrder = sortOrder

	err = database.Instance.QueryRow(
		query,
		workshop.ID,
		workshop.CreatorID,
		workshop.Name,
		workshop.Title,
		workshop.TitleAr,
		workshop.Description,
		workshop.DescriptionAr,
		workshop.Price,
		workshop.Currency,
		workshop.Duration,
		workshop.MaxStudents,
		workshop.Status,
		workshop.IsActive,
		workshop.IsFree,
		workshop.IsRecurring,
		workshop.RecurrenceType,
		workshop.WorkshopType,
		workshop.LocationName,
		workshop.LocationLink,
		workshop.SortOrder,
		false, // deleted = false for new workshops
	).Scan(&workshop.CreatedAt, &workshop.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create workshop: %w", err)
	}

	return nil
}

// CreateWorkshopSession creates a new workshop session in the database
func (s *WorkshopService) CreateWorkshopSession(session *models.WorkshopSession) error {
	query := `
		INSERT INTO workshop_sessions (
			id, workshop_id, session_date, end_date, start_time, end_time, duration, day_count,
			timezone, location, location_ar, max_attendees, current_attendees, is_completed,
			notes, notes_ar, status, status_ar, session_number, run_id, metadata,
			session_dates, total_days
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
		)
		RETURNING created_at, updated_at
	`

	// Handle nullable fields
	var endDate *time.Time
	if session.EndDate != nil {
		endDate = session.EndDate
	}
	
	var endTime *string
	if session.EndTime != nil {
		endTime = session.EndTime
	}

	var location *string
	if session.Location != nil && *session.Location != "" {
		location = session.Location
	}
	var locationAr *string
	if session.LocationAr != nil && *session.LocationAr != "" {
		locationAr = session.LocationAr
	}
	
	var notes *string
	if session.Notes != nil && *session.Notes != "" {
		notes = session.Notes
	}
	var notesAr *string
	if session.NotesAr != nil && *session.NotesAr != "" {
		notesAr = session.NotesAr
	}

	var runID *uuid.UUID
	if session.RunID != nil {
		runID = session.RunID
	}

	// Convert SessionDates slice to JSON for database storage
	sessionDatesJSON, err := json.Marshal(session.SessionDates)
	if err != nil {
		return fmt.Errorf("failed to marshal session dates: %w", err)
	}

	err = database.Instance.QueryRow(
		query,
		session.ID,
		session.WorkshopID,
		session.SessionDate,
		endDate,
		session.StartTime,
		endTime,
		session.Duration,
		session.DayCount,
		session.Timezone,
		location,
		locationAr,
		session.MaxAttendees,
		session.CurrentAttendees,
		session.IsCompleted,
		notes,
		notesAr,
		session.Status,
		session.StatusAr,
		session.SessionNumber,
		runID,
		nil, // metadata will be added later if needed
		sessionDatesJSON,
		session.TotalDays,
	).Scan(&session.CreatedAt, &session.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create workshop session: %w", err)
	}

	return nil
}

// GetWorkshopsByCreatorID retrieves workshops from database for a creator
func (s *WorkshopService) GetWorkshopsByCreatorID(creatorID uuid.UUID) []models.Workshop {
	query := `
		SELECT w.id, w.creator_id, w.name, w.title, w.title_ar, w.description, w.description_ar,
			   w.price, w.currency, w.duration, w.max_students,
			   COALESCE(c.name, '') as category, COALESCE(c.name_ar, '') as category_ar,
			   w.status, w.is_active, w.is_free, w.is_recurring, w.recurrence_type,
			   COALESCE(w.workshop_type, 'single') as workshop_type,
			   w.sort_order, w.view_count, w.enrollment_count, w.created_at, w.updated_at, w.deleted
		FROM workshops w
		LEFT JOIN categories c ON w.category_id = c.id
		WHERE w.creator_id = $1 AND w.deleted = false
		ORDER BY w.sort_order ASC, w.created_at DESC
	`

	rows, err := database.Instance.Query(query, creatorID)
	if err != nil {
		fmt.Printf("Error querying workshops: %v\n", err)
		// Return demo data as fallback
		return s.getWorkshopsByCreatorIDFallback(creatorID)
	}
	defer rows.Close()

	var workshops []models.Workshop
	for rows.Next() {
		var workshop models.Workshop
		err := rows.Scan(
			&workshop.ID,
			&workshop.CreatorID,
			&workshop.Name,
			&workshop.Title,
			&workshop.TitleAr,
			&workshop.Description,
			&workshop.DescriptionAr,
			&workshop.Price,
			&workshop.Currency,
			&workshop.Duration,
			&workshop.MaxStudents,
			&workshop.Category,
			&workshop.CategoryAr,
			&workshop.Status,
			&workshop.IsActive,
			&workshop.IsFree,
			&workshop.IsRecurring,
			&workshop.RecurrenceType,
			&workshop.WorkshopType,
			&workshop.SortOrder,
			&workshop.ViewCount,
			&workshop.EnrollmentCount,
			&workshop.CreatedAt,
			&workshop.UpdatedAt,
			&workshop.Deleted,
		)
		if err != nil {
			fmt.Printf("Error scanning workshop: %v\n", err)
			continue
		}
		workshops = append(workshops, workshop)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating workshops: %v\n", err)
		return s.getWorkshopsByCreatorIDFallback(creatorID)
	}

	// Enhance workshops with session information
	sessionService := NewWorkshopSessionService()
	for i := range workshops {
		sessions, err := sessionService.GetAvailableSessions(workshops[i].ID)
		if err == nil {
			workshops[i].Sessions = sessions
		}
		
		// Calculate enrollment count from sessions
		if len(sessions) > 0 {
			totalEnrollments := 0
			for _, session := range sessions {
				totalEnrollments += session.CurrentAttendees
			}
			workshops[i].EnrollmentCount = totalEnrollments
		}
	}

	return workshops
}

// GetActiveWorkshopsWithUpcomingSessions retrieves only workshops that have upcoming sessions
func (s *WorkshopService) GetActiveWorkshopsWithUpcomingSessions(creatorID uuid.UUID) []models.Workshop {
	query := `
		SELECT w.id, w.creator_id, w.name, w.title, w.title_ar, w.description, w.description_ar,
			   w.price, w.currency, w.duration, w.max_students,
			   COALESCE(c.name, '') as category, COALESCE(c.name_ar, '') as category_ar,
			   w.status, w.is_active, w.is_free, w.is_recurring, w.recurrence_type,
			   COALESCE(w.workshop_type, 'single') as workshop_type,
			   w.location_name, w.location_link,
			   w.sort_order, w.view_count, w.enrollment_count, w.created_at, w.updated_at, w.deleted
		FROM workshops w
		LEFT JOIN categories c ON w.category_id = c.id
		WHERE w.creator_id = $1 AND w.is_active = true AND w.deleted = false
		ORDER BY w.sort_order ASC, w.created_at DESC
	`

	rows, err := database.Instance.Query(query, creatorID)
	if err != nil {
		fmt.Printf("Error querying active workshops: %v\n", err)
		// Return demo data as fallback
		return s.getWorkshopsByCreatorIDFallback(creatorID)
	}
	defer rows.Close()

	var workshops []models.Workshop
	currentTime := time.Now()

	for rows.Next() {
		var workshop models.Workshop
		err := rows.Scan(
			&workshop.ID,
			&workshop.CreatorID,
			&workshop.Name,
			&workshop.Title,
			&workshop.TitleAr,
			&workshop.Description,
			&workshop.DescriptionAr,
			&workshop.Price,
			&workshop.Currency,
			&workshop.Duration,
			&workshop.MaxStudents,
			&workshop.Category,
			&workshop.CategoryAr,
			&workshop.Status,
			&workshop.IsActive,
			&workshop.IsFree,
			&workshop.IsRecurring,
			&workshop.RecurrenceType,
			&workshop.WorkshopType,
			&workshop.LocationName,
			&workshop.LocationLink,
			&workshop.SortOrder,
			&workshop.ViewCount,
			&workshop.EnrollmentCount,
			&workshop.CreatedAt,
			&workshop.UpdatedAt,
			&workshop.Deleted,
		)
		if err != nil {
			fmt.Printf("Error scanning workshop: %v\n", err)
			continue
		}

		// Load sessions for this workshop
		sessions, err := s.GetWorkshopSessions(workshop.ID)
		if err != nil {
			continue // Skip workshops without sessions or with errors
		}

		// Check if this is a private workshop (has special "always available" sessions)
		isPrivateWorkshop := workshop.WorkshopType == "private"
		
		// Filter to only upcoming sessions (considering both date and time)
		var upcomingSessions []models.WorkshopSession
		var sessionsNeedingTime []models.WorkshopSession
		var privateWorkshopSessions []models.WorkshopSession

		for _, session := range sessions {
			// Handle private workshop sessions with special date (9999-12-31)
			if isPrivateWorkshop && session.SessionDate.Year() == 9999 {
				// This is a private workshop session - always include it
				session.StartTime = "Available on demand" // Mark as available on demand
				privateWorkshopSessions = append(privateWorkshopSessions, session)
				continue
			}
			
			// Check if session has invalid/unset time
			if session.StartTime == "00:00:00" || session.StartTime == "" {
				// Keep sessions that need time to be set (for future dates)
				if !session.SessionDate.Before(currentTime.Truncate(24 * time.Hour)) {
					sessionsNeedingTime = append(sessionsNeedingTime, session)
				}
				continue
			}

			// Parse session date and start time to create a complete datetime
			sessionDateTime, err := time.Parse("2006-01-02 15:04:05",
				session.SessionDate.Format("2006-01-02")+" "+session.StartTime)
			if err != nil {
				// If parsing fails, fallback to date-only comparison
				if !session.SessionDate.Before(currentTime.Truncate(24 * time.Hour)) {
					upcomingSessions = append(upcomingSessions, session)
				}
			} else {
				// Use full datetime comparison
				if sessionDateTime.After(currentTime) {
					upcomingSessions = append(upcomingSessions, session)
				}
			}
		}

		// Add workshop if it has upcoming sessions OR sessions that need times OR private sessions
		if len(privateWorkshopSessions) > 0 {
			// Private workshop - always include with special sessions
			workshop.Sessions = privateWorkshopSessions
			workshops = append(workshops, workshop)
		} else if len(upcomingSessions) > 0 {
			workshop.Sessions = upcomingSessions
			workshops = append(workshops, workshop)
		} else if len(sessionsNeedingTime) > 0 {
			// For sessions needing time, add a special marker
			for i := range sessionsNeedingTime {
				sessionsNeedingTime[i].StartTime = "TBD" // Mark as To Be Determined
			}
			workshop.Sessions = sessionsNeedingTime
			workshops = append(workshops, workshop)
		}
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating workshops: %v\n", err)
		return s.getWorkshopsByCreatorIDFallback(creatorID)
	}

	return workshops
}

// Fallback method for demo data
func (s *WorkshopService) getWorkshopsByCreatorIDFallback(creatorID uuid.UUID) []models.Workshop {
	var result []models.Workshop
	for _, workshop := range s.workshops {
		if workshop.CreatorID == creatorID {
			result = append(result, workshop)
		}
	}
	return result
}

// UpdateWorkshop updates an existing workshop
func (s *WorkshopService) UpdateWorkshop(workshop *models.Workshop) error {
	query := `
    UPDATE workshops SET
        name = $2, title = $3, title_ar = $4, description = $5, description_ar = $6,
        price = $7, currency = $8, duration = $9, max_students = $10,
        status = $11, is_active = $12, is_free = $13, is_recurring = $14,
        recurrence_type = NULLIF($15, ''), workshop_type = $16, location_name = $17, 
        location_link = $18, updated_at = CURRENT_TIMESTAMP
    WHERE id = $1 AND creator_id = $19
    RETURNING updated_at
`

	err := database.Instance.QueryRow(
		query,
		workshop.ID,
		workshop.Name,
		workshop.Title,
		workshop.TitleAr,
		workshop.Description,
		workshop.DescriptionAr,
		workshop.Price,
		workshop.Currency,
		workshop.Duration,
		workshop.MaxStudents,
		workshop.Status,
		workshop.IsActive,
		workshop.IsFree,
		workshop.IsRecurring,
		workshop.RecurrenceType,
		workshop.WorkshopType,
		workshop.LocationName,
		workshop.LocationLink,
		workshop.CreatorID,
	).Scan(&workshop.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update workshop: %w", err)
	}

	return nil
}

// DeleteWorkshop soft deletes a workshop by setting deleted=true
func (s *WorkshopService) DeleteWorkshop(id uuid.UUID) error {
	// Instead of hard delete, set deleted=true (soft delete)
	query := `UPDATE workshops SET deleted = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := database.Instance.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete workshop: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workshop not found")
	}

	return nil
}

// ReorderWorkshop changes the sort order of a workshop
func (s *WorkshopService) ReorderWorkshop(workshopID uuid.UUID, direction string) error {
	// Get current workshop details
	var currentOrder int
	var creatorID uuid.UUID

	query := `SELECT sort_order, creator_id FROM workshops WHERE id = $1`
	err := database.Instance.QueryRow(query, workshopID).Scan(&currentOrder, &creatorID)
	if err != nil {
		return fmt.Errorf("failed to get workshop details: %w", err)
	}

	var newOrder int
	if direction == "up" {
		newOrder = currentOrder - 1
	} else if direction == "down" {
		newOrder = currentOrder + 1
	} else {
		return fmt.Errorf("invalid direction: %s", direction)
	}

	// Start transaction
	tx, err := database.Instance.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Find workshop at target position
	var targetWorkshopID uuid.UUID
	query = `SELECT id FROM workshops WHERE creator_id = $1 AND sort_order = $2 AND deleted = false`
	err = tx.QueryRow(query, creatorID, newOrder).Scan(&targetWorkshopID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to find target workshop: %w", err)
	}

	// Swap positions
	if err != sql.ErrNoRows {
		// Update target workshop to current position
		query = `UPDATE workshops SET sort_order = $1 WHERE id = $2`
		_, err = tx.Exec(query, currentOrder, targetWorkshopID)
		if err != nil {
			return fmt.Errorf("failed to update target workshop: %w", err)
		}
	}

	// Update current workshop to new position
	query = `UPDATE workshops SET sort_order = $1 WHERE id = $2`
	_, err = tx.Exec(query, newOrder, workshopID)
	if err != nil {
		return fmt.Errorf("failed to update workshop position: %w", err)
	}

	return tx.Commit()
}

// ToggleWorkshopStatus toggles the active status of a workshop
func (s *WorkshopService) ToggleWorkshopStatus(workshopID uuid.UUID) error {
	query := `
		UPDATE workshops
		SET is_active = NOT is_active, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := database.Instance.Exec(query, workshopID)
	if err != nil {
		return fmt.Errorf("failed to toggle workshop status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workshop not found")
	}

	return nil
}

// GetDashboardStats calculates dashboard statistics
func (s *WorkshopService) GetDashboardStats(creatorID uuid.UUID) models.DashboardStats {
	// Try to get stats from database first
	stats, err := s.getDashboardStatsFromDB(creatorID)
	if err != nil {
		fmt.Printf("Error getting dashboard stats from DB: %v\n", err)
		// Fall back to original implementation
		return s.getDashboardStatsFallback(creatorID)
	}
	return stats
}

func (s *WorkshopService) getDashboardStatsFromDB(creatorID uuid.UUID) (models.DashboardStats, error) {
	var stats models.DashboardStats

	// Get workshop counts
	query := `
		SELECT
			COUNT(*) as total_workshops,
			COUNT(CASE WHEN is_active = true THEN 1 END) as active_workshops,
			COALESCE(SUM(CASE WHEN is_active = true AND max_students > 0 THEN max_students ELSE 0 END), 0) as total_seats
		FROM workshops
		WHERE creator_id = $1 AND deleted = false
	`

	var totalSeats int
	err := database.Instance.QueryRow(query, creatorID).Scan(
		&stats.TotalWorkshops,
		&stats.ActiveWorkshops,
		&totalSeats,
	)
	if err != nil {
		return stats, fmt.Errorf("failed to get workshop stats: %w", err)
	}

	// Get enrollment count
	enrollmentQuery := `
		SELECT COUNT(e.id)
		FROM enrollments e
		JOIN workshops w ON e.workshop_id = w.id
		WHERE w.creator_id = $1
	`
	err = database.Instance.QueryRow(enrollmentQuery, creatorID).Scan(&stats.TotalEnrollments)
	if err != nil {
		stats.TotalEnrollments = 0 // Default if error
	}

	// Get monthly revenue (current month)
	revenueQuery := `
		SELECT COALESCE(SUM(e.total_price), 0)
		FROM enrollments e
		JOIN workshops w ON e.workshop_id = w.id
		WHERE w.creator_id = $1
		AND e.status = 'successful'
		AND e.enrollment_date >= DATE_TRUNC('month', CURRENT_DATE)
	`
	err = database.Instance.QueryRow(revenueQuery, creatorID).Scan(&stats.MonthlyRevenue)
	if err != nil {
		stats.MonthlyRevenue = 0 // Default if error
	}

	// Calculate projected sales (70% of total possible revenue)
	projectedQuery := `
		SELECT COALESCE(SUM(CASE WHEN w.is_active = true AND w.max_students > 0 THEN w.price * w.max_students ELSE 0 END), 0) * 0.7
		FROM workshops w
		WHERE w.creator_id = $1 AND w.deleted = false
	`
	err = database.Instance.QueryRow(projectedQuery, creatorID).Scan(&stats.ProjectedSales)
	if err != nil {
		stats.ProjectedSales = 0 // Default if error
	}

	// Calculate remaining seats
	stats.RemainingSeats = totalSeats - stats.TotalEnrollments
	if stats.RemainingSeats < 0 {
		stats.RemainingSeats = 0
	}

	return stats, nil
}

// Fallback implementation for demo data
func (s *WorkshopService) getDashboardStatsFallback(creatorID uuid.UUID) models.DashboardStats {
	workshops := s.getWorkshopsByCreatorIDFallback(creatorID)

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

// Helper function to get next sort order
func (s *WorkshopService) getNextSortOrder(creatorID uuid.UUID) (int, error) {
	var maxOrder sql.NullInt64

	query := `SELECT MAX(sort_order) FROM workshops WHERE creator_id = $1 AND deleted = false`
	err := database.Instance.QueryRow(query, creatorID).Scan(&maxOrder)
	if err != nil {
		return 0, fmt.Errorf("failed to get max sort order: %w", err)
	}

	if maxOrder.Valid {
		return int(maxOrder.Int64) + 1, nil
	}
	return 1, nil // First workshop
}

// Add this debug version to your workshop service temporarily
func (s *WorkshopService) GetWorkshopByID(workshopID uuid.UUID, creatorID uuid.UUID) (*models.Workshop, error) {
	fmt.Printf("DEBUG: GetWorkshopByID called with workshopID: %s, creatorID: %s\n", workshopID, creatorID)

	query := `
		SELECT w.id, w.creator_id, w.name, w.title, w.title_ar, w.description, w.description_ar,
			   w.price, w.currency, w.duration, w.max_students,
			   COALESCE(c.name, '') as category, COALESCE(c.name_ar, '') as category_ar,
			   w.status, w.is_active, w.is_free, w.is_recurring, w.recurrence_type,
			   COALESCE(w.workshop_type, 'single') as workshop_type,
			   w.location_name, w.location_link,
			   w.sort_order, w.view_count, w.enrollment_count, w.created_at, w.updated_at, w.deleted
		FROM workshops w
		LEFT JOIN categories c ON w.category_id = c.id
		WHERE w.id = $1 AND w.creator_id = $2 AND w.deleted = false
	`

	fmt.Printf("DEBUG: Executing query with parameters: workshopID=%s, creatorID=%s\n", workshopID, creatorID)

	var workshop models.Workshop
	err := database.Instance.QueryRow(query, workshopID, creatorID).Scan(
		&workshop.ID,
		&workshop.CreatorID,
		&workshop.Name,
		&workshop.Title,
		&workshop.TitleAr,
		&workshop.Description,
		&workshop.DescriptionAr,
		&workshop.Price,
		&workshop.Currency,
		&workshop.Duration,
		&workshop.MaxStudents,
		&workshop.Category,
		&workshop.CategoryAr,
		&workshop.Status,
		&workshop.IsActive,
		&workshop.IsFree,
		&workshop.IsRecurring,
		&workshop.RecurrenceType,
		&workshop.WorkshopType,
		&workshop.LocationName,
		&workshop.LocationLink,
		&workshop.SortOrder,
		&workshop.ViewCount,
		&workshop.EnrollmentCount,
		&workshop.CreatedAt,
		&workshop.UpdatedAt,
		&workshop.Deleted,
	)

	if err == sql.ErrNoRows {
		fmt.Printf("DEBUG: No workshop found with ID %s for creator %s\n", workshopID, creatorID)
		return nil, nil // Workshop not found
	}
	if err != nil {
		fmt.Printf("DEBUG: Database error: %v\n", err)
		return nil, fmt.Errorf("failed to get workshop by ID: %w", err)
	}

	fmt.Printf("DEBUG: Workshop found successfully - ID: %s, Name: %s, Price: %.2f\n",
		workshop.ID, workshop.Name, workshop.Price)

	return &workshop, nil
}

// GetWorkshopSessions retrieves all sessions for a workshop
func (s *WorkshopService) GetWorkshopSessions(workshopID uuid.UUID) ([]models.WorkshopSession, error) {
	query := `
		SELECT id, workshop_id, session_date, end_date, start_time::text, end_time, duration, day_count,
			   timezone, location, location_ar, max_attendees, current_attendees,
			   is_completed, notes, notes_ar, status, status_ar, session_number, run_id,
			   session_dates, total_days, created_at, updated_at
		FROM workshop_sessions
		WHERE workshop_id = $1
		ORDER BY session_date ASC, start_time ASC
	`

	rows, err := database.Instance.Query(query, workshopID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workshop sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.WorkshopSession
	for rows.Next() {
		var session models.WorkshopSession
		var endDate sql.NullTime
		var endTime sql.NullString
		var location sql.NullString
		var locationAr sql.NullString
		var notes sql.NullString
		var notesAr sql.NullString
		var runID sql.NullString
		var sessionDatesJSON sql.NullString

		err := rows.Scan(
			&session.ID,
			&session.WorkshopID,
			&session.SessionDate,
			&endDate,
			&session.StartTime,
			&endTime,
			&session.Duration,
			&session.DayCount,
			&session.Timezone,
			&location,
			&locationAr,
			&session.MaxAttendees,
			&session.CurrentAttendees,
			&session.IsCompleted,
			&notes,
			&notesAr,
			&session.Status,
			&session.StatusAr,
			&session.SessionNumber,
			&runID,
			&sessionDatesJSON,
			&session.TotalDays,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workshop session: %w", err)
		}

		// Handle nullable fields
		if endDate.Valid {
			session.EndDate = &endDate.Time
		}
		if endTime.Valid {
			session.EndTime = &endTime.String
		}
		if location.Valid {
			session.Location = &location.String
		}
		if locationAr.Valid {
			session.LocationAr = &locationAr.String
		}
		if notes.Valid {
			session.Notes = &notes.String
		}
		if notesAr.Valid {
			session.NotesAr = &notesAr.String
		}
		if runID.Valid {
			if parsedRunID, err := uuid.Parse(runID.String); err == nil {
				session.RunID = &parsedRunID
			}
		}

		// Parse session dates JSON
		if sessionDatesJSON.Valid && sessionDatesJSON.String != "" {
			var dates []time.Time
			if err := json.Unmarshal([]byte(sessionDatesJSON.String), &dates); err == nil {
				session.SessionDates = dates
			}
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workshop sessions: %w", err)
	}

	return sessions, nil
}

// GetWorkshopSessionByID retrieves a single workshop session by ID
func (s *WorkshopService) GetWorkshopSessionByID(sessionID uuid.UUID) (*models.WorkshopSession, error) {
	query := `
		SELECT id, workshop_id, session_date, start_time::text, end_time, duration,
			   timezone, location, location_ar, max_attendees, current_attendees,
			   is_completed, notes, notes_ar, session_dates, total_days, created_at, updated_at
		FROM workshop_sessions
		WHERE id = $1
	`

	var session models.WorkshopSession
	var endTime sql.NullString
	var location sql.NullString
	var locationAr sql.NullString
	var notes sql.NullString
	var notesAr sql.NullString
	var sessionDatesJSON sql.NullString

	err := database.Instance.QueryRow(query, sessionID).Scan(
		&session.ID,
		&session.WorkshopID,
		&session.SessionDate,
		&session.StartTime,
		&endTime,
		&session.Duration,
		&session.Timezone,
		&location,
		&locationAr,
		&session.MaxAttendees,
		&session.CurrentAttendees,
		&session.IsCompleted,
		&notes,
		&notesAr,
		&sessionDatesJSON,
		&session.TotalDays,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get workshop session by ID: %w", err)
	}

	// Handle nullable fields
	if endTime.Valid {
		session.EndTime = &endTime.String
	}
	if location.Valid {
		session.Location = &location.String
	}
	if locationAr.Valid {
		session.LocationAr = &locationAr.String
	}
	if notes.Valid {
		session.Notes = &notes.String
	}
	if notesAr.Valid {
		session.NotesAr = &notesAr.String
	}

	// Parse session dates JSON
	if sessionDatesJSON.Valid && sessionDatesJSON.String != "" {
		var dates []time.Time
		if err := json.Unmarshal([]byte(sessionDatesJSON.String), &dates); err == nil {
			session.SessionDates = dates
		}
	}

	return &session, nil
}

// DeleteWorkshopSessions deletes all sessions for a workshop
func (s *WorkshopService) DeleteWorkshopSessions(workshopID uuid.UUID) error {
	query := `DELETE FROM workshop_sessions WHERE workshop_id = $1`

	_, err := database.Instance.Exec(query, workshopID)
	if err != nil {
		return fmt.Errorf("failed to delete workshop sessions: %w", err)
	}

	return nil
}

// UpdateWorkshopSessionsSafely updates sessions while protecting those with orders
func (s *WorkshopService) UpdateWorkshopSessionsSafely(workshopID uuid.UUID, newSessions []models.WorkshopSession) error {
	// Get existing sessions for this workshop
	existingSessions, err := s.GetWorkshopSessions(workshopID)
	if err != nil {
		return fmt.Errorf("failed to get existing sessions: %w", err)
	}

	fmt.Printf("Safe update mode: Found %d existing sessions, %d new sessions\n", len(existingSessions), len(newSessions))

	// Track which existing sessions have been matched
	matchedSessions := make(map[uuid.UUID]bool)

	// Process new sessions - either update existing or create new
	for _, newSession := range newSessions {
		// Find the best matching existing session
		var bestMatch *models.WorkshopSession
		
		for i := range existingSessions {
			existing := &existingSessions[i]
			
			// Try to match by primary session date and start time
			if existing.SessionDate.Format("2006-01-02") == newSession.SessionDate.Format("2006-01-02") && 
			   existing.StartTime == newSession.StartTime {
				bestMatch = existing
				break
			}
		}
		
		if bestMatch != nil {
			// Mark this session as matched
			matchedSessions[bestMatch.ID] = true
			
			// Update only safe attributes that don't affect order relationships
			updateQuery := `
				UPDATE workshop_sessions 
				SET max_attendees = $1, 
				    duration = $2,
				    session_dates = $3,
				    total_days = $4,
				    updated_at = NOW()
				WHERE id = $5`
			
			// Convert SessionDates to JSON for database storage
			sessionDatesJSON, err := json.Marshal(newSession.SessionDates)
			if err != nil {
				fmt.Printf("Warning: Failed to marshal session dates for session %s: %v\n", bestMatch.ID.String(), err)
				sessionDatesJSON = []byte("[]")
			}
			
			_, err = database.Instance.Exec(updateQuery, 
				newSession.MaxAttendees, 
				newSession.Duration,
				sessionDatesJSON,
				newSession.TotalDays,
				bestMatch.ID)
			if err != nil {
				return fmt.Errorf("failed to update session %s: %w", bestMatch.ID.String(), err)
			}
			
			fmt.Printf("Updated session %s: max_attendees=%d, duration=%.1f, total_days=%d\n", 
				bestMatch.ID.String(), newSession.MaxAttendees, newSession.Duration, newSession.TotalDays)
		} else {
			// This is a new session that doesn't exist yet - create it
			fmt.Printf("Creating new session for date %s time %s in safe update mode\n", 
				newSession.SessionDate.Format("2006-01-02"), newSession.StartTime)
			
			// Set up the new session properly
			newSession.ID = uuid.New()
			newSession.WorkshopID = workshopID
			newSession.Status = "upcoming"
			newSession.StatusAr = "قادم"
			// Calculate next session number
			maxSessionNumber := 0
			for _, existing := range existingSessions {
				if existing.SessionNumber > maxSessionNumber {
					maxSessionNumber = existing.SessionNumber
				}
			}
			newSession.SessionNumber = maxSessionNumber + 1
			newSession.Timezone = "Asia/Kuwait"
			newSession.CurrentAttendees = 0
			newSession.IsCompleted = false
			newSession.CreatedAt = time.Now()
			newSession.UpdatedAt = time.Now()
			
			// Create a run_id for this session
			runID := uuid.New()
			newSession.RunID = &runID
			
			// Create the new session
			err := s.CreateWorkshopSession(&newSession)
			if err != nil {
				fmt.Printf("Warning: Failed to create new session: %v\n", err)
				continue
			}
			
			fmt.Printf("Successfully created new session %s for date %s time %s\n", 
				newSession.ID.String(), newSession.SessionDate.Format("2006-01-02"), newSession.StartTime)
		}
	}

	// Handle existing sessions that weren't matched (deleted by user)
	for _, existing := range existingSessions {
		if !matchedSessions[existing.ID] {
			// Check if this session has orders
			hasOrders, err := s.SessionHasOrders(existing.ID)
			if err != nil {
				fmt.Printf("Warning: Could not check orders for session %s: %v\n", existing.ID.String(), err)
				continue // Skip deletion if we can't verify
			}
			
			if hasOrders {
				fmt.Printf("Warning: Cannot delete session %s (date: %s, time: %s) because it has orders\n", 
					existing.ID.String(), existing.SessionDate.Format("2006-01-02"), existing.StartTime)
			} else {
				// Safe to delete this session
				err = s.DeleteWorkshopSession(existing.ID)
				if err != nil {
					fmt.Printf("Warning: Failed to delete session %s: %v\n", existing.ID.String(), err)
				} else {
					fmt.Printf("Deleted session %s (date: %s, time: %s) - no orders found\n", 
						existing.ID.String(), existing.SessionDate.Format("2006-01-02"), existing.StartTime)
				}
			}
		}
	}

	return nil
}

// DeleteWorkshopSession deletes a single workshop session
func (s *WorkshopService) DeleteWorkshopSession(sessionID uuid.UUID) error {
	query := `DELETE FROM workshop_sessions WHERE id = $1`

	result, err := database.Instance.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete workshop session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// GetWorkshopImages retrieves all images for a workshop
func (s *WorkshopService) GetWorkshopImages(workshopID uuid.UUID) ([]models.WorkshopImage, error) {
	query := `
		SELECT id, workshop_id, image_url, is_cover, sort_order, alt_text, alt_text_ar, created_at
		FROM workshop_images
		WHERE workshop_id = $1
		ORDER BY sort_order ASC, created_at ASC
	`

	rows, err := database.Instance.Query(query, workshopID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workshop images: %w", err)
	}
	defer rows.Close()

	var images []models.WorkshopImage
	for rows.Next() {
		var image models.WorkshopImage
		var altText sql.NullString
		var altTextAr sql.NullString

		err := rows.Scan(
			&image.ID,
			&image.WorkshopID,
			&image.ImageURL,
			&image.IsCover,
			&image.SortOrder,
			&altText,
			&altTextAr,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workshop image: %w", err)
		}

		// Handle nullable fields
		if altText.Valid {
			image.AltText = altText.String
		}
		if altTextAr.Valid {
			image.AltTextAr = altTextAr.String
		}

		images = append(images, image)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workshop images: %w", err)
	}

	return images, nil
}

// Legacy methods for compatibility
func (s *WorkshopService) ToJSON(workshops []models.Workshop) string {
	data, _ := json.Marshal(workshops)
	return string(data)
}

func (s *WorkshopService) GetWorkshops(creatorID uuid.UUID) []models.Workshop {
	return s.GetWorkshopsByCreatorID(creatorID)
}

// CreateWorkshopImage creates a new workshop image record
func (s *WorkshopService) CreateWorkshopImage(image *models.WorkshopImage) error {
	query := `
		INSERT INTO workshop_images (
			id, workshop_id, image_url, is_cover, sort_order, alt_text, alt_text_ar
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING created_at
	`

	err := database.Instance.QueryRow(
		query,
		image.ID,
		image.WorkshopID,
		image.ImageURL,
		image.IsCover,
		image.SortOrder,
		image.AltText,
		image.AltTextAr,
	).Scan(&image.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create workshop image: %w", err)
	}

	return nil
}

// GetWorkshopImagesByWorkshopID retrieves all images for a workshop
func (s *WorkshopService) GetWorkshopImagesByWorkshopID(workshopID uuid.UUID) ([]models.WorkshopImage, error) {
	query := `
		SELECT id, workshop_id, image_url, is_cover, sort_order, alt_text, alt_text_ar, created_at
		FROM workshop_images
		WHERE workshop_id = $1
		ORDER BY sort_order ASC, created_at ASC
	`

	rows, err := database.Instance.Query(query, workshopID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workshop images: %w", err)
	}
	defer rows.Close()

	var images []models.WorkshopImage
	for rows.Next() {
		var image models.WorkshopImage
		var altText sql.NullString
		var altTextAr sql.NullString

		err := rows.Scan(
			&image.ID,
			&image.WorkshopID,
			&image.ImageURL,
			&image.IsCover,
			&image.SortOrder,
			&altText,
			&altTextAr,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workshop image: %w", err)
		}

		// Handle nullable fields
		if altText.Valid {
			image.AltText = altText.String
		}
		if altTextAr.Valid {
			image.AltTextAr = altTextAr.String
		}

		images = append(images, image)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workshop images: %w", err)
	}

	return images, nil
}

// DeleteWorkshopImages deletes all images for a workshop
func (s *WorkshopService) DeleteWorkshopImages(workshopID uuid.UUID) error {
	query := `DELETE FROM workshop_images WHERE workshop_id = $1`

	_, err := database.Instance.Exec(query, workshopID)
	if err != nil {
		return fmt.Errorf("failed to delete workshop images: %w", err)
	}

	return nil
}

// DeleteWorkshopImage deletes a single workshop image
func (s *WorkshopService) DeleteWorkshopImage(imageID uuid.UUID) error {
	query := `DELETE FROM workshop_images WHERE id = $1`

	result, err := database.Instance.Exec(query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete workshop image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("image not found")
	}

	return nil
}

// UpdateWorkshopImageCover updates which image is the cover
func (s *WorkshopService) UpdateWorkshopImageCover(workshopID uuid.UUID, newCoverImageID uuid.UUID) error {
	// Start transaction
	tx, err := database.Instance.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Remove cover flag from all images for this workshop
	query1 := `UPDATE workshop_images SET is_cover = false WHERE workshop_id = $1`
	_, err = tx.Exec(query1, workshopID)
	if err != nil {
		return fmt.Errorf("failed to remove cover flags: %w", err)
	}

	// Set new cover image
	query2 := `UPDATE workshop_images SET is_cover = true WHERE id = $1 AND workshop_id = $2`
	result, err := tx.Exec(query2, newCoverImageID, workshopID)
	if err != nil {
		return fmt.Errorf("failed to set new cover: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("image not found or doesn't belong to workshop")
	}

	return tx.Commit()
}

// ProcessWorkshopImages handles the creation/update of workshop images
func (s *WorkshopService) ProcessWorkshopImages(workshopID uuid.UUID, imageURLs []string, coverImageIndex int) error {
	// Delete existing images for this workshop
	err := s.DeleteWorkshopImages(workshopID)
	if err != nil {
		return fmt.Errorf("failed to delete existing images: %w", err)
	}

	// Create new image records
	for i, imageURL := range imageURLs {
		if imageURL == "" {
			continue
		}

		image := &models.WorkshopImage{
			ID:         uuid.New(),
			WorkshopID: workshopID,
			ImageURL:   imageURL,
			IsCover:    i == coverImageIndex,
			SortOrder:  i,
			CreatedAt:  time.Now(),
		}

		err = s.CreateWorkshopImage(image)
		if err != nil {
			// Log error but continue with other images
			fmt.Printf("Error creating workshop image: %v\n", err)
		}
	}

	return nil
}

// SessionHasOrders checks if a workshop session has any associated orders
func (s *WorkshopService) SessionHasOrders(sessionID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM order_items WHERE session_id = $1`
	
	var count int
	err := database.Instance.QueryRow(query, sessionID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check session orders: %w", err)
	}
	
	return count > 0, nil
}

// UpdatePrivateWorkshopSession updates the duration and capacity of a private workshop session
func (s *WorkshopService) UpdatePrivateWorkshopSession(sessionID uuid.UUID, duration float64, capacity int) error {
	query := `
		UPDATE workshop_sessions 
		SET duration = $1, max_attendees = $2, updated_at = NOW()
		WHERE id = $3`
	
	_, err := database.Instance.Exec(query, duration, capacity, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update private workshop session: %w", err)
	}
	
	return nil
}
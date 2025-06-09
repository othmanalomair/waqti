package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"waqti/internal/database"
	"waqti/internal/middleware"
	"waqti/internal/models"
	"waqti/internal/services"
	"waqti/web/templates"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WorkshopHandler struct {
	workshopService *services.WorkshopService
}

func NewWorkshopHandler(workshopService *services.WorkshopService) *WorkshopHandler {
	return &WorkshopHandler{
		workshopService: workshopService,
	}
}

func (h *WorkshopHandler) ShowAddWorkshop(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	// Convert to models.Creator for template compatibility
	creator := &models.Creator{
		ID:       dbCreator.ID,
		Name:     dbCreator.Name,
		NameAr:   dbCreator.NameAr,
		Username: dbCreator.Username,
		Email:    dbCreator.Email,
		Plan:     dbCreator.Plan,
		PlanAr:   dbCreator.PlanAr,
		IsActive: dbCreator.IsActive,
	}

	return templates.AddWorkshopPage(creator, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) CreateWorkshop(c echo.Context) error {
	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	// Parse form data
	name := strings.TrimSpace(c.FormValue("name"))
	description := strings.TrimSpace(c.FormValue("description"))
	priceStr := c.FormValue("price")
	currency := c.FormValue("currency")
	isFree := c.FormValue("is_free") == "on" || c.FormValue("is_free") == "true"
	status := c.FormValue("status")

	// Validate required fields
	if name == "" {
		return c.String(http.StatusBadRequest, "Workshop name is required")
	}

	// Parse price
	var price float64 = 0
	if !isFree && priceStr != "" {
		var err error
		price, err = strconv.ParseFloat(priceStr, 64)
		if err != nil || price < 0 {
			return c.String(http.StatusBadRequest, "Invalid price")
		}
	}


	// Set default currency
	if currency == "" {
		currency = "KWD"
	}

	// Validate status
	if status != "published" && status != "draft" {
		status = "draft"
	}

	// Create workshop object
	workshop := &models.Workshop{
		ID:             uuid.New(),
		CreatorID:      dbCreator.ID,
		Name:           name,
		Title:          name,
		TitleAr:        "",
		Description:    description,
		DescriptionAr:  "",
		Price:          price,
		Currency:       currency,
		Duration:       120,
		MaxStudents:    0,
		Status:         status,
		IsActive:       status == "published",
		IsFree:         isFree,
		IsRecurring:    false,           // Set to false for new workshop_type system
		RecurrenceType: nil,             // Set to nil for new workshop_type system
		WorkshopType:   c.FormValue("workshop_type"), // Get from form
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Parse sessions data
	sessions, err := h.parseSessions(c)
	if err != nil {
		c.Logger().Error("Error parsing sessions:", err)
		return c.String(http.StatusBadRequest, "Invalid session data: "+err.Error())
	}


	// Create workshop in database
	err = h.workshopService.CreateWorkshop(workshop)
	if err != nil {
		c.Logger().Error("Error creating workshop:", err)
		return c.String(http.StatusInternalServerError, "Failed to create workshop")
	}

	// Create sessions if provided
	if len(sessions) > 0 {
		// Create a workshop run for these sessions
		runName := fmt.Sprintf("%s - %s", workshop.Name, time.Now().Format("Jan 2006"))
		
		// Create run record first
		runID := uuid.New()
		startDate := sessions[0].SessionDate
		endDate := sessions[len(sessions)-1].SessionDate
		
		for i, session := range sessions {
			session.WorkshopID = workshop.ID
			session.RunID = &runID
			session.SessionNumber = i + 1
			session.Status = "upcoming"
			session.StatusAr = "قادم"
			
			err = h.workshopService.CreateWorkshopSession(&session)
			if err != nil {
				c.Logger().Error("Error creating workshop session:", err)
			}
		}
		
		// Create the run record
		err = h.createWorkshopRun(runID, workshop.ID, runName, startDate, endDate)
		if err != nil {
			c.Logger().Error("Error creating workshop run:", err)
		}
	}

	// Process workshop images
	imageURLs := c.Request().Form["image_urls[]"]
	coverImageIndexStr := c.FormValue("cover_image_index")
	if len(imageURLs) > 0 {
		coverImageIndex := 0
		if coverImageIndexStr != "" {
			if idx, err := strconv.Atoi(coverImageIndexStr); err == nil {
				coverImageIndex = idx
			}
		}

		err = h.workshopService.ProcessWorkshopImages(workshop.ID, imageURLs, coverImageIndex)
		if err != nil {
			c.Logger().Error("Error processing workshop images:", err)
			// Continue even if images fail - workshop is already created
		}
	}

	c.Logger().Infof("Workshop created successfully: %s (ID: %s)", workshop.Name, workshop.ID)

	// Redirect to workshop reorder page
	return c.Redirect(http.StatusSeeOther, "/workshops/reorder")
}

func (h *WorkshopHandler) parseSessions(c echo.Context) ([]models.WorkshopSession, error) {
	var rawSessions []struct {
		Date     time.Time
		Time     string
		Duration float64
		Capacity int
	}

	// Get all form values
	form, err := c.FormParams()
	if err != nil {
		return nil, err
	}

	// Parse session data
	sessionIndex := 0
	for {
		dateKey := fmt.Sprintf("session_date_%d", sessionIndex)
		timeKey := fmt.Sprintf("session_time_%d", sessionIndex)
		durationKey := fmt.Sprintf("session_duration_%d", sessionIndex)
		capacityKey := fmt.Sprintf("session_capacity_%d", sessionIndex)

		dateStr := ""
		timeStr := ""
		durationStr := ""
		capacityStr := ""

		// Check if form has these keys
		if values, exists := form[dateKey]; exists && len(values) > 0 {
			dateStr = values[0]
		}
		if values, exists := form[timeKey]; exists && len(values) > 0 {
			timeStr = values[0]
		}
		if values, exists := form[durationKey]; exists && len(values) > 0 {
			durationStr = values[0]
		}
		if values, exists := form[capacityKey]; exists && len(values) > 0 {
			capacityStr = values[0]
		}

		// If no date found, we've reached the end
		if dateStr == "" {
			break
		}

		// Parse date
		sessionDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid session date format for session %d: %v", sessionIndex, err)
		}

		// Validate that time is not empty
		if timeStr == "" {
			return nil, fmt.Errorf("session time is required for session %d", sessionIndex)
		}

		// Parse duration
		duration := 2.0 // Default 2 hours
		if durationStr != "" {
			duration, err = strconv.ParseFloat(durationStr, 64)
			if err != nil || duration <= 0 {
				return nil, fmt.Errorf("invalid duration for session %d: %v", sessionIndex, err)
			}
		}

		// Parse capacity
		capacity := 20 // Default capacity
		if capacityStr != "" {
			capacity, err = strconv.Atoi(capacityStr)
			if err != nil || capacity < 0 {
				capacity = 20 // Fallback to default if invalid
			}
		}
		
		// Check for unlimited capacity (capacity of 0 means unlimited)
		if capacity == 0 {
			capacity = 0 // Explicitly set to 0 for unlimited
		}

		rawSessions = append(rawSessions, struct {
			Date     time.Time
			Time     string
			Duration float64
			Capacity int
		}{
			Date:     sessionDate,
			Time:     timeStr,
			Duration: duration,
			Capacity: capacity,
		})

		sessionIndex++
	}

	// Group consecutive days into single sessions
	return h.groupConsecutiveDays(rawSessions)
}

func (h *WorkshopHandler) groupConsecutiveDays(rawSessions []struct {
	Date     time.Time
	Time     string
	Duration float64
	Capacity int
}) ([]models.WorkshopSession, error) {
	if len(rawSessions) == 0 {
		return nil, nil
	}

	// Sort sessions by date
	sort.Slice(rawSessions, func(i, j int) bool {
		return rawSessions[i].Date.Before(rawSessions[j].Date)
	})

	var sessions []models.WorkshopSession
	sessionNumber := 1

	// Group sessions by identical properties (time, duration, capacity)
	// Each group becomes ONE session regardless of date gaps
	sessionGroups := make(map[string][]time.Time)
	sessionProps := make(map[string]struct {
		Time     string
		Duration float64
		Capacity int
	})

	for _, raw := range rawSessions {
		// Create a key based on session properties
		key := fmt.Sprintf("%s_%f_%d", raw.Time, raw.Duration, raw.Capacity)
		
		// Group dates by identical properties
		sessionGroups[key] = append(sessionGroups[key], raw.Date)
		sessionProps[key] = struct {
			Time     string
			Duration float64
			Capacity int
		}{
			Time:     raw.Time,
			Duration: raw.Duration,
			Capacity: raw.Capacity,
		}
	}

	// Create ONE session for each group of dates with identical properties
	for key, dates := range sessionGroups {
		props := sessionProps[key]
		
		// Sort dates within the group
		sort.Slice(dates, func(i, j int) bool {
			return dates[i].Before(dates[j])
		})

		// Parse time
		startTime, err := time.Parse("15:04", props.Time)
		if err != nil {
			return nil, fmt.Errorf("invalid session time format: %v", err)
		}

		// Calculate end time
		endTime := startTime.Add(time.Duration(props.Duration * float64(time.Hour)))
		endTimeStr := endTime.Format("15:04:05")

		// Get first and last dates for legacy compatibility
		firstDate := dates[0]
		var lastDate *time.Time
		if len(dates) > 1 {
			lastDate = &dates[len(dates)-1]
		}

		// Create ONE session with ALL dates (regardless of gaps)
		session := models.WorkshopSession{
			ID:               uuid.New(),
			SessionDate:      firstDate,              // Primary start date for compatibility
			EndDate:          lastDate,               // Legacy end date
			SessionDates:     dates,                  // ALL dates for this session
			TotalDays:        len(dates),             // Total days including gaps
			StartTime:        startTime.Format("15:04:05"),
			EndTime:          &endTimeStr,
			Duration:         props.Duration,
			DayCount:         len(dates),             // Legacy field
			Timezone:         "Asia/Kuwait",
			MaxAttendees:     props.Capacity,
			CurrentAttendees: 0,
			IsCompleted:      false,
			Status:           "upcoming",
			StatusAr:         "قادم",
			SessionNumber:    sessionNumber,
			Metadata:         make(map[string]interface{}),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		sessions = append(sessions, session)
		sessionNumber++
	}

	// Sort sessions by first date
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].SessionDate.Before(sessions[j].SessionDate)
	})

	return sessions, nil
}

// createWorkshopRun creates a workshop run record in the database
func (h *WorkshopHandler) createWorkshopRun(runID, workshopID uuid.UUID, runName string, startDate, endDate time.Time) error {
	// For now, return nil - this will be implemented when we have proper DB access
	// or could be moved to the workshop service
	return nil
}

func (h *WorkshopHandler) ShowReorderWorkshops(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	// Convert to models.Creator for template compatibility
	creator := &models.Creator{
		ID:       dbCreator.ID,
		Name:     dbCreator.Name,
		NameAr:   dbCreator.NameAr,
		Username: dbCreator.Username,
		Email:    dbCreator.Email,
		Plan:     dbCreator.Plan,
		PlanAr:   dbCreator.PlanAr,
		IsActive: dbCreator.IsActive,
	}

	workshops := h.workshopService.GetWorkshopsByCreatorID(dbCreator.ID)
	stats := h.workshopService.GetDashboardStats(dbCreator.ID)

	return templates.ReorderWorkshopsPage(creator, workshops, stats, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) ReorderWorkshop(c echo.Context) error {
	// Check authentication
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	workshopIDStr := c.FormValue("workshop_id")
	direction := c.FormValue("direction")

	workshopID, err := uuid.Parse(workshopIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid workshop ID")
	}

	err = h.workshopService.ReorderWorkshop(workshopID, direction)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to reorder workshop")
	}

	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	workshops := h.workshopService.GetWorkshopsByCreatorID(dbCreator.ID)

	return templates.WorkshopsListFixed(workshops, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) ToggleWorkshopStatus(c echo.Context) error {
	// Check authentication
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	workshopIDStr := c.FormValue("workshop_id")

	workshopID, err := uuid.Parse(workshopIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid workshop ID")
	}

	err = h.workshopService.ToggleWorkshopStatus(workshopID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to toggle workshop status")
	}

	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	workshops := h.workshopService.GetWorkshopsByCreatorID(dbCreator.ID)

	return templates.WorkshopsListFixed(workshops, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) ShowEditWorkshop(c echo.Context) error {
	lang := c.Get("lang").(string)
	isRTL := lang == "ar"

	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	// Get workshop ID from URL parameter
	workshopIDStr := c.Param("id")
	c.Logger().Infof("DEBUG: Attempting to edit workshop ID: %s", workshopIDStr)

	workshopID, err := uuid.Parse(workshopIDStr)
	if err != nil {
		c.Logger().Error("DEBUG: Invalid workshop ID:", err)
		return c.String(http.StatusBadRequest, "Invalid workshop ID")
	}

	// Get workshop from database
	workshop, err := h.workshopService.GetWorkshopByID(workshopID, dbCreator.ID)
	if err != nil {
		c.Logger().Error("DEBUG: Error getting workshop:", err)
		return c.String(http.StatusInternalServerError, "Workshop not found")
	}
	if workshop == nil {
		c.Logger().Error("DEBUG: Workshop not found in database")
		return c.String(http.StatusNotFound, "Workshop not found")
	}

	// Convert to models.Creator for template compatibility
	creator := &models.Creator{
		ID:       dbCreator.ID,
		Name:     dbCreator.Name,
		NameAr:   dbCreator.NameAr,
		Username: dbCreator.Username,
		Email:    dbCreator.Email,
		Plan:     dbCreator.Plan,
		PlanAr:   dbCreator.PlanAr,
		IsActive: dbCreator.IsActive,
	}

	// Get workshop sessions
	sessions, err := h.workshopService.GetWorkshopSessions(workshopID)
	if err != nil {
		c.Logger().Error("DEBUG: Error getting workshop sessions:", err)
		sessions = []models.WorkshopSession{}
	}

	c.Logger().Infof("DEBUG: Found %d sessions for workshop", len(sessions))

	return templates.EditWorkshopPage(creator, workshop, sessions, lang, isRTL).Render(c.Request().Context(), c.Response().Writer)
}

func (h *WorkshopHandler) UpdateWorkshop(c echo.Context) error {
	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.Redirect(http.StatusSeeOther, "/signin")
	}

	// Get workshop ID from URL parameter
	workshopIDStr := c.Param("id")
	workshopID, err := uuid.Parse(workshopIDStr)
	if err != nil {
		c.Logger().Error("Invalid workshop ID:", err)
		return c.Redirect(http.StatusSeeOther, "/workshops/reorder?error=invalid_id")
	}

	// Parse form data
	name := strings.TrimSpace(c.FormValue("name"))
	description := strings.TrimSpace(c.FormValue("description"))
	priceStr := c.FormValue("price")
	currency := c.FormValue("currency")
	durationStr := c.FormValue("duration")
	maxStudentsStr := c.FormValue("max_students")
	isFree := c.FormValue("is_free") == "on" || c.FormValue("is_free") == "true"
	status := c.FormValue("status")

	// Validate required fields
	if name == "" {
		c.Logger().Error("Workshop name is required")
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/workshops/edit/%s?error=name_required", workshopIDStr))
	}

	// Parse price
	var price float64 = 0
	if !isFree && priceStr != "" {
		var err error
		price, err = strconv.ParseFloat(priceStr, 64)
		if err != nil || price < 0 {
			c.Logger().Error("Invalid price:", err)
			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/workshops/edit/%s?error=invalid_price", workshopIDStr))
		}
	}

	// Parse duration (default 120 minutes)
	var duration int = 120
	if durationStr != "" {
		duration, err = strconv.Atoi(durationStr)
		if err != nil || duration <= 0 {
			duration = 120
		}
	}

	// Parse max students (default 0 = unlimited)
	var maxStudents int = 0
	if maxStudentsStr != "" {
		maxStudents, err = strconv.Atoi(maxStudentsStr)
		if err != nil || maxStudents < 0 {
			maxStudents = 0
		}
	}

	// Set default currency
	if currency == "" {
		currency = "KWD"
	}

	// Validate status
	if status != "published" && status != "draft" {
		status = "draft"
	}

	// Get existing workshop to update
	existingWorkshop, err := h.workshopService.GetWorkshopByID(workshopID, dbCreator.ID)
	if err != nil || existingWorkshop == nil {
		c.Logger().Error("Workshop not found for update:", err)
		return c.Redirect(http.StatusSeeOther, "/workshops/reorder?error=workshop_not_found")
	}

	// Update workshop object
	existingWorkshop.Name = name
	existingWorkshop.Title = name
	existingWorkshop.Description = description
	existingWorkshop.Price = price
	existingWorkshop.Currency = currency
	existingWorkshop.Duration = duration
	existingWorkshop.MaxStudents = maxStudents
	existingWorkshop.Status = status
	existingWorkshop.IsActive = status == "published"
	existingWorkshop.IsFree = isFree
	existingWorkshop.IsRecurring = false     // Set to false for new workshop_type system
	existingWorkshop.RecurrenceType = nil    // Set to nil for new workshop_type system
	existingWorkshop.WorkshopType = c.FormValue("workshop_type") // Get from form
	existingWorkshop.UpdatedAt = time.Now()

	// Update workshop in database
	err = h.workshopService.UpdateWorkshop(existingWorkshop)
	if err != nil {
		c.Logger().Error("Error updating workshop:", err)
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/workshops/edit/%s?error=update_failed", workshopIDStr))
	}

	// Handle session updates
	err = h.updateWorkshopSessions(c, workshopID)
	if err != nil {
		c.Logger().Error("Error updating workshop sessions:", err)
	}

	// Process workshop images
	imageURLs := c.Request().Form["image_urls[]"]
	coverImageIndexStr := c.FormValue("cover_image_index")
	if len(imageURLs) > 0 {
		coverImageIndex := 0
		if coverImageIndexStr != "" {
			if idx, err := strconv.Atoi(coverImageIndexStr); err == nil {
				coverImageIndex = idx
			}
		}

		err = h.workshopService.ProcessWorkshopImages(workshopID, imageURLs, coverImageIndex)
		if err != nil {
			c.Logger().Error("Error processing workshop images:", err)
		}
	}

	c.Logger().Infof("Workshop updated successfully: %s (ID: %s)", existingWorkshop.Name, existingWorkshop.ID)

	// Determine success message based on status
	successMsg := "workshop_updated"
	if status == "draft" {
		successMsg = "draft_saved"
	} else if status == "published" {
		successMsg = "workshop_published"
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/workshops/reorder?success=%s", successMsg))
}

func (h *WorkshopHandler) DeleteWorkshop(c echo.Context) error {
	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	// Get workshop ID from URL parameter
	workshopIDStr := c.Param("id")
	workshopID, err := uuid.Parse(workshopIDStr)
	if err != nil {
		c.Logger().Error("Invalid workshop ID:", err)
		return c.String(http.StatusBadRequest, "Invalid workshop ID")
	}

	// Verify workshop belongs to current creator
	workshop, err := h.workshopService.GetWorkshopByID(workshopID, dbCreator.ID)
	if err != nil || workshop == nil {
		c.Logger().Error("Workshop not found for deletion:", err)
		return c.String(http.StatusNotFound, "Workshop not found")
	}

	// Delete workshop from database (this should cascade delete sessions)
	err = h.workshopService.DeleteWorkshop(workshopID)
	if err != nil {
		c.Logger().Error("Error deleting workshop:", err)
		return c.String(http.StatusInternalServerError, "Failed to delete workshop")
	}

	c.Logger().Infof("Workshop deleted successfully: %s (ID: %s)", workshop.Name, workshop.ID)

	return c.String(http.StatusOK, "Workshop deleted successfully")
}

func (h *WorkshopHandler) updateWorkshopSessions(c echo.Context, workshopID uuid.UUID) error {
	// Parse sessions data similar to creation
	newSessions, err := h.parseSessions(c)
	if err != nil {
		return fmt.Errorf("failed to parse sessions: %w", err)
	}

	// Get existing sessions for comparison
	existingSessions, err := h.workshopService.GetWorkshopSessions(workshopID)
	if err != nil {
		return fmt.Errorf("failed to get existing sessions: %w", err)
	}

	// Create a map of existing sessions by date+time for easier matching
	existingMap := make(map[string]*models.WorkshopSession)
	for i := range existingSessions {
		key := fmt.Sprintf("%s_%s", existingSessions[i].SessionDate.Format("2006-01-02"), existingSessions[i].StartTime)
		existingMap[key] = &existingSessions[i]
	}

	// Process each new session
	for _, newSession := range newSessions {
		key := fmt.Sprintf("%s_%s", newSession.SessionDate.Format("2006-01-02"), newSession.StartTime)
		existingSession, exists := existingMap[key]
		
		if exists {
			// Calculate new end time based on start time and duration
			startTime, err := time.Parse("15:04:05", existingSession.StartTime)
			if err != nil {
				return fmt.Errorf("failed to parse existing start time %s: %w", existingSession.StartTime, err)
			}
			
			// Add duration (in hours) to start time to get end time
			endTime := startTime.Add(time.Duration(newSession.Duration * float64(time.Hour)))
			endTimeStr := endTime.Format("15:04:05")
			
			// Update existing session with recalculated end_time
			updateQuery := `
				UPDATE workshop_sessions 
				SET max_attendees = $1, 
				    notes = $2, 
				    notes_ar = $3,
				    duration = $4,
				    end_time = $5,
				    updated_at = NOW()
				WHERE id = $6`
			
			_, err = database.Instance.Exec(updateQuery, 
				newSession.MaxAttendees, 
				newSession.Notes, 
				newSession.NotesAr, 
				newSession.Duration,
				endTimeStr,
				existingSession.ID)
			if err != nil {
				return fmt.Errorf("failed to update session %s: %w", existingSession.ID.String(), err)
			}
			
			fmt.Printf("Updated existing session %s: max_attendees=%d, duration=%.1f, end_time=%s\n", 
				existingSession.ID.String(), newSession.MaxAttendees, newSession.Duration, endTimeStr)
		} else {
			// This is a new session - create it
			newSession.ID = uuid.New()
			newSession.WorkshopID = workshopID
			newSession.Status = "upcoming"
			newSession.StatusAr = "قادم"
			newSession.SessionNumber = len(existingSessions) + 1  // Assign next session number
			
			// Create a run_id for this session (can be same as session ID for simplicity)
			runID := uuid.New()
			newSession.RunID = &runID
			
			err = h.workshopService.CreateWorkshopSession(&newSession)
			if err != nil {
				return fmt.Errorf("failed to create new session: %w", err)
			}
			
			fmt.Printf("Created new session %s for date %s time %s\n", 
				newSession.ID.String(), newSession.SessionDate.Format("2006-01-02"), newSession.StartTime)
		}
	}

	return nil
}

func (h *WorkshopHandler) DeleteWorkshopSession(c echo.Context) error {
	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	// Get session ID from URL parameter
	sessionIDStr := c.Param("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid session ID")
	}

	// Verify session belongs to creator's workshop
	session, err := h.workshopService.GetWorkshopSessionByID(sessionID)
	if err != nil || session == nil {
		return c.String(http.StatusNotFound, "Session not found")
	}

	// Verify workshop belongs to creator
	workshop, err := h.workshopService.GetWorkshopByID(session.WorkshopID, dbCreator.ID)
	if err != nil || workshop == nil {
		return c.String(http.StatusNotFound, "Workshop not found")
	}

	// Delete session
	err = h.workshopService.DeleteWorkshopSession(sessionID)
	if err != nil {
		c.Logger().Error("Error deleting session:", err)
		return c.String(http.StatusInternalServerError, "Failed to delete session")
	}

	return c.String(http.StatusOK, "Session deleted successfully")
}

// GetWorkshopImages returns images for a specific workshop
func (h *WorkshopHandler) GetWorkshopImages(c echo.Context) error {
	// Check authentication
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Unauthorized",
		})
	}

	// Get workshop ID from URL parameter
	workshopIDStr := c.Param("id")
	workshopID, err := uuid.Parse(workshopIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid workshop ID",
		})
	}

	// Verify workshop belongs to current creator
	workshop, err := h.workshopService.GetWorkshopByID(workshopID, dbCreator.ID)
	if err != nil || workshop == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   "Workshop not found",
		})
	}

	// Get workshop images
	images, err := h.workshopService.GetWorkshopImagesByWorkshopID(workshopID)
	if err != nil {
		c.Logger().Error("Error getting workshop images:", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to get workshop images",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"images":  images,
	})
}

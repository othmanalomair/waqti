package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
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

	// Handle session creation based on workshop type
	if workshop.WorkshopType == "private" {
		// For private workshops, create a special "always available" session
		c.Logger().Infof("Creating special session for private workshop %s", workshop.ID)
		
		// Parse duration and capacity from session data (private workshops send this as session_0)
		duration := 2.0 // Default 2 hours
		capacity := 1   // Default 1 person
		
		durationStr := c.FormValue("session_duration_0")
		capacityStr := c.FormValue("session_capacity_0")
		c.Logger().Infof("DEBUG: Creating private workshop - duration: '%s', capacity: '%s'", durationStr, capacityStr)
		
		if durationStr != "" {
			if parsedDuration, err := strconv.ParseFloat(durationStr, 64); err == nil && parsedDuration > 0 {
				duration = parsedDuration
				c.Logger().Infof("DEBUG: Parsed creation duration: %.2f", duration)
			} else {
				c.Logger().Warnf("DEBUG: Failed to parse creation duration '%s': %v", durationStr, err)
			}
		}
		
		if capacityStr != "" {
			if parsedCapacity, err := strconv.Atoi(capacityStr); err == nil && parsedCapacity >= 0 {
				capacity = parsedCapacity
				c.Logger().Infof("DEBUG: Parsed creation capacity: %d", capacity)
			} else {
				c.Logger().Warnf("DEBUG: Failed to parse creation capacity '%s': %v", capacityStr, err)
			}
		}
		
		// Create special session with far future date to indicate "always available"
		specialDate, _ := time.Parse("2006-01-02", "9999-12-31")
		specialSession := models.WorkshopSession{
			ID:               uuid.New(),
			WorkshopID:       workshop.ID,
			SessionDate:      specialDate,
			EndDate:          nil,
			SessionDates:     []time.Time{specialDate},
			TotalDays:        1,
			StartTime:        "00:00:00", // No specific time for private workshops
			EndTime:          nil,
			Duration:         duration,
			DayCount:         1,
			Timezone:         "Asia/Kuwait",
			MaxAttendees:     capacity,
			CurrentAttendees: 0,
			IsCompleted:      false,
			Status:           "upcoming",
			StatusAr:         "قادم",
			SessionNumber:    1,
			Metadata:         make(map[string]interface{}),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		
		// Create run for private workshop
		runID := uuid.New()
		specialSession.RunID = &runID
		
		err = h.workshopService.CreateWorkshopSession(&specialSession)
		if err != nil {
			c.Logger().Error("Error creating private workshop session:", err)
		} else {
			c.Logger().Infof("Created special session %s for private workshop", specialSession.ID)
		}
		
		// Create the run record
		err = h.createWorkshopRun(runID, workshop.ID, fmt.Sprintf("%s - Private", workshop.Name), specialDate, specialDate)
		if err != nil {
			c.Logger().Error("Error creating workshop run:", err)
		}
	} else if len(sessions) > 0 {
		// For regular workshops, create sessions as before
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
		Date         time.Time
		Time         string
		Duration     float64
		Capacity     int
		SessionDates []time.Time // For consecutive sessions
		TotalDays    int         // For consecutive sessions
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
		sessionDatesKey := fmt.Sprintf("session_dates_%d", sessionIndex)
		totalDaysKey := fmt.Sprintf("session_total_days_%d", sessionIndex)

		dateStr := ""
		timeStr := ""
		durationStr := ""
		capacityStr := ""
		sessionDatesStr := ""
		totalDaysStr := ""

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
		if values, exists := form[sessionDatesKey]; exists && len(values) > 0 {
			sessionDatesStr = values[0]
		}
		if values, exists := form[totalDaysKey]; exists && len(values) > 0 {
			totalDaysStr = values[0]
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

		// Parse consecutive session dates if provided
		var sessionDates []time.Time
		var totalDays int
		if sessionDatesStr != "" {
			dateStrings := strings.Split(sessionDatesStr, ",")
			for _, dateStr := range dateStrings {
				if dateStr = strings.TrimSpace(dateStr); dateStr != "" {
					if date, err := time.Parse("2006-01-02", dateStr); err == nil {
						sessionDates = append(sessionDates, date)
					}
				}
			}
		}
		if totalDaysStr != "" {
			if days, err := strconv.Atoi(totalDaysStr); err == nil {
				totalDays = days
			}
		}

		rawSessions = append(rawSessions, struct {
			Date         time.Time
			Time         string
			Duration     float64
			Capacity     int
			SessionDates []time.Time
			TotalDays    int
		}{
			Date:         sessionDate,
			Time:         timeStr,
			Duration:     duration,
			Capacity:     capacity,
			SessionDates: sessionDates,
			TotalDays:    totalDays,
		})

		sessionIndex++
	}

	// Group consecutive days into single sessions
	return h.groupConsecutiveDays(rawSessions)
}

func (h *WorkshopHandler) groupConsecutiveDays(rawSessions []struct {
	Date         time.Time
	Time         string
	Duration     float64
	Capacity     int
	SessionDates []time.Time
	TotalDays    int
}) ([]models.WorkshopSession, error) {
	if len(rawSessions) == 0 {
		return nil, nil
	}

	var sessions []models.WorkshopSession
	sessionNumber := 1

	// Check if we have sessions with pre-populated SessionDates (consecutive sessions from frontend)
	for _, raw := range rawSessions {
		if len(raw.SessionDates) > 0 {
			// This is a consecutive session with multiple dates already populated
			// Parse time
			startTime, err := time.Parse("15:04", raw.Time)
			if err != nil {
				return nil, fmt.Errorf("invalid session time format: %v", err)
			}

			// Calculate end time
			endTime := startTime.Add(time.Duration(raw.Duration * float64(time.Hour)))
			endTimeStr := endTime.Format("15:04:05")

			// Create session with all dates in SessionDates
			firstDate := raw.SessionDates[0]
			var lastDate *time.Time
			if len(raw.SessionDates) > 1 {
				lastDate = &raw.SessionDates[len(raw.SessionDates)-1]
			}

			session := models.WorkshopSession{
				ID:               uuid.New(),
				SessionDate:      firstDate,              // Primary start date for compatibility
				EndDate:          lastDate,               // Legacy end date
				SessionDates:     raw.SessionDates,       // All dates for this consecutive session
				TotalDays:        len(raw.SessionDates),  // Total days in this range
				StartTime:        startTime.Format("15:04:05"),
				EndTime:          &endTimeStr,
				Duration:         raw.Duration,
				DayCount:         len(raw.SessionDates),  // Legacy field
				Timezone:         "Asia/Kuwait",
				MaxAttendees:     raw.Capacity,
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
			continue
		}
	}

	// For non-consecutive sessions, create one session per raw session (don't group)
	// This ensures each date gets its own database row
	for _, raw := range rawSessions {
		// Skip if this was already handled above (has SessionDates)
		if len(raw.SessionDates) > 0 {
			continue
		}

		// Parse time
		startTime, err := time.Parse("15:04", raw.Time)
		if err != nil {
			return nil, fmt.Errorf("invalid session time format: %v", err)
		}

		// Calculate end time
		endTime := startTime.Add(time.Duration(raw.Duration * float64(time.Hour)))
		endTimeStr := endTime.Format("15:04:05")

		// Create individual session for each date (no grouping)
		session := models.WorkshopSession{
			ID:               uuid.New(),
			SessionDate:      raw.Date,                    // Individual date
			EndDate:          nil,                         // No end date for single sessions
			SessionDates:     []time.Time{raw.Date},       // Single date in array
			TotalDays:        1,                           // Always 1 for individual sessions
			StartTime:        startTime.Format("15:04:05"),
			EndTime:          &endTimeStr,
			Duration:         raw.Duration,
			DayCount:         1,                           // Always 1 for individual sessions
			Timezone:         "Asia/Kuwait",
			MaxAttendees:     raw.Capacity,
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

// splitIntoConsecutiveRanges splits a list of dates into consecutive date ranges
func (h *WorkshopHandler) splitIntoConsecutiveRanges(dates []time.Time) [][]time.Time {
	if len(dates) == 0 {
		return nil
	}

	var ranges [][]time.Time
	var currentRange []time.Time
	
	for i, date := range dates {
		if i == 0 {
			// First date - start new range
			currentRange = []time.Time{date}
		} else {
			// Check if current date is consecutive to previous date
			prevDate := dates[i-1]
			if date.Sub(prevDate) == 24*time.Hour {
				// Consecutive - add to current range
				currentRange = append(currentRange, date)
			} else {
				// Not consecutive - finish current range and start new one
				ranges = append(ranges, currentRange)
				currentRange = []time.Time{date}
			}
		}
	}
	
	// Add the last range
	if len(currentRange) > 0 {
		ranges = append(ranges, currentRange)
	}
	
	return ranges
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
	fmt.Printf("================================\n")
	fmt.Printf("DEBUG: UpdateWorkshop called - starting update process\n")
	fmt.Printf("================================\n")
	c.Logger().Infof("DEBUG: UpdateWorkshop called - starting update process")
	
	// Get current authenticated creator
	dbCreator := middleware.GetCurrentCreator(c)
	if dbCreator == nil {
		c.Logger().Infof("DEBUG: No authenticated creator found, redirecting to signin")
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
	fmt.Printf("DEBUG: Parsing form data...\n")
	name := strings.TrimSpace(c.FormValue("name"))
	description := strings.TrimSpace(c.FormValue("description"))
	priceStr := c.FormValue("price")
	currency := c.FormValue("currency")
	durationStr := c.FormValue("duration")
	maxStudentsStr := c.FormValue("max_students")
	isFree := c.FormValue("is_free") == "on" || c.FormValue("is_free") == "true"
	status := c.FormValue("status")
	fmt.Printf("DEBUG: Form data parsed - name: %s, price: %s, status: %s\n", name, priceStr, status)

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
	fmt.Printf("DEBUG: Getting existing workshop for update...\n")
	existingWorkshop, err := h.workshopService.GetWorkshopByID(workshopID, dbCreator.ID)
	if err != nil || existingWorkshop == nil {
		fmt.Printf("DEBUG: Workshop not found - err: %v, workshop: %v\n", err, existingWorkshop)
		c.Logger().Error("Workshop not found for update:", err)
		return c.Redirect(http.StatusSeeOther, "/workshops/reorder?error=workshop_not_found")
	}
	fmt.Printf("DEBUG: Workshop found for update - Type: %s\n", existingWorkshop.WorkshopType)

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
	fmt.Printf("DEBUG: Updating workshop in database...\n")
	err = h.workshopService.UpdateWorkshop(existingWorkshop)
	if err != nil {
		fmt.Printf("DEBUG: Error updating workshop: %v\n", err)
		c.Logger().Error("Error updating workshop:", err)
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/workshops/edit/%s?error=update_failed", workshopIDStr))
	}
	fmt.Printf("DEBUG: Workshop updated successfully in database\n")

	// Handle session updates
	fmt.Printf("DEBUG: About to call updateWorkshopSessions for workshopID: %s\n", workshopID)
	c.Logger().Infof("DEBUG: About to call updateWorkshopSessions for workshopID: %s", workshopID)
	err = h.updateWorkshopSessions(c, workshopID)
	if err != nil {
		c.Logger().Error("Error updating workshop sessions:", err)
	} else {
		c.Logger().Infof("DEBUG: updateWorkshopSessions completed successfully")
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
	// Check if this is a private workshop first
	workshopType := c.FormValue("workshop_type")
	fmt.Printf("DEBUG: updateWorkshopSessions called for workshopID: %s\n", workshopID)
	fmt.Printf("DEBUG: updateWorkshopSessions - workshop_type received: '%s'\n", workshopType)
	
	// Debug: Print ALL form values
	form, _ := c.FormParams()
	fmt.Printf("DEBUG: ALL FORM VALUES:\n")
	for key, values := range form {
		fmt.Printf("  %s: %v\n", key, values)
	}
	fmt.Printf("DEBUG: END FORM VALUES\n")
	
	c.Logger().Infof("DEBUG: updateWorkshopSessions - workshop_type received: '%s'", workshopType)
	
	if workshopType == "private" {
		// For private workshops, update the special session's duration and capacity
		fmt.Printf("DEBUG: Processing private workshop session for workshop %s\n", workshopID)
		c.Logger().Infof("Updating private workshop session for workshop %s", workshopID)
		
		// Parse duration and capacity from form
		duration := 2.0 // Default 2 hours
		capacity := 1   // Default 1 person
		
		durationStr := c.FormValue("session_duration_0")
		capacityStr := c.FormValue("session_capacity_0")
		fmt.Printf("DEBUG: Private workshop form values - duration: '%s', capacity: '%s'\n", durationStr, capacityStr)
		c.Logger().Infof("DEBUG: Private workshop form values - duration: '%s', capacity: '%s'", durationStr, capacityStr)
		
		if durationStr != "" {
			if parsedDuration, err := strconv.ParseFloat(durationStr, 64); err == nil && parsedDuration > 0 {
				duration = parsedDuration
				c.Logger().Infof("DEBUG: Parsed duration: %.2f", duration)
			} else {
				c.Logger().Warnf("DEBUG: Failed to parse duration '%s': %v", durationStr, err)
			}
		}
		
		if capacityStr != "" {
			if parsedCapacity, err := strconv.Atoi(capacityStr); err == nil && parsedCapacity >= 0 {
				capacity = parsedCapacity
				c.Logger().Infof("DEBUG: Parsed capacity: %d", capacity)
			} else {
				c.Logger().Warnf("DEBUG: Failed to parse capacity '%s': %v", capacityStr, err)
			}
		}
		
		// Get existing sessions for this workshop
		existingSessions, err := h.workshopService.GetWorkshopSessions(workshopID)
		if err != nil {
			return fmt.Errorf("failed to get existing sessions: %w", err)
		}
		
		// Find the private workshop session (should have year 9999)
		for _, session := range existingSessions {
			if session.SessionDate.Year() == 9999 {
				// Update the special session with new duration and capacity
				err = h.workshopService.UpdatePrivateWorkshopSession(session.ID, duration, capacity)
				if err != nil {
					return fmt.Errorf("failed to update private workshop session: %w", err)
				}
				
				c.Logger().Infof("Updated private workshop session %s: duration=%.1f, capacity=%d", 
					session.ID.String(), duration, capacity)
				break
			}
		}
		
		return nil
	}
	
	// For regular workshops, parse sessions data similar to creation
	newSessions, err := h.parseSessions(c)
	if err != nil {
		return fmt.Errorf("failed to parse sessions: %w", err)
	}

	// Check if any sessions have orders (foreign key constraints)
	hasOrders, err := h.checkWorkshopHasOrders(workshopID)
	if err != nil {
		c.Logger().Warnf("Could not check for orders, using safe update: %v", err)
		hasOrders = true // Default to safe mode if check fails
	}

	if hasOrders {
		// Use safe update method when orders exist
		c.Logger().Infof("Workshop has orders, using safe session update for workshop %s", workshopID)
		err = h.workshopService.UpdateWorkshopSessionsSafely(workshopID, newSessions)
		if err != nil {
			return fmt.Errorf("failed to safely update sessions: %w", err)
		}
	} else {
		// No orders exist, we can safely delete and recreate all sessions
		c.Logger().Infof("Workshop has no orders, using full session replacement for workshop %s", workshopID)
		
		// Get existing sessions
		existingSessions, err := h.workshopService.GetWorkshopSessions(workshopID)
		if err != nil {
			return fmt.Errorf("failed to get existing sessions: %w", err)
		}

		// Delete all existing sessions
		for _, existingSession := range existingSessions {
			err = h.workshopService.DeleteWorkshopSession(existingSession.ID)
			if err != nil {
				return fmt.Errorf("failed to delete session %s: %w", existingSession.ID.String(), err)
			}
		}

		// Create all new sessions
		for i, newSession := range newSessions {
			newSession.ID = uuid.New()
			newSession.WorkshopID = workshopID
			newSession.Status = "upcoming"
			newSession.StatusAr = "قادم"
			newSession.SessionNumber = i + 1
			
			// Create a run_id for this session
			runID := uuid.New()
			newSession.RunID = &runID
			
			err = h.workshopService.CreateWorkshopSession(&newSession)
			if err != nil {
				return fmt.Errorf("failed to create new session: %w", err)
			}
		}
	}

	return nil
}

// checkWorkshopHasOrders checks if a workshop has any orders associated with its sessions
func (h *WorkshopHandler) checkWorkshopHasOrders(workshopID uuid.UUID) (bool, error) {
	// Check if any sessions of this workshop have orders
	sessions, err := h.workshopService.GetWorkshopSessions(workshopID)
	if err != nil {
		return false, fmt.Errorf("failed to get workshop sessions: %w", err)
	}

	for _, session := range sessions {
		hasOrders, err := h.workshopService.SessionHasOrders(session.ID)
		if err != nil {
			return false, fmt.Errorf("failed to check orders for session %s: %w", session.ID, err)
		}
		if hasOrders {
			return true, nil
		}
	}

	return false, nil
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

package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Workshop struct {
	ID              uuid.UUID `json:"id"`
	CreatorID       uuid.UUID `json:"creator_id"`
	Name            string    `json:"name"`
	Title           string    `json:"title"`
	TitleAr         string    `json:"title_ar"`
	Description     string    `json:"description"`
	DescriptionAr   string    `json:"description_ar"`
	Price           float64   `json:"price"`
	Currency        string    `json:"currency"`
	Duration        int       `json:"duration"` // in minutes
	MaxStudents     int       `json:"max_students"`
	Category        string    `json:"category"`
	CategoryAr      string    `json:"category_ar"`
	Status          string    `json:"status"`
	IsActive        bool      `json:"is_active"`
	IsFree          bool      `json:"is_free"`
	IsRecurring     bool      `json:"is_recurring"`
	RecurrenceType  *string   `json:"recurrence_type"` // Change to pointer
	WorkshopType    string    `json:"workshop_type"` // single, consecutive, spread, custom
	SortOrder       int       `json:"sort_order"`
	ViewCount       int       `json:"view_count"`
	EnrollmentCount int       `json:"enrollment_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	// Related data for enhanced display
	Images          []WorkshopImage   `json:"images,omitempty"`
	Sessions        []WorkshopSession `json:"sessions,omitempty"`
}

type WorkshopSession struct {
	ID               uuid.UUID              `json:"id"`
	WorkshopID       uuid.UUID              `json:"workshop_id"`
	SessionDate      time.Time              `json:"session_date"`      // Primary start date (for compatibility)
	EndDate          *time.Time             `json:"end_date"`          // End date for multi-day sessions (may be deprecated)
	SessionDates     []time.Time            `json:"session_dates"`     // All dates this session spans (may have gaps)
	TotalDays        int                    `json:"total_days"`        // Total number of days including gaps
	StartTime        string                 `json:"start_time"`
	EndTime          *string                `json:"end_time"`
	Duration         float64                `json:"duration"`          // in hours per day
	DayCount         int                    `json:"day_count"`         // Number of consecutive days (legacy)
	Timezone         string                 `json:"timezone"`
	Location         *string                `json:"location"`
	LocationAr       *string                `json:"location_ar"`
	MaxAttendees     int                    `json:"max_attendees"`
	CurrentAttendees int                    `json:"current_attendees"`
	IsCompleted      bool                   `json:"is_completed"`
	Notes            *string                `json:"notes"`
	NotesAr          *string                `json:"notes_ar"`
	Status           string                 `json:"status"`            // upcoming, active, full, completed, cancelled
	StatusAr         string                 `json:"status_ar"`         // Arabic status
	SessionNumber    int                    `json:"session_number"`    // For numbering sessions
	RunID            *uuid.UUID             `json:"run_id"`            // Groups sessions in same run
	Metadata         map[string]interface{} `json:"metadata"`          // Flexible additional data
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// GetDateRange returns a formatted string of the session date range
func (ws *WorkshopSession) GetDateRange(lang string) string {
	if len(ws.SessionDates) == 0 {
		// Fallback to old format
		if ws.EndDate != nil {
			if lang == "ar" {
				return fmt.Sprintf("%s - %s", ws.SessionDate.Format("2 Jan"), ws.EndDate.Format("2 Jan 2006"))
			}
			return fmt.Sprintf("%s - %s", ws.SessionDate.Format("Jan 2"), ws.EndDate.Format("Jan 2, 2006"))
		}
		if lang == "ar" {
			return ws.SessionDate.Format("2 Jan 2006")
		}
		return ws.SessionDate.Format("Jan 2, 2006")
	}

	if len(ws.SessionDates) == 1 {
		if lang == "ar" {
			return ws.SessionDates[0].Format("2 Jan 2006")
		}
		return ws.SessionDates[0].Format("Jan 2, 2006")
	}

	// Multiple dates
	firstDate := ws.SessionDates[0]
	lastDate := ws.SessionDates[len(ws.SessionDates)-1]
	
	if lang == "ar" {
		return fmt.Sprintf("%d أيام: %s - %s", len(ws.SessionDates), 
			firstDate.Format("2 Jan"), lastDate.Format("2 Jan 2006"))
	}
	return fmt.Sprintf("%d days: %s - %s", len(ws.SessionDates),
		firstDate.Format("Jan 2"), lastDate.Format("Jan 2, 2006"))
}

type WorkshopRun struct {
	ID         uuid.UUID `json:"id"`
	WorkshopID uuid.UUID `json:"workshop_id"`
	RunName    string    `json:"run_name"`
	RunNameAr  string    `json:"run_name_ar"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Status     string    `json:"status"` // upcoming, active, full, completed, cancelled
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	// Related data
	Sessions []WorkshopSession `json:"sessions,omitempty"`
}

type WorkshopImage struct {
	ID         uuid.UUID `json:"id"`
	WorkshopID uuid.UUID `json:"workshop_id"`
	ImageURL   string    `json:"image_url"`
	IsCover    bool      `json:"is_cover"`
	SortOrder  int       `json:"sort_order"`
	AltText    string    `json:"alt_text"`
	AltTextAr  string    `json:"alt_text_ar"`
	CreatedAt  time.Time `json:"created_at"`
}

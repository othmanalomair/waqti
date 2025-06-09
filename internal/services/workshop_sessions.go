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

// WorkshopSessionService handles session-specific operations
type WorkshopSessionService struct {
	db *database.DB
}

func NewWorkshopSessionService() *WorkshopSessionService {
	return &WorkshopSessionService{
		db: database.Instance,
	}
}

// GetAvailableSessions returns all available sessions for a workshop
func (s *WorkshopSessionService) GetAvailableSessions(workshopID uuid.UUID) ([]models.WorkshopSession, error) {
	query := `
		SELECT 
			ws.id, ws.workshop_id, ws.session_date, ws.end_date, ws.start_time, ws.end_time,
			ws.duration, ws.day_count, ws.timezone, ws.location, ws.location_ar,
			ws.max_attendees, ws.current_attendees, ws.is_completed,
			ws.notes, ws.notes_ar, ws.status, ws.status_ar,
			ws.session_number, ws.run_id, ws.metadata::text,
			ws.session_dates, ws.total_days,
			ws.created_at, ws.updated_at
		FROM workshop_sessions ws
		WHERE ws.workshop_id = $1
		AND ws.status NOT IN ('cancelled', 'completed')
		AND (ws.max_attendees = 0 OR ws.current_attendees < ws.max_attendees)
		AND ws.session_date >= CURRENT_DATE
		ORDER BY ws.session_date, ws.start_time
	`

	rows, err := s.db.Query(query, workshopID)
	if err != nil {
		return nil, fmt.Errorf("failed to query available sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.WorkshopSession
	for rows.Next() {
		var session models.WorkshopSession
		var metadataStr, location, locationAr, notes, notesAr sql.NullString
		var endDate sql.NullTime
		var sessionDatesJSON sql.NullString

		err := rows.Scan(
			&session.ID, &session.WorkshopID, &session.SessionDate, &endDate,
			&session.StartTime, &session.EndTime, &session.Duration, &session.DayCount,
			&session.Timezone, &location, &locationAr,
			&session.MaxAttendees, &session.CurrentAttendees, &session.IsCompleted,
			&notes, &notesAr, &session.Status, &session.StatusAr,
			&session.SessionNumber, &session.RunID, &metadataStr,
			&sessionDatesJSON, &session.TotalDays,
			&session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		// Convert nullable fields to pointers
		if endDate.Valid {
			session.EndDate = &endDate.Time
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

		// Parse metadata if present
		if metadataStr.Valid && metadataStr.String != "" {
			// Handle JSON parsing for metadata
			session.Metadata = make(map[string]interface{})
			// You might want to add JSON parsing here
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetNextAvailableSession returns the next available session for enrollment
func (s *WorkshopSessionService) GetNextAvailableSession(workshopID uuid.UUID) (*models.WorkshopSession, error) {
	query := `
		SELECT 
			ws.id, ws.workshop_id, ws.session_date, ws.end_date, ws.start_time, ws.end_time,
			ws.duration, ws.day_count, ws.timezone, ws.location, ws.location_ar,
			ws.max_attendees, ws.current_attendees, ws.is_completed,
			ws.notes, ws.notes_ar, ws.status, ws.status_ar,
			ws.session_number, ws.run_id, ws.metadata::text,
			ws.created_at, ws.updated_at
		FROM workshop_sessions ws
		WHERE ws.workshop_id = $1
		AND ws.status NOT IN ('cancelled', 'completed', 'full')
		AND (ws.max_attendees = 0 OR ws.current_attendees < ws.max_attendees)
		AND ws.session_date >= CURRENT_DATE
		AND ws.run_id IS NOT NULL
		AND EXISTS (SELECT 1 FROM workshop_runs wr WHERE wr.id = ws.run_id)
		ORDER BY ws.session_date, ws.start_time
		LIMIT 1
	`

	var session models.WorkshopSession
	var metadataStr, location, locationAr, notes, notesAr sql.NullString
	var endDate sql.NullTime

	err := s.db.QueryRow(query, workshopID).Scan(
		&session.ID, &session.WorkshopID, &session.SessionDate, &endDate,
		&session.StartTime, &session.EndTime, &session.Duration, &session.DayCount,
		&session.Timezone, &location, &locationAr,
		&session.MaxAttendees, &session.CurrentAttendees, &session.IsCompleted,
		&notes, &notesAr, &session.Status, &session.StatusAr,
		&session.SessionNumber, &session.RunID, &metadataStr,
		&session.CreatedAt, &session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Try to create a default run for sessions without run_id
		err = s.ensureDefaultWorkshopRun(workshopID)
		if err != nil {
			return nil, fmt.Errorf("failed to create default workshop run: %w", err)
		}
		
		// Retry the query after creating default run
		err = s.db.QueryRow(query, workshopID).Scan(
			&session.ID, &session.WorkshopID, &session.SessionDate, &endDate,
			&session.StartTime, &session.EndTime, &session.Duration, &session.DayCount,
			&session.Timezone, &location, &locationAr,
			&session.MaxAttendees, &session.CurrentAttendees, &session.IsCompleted,
			&notes, &notesAr, &session.Status, &session.StatusAr,
			&session.SessionNumber, &session.RunID, &metadataStr,
			&session.CreatedAt, &session.UpdatedAt,
		)
		
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no available sessions for this workshop")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get next available session after creating default run: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get next available session: %w", err)
	}

	// Convert nullable fields to pointers
	if endDate.Valid {
		session.EndDate = &endDate.Time
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

	return &session, nil
}

// ensureDefaultWorkshopRun creates a default workshop run for sessions without run_id
func (s *WorkshopSessionService) ensureDefaultWorkshopRun(workshopID uuid.UUID) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if there are sessions without run_id
	var sessionCount int
	countQuery := `
		SELECT COUNT(*) 
		FROM workshop_sessions 
		WHERE workshop_id = $1 
		AND (run_id IS NULL OR run_id NOT IN (SELECT id FROM workshop_runs WHERE workshop_id = $1))
	`
	err = tx.QueryRow(countQuery, workshopID).Scan(&sessionCount)
	if err != nil {
		return fmt.Errorf("failed to count sessions without run_id: %w", err)
	}

	if sessionCount == 0 {
		return tx.Commit() // No sessions need fixing
	}

	// Get workshop details for run naming
	var workshopName, workshopNameAr string
	workshopQuery := `SELECT name, COALESCE(title_ar, name) FROM workshops WHERE id = $1`
	err = tx.QueryRow(workshopQuery, workshopID).Scan(&workshopName, &workshopNameAr)
	if err != nil {
		return fmt.Errorf("failed to get workshop details: %w", err)
	}

	// Get date range for the run
	var startDate, endDate time.Time
	dateQuery := `
		SELECT MIN(session_date), COALESCE(MAX(session_date), MIN(session_date))
		FROM workshop_sessions 
		WHERE workshop_id = $1 
		AND (run_id IS NULL OR run_id NOT IN (SELECT id FROM workshop_runs WHERE workshop_id = $1))
	`
	err = tx.QueryRow(dateQuery, workshopID).Scan(&startDate, &endDate)
	if err != nil {
		return fmt.Errorf("failed to get session date range: %w", err)
	}

	// Create default workshop run
	runID := uuid.New()
	runName := fmt.Sprintf("Default Run - %s", workshopName)
	runNameAr := fmt.Sprintf("الدورة الافتراضية - %s", workshopNameAr)

	insertRunQuery := `
		INSERT INTO workshop_runs (id, workshop_id, run_name, run_name_ar, start_date, end_date, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'upcoming')
	`
	_, err = tx.Exec(insertRunQuery, runID, workshopID, runName, runNameAr, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to create default workshop run: %w", err)
	}

	// Update sessions to use the new run_id
	updateSessionsQuery := `
		UPDATE workshop_sessions 
		SET run_id = $1
		WHERE workshop_id = $2 
		AND (run_id IS NULL OR run_id NOT IN (SELECT id FROM workshop_runs WHERE workshop_id = $2))
	`
	_, err = tx.Exec(updateSessionsQuery, runID, workshopID)
	if err != nil {
		return fmt.Errorf("failed to update sessions with run_id: %w", err)
	}

	return tx.Commit()
}

// IncrementSessionAttendance increments the attendance count for a session
func (s *WorkshopSessionService) IncrementSessionAttendance(sessionID uuid.UUID, count int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Lock the session row for update
	query := `
		UPDATE workshop_sessions
		SET current_attendees = current_attendees + $1
		WHERE id = $2
		AND (max_attendees = 0 OR current_attendees + $1 <= max_attendees)
		RETURNING current_attendees, max_attendees
	`

	var currentAttendees, maxAttendees int
	err = tx.QueryRow(query, count, sessionID).Scan(&currentAttendees, &maxAttendees)
	if err == sql.ErrNoRows {
		return fmt.Errorf("session is full or does not exist")
	}
	if err != nil {
		return fmt.Errorf("failed to update attendance: %w", err)
	}

	// The trigger will automatically update the status if needed
	return tx.Commit()
}

// CloneWorkshopRun creates a new run of a workshop with all its sessions
func (s *WorkshopSessionService) CloneWorkshopRun(workshopID uuid.UUID, startDate time.Time, runName string) (*models.WorkshopRun, error) {
	query := `SELECT * FROM clone_workshop_sessions($1, $2, $3)`
	
	var runID uuid.UUID
	var sessionsCreated int
	
	err := s.db.QueryRow(query, workshopID, startDate, runName).Scan(&runID, &sessionsCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to clone workshop sessions: %w", err)
	}

	// Now fetch the created run details
	runQuery := `
		SELECT id, workshop_id, run_name, run_name_ar, start_date, end_date, status, created_at, updated_at
		FROM workshop_runs
		WHERE id = $1
	`
	
	var run models.WorkshopRun
	err = s.db.QueryRow(runQuery, runID).Scan(
		&run.ID, &run.WorkshopID, &run.RunName, &run.RunNameAr,
		&run.StartDate, &run.EndDate, &run.Status,
		&run.CreatedAt, &run.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created run: %w", err)
	}

	return &run, nil
}

// GetWorkshopRuns returns all runs for a workshop
func (s *WorkshopSessionService) GetWorkshopRuns(workshopID uuid.UUID) ([]models.WorkshopRun, error) {
	query := `
		SELECT 
			wr.id, wr.workshop_id, wr.run_name, wr.run_name_ar,
			wr.start_date, wr.end_date, wr.status,
			wr.created_at, wr.updated_at
		FROM workshop_runs wr
		WHERE wr.workshop_id = $1
		ORDER BY wr.start_date DESC
	`

	rows, err := s.db.Query(query, workshopID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workshop runs: %w", err)
	}
	defer rows.Close()

	var runs []models.WorkshopRun
	for rows.Next() {
		var run models.WorkshopRun
		err := rows.Scan(
			&run.ID, &run.WorkshopID, &run.RunName, &run.RunNameAr,
			&run.StartDate, &run.EndDate, &run.Status,
			&run.CreatedAt, &run.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan run: %w", err)
		}
		runs = append(runs, run)
	}

	return runs, nil
}

// GetSessionsByRunID returns all sessions for a specific run
func (s *WorkshopSessionService) GetSessionsByRunID(runID uuid.UUID) ([]models.WorkshopSession, error) {
	query := `
		SELECT 
			ws.id, ws.workshop_id, ws.session_date, ws.end_date, ws.start_time, ws.end_time,
			ws.duration, ws.day_count, ws.timezone, ws.location, ws.location_ar,
			ws.max_attendees, ws.current_attendees, ws.is_completed,
			ws.notes, ws.notes_ar, ws.status, ws.status_ar,
			ws.session_number, ws.run_id, ws.metadata::text,
			ws.created_at, ws.updated_at
		FROM workshop_sessions ws
		WHERE ws.run_id = $1
		ORDER BY ws.session_date, ws.start_time
	`

	rows, err := s.db.Query(query, runID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions by run: %w", err)
	}
	defer rows.Close()

	var sessions []models.WorkshopSession
	for rows.Next() {
		var session models.WorkshopSession
		var metadataStr, location, locationAr, notes, notesAr sql.NullString
		var endDate sql.NullTime
		var sessionDatesJSON sql.NullString

		err := rows.Scan(
			&session.ID, &session.WorkshopID, &session.SessionDate, &endDate,
			&session.StartTime, &session.EndTime, &session.Duration, &session.DayCount,
			&session.Timezone, &location, &locationAr,
			&session.MaxAttendees, &session.CurrentAttendees, &session.IsCompleted,
			&notes, &notesAr, &session.Status, &session.StatusAr,
			&session.SessionNumber, &session.RunID, &metadataStr,
			&sessionDatesJSON, &session.TotalDays,
			&session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		// Convert nullable fields to pointers
		if endDate.Valid {
			session.EndDate = &endDate.Time
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

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetSessionAvailability returns availability information for all sessions
func (s *WorkshopSessionService) GetSessionAvailability(creatorID uuid.UUID) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			session_id, workshop_id, workshop_name, workshop_name_ar,
			session_date, start_time, end_time,
			max_attendees, current_attendees, available_seats,
			calculated_status, run_id, run_name,
			creator_name, creator_username
		FROM session_availability
		WHERE creator_username = (SELECT username FROM creators WHERE id = $1)
		ORDER BY session_date, start_time
	`

	rows, err := s.db.Query(query, creatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query session availability: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var (
			sessionID, workshopID       uuid.UUID
			workshopName, workshopNameAr string
			sessionDate                 time.Time
			startTime, endTime          string
			maxAttendees, currentAttendees, availableSeats int
			status                      string
			runID                       *uuid.UUID
			runName                     sql.NullString
			creatorName, creatorUsername string
		)

		err := rows.Scan(
			&sessionID, &workshopID, &workshopName, &workshopNameAr,
			&sessionDate, &startTime, &endTime,
			&maxAttendees, &currentAttendees, &availableSeats,
			&status, &runID, &runName,
			&creatorName, &creatorUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan availability: %w", err)
		}

		result := map[string]interface{}{
			"session_id":        sessionID,
			"workshop_id":       workshopID,
			"workshop_name":     workshopName,
			"workshop_name_ar":  workshopNameAr,
			"session_date":      sessionDate,
			"start_time":        startTime,
			"end_time":          endTime,
			"max_attendees":     maxAttendees,
			"current_attendees": currentAttendees,
			"available_seats":   availableSeats,
			"status":            status,
			"creator_name":      creatorName,
			"creator_username":  creatorUsername,
		}

		if runID != nil {
			result["run_id"] = *runID
		}
		if runName.Valid {
			result["run_name"] = runName.String
		}

		results = append(results, result)
	}

	return results, nil
}
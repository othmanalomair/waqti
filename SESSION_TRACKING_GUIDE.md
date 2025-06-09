# Session-Based Workshop Tracking Implementation Guide

## Overview

This guide explains the new session-based tracking system for workshops that solves the following problems:
- Difficulty tracking if sessions are full
- Manual deletion of old sessions when repeating workshops
- No clear representation of multi-day workshops
- No link between enrollments and specific sessions

## Key Changes

### 1. Enhanced Workshop Sessions
Each workshop session now includes:
- **Status tracking**: upcoming, active, full, completed, cancelled
- **Session numbering**: For multi-day workshops (Day 1, Day 2, etc.)
- **Run ID**: Groups sessions that belong to the same "run" of a workshop
- **Metadata**: Flexible JSONB field for additional information

### 2. Workshop Runs
A new concept that groups multiple sessions of the same workshop:
```sql
workshop_runs
├── id (UUID)
├── workshop_id (references workshops)
├── run_name (e.g., "July 2025 Batch")
├── start_date
├── end_date
└── status
```

### 3. Session-Specific Enrollments
Enrollments now link to specific sessions:
- `enrollments.session_id` → `workshop_sessions.id`
- `order_items.session_id` → `workshop_sessions.id`
- `order_items.run_id` → `workshop_runs.id`

## Migration Steps

1. **Run the migration script**:
   ```bash
   psql -d waqti -f migrations/001_session_tracking.sql
   ```

2. **Update your Go code** to use the new models (already done in `/internal/models/`)

## Usage Examples

### Creating a New Workshop Run (Repeating a Workshop)

```sql
-- Clone an existing workshop for a new date
SELECT * FROM clone_workshop_sessions(
    'workshop-uuid-here'::uuid,  -- workshop_id
    '2025-07-01'::date,          -- new start date
    'Summer 2025 Session'        -- optional run name
);
```

### Checking Session Availability

```sql
-- Use the session_availability view
SELECT * FROM session_availability
WHERE workshop_id = 'workshop-uuid-here'
AND calculated_status != 'full'
ORDER BY session_date;
```

### Creating an Enrollment for a Specific Session

```go
// In your enrollment service
enrollment := &models.Enrollment{
    WorkshopID: workshopID,
    SessionID:  &sessionID,  // Now links to specific session
    OrderID:    &orderID,
    // ... other fields
}
```

### Updating Session Attendance

```sql
-- The trigger automatically updates status when attendance changes
UPDATE workshop_sessions
SET current_attendees = current_attendees + 1
WHERE id = 'session-uuid-here';
-- Status automatically changes to 'full' if max_attendees reached
```

## Implementation in Services

### Workshop Service Updates

```go
// Example: Get available sessions for a workshop
func (s *WorkshopService) GetAvailableSessions(workshopID uuid.UUID) ([]models.WorkshopSession, error) {
    query := `
        SELECT * FROM workshop_sessions
        WHERE workshop_id = $1
        AND status != 'cancelled'
        AND (max_attendees = 0 OR current_attendees < max_attendees)
        AND session_date >= CURRENT_DATE
        ORDER BY session_date, start_time
    `
    // ... implementation
}

// Example: Clone workshop for new run
func (s *WorkshopService) CloneWorkshopRun(workshopID uuid.UUID, startDate time.Time, runName string) (*models.WorkshopRun, error) {
    query := `SELECT * FROM clone_workshop_sessions($1, $2, $3)`
    // ... implementation
}
```

### Order Service Updates

```go
// When creating an order, specify the session
func (s *OrderService) CreateOrder(req models.CreateOrderRequest) (*models.Order, error) {
    // For each item, find the next available session
    for _, item := range req.Items {
        session, err := s.findNextAvailableSession(item.WorkshopID)
        if err != nil {
            return nil, err
        }
        
        orderItem := models.OrderItem{
            WorkshopID: item.WorkshopID,
            SessionID:  &session.ID,
            RunID:      session.RunID,
            // ... other fields
        }
    }
}
```

## UI/Template Updates Needed

### 1. Store Page (`store.templ`)
- Show session dates and availability
- Allow selection of specific sessions if multiple are available
- Display "Full" badge for full sessions

### 2. Add Workshop Page (`add_workshop.templ`)
- Add option to set max attendees per session
- Show session capacity configuration

### 3. Workshop Management
- Add "Repeat Workshop" button that uses `clone_workshop_sessions`
- Show workshop runs with their status
- Allow viewing attendees by session

### 4. Order Tracking (`order_tracking.templ`)
- Display which session the order is for
- Show session date alongside workshop name

## Benefits

1. **Automatic Capacity Management**: Sessions automatically marked as full
2. **Easy Workshop Repetition**: One function call to create a new run
3. **Better Tracking**: Know exactly which session each student is in
4. **Historical Data**: Keep old sessions for records while creating new ones
5. **Flexible Scheduling**: Support for single-day, multi-day, and recurring workshops

## Views Available

1. **session_availability**: Shows all sessions with calculated availability
2. **workshop_run_summary**: Summary of all runs for a workshop

## Status Management

Sessions automatically transition through statuses:
- `upcoming` → `active` (when first enrollment)
- `active` → `full` (when max capacity reached)
- Any → `completed` (when date passes)
- Manual → `cancelled` (by creator action)

## Next Steps

1. Update the workshop service to use session-specific logic
2. Modify the enrollment flow to select specific sessions
3. Update the UI to show session availability
4. Add "Repeat Workshop" functionality to the dashboard
5. Update order tracking to show session information

## Example Queries

### Get workshops with available sessions
```sql
SELECT DISTINCT w.*, 
       COUNT(DISTINCT s.id) as available_sessions
FROM workshops w
JOIN workshop_sessions s ON w.id = s.workshop_id
WHERE w.is_active = true
AND s.status NOT IN ('full', 'cancelled', 'completed')
AND s.session_date >= CURRENT_DATE
GROUP BY w.id;
```

### Get enrollment count by session
```sql
SELECT ws.*, 
       COUNT(e.id) as enrollment_count
FROM workshop_sessions ws
LEFT JOIN enrollments e ON ws.id = e.session_id
WHERE ws.workshop_id = 'workshop-uuid'
GROUP BY ws.id
ORDER BY ws.session_date;
```
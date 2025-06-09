# Fix for GetNextAvailableSession Foreign Key Constraint Issue

## Problem Description

The `GetNextAvailableSession` method in `WorkshopSessionService` was causing foreign key constraint violations when creating orders. The issue occurred because:

1. **Missing or Invalid run_id Values**: Some workshop sessions had NULL or invalid `run_id` values that didn't reference existing entries in the `workshop_runs` table.

2. **Schema Evolution**: The database schema was evolved through migrations to support workshop runs, but existing sessions weren't properly migrated to have valid `run_id` references.

3. **Order Creation Failure**: When creating orders, the system tried to insert order items with invalid `run_id` values, causing foreign key constraint violations.

## Root Cause Analysis

### Database State Issues
- Workshop sessions existed without valid `run_id` references
- Some sessions had `run_id` values pointing to non-existent workshop runs
- Order items inherited invalid `run_id` values from sessions

### Code Logic Issues
- `GetNextAvailableSession` didn't validate `run_id` existence
- Order creation didn't handle NULL/invalid `run_id` gracefully
- No fallback mechanism for creating default workshop runs

## Solution Implemented

### 1. Database Fixes (`migrations/004_fix_run_id_foreign_keys.sql`)

```sql
-- Create default workshop runs for orphaned sessions
INSERT INTO workshop_runs (id, workshop_id, run_name, run_name_ar, start_date, end_date, status)
SELECT 
    gen_random_uuid() as id,
    ws.workshop_id,
    'Default Run - ' || w.name as run_name,
    'الدورة الافتراضية - ' || COALESCE(w.title_ar, w.name) as run_name_ar,
    MIN(ws.session_date) as start_date,
    COALESCE(MAX(ws.session_date), MIN(ws.session_date)) as end_date,
    'upcoming' as status
FROM workshop_sessions ws
JOIN workshops w ON ws.workshop_id = w.id
WHERE ws.run_id IS NULL 
   OR ws.run_id NOT IN (SELECT id FROM workshop_runs WHERE workshop_id = ws.workshop_id)
GROUP BY ws.workshop_id, w.name, w.title_ar;

-- Update sessions to reference valid runs
UPDATE workshop_sessions 
SET run_id = (SELECT wr.id FROM workshop_runs wr WHERE wr.workshop_id = workshop_sessions.workshop_id AND wr.run_name LIKE 'Default Run -%' LIMIT 1)
WHERE run_id IS NULL OR run_id NOT IN (SELECT id FROM workshop_runs WHERE workshop_id = workshop_sessions.workshop_id);
```

### 2. Code Improvements

#### Enhanced `GetNextAvailableSession` Method
```go
// Added validation for run_id existence
query := `
    SELECT ... FROM workshop_sessions ws
    WHERE ws.workshop_id = $1
    AND ws.run_id IS NOT NULL
    AND EXISTS (SELECT 1 FROM workshop_runs wr WHERE wr.id = ws.run_id)
    ...
`

// Added fallback to create default run if needed
if err == sql.ErrNoRows {
    err = s.ensureDefaultWorkshopRun(workshopID)
    // Retry query after creating default run
}
```

#### New Helper Method `ensureDefaultWorkshopRun`
```go
func (s *WorkshopSessionService) ensureDefaultWorkshopRun(workshopID uuid.UUID) error {
    // Creates a default workshop run for sessions without valid run_id
    // Updates orphaned sessions to reference the new run
}
```

#### Improved Order Creation
```go
// Added run_id retrieval for specific sessions
if itemReq.SessionID != nil {
    sessionID = itemReq.SessionID
    // Get the run_id for this session
    var sessionRunID *uuid.UUID
    runQuery := `SELECT run_id FROM workshop_sessions WHERE id = $1`
    err := database.Instance.QueryRow(runQuery, sessionID).Scan(&sessionRunID)
    if err == nil {
        runID = sessionRunID
    }
}
```

### 3. Data Validation Trigger
```sql
-- Added trigger to prevent future run_id issues
CREATE OR REPLACE FUNCTION validate_session_run_id() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.run_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM workshop_runs wr 
            WHERE wr.id = NEW.run_id 
            AND wr.workshop_id = NEW.workshop_id
        ) THEN
            RAISE EXCEPTION 'run_id % does not belong to workshop %', NEW.run_id, NEW.workshop_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

## How to Apply the Fix

### 1. Apply Database Migration
```bash
# Apply the SQL migration to fix existing data
psql -d your_database -f migrations/004_fix_run_id_foreign_keys.sql
```

### 2. Deploy Code Changes
The code changes are already in place in:
- `/home/most3mr/mytrash/waqti/internal/services/workshop_sessions.go`
- `/home/most3mr/mytrash/waqti/internal/services/order.go`

### 3. Verify the Fix
```bash
# Run the application and test order creation
make run
```

## Prevention Measures

1. **Data Validation**: The new trigger prevents inserting sessions with invalid `run_id`
2. **Automatic Fallback**: `ensureDefaultWorkshopRun` creates runs when needed
3. **Strict Queries**: `GetNextAvailableSession` only returns sessions with valid `run_id`
4. **Error Handling**: Better error messages for debugging

## Testing Recommendations

1. **Test Order Creation**: Create orders for workshops to ensure no foreign key errors
2. **Test Session Retrieval**: Verify `GetNextAvailableSession` returns valid sessions
3. **Test Edge Cases**: Try creating orders for workshops without existing sessions
4. **Monitor Logs**: Check for any remaining foreign key constraint warnings

## Files Modified

- `internal/services/workshop_sessions.go` - Enhanced session retrieval and run creation
- `internal/services/order.go` - Improved order creation with run_id handling
- `migrations/004_fix_run_id_foreign_keys.sql` - Database migration to fix existing data
- `fix_run_id_constraints.sql` - Standalone fix script (alternative approach)

This fix ensures that the foreign key constraint issue with `run_id` is resolved both for existing data and future operations.
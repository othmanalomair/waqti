-- Migration to store session dates as JSON array
-- This allows storing multiple non-consecutive dates in a single session row
-- =====================================================================

-- 1. Add new columns for JSON-based date storage
ALTER TABLE workshop_sessions 
ADD COLUMN session_dates JSONB,  -- Array of dates for this session
ADD COLUMN total_days INTEGER DEFAULT 1;  -- Total number of days (may have gaps)

-- 2. Migrate existing data to new format
UPDATE workshop_sessions 
SET session_dates = CASE 
    WHEN end_date IS NOT NULL THEN 
        -- Multi-day sessions: create array of all dates from start to end
        (SELECT jsonb_agg(date_val::date)
         FROM generate_series(session_date::date, end_date::date, '1 day'::interval) AS date_val)
    ELSE 
        -- Single day sessions: create array with just one date
        jsonb_build_array(session_date::date)
    END,
total_days = CASE 
    WHEN day_count > 0 THEN day_count 
    ELSE 1 
    END;

-- 3. Update the session_availability view to work with new structure
DROP VIEW IF EXISTS session_availability;
CREATE VIEW session_availability AS
SELECT 
    ws.id as session_id,
    ws.workshop_id,
    w.name as workshop_name,
    COALESCE(w.title_ar, w.name) as workshop_name_ar,
    ws.session_date,  -- Keep for compatibility
    ws.session_dates, -- New JSON array of dates
    ws.total_days,    -- Total days including gaps
    ws.start_time,
    ws.end_time,
    ws.max_attendees,
    ws.current_attendees,
    GREATEST(0, ws.max_attendees - ws.current_attendees) as available_seats,
    CASE 
        WHEN ws.status = 'cancelled' THEN 'cancelled'
        WHEN ws.session_date < CURRENT_DATE THEN 'completed'
        WHEN ws.current_attendees >= ws.max_attendees AND ws.max_attendees > 0 THEN 'full'
        ELSE COALESCE(ws.status, 'upcoming')
    END as calculated_status,
    ws.run_id,
    wr.run_name,
    c.name as creator_name,
    c.username as creator_username
FROM workshop_sessions ws
JOIN workshops w ON ws.workshop_id = w.id
JOIN creators c ON w.creator_id = c.id
LEFT JOIN workshop_runs wr ON ws.run_id = wr.id
WHERE w.is_active = true;

-- 4. Create helper function to get formatted date range for display
CREATE OR REPLACE FUNCTION get_session_date_display(session_dates JSONB, lang TEXT DEFAULT 'en')
RETURNS TEXT AS $$
DECLARE
    dates_array DATE[];
    first_date DATE;
    last_date DATE;
    total_count INTEGER;
BEGIN
    -- Convert JSONB array to PostgreSQL array
    SELECT array_agg(value::date ORDER BY value::date)
    INTO dates_array
    FROM jsonb_array_elements_text(session_dates);
    
    total_count := array_length(dates_array, 1);
    
    IF total_count IS NULL OR total_count = 0 THEN
        RETURN '';
    ELSIF total_count = 1 THEN
        -- Single day
        IF lang = 'ar' THEN
            RETURN to_char(dates_array[1], 'DD Mon YYYY');
        ELSE
            RETURN to_char(dates_array[1], 'Mon DD, YYYY');
        END IF;
    ELSE
        -- Multiple days
        first_date := dates_array[1];
        last_date := dates_array[total_count];
        
        IF lang = 'ar' THEN
            RETURN format('%s أيام: %s - %s', 
                total_count,
                to_char(first_date, 'DD Mon'),
                to_char(last_date, 'DD Mon YYYY')
            );
        ELSE
            RETURN format('%s days: %s - %s', 
                total_count,
                to_char(first_date, 'Mon DD'),
                to_char(last_date, 'Mon DD, YYYY')
            );
        END IF;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- 5. Add index for better performance on JSON queries
CREATE INDEX IF NOT EXISTS idx_workshop_sessions_dates_gin ON workshop_sessions USING GIN (session_dates);

-- 6. Test the migration
SELECT 'Testing new session dates structure:' as test_message;
SELECT 
    id,
    session_date as old_start,
    end_date as old_end,
    day_count as old_day_count,
    session_dates as new_dates_json,
    total_days as new_total_days,
    get_session_date_display(session_dates, 'en') as display_en,
    get_session_date_display(session_dates, 'ar') as display_ar
FROM workshop_sessions 
LIMIT 5;
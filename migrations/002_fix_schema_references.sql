-- Fix schema references in views and functions
-- ===============================================

-- 1. Fix session_availability view to use correct column names
CREATE OR REPLACE VIEW session_availability AS
SELECT 
    ws.id as session_id,
    ws.workshop_id,
    w.name as workshop_name,
    COALESCE(w.title_ar, w.name) as workshop_name_ar,
    ws.session_date,
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

-- 2. Fix workshop_run_summary view to use correct column names
CREATE OR REPLACE VIEW workshop_run_summary AS
SELECT 
    wr.id as run_id,
    wr.workshop_id,
    w.name as workshop_name,
    COALESCE(w.title_ar, w.name) as workshop_name_ar,
    wr.run_name,
    wr.start_date,
    wr.end_date,
    COUNT(ws.id) as total_sessions,
    SUM(ws.max_attendees) as total_capacity,
    SUM(ws.current_attendees) as total_enrolled,
    SUM(CASE WHEN ws.status = 'full' THEN 1 ELSE 0 END) as full_sessions,
    wr.status,
    w.creator_id
FROM workshop_runs wr
JOIN workshops w ON wr.workshop_id = w.id
LEFT JOIN workshop_sessions ws ON ws.run_id = wr.id
GROUP BY wr.id, w.id;

-- 3. Add some initial test data to verify the system works
-- Insert some sample sessions for existing workshops if they don't have any
INSERT INTO workshop_sessions (
    id, workshop_id, session_date, start_time, end_time, duration, 
    timezone, max_attendees, current_attendees, status, status_ar, 
    session_number, created_at, updated_at
)
SELECT 
    gen_random_uuid() as id,
    w.id as workshop_id,
    CURRENT_DATE + INTERVAL '7 days' as session_date,
    '10:00:00' as start_time,
    '12:00:00' as end_time,
    2.0 as duration,
    'Asia/Kuwait' as timezone,
    CASE WHEN w.max_students > 0 THEN w.max_students ELSE 20 END as max_attendees,
    0 as current_attendees,
    'upcoming' as status,
    'قادم' as status_ar,
    1 as session_number,
    NOW() as created_at,
    NOW() as updated_at
FROM workshops w
WHERE w.is_active = true 
AND NOT EXISTS (
    SELECT 1 FROM workshop_sessions ws WHERE ws.workshop_id = w.id
)
LIMIT 5; -- Only add sessions for first 5 workshops without sessions

-- 4. Test the views work correctly
SELECT 'Testing session_availability view:' as test_name;
SELECT COUNT(*) as total_available_sessions FROM session_availability;

SELECT 'Testing workshop_run_summary view:' as test_name;
SELECT COUNT(*) as total_runs FROM workshop_run_summary;
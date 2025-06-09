-- Database Redesign for Session-Based Workshop Management
-- ========================================================

-- 1. Enhanced workshop_sessions table
ALTER TABLE workshop_sessions
ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'active', 'full', 'completed', 'cancelled')),
ADD COLUMN IF NOT EXISTS status_ar VARCHAR(50) DEFAULT 'قادم',
ADD COLUMN IF NOT EXISTS session_number INTEGER DEFAULT 1, -- For multi-day workshops (Day 1, Day 2, etc.)
ADD COLUMN IF NOT EXISTS parent_run_id UUID, -- Groups sessions of the same workshop run
ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}'; -- Flexible storage for additional session info

-- Add index for better performance
CREATE INDEX IF NOT EXISTS idx_workshop_sessions_status ON workshop_sessions(status);
CREATE INDEX IF NOT EXISTS idx_workshop_sessions_parent_run ON workshop_sessions(parent_run_id);
CREATE INDEX IF NOT EXISTS idx_workshop_sessions_date ON workshop_sessions(session_date);

-- 2. Create workshop_runs table to group multi-day sessions
CREATE TABLE IF NOT EXISTS workshop_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workshop_id UUID NOT NULL REFERENCES workshops(id) ON DELETE CASCADE,
    run_name VARCHAR(255), -- e.g., "July 2025 Batch", "Summer Session"
    run_name_ar VARCHAR(255),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_sessions INTEGER NOT NULL DEFAULT 1,
    max_attendees INTEGER NOT NULL DEFAULT 0,
    current_attendees INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'active', 'full', 'completed', 'cancelled')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Link enrollments to specific sessions
ALTER TABLE enrollments
ADD CONSTRAINT fk_enrollment_session FOREIGN KEY (session_id) REFERENCES workshop_sessions(id);

-- 4. Update order_items to track session_id
ALTER TABLE order_items
ADD COLUMN IF NOT EXISTS session_id UUID REFERENCES workshop_sessions(id),
ADD COLUMN IF NOT EXISTS workshop_run_id UUID REFERENCES workshop_runs(id);

-- 5. Create a view for session availability
CREATE OR REPLACE VIEW session_availability AS
SELECT 
    ws.id as session_id,
    ws.workshop_id,
    w.name as workshop_name,
    w.name_ar as workshop_name_ar,
    ws.session_date,
    ws.start_time,
    ws.end_time,
    ws.max_attendees,
    ws.current_attendees,
    ws.max_attendees - ws.current_attendees as available_seats,
    CASE 
        WHEN ws.current_attendees >= ws.max_attendees THEN 'full'
        WHEN ws.session_date < CURRENT_DATE THEN 'completed'
        ELSE ws.status
    END as calculated_status,
    ws.parent_run_id,
    wr.run_name
FROM workshop_sessions ws
JOIN workshops w ON ws.workshop_id = w.id
LEFT JOIN workshop_runs wr ON ws.parent_run_id = wr.id;

-- 6. Function to automatically update session status
CREATE OR REPLACE FUNCTION update_session_status() RETURNS TRIGGER AS $$
BEGIN
    -- Update status to 'full' when max capacity is reached
    IF NEW.current_attendees >= NEW.max_attendees AND NEW.max_attendees > 0 THEN
        NEW.status = 'full';
        NEW.status_ar = 'ممتلئ';
    END IF;
    
    -- Update updated_at timestamp
    NEW.updated_at = NOW();
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_session_status
BEFORE UPDATE ON workshop_sessions
FOR EACH ROW
EXECUTE FUNCTION update_session_status();

-- 7. Function to clone workshop sessions for repeat runs
CREATE OR REPLACE FUNCTION clone_workshop_sessions(
    p_workshop_id UUID,
    p_new_start_date DATE,
    p_run_name VARCHAR DEFAULT NULL
) RETURNS TABLE(new_run_id UUID, sessions_created INTEGER) AS $$
DECLARE
    v_run_id UUID;
    v_date_offset INTERVAL;
    v_sessions_count INTEGER;
    v_first_session_date DATE;
BEGIN
    -- Generate new run ID
    v_run_id := gen_random_uuid();
    
    -- Get the date offset from the original sessions
    SELECT MIN(session_date) INTO v_first_session_date
    FROM workshop_sessions
    WHERE workshop_id = p_workshop_id
    AND parent_run_id IS NULL OR parent_run_id = (
        SELECT id FROM workshop_runs WHERE workshop_id = p_workshop_id ORDER BY created_at DESC LIMIT 1
    );
    
    v_date_offset := p_new_start_date - v_first_session_date;
    
    -- Create workshop run record
    INSERT INTO workshop_runs (id, workshop_id, run_name, start_date, end_date)
    SELECT 
        v_run_id,
        p_workshop_id,
        COALESCE(p_run_name, 'Session ' || TO_CHAR(p_new_start_date, 'Month YYYY')),
        MIN(session_date + v_date_offset),
        MAX(session_date + v_date_offset)
    FROM workshop_sessions
    WHERE workshop_id = p_workshop_id;
    
    -- Clone sessions with new dates
    INSERT INTO workshop_sessions (
        workshop_id, session_date, start_time, end_time, duration,
        timezone, location, location_ar, max_attendees, session_number,
        parent_run_id, metadata
    )
    SELECT 
        workshop_id,
        session_date + v_date_offset,
        start_time,
        end_time,
        duration,
        timezone,
        location,
        location_ar,
        max_attendees,
        session_number,
        v_run_id,
        metadata
    FROM workshop_sessions
    WHERE workshop_id = p_workshop_id
    AND (parent_run_id IS NULL OR parent_run_id = (
        SELECT id FROM workshop_runs WHERE workshop_id = p_workshop_id ORDER BY created_at DESC LIMIT 1
    ))
    ORDER BY session_date, start_time;
    
    GET DIAGNOSTICS v_sessions_count = ROW_COUNT;
    
    RETURN QUERY SELECT v_run_id, v_sessions_count;
END;
$$ LANGUAGE plpgsql;

-- 8. Update trigger for workshop_runs attendance
CREATE OR REPLACE FUNCTION update_workshop_run_attendance() RETURNS TRIGGER AS $$
BEGIN
    -- Update the workshop run's total attendance
    UPDATE workshop_runs
    SET current_attendees = (
        SELECT SUM(current_attendees) 
        FROM workshop_sessions 
        WHERE parent_run_id = NEW.parent_run_id
    )
    WHERE id = NEW.parent_run_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_run_attendance
AFTER UPDATE OF current_attendees ON workshop_sessions
FOR EACH ROW
WHEN (NEW.parent_run_id IS NOT NULL)
EXECUTE FUNCTION update_workshop_run_attendance();

-- 9. Helper view for creators to see all their workshop runs
CREATE OR REPLACE VIEW creator_workshop_runs AS
SELECT 
    wr.*,
    w.name as workshop_name,
    w.name_ar as workshop_name_ar,
    w.creator_id,
    COUNT(DISTINCT ws.id) as total_sessions,
    SUM(ws.max_attendees) as total_capacity,
    SUM(ws.current_attendees) as total_enrolled,
    MIN(ws.session_date) as first_session,
    MAX(ws.session_date) as last_session
FROM workshop_runs wr
JOIN workshops w ON wr.workshop_id = w.id
JOIN workshop_sessions ws ON ws.parent_run_id = wr.id
GROUP BY wr.id, w.id;

-- 10. Migration to link existing enrollments to sessions
-- This attempts to match enrollments to sessions based on workshop and date proximity
UPDATE enrollments e
SET session_id = (
    SELECT ws.id
    FROM workshop_sessions ws
    WHERE ws.workshop_id = e.workshop_id
    AND ws.session_date >= e.enrollment_date::date
    ORDER BY ws.session_date, ws.start_time
    LIMIT 1
)
WHERE e.session_id IS NULL;

-- 11. Update session attendance counts based on enrollments
UPDATE workshop_sessions ws
SET current_attendees = (
    SELECT COUNT(*)
    FROM enrollments e
    WHERE e.session_id = ws.id
    AND e.status IN ('successful', 'pending')
);
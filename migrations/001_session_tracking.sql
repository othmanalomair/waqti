-- Migration: Enhanced Session-Based Enrollment Tracking
-- =====================================================
-- This migration enhances the existing workshop_sessions table
-- to better track enrollment capacity and session status

-- 1. Add status tracking to workshop_sessions
ALTER TABLE workshop_sessions
ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'active', 'full', 'completed', 'cancelled')),
ADD COLUMN IF NOT EXISTS status_ar VARCHAR(50) DEFAULT 'قادم',
ADD COLUMN IF NOT EXISTS session_number INTEGER DEFAULT 1, -- For multi-day workshops (Day 1, Day 2, etc.)
ADD COLUMN IF NOT EXISTS run_id UUID, -- Groups sessions that belong to the same "run" of a workshop
ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}'; -- Flexible storage for additional info

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_workshop_sessions_status ON workshop_sessions(status);
CREATE INDEX IF NOT EXISTS idx_workshop_sessions_run_id ON workshop_sessions(run_id);
CREATE INDEX IF NOT EXISTS idx_workshop_sessions_date ON workshop_sessions(session_date);
CREATE INDEX IF NOT EXISTS idx_workshop_sessions_workshop ON workshop_sessions(workshop_id);

-- 2. Create workshop_runs table to group multi-day sessions
CREATE TABLE IF NOT EXISTS workshop_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workshop_id UUID NOT NULL REFERENCES workshops(id) ON DELETE CASCADE,
    run_name VARCHAR(255), -- e.g., "July 2025 Batch", "Summer Session"
    run_name_ar VARCHAR(255),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'upcoming' CHECK (status IN ('upcoming', 'active', 'full', 'completed', 'cancelled')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_workshop_runs_workshop ON workshop_runs(workshop_id);
CREATE INDEX IF NOT EXISTS idx_workshop_runs_dates ON workshop_runs(start_date, end_date);

-- 3. Fix enrollments table to properly link to sessions
-- First ensure we have the session_id column
ALTER TABLE enrollments
ADD COLUMN IF NOT EXISTS session_id UUID REFERENCES workshop_sessions(id);

-- Add index for performance
CREATE INDEX IF NOT EXISTS idx_enrollments_session ON enrollments(session_id);

-- 4. Update order_items to track session information
ALTER TABLE order_items
ADD COLUMN IF NOT EXISTS session_id UUID REFERENCES workshop_sessions(id),
ADD COLUMN IF NOT EXISTS run_id UUID REFERENCES workshop_runs(id);

CREATE INDEX IF NOT EXISTS idx_order_items_session ON order_items(session_id);

-- 5. Create a view for easy session availability checking
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

-- 6. Function to update session status automatically
CREATE OR REPLACE FUNCTION update_session_status() 
RETURNS TRIGGER AS $$
BEGIN
    -- Update status based on attendance
    IF NEW.max_attendees > 0 AND NEW.current_attendees >= NEW.max_attendees THEN
        NEW.status = 'full';
        NEW.status_ar = 'ممتلئ';
    ELSIF NEW.current_attendees > 0 AND NEW.status = 'upcoming' THEN
        NEW.status = 'active';
        NEW.status_ar = 'نشط';
    END IF;
    
    -- Update timestamp
    NEW.updated_at = NOW();
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger if it doesn't exist
DROP TRIGGER IF EXISTS trigger_update_session_status ON workshop_sessions;
CREATE TRIGGER trigger_update_session_status
BEFORE UPDATE OF current_attendees ON workshop_sessions
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
    v_last_session_date DATE;
BEGIN
    -- Generate new run ID
    v_run_id := gen_random_uuid();
    
    -- Get the date range of existing sessions
    SELECT MIN(session_date), MAX(session_date) 
    INTO v_first_session_date, v_last_session_date
    FROM workshop_sessions
    WHERE workshop_id = p_workshop_id
    AND (run_id IS NULL OR run_id = (
        SELECT id FROM workshop_runs 
        WHERE workshop_id = p_workshop_id 
        ORDER BY created_at DESC LIMIT 1
    ));
    
    -- Calculate date offset
    v_date_offset := p_new_start_date - v_first_session_date;
    
    -- Create workshop run record
    INSERT INTO workshop_runs (id, workshop_id, run_name, start_date, end_date)
    VALUES (
        v_run_id,
        p_workshop_id,
        COALESCE(p_run_name, 'Session ' || TO_CHAR(p_new_start_date, 'Month YYYY')),
        p_new_start_date,
        v_last_session_date + v_date_offset
    );
    
    -- Clone sessions with new dates
    INSERT INTO workshop_sessions (
        workshop_id, session_date, start_time, end_time, duration,
        timezone, location, location_ar, max_attendees, session_number,
        run_id, metadata, status, current_attendees
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
        ROW_NUMBER() OVER (ORDER BY session_date, start_time),
        v_run_id,
        metadata,
        'upcoming',
        0
    FROM workshop_sessions
    WHERE workshop_id = p_workshop_id
    AND (run_id IS NULL OR run_id = (
        SELECT id FROM workshop_runs 
        WHERE workshop_id = p_workshop_id 
        ORDER BY created_at DESC LIMIT 1
    ))
    ORDER BY session_date, start_time;
    
    GET DIAGNOSTICS v_sessions_count = ROW_COUNT;
    
    RETURN QUERY SELECT v_run_id, v_sessions_count;
END;
$$ LANGUAGE plpgsql;

-- 8. Update existing data to use the new structure
-- Set session numbers for existing sessions
WITH numbered_sessions AS (
    SELECT 
        id,
        ROW_NUMBER() OVER (PARTITION BY workshop_id ORDER BY session_date, start_time) as session_num
    FROM workshop_sessions
    WHERE session_number IS NULL OR session_number = 0
)
UPDATE workshop_sessions ws
SET session_number = ns.session_num
FROM numbered_sessions ns
WHERE ws.id = ns.id;

-- Update session status based on dates
UPDATE workshop_sessions
SET status = CASE 
    WHEN session_date < CURRENT_DATE THEN 'completed'
    WHEN current_attendees >= max_attendees AND max_attendees > 0 THEN 'full'
    ELSE 'upcoming'
END,
status_ar = CASE 
    WHEN session_date < CURRENT_DATE THEN 'مكتمل'
    WHEN current_attendees >= max_attendees AND max_attendees > 0 THEN 'ممتلئ'
    ELSE 'قادم'
END
WHERE status IS NULL;

-- 9. Create a helper view for creators to see workshop run summaries
CREATE OR REPLACE VIEW workshop_run_summary AS
SELECT 
    wr.id as run_id,
    wr.workshop_id,
    w.name as workshop_name,
    w.name_ar as workshop_name_ar,
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

-- 10. Add comments for documentation
COMMENT ON TABLE workshop_runs IS 'Groups multiple sessions of the same workshop run together';
COMMENT ON COLUMN workshop_sessions.run_id IS 'Links session to a specific run/batch of the workshop';
COMMENT ON COLUMN workshop_sessions.session_number IS 'Sequential number for multi-day workshops (Day 1, 2, etc.)';
COMMENT ON COLUMN workshop_sessions.status IS 'Current status of the session';
COMMENT ON FUNCTION clone_workshop_sessions IS 'Creates a new run of a workshop by cloning existing sessions with new dates';
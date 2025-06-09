-- Migration: Fix run_id foreign key constraint issues
-- This addresses the GetNextAvailableSession foreign key constraint errors

-- Step 1: Create default workshop runs for sessions that don't have valid run_id
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
GROUP BY ws.workshop_id, w.name, w.title_ar
ON CONFLICT DO NOTHING;

-- Step 2: Update sessions with NULL run_id to reference the default run
UPDATE workshop_sessions 
SET run_id = (
    SELECT wr.id 
    FROM workshop_runs wr 
    WHERE wr.workshop_id = workshop_sessions.workshop_id 
    AND wr.run_name LIKE 'Default Run -%'
    ORDER BY wr.created_at DESC
    LIMIT 1
)
WHERE run_id IS NULL;

-- Step 3: Update sessions with invalid run_id to reference the default run  
UPDATE workshop_sessions 
SET run_id = (
    SELECT wr.id 
    FROM workshop_runs wr 
    WHERE wr.workshop_id = workshop_sessions.workshop_id 
    AND wr.run_name LIKE 'Default Run -%'
    ORDER BY wr.created_at DESC
    LIMIT 1
)
WHERE run_id NOT IN (SELECT id FROM workshop_runs WHERE workshop_id = workshop_sessions.workshop_id);

-- Step 4: Ensure all order_items have valid run_id values
UPDATE order_items oi
SET run_id = (
    SELECT ws.run_id 
    FROM workshop_sessions ws 
    WHERE ws.id = oi.session_id
)
WHERE oi.session_id IS NOT NULL 
  AND (oi.run_id IS NULL OR oi.run_id NOT IN (SELECT id FROM workshop_runs));

-- Step 5: Add a function to validate run_id relationships
CREATE OR REPLACE FUNCTION validate_session_run_id() RETURNS TRIGGER AS $$
BEGIN
    -- Ensure run_id references a valid workshop run for the same workshop
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

-- Step 6: Create trigger to validate run_id on insert/update
DROP TRIGGER IF EXISTS validate_session_run_id_trigger ON workshop_sessions;
CREATE TRIGGER validate_session_run_id_trigger
    BEFORE INSERT OR UPDATE OF run_id ON workshop_sessions
    FOR EACH ROW
    EXECUTE FUNCTION validate_session_run_id();

-- Step 7: Verification queries (output will be shown in logs)
DO $$
DECLARE
    sessions_without_run_id INTEGER;
    invalid_order_items INTEGER;
    total_runs INTEGER;
BEGIN
    -- Count sessions without valid run_id
    SELECT COUNT(*) INTO sessions_without_run_id
    FROM workshop_sessions ws
    WHERE ws.run_id IS NULL 
       OR ws.run_id NOT IN (SELECT id FROM workshop_runs WHERE workshop_id = ws.workshop_id);
    
    -- Count order items with invalid run_id
    SELECT COUNT(*) INTO invalid_order_items
    FROM order_items oi
    WHERE oi.session_id IS NOT NULL 
      AND (oi.run_id IS NULL OR oi.run_id NOT IN (SELECT id FROM workshop_runs));
    
    -- Count total runs
    SELECT COUNT(*) INTO total_runs FROM workshop_runs;
    
    RAISE NOTICE 'Migration 004 Results:';
    RAISE NOTICE 'Sessions without valid run_id: %', sessions_without_run_id;
    RAISE NOTICE 'Order items with invalid run_id: %', invalid_order_items;
    RAISE NOTICE 'Total workshop runs: %', total_runs;
    
    IF sessions_without_run_id > 0 OR invalid_order_items > 0 THEN
        RAISE WARNING 'Some run_id issues remain - manual investigation may be needed';
    ELSE
        RAISE NOTICE 'All run_id foreign key issues have been resolved';
    END IF;
END;
$$;
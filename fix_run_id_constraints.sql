-- Fix foreign key constraint issues with run_id
-- This script addresses the issue where workshop_sessions have NULL or invalid run_id values
-- that cause foreign key constraint violations when creating orders

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

-- Step 5: Add constraints to prevent future issues (optional - remove if causing problems)
-- ALTER TABLE workshop_sessions 
-- ADD CONSTRAINT workshop_sessions_run_id_not_null 
-- CHECK (run_id IS NOT NULL);

-- Step 6: Verify the fix
SELECT 'Verification Results:' as status;

SELECT 
    'Sessions without valid run_id:' as check_type,
    COUNT(*) as count
FROM workshop_sessions ws
WHERE ws.run_id IS NULL 
   OR ws.run_id NOT IN (SELECT id FROM workshop_runs WHERE workshop_id = ws.workshop_id);

SELECT 
    'Order items without valid run_id:' as check_type,
    COUNT(*) as count
FROM order_items oi
WHERE oi.session_id IS NOT NULL 
  AND (oi.run_id IS NULL OR oi.run_id NOT IN (SELECT id FROM workshop_runs));

SELECT 
    'Total workshop runs created:' as check_type,
    COUNT(*) as count
FROM workshop_runs;

SELECT 
    'Sessions per workshop run:' as check_type,
    wr.run_name,
    COUNT(ws.id) as session_count
FROM workshop_runs wr
LEFT JOIN workshop_sessions ws ON wr.id = ws.run_id
GROUP BY wr.id, wr.run_name
ORDER BY wr.run_name;
-- Remove role column from creators table since all users are the same
ALTER TABLE creators DROP COLUMN IF EXISTS role;

-- Remove any role-related constraints or indexes
DROP INDEX IF EXISTS idx_creators_role;
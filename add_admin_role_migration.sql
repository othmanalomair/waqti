-- Migration: Add admin role support to creators table
-- Description: Adds role field to creators table and creates first admin user

-- Add role column to creators table
ALTER TABLE creators ADD COLUMN role VARCHAR(20) DEFAULT 'creator' NOT NULL;

-- Add index for role lookups
CREATE INDEX idx_creators_role ON creators(role);

-- Add constraint to ensure valid roles
ALTER TABLE creators ADD CONSTRAINT chk_creators_role 
    CHECK (role IN ('creator', 'admin', 'super_admin'));

-- Create a comment to document the role values
COMMENT ON COLUMN creators.role IS 'User role: creator (default), admin (platform admin), super_admin (full access)';

-- Update existing creators to have 'creator' role (they already have DEFAULT 'creator')
-- This is just to be explicit - no action needed due to DEFAULT value

-- Example: Create first admin user (commented out - should be done manually)
-- INSERT INTO creators (name, name_ar, username, email, password_hash, role, plan, plan_ar, is_active, email_verified)
-- VALUES (
--     'Admin User',
--     'مدير النظام',
--     'admin',
--     'admin@waqti.me',
--     '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: 'password'
--     'super_admin',
--     'unlimited',
--     'غير محدود',
--     true,
--     true
-- );
-- Migration: Add admin role support to creators table
-- Date: 2025-06-15
-- Description: Adds role column to creators table and creates admin analytics table

BEGIN;

-- Add role column to creators table
ALTER TABLE creators ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'creator';

-- Add constraint to validate role values
ALTER TABLE creators ADD CONSTRAINT creators_role_check 
    CHECK (role IN ('creator', 'admin', 'super_admin'));

-- Create index for role-based queries
CREATE INDEX IF NOT EXISTS idx_creators_role ON creators(role);

-- Update existing creators to have 'creator' role if null
UPDATE creators SET role = 'creator' WHERE role IS NULL;

-- Create admin analytics table for tracking system traffic
CREATE TABLE IF NOT EXISTS admin_analytics (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    page_type VARCHAR(50) NOT NULL, -- 'landing', 'signin', 'signup'
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    country VARCHAR(100),
    device VARCHAR(50),
    browser VARCHAR(50),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Index for efficient queries
    CONSTRAINT admin_analytics_page_type_check 
        CHECK (page_type IN ('landing', 'signin', 'signup', 'store_visit'))
);

-- Indexes for admin analytics
CREATE INDEX IF NOT EXISTS idx_admin_analytics_page_type ON admin_analytics(page_type);
CREATE INDEX IF NOT EXISTS idx_admin_analytics_timestamp ON admin_analytics(timestamp);
CREATE INDEX IF NOT EXISTS idx_admin_analytics_page_date ON admin_analytics(page_type, date(timestamp));

-- Create view for admin analytics summary
CREATE OR REPLACE VIEW admin_analytics_summary AS
SELECT 
    page_type,
    DATE(timestamp) as date,
    COUNT(*) as total_visits,
    COUNT(DISTINCT ip_address) as unique_visitors,
    COUNT(CASE WHEN device = 'Mobile' THEN 1 END) as mobile_visits,
    COUNT(CASE WHEN device = 'Desktop' THEN 1 END) as desktop_visits
FROM admin_analytics 
GROUP BY page_type, DATE(timestamp)
ORDER BY date DESC, page_type;

-- Example: Create first super admin user (update credentials as needed)
-- Note: This is commented out - use the CLI tool instead for security
/*
INSERT INTO creators (
    name, 
    name_ar, 
    username, 
    email, 
    password_hash, 
    role, 
    is_active, 
    email_verified
) VALUES (
    'System Administrator',
    'مدير النظام',
    'admin',
    'admin@waqti.me',
    '$2a$10$example_hash_replace_with_real_hash', -- Replace with actual bcrypt hash
    'super_admin',
    true,
    true
) ON CONFLICT (username) DO NOTHING;
*/

COMMIT;

-- Success message
SELECT 'Admin role support added successfully!' as status;
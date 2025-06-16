-- Add store_username column to admin_analytics table for better store visit tracking

ALTER TABLE admin_analytics 
ADD COLUMN store_username VARCHAR(100);

-- Add index for store_username filtering
CREATE INDEX idx_admin_analytics_store_username 
ON admin_analytics(store_username) 
WHERE store_username IS NOT NULL;

-- Add composite index for store analytics
CREATE INDEX idx_admin_analytics_store_page_date 
ON admin_analytics(store_username, page_type, DATE(timestamp)) 
WHERE store_username IS NOT NULL;

COMMENT ON COLUMN admin_analytics.store_username IS 'Username of the store being visited (only for store_visit page type)';

SELECT 'Store username column added to admin_analytics table successfully!' as result;
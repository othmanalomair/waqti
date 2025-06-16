-- Add enhanced analytics test data with actual store usernames for better filtering testing

INSERT INTO admin_analytics (page_type, ip_address, user_agent, referrer, country, device, browser, store_username, timestamp) VALUES
-- Store visits with actual usernames from the last hour
('store_visit', '192.168.1.21', 'Mozilla/5.0 Chrome/120.0', 'https://instagram.com', 'Kuwait', 'Mobile', 'Chrome', 'othman', NOW() - INTERVAL '15 minutes'),
('store_visit', '192.168.1.22', 'Mozilla/5.0 Safari/17.0', 'https://whatsapp.com', 'UAE', 'Mobile', 'Safari', 'ahmed123', NOW() - INTERVAL '25 minutes'),
('store_visit', '192.168.1.23', 'Mozilla/5.0 Firefox/119.0', 'waqti.me', 'Saudi Arabia', 'Desktop', 'Firefox', 'othmanalomair11', NOW() - INTERVAL '35 minutes'),
('store_visit', '192.168.1.24', 'Mozilla/5.0 Chrome/120.0', 'https://telegram.com', 'Qatar', 'Mobile', 'Chrome', 'othmanalothman', NOW() - INTERVAL '45 minutes'),
('store_visit', '192.168.1.25', 'Mozilla/5.0 Safari/17.0', 'waqti.me', 'Kuwait', 'Tablet', 'Safari', 'waqtipartial', NOW() - INTERVAL '55 minutes'),

-- More store visits from yesterday
('store_visit', '192.168.1.26', 'Mozilla/5.0 Chrome/120.0', 'https://google.com', 'UAE', 'Desktop', 'Chrome', 'othman', NOW() - INTERVAL '1 day 2 hours'),
('store_visit', '192.168.1.27', 'Mozilla/5.0 Safari/17.0', 'https://facebook.com', 'Kuwait', 'Mobile', 'Safari', 'ahmed123', NOW() - INTERVAL '1 day 4 hours'),
('store_visit', '192.168.1.28', 'Mozilla/5.0 Edge/119.0', 'waqti.me', 'Bahrain', 'Desktop', 'Edge', 'othman', NOW() - INTERVAL '1 day 6 hours'),

-- Store visits from a week ago
('store_visit', '192.168.1.29', 'Mozilla/5.0 Chrome/120.0', 'https://twitter.com', 'Kuwait', 'Mobile', 'Chrome', 'othmanalomair11', NOW() - INTERVAL '7 days 1 hour'),
('store_visit', '192.168.1.30', 'Mozilla/5.0 Safari/17.0', 'waqti.me', 'UAE', 'Tablet', 'Safari', 'waqtipartial', NOW() - INTERVAL '7 days 3 hours'),

-- Recent landing and signin pages (for comparison)
('landing', '192.168.1.31', 'Mozilla/5.0 Chrome/120.0', 'https://google.com', 'Kuwait', 'Desktop', 'Chrome', NULL, NOW() - INTERVAL '10 minutes'),
('signin', '192.168.1.32', 'Mozilla/5.0 Safari/17.0', 'waqti.me', 'UAE', 'Mobile', 'Safari', NULL, NOW() - INTERVAL '20 minutes'),
('signup', '192.168.1.33', 'Mozilla/5.0 Firefox/119.0', 'https://instagram.com', 'Saudi Arabia', 'Desktop', 'Firefox', NULL, NOW() - INTERVAL '30 minutes');

-- Display enhanced analytics summary
SELECT 'Enhanced Analytics Test Data Added Successfully!' as message;

SELECT 
    'Store Visit Summary:' as summary_type,
    store_username,
    COUNT(*) as visits,
    COUNT(DISTINCT ip_address) as unique_visitors
FROM admin_analytics 
WHERE page_type = 'store_visit' AND store_username IS NOT NULL
GROUP BY store_username 
ORDER BY visits DESC;

SELECT 
    'Recent Activity with Store Names:' as summary_type,
    page_type,
    CASE 
        WHEN page_type = 'store_visit' AND store_username IS NOT NULL 
        THEN store_username 
        ELSE 'N/A' 
    END as store,
    device,
    timestamp
FROM admin_analytics 
WHERE timestamp >= NOW() - INTERVAL '2 hours'
ORDER BY timestamp DESC 
LIMIT 10;
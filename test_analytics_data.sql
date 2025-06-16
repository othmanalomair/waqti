-- Add some test analytics data for better demo
INSERT INTO admin_analytics (page_type, ip_address, user_agent, referrer, country, device, browser, timestamp) VALUES
-- Landing page visits
('landing', '192.168.1.1', 'Mozilla/5.0 Chrome/91.0', 'https://google.com', 'Kuwait', 'Desktop', 'Chrome', NOW() - INTERVAL '1 hour'),
('landing', '192.168.1.2', 'Mozilla/5.0 Safari/14.0', 'https://instagram.com', 'UAE', 'Mobile', 'Safari', NOW() - INTERVAL '2 hours'),
('landing', '192.168.1.3', 'Mozilla/5.0 Firefox/89.0', 'https://twitter.com', 'Saudi Arabia', 'Desktop', 'Firefox', NOW() - INTERVAL '3 hours'),
('landing', '192.168.1.4', 'Mozilla/5.0 Chrome/91.0', 'https://facebook.com', 'Kuwait', 'Mobile', 'Chrome', NOW() - INTERVAL '4 hours'),
('landing', '192.168.1.5', 'Mozilla/5.0 Safari/14.0', 'https://linkedin.com', 'Qatar', 'Tablet', 'Safari', NOW() - INTERVAL '5 hours'),

-- Sign in page visits
('signin', '192.168.1.1', 'Mozilla/5.0 Chrome/91.0', 'waqti.me', 'Kuwait', 'Desktop', 'Chrome', NOW() - INTERVAL '1 hour'),
('signin', '192.168.1.2', 'Mozilla/5.0 Safari/14.0', 'waqti.me', 'UAE', 'Mobile', 'Safari', NOW() - INTERVAL '2 hours'),
('signin', '192.168.1.6', 'Mozilla/5.0 Edge/91.0', 'waqti.me', 'Bahrain', 'Desktop', 'Edge', NOW() - INTERVAL '6 hours'),
('signin', '192.168.1.7', 'Mozilla/5.0 Chrome/91.0', 'waqti.me', 'Kuwait', 'Mobile', 'Chrome', NOW() - INTERVAL '7 hours'),

-- Sign up page visits
('signup', '192.168.1.8', 'Mozilla/5.0 Chrome/91.0', 'https://google.com', 'Kuwait', 'Desktop', 'Chrome', NOW() - INTERVAL '8 hours'),
('signup', '192.168.1.9', 'Mozilla/5.0 Safari/14.0', 'https://instagram.com', 'UAE', 'Mobile', 'Safari', NOW() - INTERVAL '9 hours'),
('signup', '192.168.1.10', 'Mozilla/5.0 Firefox/89.0', 'waqti.me', 'Saudi Arabia', 'Desktop', 'Firefox', NOW() - INTERVAL '10 hours'),

-- Store visits
('store_visit', '192.168.1.11', 'Mozilla/5.0 Chrome/91.0', 'waqti.me', 'Kuwait', 'Desktop', 'Chrome', NOW() - INTERVAL '30 minutes'),
('store_visit', '192.168.1.12', 'Mozilla/5.0 Safari/14.0', 'https://whatsapp.com', 'UAE', 'Mobile', 'Safari', NOW() - INTERVAL '45 minutes'),
('store_visit', '192.168.1.13', 'Mozilla/5.0 Chrome/91.0', 'https://telegram.com', 'Kuwait', 'Mobile', 'Chrome', NOW() - INTERVAL '1 hour 15 minutes'),
('store_visit', '192.168.1.14', 'Mozilla/5.0 Firefox/89.0', 'waqti.me', 'Qatar', 'Desktop', 'Firefox', NOW() - INTERVAL '2 hours 30 minutes'),
('store_visit', '192.168.1.15', 'Mozilla/5.0 Safari/14.0', 'waqti.me', 'Saudi Arabia', 'Tablet', 'Safari', NOW() - INTERVAL '3 hours 45 minutes'),

-- Previous day data
('landing', '192.168.1.16', 'Mozilla/5.0 Chrome/91.0', 'https://google.com', 'Kuwait', 'Desktop', 'Chrome', NOW() - INTERVAL '1 day 2 hours'),
('landing', '192.168.1.17', 'Mozilla/5.0 Safari/14.0', 'https://instagram.com', 'UAE', 'Mobile', 'Safari', NOW() - INTERVAL '1 day 4 hours'),
('signin', '192.168.1.18', 'Mozilla/5.0 Firefox/89.0', 'waqti.me', 'Kuwait', 'Desktop', 'Firefox', NOW() - INTERVAL '1 day 6 hours'),
('store_visit', '192.168.1.19', 'Mozilla/5.0 Chrome/91.0', 'waqti.me', 'UAE', 'Mobile', 'Chrome', NOW() - INTERVAL '1 day 8 hours'),
('store_visit', '192.168.1.20', 'Mozilla/5.0 Safari/14.0', 'waqti.me', 'Qatar', 'Tablet', 'Safari', NOW() - INTERVAL '1 day 10 hours');

-- Display current analytics summary
SELECT 'Analytics Summary After Adding Test Data:' as message;
SELECT page_type, COUNT(*) as total_visits, COUNT(DISTINCT ip_address) as unique_visitors 
FROM admin_analytics 
GROUP BY page_type 
ORDER BY total_visits DESC;
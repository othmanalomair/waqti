-- Waqti.me PostgreSQL Database Schema
-- Created for a mobile-first platform for Gulf creators to monetize their time

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Creators table - Main user accounts
CREATE TABLE creators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    name_ar VARCHAR(255),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar TEXT,
    plan VARCHAR(50) DEFAULT 'free' CHECK (plan IN ('free', 'pro')),
    plan_ar VARCHAR(50) DEFAULT 'مجاني',
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    email_verified_at TIMESTAMP,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast username and email lookups
CREATE INDEX idx_creators_username ON creators(username);
CREATE INDEX idx_creators_email ON creators(email);
CREATE INDEX idx_creators_active ON creators(is_active);

-- URL settings for customizable creator links
CREATE TABLE url_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    username VARCHAR(50) NOT NULL, -- Matches creators.username but tracked separately for change history
    changes_used INTEGER DEFAULT 0,
    max_changes INTEGER DEFAULT 5,
    last_changed TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_url_settings_creator ON url_settings(creator_id);

-- Shop settings for branding and customization
CREATE TABLE shop_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    -- Branding
    logo_url TEXT,
    creator_name VARCHAR(255),
    creator_name_ar VARCHAR(255),
    sub_header TEXT,
    sub_header_ar TEXT,
    enrollment_whatsapp VARCHAR(20),
    contact_whatsapp VARCHAR(20),
    -- Checkout preferences
    checkout_language VARCHAR(10) DEFAULT 'both' CHECK (checkout_language IN ('ar', 'en', 'both')),
    greeting_message TEXT,
    greeting_message_ar TEXT,
    currency_symbol VARCHAR(10) DEFAULT 'KD',
    currency_symbol_ar VARCHAR(10) DEFAULT 'د.ك',
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_shop_settings_creator ON shop_settings(creator_id);

-- Categories for organizing workshops
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    name_ar VARCHAR(255),
    description TEXT,
    description_ar TEXT,
    color VARCHAR(7) DEFAULT '#2DD4BF', -- Hex color
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_creator ON categories(creator_id);
CREATE INDEX idx_categories_active ON categories(creator_id, is_active);

-- Main workshops/courses table
CREATE TABLE workshops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL, -- Internal name, added NOT NULL constraint
    title VARCHAR(255) NOT NULL,
    title_ar VARCHAR(255),
    description TEXT,
    description_ar TEXT,
    price DECIMAL(10,3) DEFAULT 0.000, -- Support 3 decimal places for KWD
    currency VARCHAR(3) DEFAULT 'KWD',
    duration INTEGER, -- Duration in minutes
    max_students INTEGER DEFAULT 0, -- 0 = unlimited
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    is_active BOOLEAN DEFAULT false,
    is_free BOOLEAN DEFAULT false,
    is_recurring BOOLEAN DEFAULT false,
    recurrence_type VARCHAR(20) CHECK (recurrence_type IN ('weekly', 'monthly', 'yearly')),
    sort_order INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    enrollment_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_workshops_creator ON workshops(creator_id);
CREATE INDEX idx_workshops_active ON workshops(creator_id, is_active);
CREATE INDEX idx_workshops_status ON workshops(creator_id, status);
CREATE INDEX idx_workshops_category ON workshops(category_id);

-- Workshop images for galleries
CREATE TABLE workshop_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workshop_id UUID NOT NULL REFERENCES workshops(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    is_cover BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 0,
    alt_text VARCHAR(255),
    alt_text_ar VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_workshop_images_workshop ON workshop_images(workshop_id);
CREATE INDEX idx_workshop_images_cover ON workshop_images(workshop_id, is_cover);

-- Workshop sessions for scheduling
CREATE TABLE workshop_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workshop_id UUID NOT NULL REFERENCES workshops(id) ON DELETE CASCADE,
    session_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME,
    duration DECIMAL(4,2), -- Duration in hours (e.g., 2.5)
    timezone VARCHAR(50) DEFAULT 'Asia/Kuwait',
    location TEXT,
    location_ar TEXT,
    max_attendees INTEGER,
    current_attendees INTEGER DEFAULT 0,
    is_completed BOOLEAN DEFAULT false,
    notes TEXT,
    notes_ar TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_workshop_sessions_workshop ON workshop_sessions(workshop_id);
CREATE INDEX idx_workshop_sessions_date ON workshop_sessions(session_date);
CREATE INDEX idx_workshop_sessions_upcoming ON workshop_sessions(workshop_id, session_date, is_completed);

-- Orders from customers
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    order_number VARCHAR(50) UNIQUE, -- Human-readable order number
    customer_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(20) NOT NULL,
    customer_email VARCHAR(255),
    total_amount DECIMAL(10,3) NOT NULL,
    currency VARCHAR(3) DEFAULT 'KWD',
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'cancelled', 'refunded')),
    status_ar VARCHAR(20) DEFAULT 'قيد الانتظار',
    payment_method VARCHAR(50),
    payment_reference VARCHAR(255),
    order_source VARCHAR(20) DEFAULT 'whatsapp' CHECK (order_source IN ('whatsapp', 'direct', 'qr')),
    notes TEXT,
    is_viewed BOOLEAN DEFAULT false, -- For notification management
    viewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_creator ON orders(creator_id);
CREATE INDEX idx_orders_status ON orders(creator_id, status);
CREATE INDEX idx_orders_date ON orders(created_at);
CREATE INDEX idx_orders_customer_phone ON orders(customer_phone);
CREATE UNIQUE INDEX idx_orders_number ON orders(order_number) WHERE order_number IS NOT NULL;

-- Order items (individual workshops in an order)
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    workshop_id UUID NOT NULL REFERENCES workshops(id) ON DELETE RESTRICT,
    workshop_name VARCHAR(255) NOT NULL, -- Snapshot at time of order
    workshop_name_ar VARCHAR(255),
    price DECIMAL(10,3) NOT NULL, -- Price at time of order
    quantity INTEGER DEFAULT 1,
    subtotal DECIMAL(10,3) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_order_items_workshop ON order_items(workshop_id);

-- Enrollments (successful registrations)
CREATE TABLE enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workshop_id UUID NOT NULL REFERENCES workshops(id) ON DELETE RESTRICT,
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    session_id UUID REFERENCES workshop_sessions(id) ON DELETE SET NULL,
    student_name VARCHAR(255) NOT NULL,
    student_email VARCHAR(255),
    student_phone VARCHAR(20),
    total_price DECIMAL(10,3) NOT NULL,
    status VARCHAR(20) DEFAULT 'successful' CHECK (status IN ('successful', 'pending', 'rejected', 'cancelled')),
    status_ar VARCHAR(20) DEFAULT 'مكتمل',
    enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completion_status VARCHAR(20) DEFAULT 'enrolled' CHECK (completion_status IN ('enrolled', 'in_progress', 'completed', 'dropped')),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_enrollments_workshop ON enrollments(workshop_id);
CREATE INDEX idx_enrollments_student_email ON enrollments(student_email);
CREATE INDEX idx_enrollments_status ON enrollments(status);
CREATE INDEX idx_enrollments_date ON enrollments(enrollment_date);

-- Analytics clicks for tracking page visits
CREATE TABLE analytics_clicks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    country VARCHAR(100),
    country_ar VARCHAR(100),
    city VARCHAR(100),
    city_ar VARCHAR(100),
    device VARCHAR(50), -- Mobile, Desktop, Tablet
    device_ar VARCHAR(50),
    os VARCHAR(50), -- iOS, Android, Windows, etc.
    os_ar VARCHAR(50),
    browser VARCHAR(50),
    browser_ar VARCHAR(50),
    platform VARCHAR(50), -- Instagram, WhatsApp, Direct, etc.
    platform_ar VARCHAR(50),
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_analytics_clicks_creator ON analytics_clicks(creator_id);
CREATE INDEX idx_analytics_clicks_date ON analytics_clicks(clicked_at);
CREATE INDEX idx_analytics_clicks_country ON analytics_clicks(creator_id, country);
CREATE INDEX idx_analytics_clicks_device ON analytics_clicks(creator_id, device);
CREATE INDEX idx_analytics_clicks_platform ON analytics_clicks(creator_id, platform);

-- Email verification tokens
CREATE TABLE email_verification_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email_tokens_creator ON email_verification_tokens(creator_id);
CREATE INDEX idx_email_tokens_token ON email_verification_tokens(token);

-- Password reset tokens
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_password_tokens_creator ON password_reset_tokens(creator_id);
CREATE INDEX idx_password_tokens_token ON password_reset_tokens(token);

-- Creator sessions for authentication
CREATE TABLE creator_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    session_token VARCHAR(255) NOT NULL UNIQUE,
    device_info TEXT,
    ip_address INET,
    expires_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_creator_sessions_creator ON creator_sessions(creator_id);
CREATE INDEX idx_creator_sessions_token ON creator_sessions(session_token);
CREATE INDEX idx_creator_sessions_expires ON creator_sessions(expires_at);

-- Notifications for creators
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- 'new_order', 'payment_received', 'new_enrollment', etc.
    title VARCHAR(255) NOT NULL,
    title_ar VARCHAR(255),
    message TEXT NOT NULL,
    message_ar TEXT,
    data JSONB, -- Additional structured data
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_creator ON notifications(creator_id);
CREATE INDEX idx_notifications_unread ON notifications(creator_id, is_read);
CREATE INDEX idx_notifications_type ON notifications(creator_id, type);

-- Subscription/billing information for pro users
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    plan VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'cancelled', 'past_due', 'trialing')),
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    cancel_at_period_end BOOLEAN DEFAULT false,
    stripe_subscription_id VARCHAR(255),
    stripe_customer_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_subscriptions_creator ON subscriptions(creator_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_stripe ON subscriptions(stripe_subscription_id);

-- Promo codes and discounts
CREATE TABLE promo_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES creators(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    type VARCHAR(20) DEFAULT 'percentage' CHECK (type IN ('percentage', 'fixed')),
    value DECIMAL(10,3) NOT NULL, -- Percentage (0-100) or fixed amount
    min_order_amount DECIMAL(10,3) DEFAULT 0,
    max_uses INTEGER DEFAULT 0, -- 0 = unlimited
    current_uses INTEGER DEFAULT 0,
    starts_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_promo_codes_creator_code ON promo_codes(creator_id, code);
CREATE INDEX idx_promo_codes_active ON promo_codes(is_active, expires_at);

-- Promo code usage tracking
CREATE TABLE promo_code_uses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    promo_code_id UUID NOT NULL REFERENCES promo_codes(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    discount_amount DECIMAL(10,3) NOT NULL,
    used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_promo_uses_code ON promo_code_uses(promo_code_id);
CREATE INDEX idx_promo_uses_order ON promo_code_uses(order_id);

-- Function to automatically update 'updated_at' timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply the update trigger to relevant tables
CREATE TRIGGER update_creators_updated_at BEFORE UPDATE ON creators FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_url_settings_updated_at BEFORE UPDATE ON url_settings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_shop_settings_updated_at BEFORE UPDATE ON shop_settings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_workshops_updated_at BEFORE UPDATE ON workshops FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_workshop_sessions_updated_at BEFORE UPDATE ON workshop_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_enrollments_updated_at BEFORE UPDATE ON enrollments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON subscriptions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_promo_codes_updated_at BEFORE UPDATE ON promo_codes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to generate order numbers
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.order_number IS NULL THEN
        NEW.order_number := 'WQ' || TO_CHAR(NEW.created_at, 'YYYYMMDD') || '-' || LPAD(CAST((SELECT COALESCE(MAX(CAST(RIGHT(order_number, 4) AS INTEGER)), 0) + 1 FROM orders WHERE order_number LIKE 'WQ' || TO_CHAR(NEW.created_at, 'YYYYMMDD') || '-%') AS TEXT), 4, '0');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply order number generation trigger
CREATE TRIGGER generate_order_number_trigger BEFORE INSERT ON orders FOR EACH ROW EXECUTE FUNCTION generate_order_number();

-- Insert initial data for the demo creator
INSERT INTO creators (
    id, name, name_ar, username, email, password_hash, plan, plan_ar, is_active, email_verified
) VALUES (
    '550e8400-e29b-41d4-a716-446655440000',
    'Ahmed Al-Kuwaiti',
    'أحمد الكويتي',
    'ahmed',
    'demo@waqti.me',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- 'password'
    'free',
    'مجاني',
    true,
    true
);

-- Insert URL settings for demo creator
INSERT INTO url_settings (
    id, creator_id, username, changes_used, max_changes, last_changed
) VALUES (
    '550e8400-e29b-41d4-a716-446655440070',
    '550e8400-e29b-41d4-a716-446655440000',
    'ahmed',
    2,
    5,
    NOW() - INTERVAL '1 month'
);

-- Insert shop settings for demo creator
INSERT INTO shop_settings (
    id, creator_id, logo_url, creator_name, creator_name_ar, sub_header, sub_header_ar,
    enrollment_whatsapp, contact_whatsapp, greeting_message, greeting_message_ar
) VALUES (
    '550e8400-e29b-41d4-a716-446655440060',
    '550e8400-e29b-41d4-a716-446655440000',
    '/static/images/default.jpg',
    'Ahmed Al-Kuwaiti',
    'أحمد الكويتي',
    'Certified Design Trainer',
    'مدرب معتمد في التصميم',
    '+965-9999-8888',
    '+965-9999-7777',
    'Welcome to my workshop! Ready to learn?',
    'مرحباً بك في ورشتي! هل أنت مستعد للتعلم؟'
);

-- Insert sample workshops
-- Added 'name' column to INSERT statement, using 'title' value for 'name'
INSERT INTO workshops (
    id, creator_id, name, title, title_ar, description, description_ar, price, duration, max_students, is_active, sort_order
) VALUES
('550e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', 'Photography Basics', 'Photography Basics', 'أساسيات التصوير', 'Learn the fundamentals of photography', 'تعلم أساسيات التصوير الفوتوغرافي', 25.000, 120, 15, true, 1),
('550e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', 'Digital Marketing', 'Digital Marketing', 'التسويق الرقمي', 'Master social media marketing strategies', 'إتقن استراتيجيات التسويق عبر وسائل التواصل', 35.000, 90, 20, true, 2),
('550e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440000', 'Arabic Calligraphy', 'Arabic Calligraphy', 'الخط العربي', 'Traditional Arabic calligraphy techniques', 'تقنيات الخط العربي التقليدية', 20.000, 150, 10, false, 3),
('550e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440000', 'Business English', 'Business English', 'الإنجليزية التجارية', 'Professional English for business communication', 'الإنجليزية المهنية للتواصل التجاري', 30.000, 60, 12, true, 4);

-- Insert sample orders
INSERT INTO orders (
    id, creator_id, customer_name, customer_phone, total_amount, status, status_ar, order_source
) VALUES
('550e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440000', 'أحمد محمد', '+965-9999-1234', 25.000, 'pending', 'قيد الانتظار', 'whatsapp'),
('550e8400-e29b-41d4-a716-446655440031', '550e8400-e29b-41d4-a716-446655440000', 'سارة أحمد', '+965-9999-5678', 35.000, 'pending', 'قيد الانتظار', 'whatsapp');

-- Insert sample order items
INSERT INTO order_items (
    id, order_id, workshop_id, workshop_name, workshop_name_ar, price, quantity, subtotal
) VALUES
('550e8400-e29b-41d4-a716-446655440040', '550e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440001', 'Photography Basics', 'أساسيات التصوير', 25.000, 1, 25.000),
('550e8400-e29b-41d4-a716-446655440041', '550e8400-e29b-41d4-a716-446655440031', '550e8400-e29b-41d4-a716-446655440002', 'Digital Marketing', 'التسويق الرقمي', 35.000, 1, 35.000);

-- Insert sample enrollments
INSERT INTO enrollments (
    id, workshop_id, student_name, student_email, total_price, status, status_ar, enrollment_date
) VALUES
('550e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440001', 'سارة أحمد', 'sara@example.com', 25.000, 'successful', 'مكتمل', NOW() - INTERVAL '2 days'),
('550e8400-e29b-41d4-a716-446655440021', '550e8400-e29b-41d4-a716-446655440002', 'محمد الكويتي', 'mohammed@example.com', 35.000, 'successful', 'مكتمل', NOW() - INTERVAL '5 days'),
('550e8400-e29b-41d4-a716-446655440022', '550e8400-e29b-41d4-a716-446655440004', 'فاطمة الزهراء', 'fatima@example.com', 30.000, 'rejected', 'مرفوض', NOW() - INTERVAL '1 day'),
('550e8400-e29b-41d4-a716-446655440023', '550e8400-e29b-41d4-a716-446655440001', 'أحمد عبدالله', 'ahmed@example.com', 25.000, 'successful', 'مكتمل', NOW() - INTERVAL '7 days'),
('550e8400-e29b-41d4-a716-446655440024', '550e8400-e29b-41d4-a716-446655440002', 'نورا السالم', 'nora@example.com', 35.000, 'pending', 'قيد المراجعة', NOW()),
('550e8400-e29b-41d4-a716-446655440025', '550e8400-e29b-41d4-a716-446655440003', 'يوسف المطيري', 'youssef@example.com', 20.000, 'successful', 'مكتمل', NOW() - INTERVAL '10 days');

-- Insert sample analytics clicks
INSERT INTO analytics_clicks (
    id, creator_id, country, country_ar, device, device_ar, os, os_ar, platform, platform_ar, clicked_at
) VALUES
('550e8400-e29b-41d4-a716-446655440050', '550e8400-e29b-41d4-a716-446655440000', 'Kuwait', 'الكويت', 'Mobile', 'جوال', 'iOS', 'آي أو إس', 'Instagram', 'إنستغرام', NOW() - INTERVAL '1 day'),
('550e8400-e29b-41d4-a716-446655440051', '550e8400-e29b-41d4-a716-446655440000', 'Saudi Arabia', 'السعودية', 'Desktop', 'سطح المكتب', 'Windows', 'ويندوز', 'WhatsApp', 'واتساب', NOW() - INTERVAL '2 days'),
('550e8400-e29b-41d4-a716-446655440052', '550e8400-e29b-41d4-a716-446655440000', 'UAE', 'الإمارات', 'Mobile', 'جوال', 'Android', 'أندرويد', 'Snapchat', 'سناب شات', NOW() - INTERVAL '3 days');

-- Comments for future development
COMMENT ON TABLE creators IS 'Main creators/users who create workshops and manage their stores';
COMMENT ON TABLE workshops IS 'Individual workshops or courses offered by creators';
COMMENT ON TABLE orders IS 'Customer orders placed through WhatsApp or direct booking';
COMMENT ON TABLE enrollments IS 'Successful enrollments/registrations for workshops';
COMMENT ON TABLE analytics_clicks IS 'Track page visits and click analytics for creators';
COMMENT ON TABLE shop_settings IS 'Customizable branding and settings for each creator''s store';
COMMENT ON TABLE url_settings IS 'Manage custom usernames and URL changes for creators';

-- Create views for common queries

-- Creator dashboard stats view
CREATE OR REPLACE VIEW creator_dashboard_stats AS
SELECT
    c.id as creator_id,
    COUNT(DISTINCT w.id) as total_workshops,
    COUNT(DISTINCT CASE WHEN w.is_active = true THEN w.id END) as active_workshops,
    COUNT(DISTINCT e.id) as total_enrollments,
    COALESCE(SUM(CASE WHEN e.status = 'successful' AND e.enrollment_date >= DATE_TRUNC('month', CURRENT_DATE) THEN e.total_price END), 0) as monthly_revenue,
    COALESCE(SUM(CASE WHEN w.is_active = true AND w.max_students > 0 THEN w.price * w.max_students ELSE 0 END), 0) * 0.7 as projected_sales, -- Added check for max_students > 0
    COALESCE(SUM(CASE WHEN w.is_active = true THEN w.max_students END), 0) - COUNT(DISTINCT e.id) as remaining_seats
FROM creators c
LEFT JOIN workshops w ON c.id = w.creator_id
LEFT JOIN enrollments e ON w.id = e.workshop_id
GROUP BY c.id;

COMMENT ON VIEW creator_dashboard_stats IS 'Aggregated statistics for creator dashboard display';

-- Additional useful views for analytics and reporting

-- Monthly revenue trends view
-- Corrected to join with workshops to get creator_id
CREATE OR REPLACE VIEW monthly_revenue_trends AS
SELECT
    w.creator_id,
    DATE_TRUNC('month', e.enrollment_date) as month,
    COUNT(e.id) as enrollments_count,
    SUM(e.total_price) as total_revenue,
    AVG(e.total_price) as avg_order_value
FROM enrollments e
JOIN workshops w ON e.workshop_id = w.id
WHERE e.status = 'successful'
GROUP BY w.creator_id, DATE_TRUNC('month', e.enrollment_date)
ORDER BY w.creator_id, month;

-- Popular workshops view
CREATE OR REPLACE VIEW popular_workshops AS
SELECT
    w.id,
    w.creator_id,
    w.title,
    w.title_ar,
    w.price,
    COUNT(e.id) as enrollment_count,
    AVG(w.price) as avg_price, -- This is avg price of workshop, not avg price per enrollment
    SUM(e.total_price) as total_revenue
FROM workshops w
LEFT JOIN enrollments e ON w.id = e.workshop_id AND e.status = 'successful'
WHERE w.is_active = true
GROUP BY w.id, w.creator_id, w.title, w.title_ar, w.price
ORDER BY enrollment_count DESC;

-- Order management view with customer details
CREATE OR REPLACE VIEW order_management AS
SELECT
    o.id,
    o.order_number,
    o.creator_id,
    o.customer_name,
    o.customer_phone,
    o.customer_email,
    o.total_amount,
    o.status,
    o.status_ar,
    o.order_source,
    o.created_at,
    o.is_viewed,
    COUNT(oi.id) as item_count,
    STRING_AGG(oi.workshop_name, ', ') as workshop_names,
    STRING_AGG(oi.workshop_name_ar, '، ') as workshop_names_ar
FROM orders o
LEFT JOIN order_items oi ON o.id = oi.order_id
GROUP BY o.id, o.order_number, o.creator_id, o.customer_name, o.customer_phone,
         o.customer_email, o.total_amount, o.status, o.status_ar, o.order_source,
         o.created_at, o.is_viewed
ORDER BY o.created_at DESC;

-- Analytics summary view
CREATE OR REPLACE VIEW analytics_summary AS
SELECT
    creator_id,
    DATE_TRUNC('day', clicked_at) as date,
    COUNT(*) as total_clicks,
    COUNT(DISTINCT ip_address) as unique_visitors,
    COUNT(CASE WHEN device = 'Mobile' THEN 1 END) as mobile_clicks,
    COUNT(CASE WHEN device = 'Desktop' THEN 1 END) as desktop_clicks,
    COUNT(CASE WHEN platform = 'Instagram' THEN 1 END) as instagram_clicks,
    COUNT(CASE WHEN platform = 'WhatsApp' THEN 1 END) as whatsapp_clicks,
    COUNT(CASE WHEN platform = 'Direct' THEN 1 END) as direct_clicks
FROM analytics_clicks
GROUP BY creator_id, DATE_TRUNC('day', clicked_at)
ORDER BY creator_id, date DESC;

-- Enrollment tracking view
CREATE OR REPLACE VIEW enrollment_tracking AS
SELECT
    e.id,
    e.workshop_id,
    w.title as workshop_name,
    w.title_ar as workshop_name_ar,
    e.student_name,
    e.student_email,
    e.student_phone,
    e.total_price,
    e.status,
    e.status_ar,
    e.enrollment_date,
    e.completion_status,
    w.creator_id
FROM enrollments e
JOIN workshops w ON e.workshop_id = w.id
ORDER BY e.enrollment_date DESC;

-- Create materialized view for performance on analytics
CREATE MATERIALIZED VIEW creator_analytics_summary AS
SELECT
    c.id as creator_id,
    c.username,
    COUNT(DISTINCT ac.id) as total_clicks,
    COUNT(DISTINCT DATE(ac.clicked_at)) as active_days,
    COUNT(DISTINCT ac.ip_address) as unique_visitors,
    COUNT(DISTINCT w.id) as total_workshops,
    COUNT(DISTINCT CASE WHEN w.is_active THEN w.id END) as active_workshops,
    COUNT(DISTINCT e.id) as total_enrollments,
    COALESCE(SUM(CASE WHEN e.status = 'successful' THEN e.total_price END), 0) as total_revenue,
    COUNT(DISTINCT o.id) as total_orders,
    COUNT(DISTINCT CASE WHEN o.status = 'pending' THEN o.id END) as pending_orders
FROM creators c
LEFT JOIN analytics_clicks ac ON c.id = ac.creator_id
LEFT JOIN workshops w ON c.id = w.creator_id
LEFT JOIN enrollments e ON w.id = e.workshop_id -- This join might inflate counts if a workshop has many enrollments
LEFT JOIN orders o ON c.id = o.creator_id
GROUP BY c.id, c.username;

-- Index for materialized view
CREATE UNIQUE INDEX idx_creator_analytics_summary_creator ON creator_analytics_summary(creator_id);

-- Function to refresh materialized view
CREATE OR REPLACE FUNCTION refresh_analytics_summary()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW creator_analytics_summary;
END;
$$ LANGUAGE plpgsql;

-- Additional helper functions

-- Function to calculate conversion rate
CREATE OR REPLACE FUNCTION calculate_conversion_rate(creator_uuid UUID, days_back INTEGER DEFAULT 30)
RETURNS DECIMAL(5,2) AS $$
DECLARE
    total_clicks INTEGER;
    total_orders INTEGER;
    conversion_rate DECIMAL(5,2);
BEGIN
    -- Get total clicks in the period
    SELECT COUNT(*) INTO total_clicks
    FROM analytics_clicks
    WHERE creator_id = creator_uuid
    AND clicked_at >= CURRENT_DATE - INTERVAL '1 day' * days_back;

    -- Get total orders in the period
    SELECT COUNT(*) INTO total_orders
    FROM orders
    WHERE creator_id = creator_uuid
    AND created_at >= CURRENT_DATE - INTERVAL '1 day' * days_back;

    -- Calculate conversion rate
    IF total_clicks > 0 THEN
        conversion_rate := (total_orders::DECIMAL / total_clicks::DECIMAL) * 100;
    ELSE
        conversion_rate := 0;
    END IF;

    RETURN conversion_rate;
END;
$$ LANGUAGE plpgsql;

-- Function to get top performing workshops
CREATE OR REPLACE FUNCTION get_top_workshops(creator_uuid UUID, limit_count INTEGER DEFAULT 5)
RETURNS TABLE (
    workshop_id UUID,
    title VARCHAR(255),
    title_ar VARCHAR(255),
    enrollment_count BIGINT,
    revenue DECIMAL(10,3)
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        w.id,
        w.title,
        w.title_ar,
        COUNT(e.id) as enrollment_count,
        COALESCE(SUM(e.total_price), 0) as revenue
    FROM workshops w
    LEFT JOIN enrollments e ON w.id = e.workshop_id AND e.status = 'successful'
    WHERE w.creator_id = creator_uuid
    GROUP BY w.id, w.title, w.title_ar
    ORDER BY enrollment_count DESC, revenue DESC
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql;

-- Function to generate analytics for date range
CREATE OR REPLACE FUNCTION get_analytics_for_period(
    creator_uuid UUID,
    start_date DATE,
    end_date DATE
)
RETURNS TABLE (
    total_clicks BIGINT,
    unique_visitors BIGINT,
    mobile_percentage DECIMAL(5,2),
    top_country VARCHAR(100),
    top_platform VARCHAR(50),
    conversion_rate DECIMAL(5,2)
) AS $$
DECLARE
    clicks_count BIGINT;
    visitors_count BIGINT;
    mobile_count BIGINT;
    top_country_name VARCHAR(100);
    top_platform_name VARCHAR(50);
    orders_count BIGINT;
    conv_rate DECIMAL(5,2);
    calculated_mobile_percentage DECIMAL(5,2); -- Renamed to avoid conflict with RETURNS TABLE column name
BEGIN
    -- Total clicks
    SELECT COUNT(*) INTO clicks_count
    FROM analytics_clicks ac -- Added alias
    WHERE ac.creator_id = creator_uuid -- Used alias
    AND DATE(ac.clicked_at) BETWEEN start_date AND end_date;

    -- Unique visitors
    SELECT COUNT(DISTINCT ac.ip_address) INTO visitors_count
    FROM analytics_clicks ac -- Added alias
    WHERE ac.creator_id = creator_uuid -- Used alias
    AND DATE(ac.clicked_at) BETWEEN start_date AND end_date;

    -- Mobile clicks
    SELECT COUNT(*) INTO mobile_count
    FROM analytics_clicks ac -- Added alias
    WHERE ac.creator_id = creator_uuid -- Used alias
    AND DATE(ac.clicked_at) BETWEEN start_date AND end_date
    AND ac.device = 'Mobile';

    -- Top country
    SELECT ac.country INTO top_country_name
    FROM analytics_clicks ac -- Added alias
    WHERE ac.creator_id = creator_uuid -- Used alias
    AND DATE(ac.clicked_at) BETWEEN start_date AND end_date
    AND ac.country IS NOT NULL -- Ensure country is not null before grouping
    GROUP BY ac.country
    ORDER BY COUNT(*) DESC
    LIMIT 1;

    -- Top platform
    SELECT ac.platform INTO top_platform_name
    FROM analytics_clicks ac -- Added alias
    WHERE ac.creator_id = creator_uuid -- Used alias
    AND DATE(ac.clicked_at) BETWEEN start_date AND end_date
    AND ac.platform IS NOT NULL -- Ensure platform is not null before grouping
    GROUP BY ac.platform
    ORDER BY COUNT(*) DESC
    LIMIT 1;

    -- Orders count for conversion
    SELECT COUNT(*) INTO orders_count
    FROM orders o -- Added alias
    WHERE o.creator_id = creator_uuid -- Used alias
    AND DATE(o.created_at) BETWEEN start_date AND end_date;

    -- Calculate conversion rate
    IF clicks_count > 0 THEN
        conv_rate := (orders_count::DECIMAL / clicks_count::DECIMAL) * 100;
    ELSE
        conv_rate := 0;
    END IF;

    -- Mobile percentage
    IF clicks_count > 0 THEN
        calculated_mobile_percentage := (mobile_count::DECIMAL / clicks_count::DECIMAL) * 100;
    ELSE
        calculated_mobile_percentage := 0;
    END IF;

    RETURN QUERY SELECT
        clicks_count,
        visitors_count,
        calculated_mobile_percentage, -- Use the calculated variable
        top_country_name,
        top_platform_name,
        conv_rate;
END;
$$ LANGUAGE plpgsql;

-- Security: Row Level Security (RLS) policies
ALTER TABLE creators ENABLE ROW LEVEL SECURITY;
ALTER TABLE workshops ENABLE ROW LEVEL SECURITY;
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE enrollments ENABLE ROW LEVEL SECURITY;
ALTER TABLE analytics_clicks ENABLE ROW LEVEL SECURITY;
ALTER TABLE shop_settings ENABLE ROW LEVEL SECURITY;
ALTER TABLE url_settings ENABLE ROW LEVEL SECURITY;

-- Create RLS policies (assuming you'll have a way to get current creator_id)
-- These would be used with a proper authentication system

-- Example policy for creators (they can only see their own data)
-- CREATE POLICY creator_isolation_workshops ON workshops
-- FOR ALL
-- USING (creator_id = current_setting('app.current_creator_id', true)::UUID)
-- WITH CHECK (creator_id = current_setting('app.current_creator_id', true)::UUID);
-- COMMENT ON POLICY creator_isolation_workshops ON workshops IS 'Creators can only manage their own workshops.';

-- Useful indexes for performance
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_creator_status_date ON orders(creator_id, status, created_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_enrollments_workshop_status_date ON enrollments(workshop_id, status, enrollment_date DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_analytics_creator_date ON analytics_clicks(creator_id, clicked_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_workshops_creator_active_sort ON workshops(creator_id, is_active, sort_order);

-- Performance optimization: Partial indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_pending ON orders(creator_id, created_at DESC) WHERE status = 'pending';
-- Renamed duplicate index idx_workshops_active
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_workshops_active_sorted_partial ON workshops(creator_id, sort_order) WHERE is_active = true;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_enrollments_successful ON enrollments(workshop_id, enrollment_date DESC) WHERE status = 'successful';

COMMENT ON VIEW creator_dashboard_stats IS 'Aggregated statistics for creator dashboard display';
COMMENT ON VIEW monthly_revenue_trends IS 'Monthly revenue and enrollment trends for reporting';
COMMENT ON VIEW popular_workshops IS 'Workshop popularity rankings by enrollment and revenue';
COMMENT ON VIEW order_management IS 'Comprehensive order details for management interface';
COMMENT ON VIEW analytics_summary IS 'Daily analytics summary with device and platform breakdown';
COMMENT ON VIEW enrollment_tracking IS 'Detailed enrollment tracking with workshop information';
COMMENT ON MATERIALIZED VIEW creator_analytics_summary IS 'High-level creator analytics (refresh periodically for performance)';

-- Final data integrity checks
-- Corrected DO block syntax and constraint checks
DO $$
BEGIN
    -- Ensure we have all necessary constraints
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'workshops_price_check' AND conrelid = 'workshops'::regclass) THEN
        ALTER TABLE workshops ADD CONSTRAINT workshops_price_check CHECK (price >= 0);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'orders_total_check' AND conrelid = 'orders'::regclass) THEN
        ALTER TABLE orders ADD CONSTRAINT orders_total_check CHECK (total_amount >= 0);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'url_settings_changes_check' AND conrelid = 'url_settings'::regclass) THEN
        ALTER TABLE url_settings ADD CONSTRAINT url_settings_changes_check CHECK (changes_used <= max_changes);
    END IF;

    RAISE NOTICE 'Waqti.me database schema created successfully!';
    RAISE NOTICE 'Demo creator available: demo@waqti.me / password';
    RAISE NOTICE 'Creator ID: 550e8400-e29b-41d4-a716-446655440000';
END;
$$;

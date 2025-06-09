--
-- PostgreSQL database dump
--

-- Dumped from database version 14.18 (Ubuntu 14.18-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.18 (Ubuntu 14.18-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

ALTER TABLE ONLY public.workshops DROP CONSTRAINT workshops_creator_id_fkey;
ALTER TABLE ONLY public.workshops DROP CONSTRAINT workshops_category_id_fkey;
ALTER TABLE ONLY public.workshop_sessions DROP CONSTRAINT workshop_sessions_workshop_id_fkey;
ALTER TABLE ONLY public.workshop_runs DROP CONSTRAINT workshop_runs_workshop_id_fkey;
ALTER TABLE ONLY public.workshop_images DROP CONSTRAINT workshop_images_workshop_id_fkey;
ALTER TABLE ONLY public.url_settings DROP CONSTRAINT url_settings_creator_id_fkey;
ALTER TABLE ONLY public.subscriptions DROP CONSTRAINT subscriptions_creator_id_fkey;
ALTER TABLE ONLY public.shop_settings DROP CONSTRAINT shop_settings_creator_id_fkey;
ALTER TABLE ONLY public.promo_codes DROP CONSTRAINT promo_codes_creator_id_fkey;
ALTER TABLE ONLY public.promo_code_uses DROP CONSTRAINT promo_code_uses_promo_code_id_fkey;
ALTER TABLE ONLY public.promo_code_uses DROP CONSTRAINT promo_code_uses_order_id_fkey;
ALTER TABLE ONLY public.password_reset_tokens DROP CONSTRAINT password_reset_tokens_creator_id_fkey;
ALTER TABLE ONLY public.orders DROP CONSTRAINT orders_creator_id_fkey;
ALTER TABLE ONLY public.order_items DROP CONSTRAINT order_items_workshop_run_id_fkey;
ALTER TABLE ONLY public.order_items DROP CONSTRAINT order_items_workshop_id_fkey;
ALTER TABLE ONLY public.order_items DROP CONSTRAINT order_items_session_id_fkey;
ALTER TABLE ONLY public.order_items DROP CONSTRAINT order_items_run_id_fkey;
ALTER TABLE ONLY public.order_items DROP CONSTRAINT order_items_order_id_fkey;
ALTER TABLE ONLY public.notifications DROP CONSTRAINT notifications_creator_id_fkey;
ALTER TABLE ONLY public.enrollments DROP CONSTRAINT fk_enrollment_session;
ALTER TABLE ONLY public.enrollments DROP CONSTRAINT enrollments_workshop_id_fkey;
ALTER TABLE ONLY public.enrollments DROP CONSTRAINT enrollments_session_id_fkey;
ALTER TABLE ONLY public.enrollments DROP CONSTRAINT enrollments_order_id_fkey;
ALTER TABLE ONLY public.email_verification_tokens DROP CONSTRAINT email_verification_tokens_creator_id_fkey;
ALTER TABLE ONLY public.creator_sessions DROP CONSTRAINT creator_sessions_creator_id_fkey;
ALTER TABLE ONLY public.categories DROP CONSTRAINT categories_creator_id_fkey;
ALTER TABLE ONLY public.analytics_clicks DROP CONSTRAINT analytics_clicks_creator_id_fkey;
DROP TRIGGER update_workshops_updated_at ON public.workshops;
DROP TRIGGER update_workshop_sessions_updated_at ON public.workshop_sessions;
DROP TRIGGER update_url_settings_updated_at ON public.url_settings;
DROP TRIGGER update_subscriptions_updated_at ON public.subscriptions;
DROP TRIGGER update_shop_settings_updated_at ON public.shop_settings;
DROP TRIGGER update_promo_codes_updated_at ON public.promo_codes;
DROP TRIGGER update_orders_updated_at ON public.orders;
DROP TRIGGER update_enrollments_updated_at ON public.enrollments;
DROP TRIGGER update_creators_updated_at ON public.creators;
DROP TRIGGER update_categories_updated_at ON public.categories;
DROP TRIGGER trigger_update_session_status ON public.workshop_sessions;
DROP TRIGGER trigger_update_run_attendance ON public.workshop_sessions;
DROP TRIGGER generate_order_number_trigger ON public.orders;
DROP TRIGGER ensure_workshop_run_trigger ON public.workshop_sessions;
CREATE OR REPLACE VIEW public.workshop_run_summary AS
SELECT
    NULL::uuid AS run_id,
    NULL::uuid AS workshop_id,
    NULL::character varying(255) AS workshop_name,
    NULL::character varying(255) AS workshop_name_ar,
    NULL::character varying(255) AS run_name,
    NULL::date AS start_date,
    NULL::date AS end_date,
    NULL::bigint AS total_sessions,
    NULL::bigint AS total_capacity,
    NULL::bigint AS total_enrolled,
    NULL::bigint AS full_sessions,
    NULL::character varying(20) AS status,
    NULL::uuid AS creator_id;
DROP INDEX public.idx_workshops_status;
DROP INDEX public.idx_workshops_creator_active_sort;
DROP INDEX public.idx_workshops_creator;
DROP INDEX public.idx_workshops_category;
DROP INDEX public.idx_workshops_active_sorted_partial;
DROP INDEX public.idx_workshops_active;
DROP INDEX public.idx_workshop_sessions_workshop;
DROP INDEX public.idx_workshop_sessions_upcoming;
DROP INDEX public.idx_workshop_sessions_status;
DROP INDEX public.idx_workshop_sessions_run_id;
DROP INDEX public.idx_workshop_sessions_parent_run;
DROP INDEX public.idx_workshop_sessions_dates_gin;
DROP INDEX public.idx_workshop_sessions_date;
DROP INDEX public.idx_workshop_runs_workshop;
DROP INDEX public.idx_workshop_runs_dates;
DROP INDEX public.idx_workshop_images_workshop;
DROP INDEX public.idx_workshop_images_cover;
DROP INDEX public.idx_url_settings_creator;
DROP INDEX public.idx_subscriptions_stripe;
DROP INDEX public.idx_subscriptions_status;
DROP INDEX public.idx_subscriptions_creator;
DROP INDEX public.idx_shop_settings_creator;
DROP INDEX public.idx_promo_uses_order;
DROP INDEX public.idx_promo_uses_code;
DROP INDEX public.idx_promo_codes_creator_code;
DROP INDEX public.idx_promo_codes_active;
DROP INDEX public.idx_password_tokens_token;
DROP INDEX public.idx_password_tokens_creator;
DROP INDEX public.idx_orders_status;
DROP INDEX public.idx_orders_pending;
DROP INDEX public.idx_orders_number;
DROP INDEX public.idx_orders_date;
DROP INDEX public.idx_orders_customer_phone;
DROP INDEX public.idx_orders_creator_status_date;
DROP INDEX public.idx_orders_creator;
DROP INDEX public.idx_order_items_workshop;
DROP INDEX public.idx_order_items_session;
DROP INDEX public.idx_order_items_order;
DROP INDEX public.idx_notifications_unread;
DROP INDEX public.idx_notifications_type;
DROP INDEX public.idx_notifications_creator;
DROP INDEX public.idx_enrollments_workshop_status_date;
DROP INDEX public.idx_enrollments_workshop;
DROP INDEX public.idx_enrollments_successful;
DROP INDEX public.idx_enrollments_student_email;
DROP INDEX public.idx_enrollments_status;
DROP INDEX public.idx_enrollments_session;
DROP INDEX public.idx_enrollments_date;
DROP INDEX public.idx_email_tokens_token;
DROP INDEX public.idx_email_tokens_creator;
DROP INDEX public.idx_creators_username;
DROP INDEX public.idx_creators_email;
DROP INDEX public.idx_creators_active;
DROP INDEX public.idx_creator_sessions_token;
DROP INDEX public.idx_creator_sessions_expires;
DROP INDEX public.idx_creator_sessions_creator;
DROP INDEX public.idx_creator_analytics_summary_creator;
DROP INDEX public.idx_categories_creator;
DROP INDEX public.idx_categories_active;
DROP INDEX public.idx_analytics_creator_date;
DROP INDEX public.idx_analytics_clicks_platform;
DROP INDEX public.idx_analytics_clicks_device;
DROP INDEX public.idx_analytics_clicks_date;
DROP INDEX public.idx_analytics_clicks_creator;
DROP INDEX public.idx_analytics_clicks_country;
ALTER TABLE ONLY public.workshops DROP CONSTRAINT workshops_pkey;
ALTER TABLE ONLY public.workshop_sessions DROP CONSTRAINT workshop_sessions_pkey;
ALTER TABLE ONLY public.workshop_runs DROP CONSTRAINT workshop_runs_pkey;
ALTER TABLE ONLY public.workshop_images DROP CONSTRAINT workshop_images_pkey;
ALTER TABLE ONLY public.url_settings DROP CONSTRAINT url_settings_pkey;
ALTER TABLE ONLY public.subscriptions DROP CONSTRAINT subscriptions_pkey;
ALTER TABLE ONLY public.shop_settings DROP CONSTRAINT shop_settings_pkey;
ALTER TABLE ONLY public.promo_codes DROP CONSTRAINT promo_codes_pkey;
ALTER TABLE ONLY public.promo_code_uses DROP CONSTRAINT promo_code_uses_pkey;
ALTER TABLE ONLY public.password_reset_tokens DROP CONSTRAINT password_reset_tokens_token_key;
ALTER TABLE ONLY public.password_reset_tokens DROP CONSTRAINT password_reset_tokens_pkey;
ALTER TABLE ONLY public.orders DROP CONSTRAINT orders_pkey;
ALTER TABLE ONLY public.orders DROP CONSTRAINT orders_order_number_key;
ALTER TABLE ONLY public.order_items DROP CONSTRAINT order_items_pkey;
ALTER TABLE ONLY public.notifications DROP CONSTRAINT notifications_pkey;
ALTER TABLE ONLY public.enrollments DROP CONSTRAINT enrollments_pkey;
ALTER TABLE ONLY public.email_verification_tokens DROP CONSTRAINT email_verification_tokens_token_key;
ALTER TABLE ONLY public.email_verification_tokens DROP CONSTRAINT email_verification_tokens_pkey;
ALTER TABLE ONLY public.creators DROP CONSTRAINT creators_username_key;
ALTER TABLE ONLY public.creators DROP CONSTRAINT creators_pkey;
ALTER TABLE ONLY public.creators DROP CONSTRAINT creators_email_key;
ALTER TABLE ONLY public.creator_sessions DROP CONSTRAINT creator_sessions_session_token_key;
ALTER TABLE ONLY public.creator_sessions DROP CONSTRAINT creator_sessions_pkey;
ALTER TABLE ONLY public.categories DROP CONSTRAINT categories_pkey;
ALTER TABLE ONLY public.analytics_clicks DROP CONSTRAINT analytics_clicks_pkey;
DROP VIEW public.workshop_run_summary;
DROP TABLE public.workshop_images;
DROP TABLE public.url_settings;
DROP TABLE public.subscriptions;
DROP TABLE public.shop_settings;
DROP VIEW public.session_availability;
DROP TABLE public.workshop_sessions;
DROP TABLE public.workshop_runs;
DROP TABLE public.promo_codes;
DROP TABLE public.promo_code_uses;
DROP VIEW public.popular_workshops;
DROP TABLE public.password_reset_tokens;
DROP VIEW public.order_management;
DROP TABLE public.order_items;
DROP TABLE public.notifications;
DROP VIEW public.monthly_revenue_trends;
DROP VIEW public.enrollment_tracking;
DROP TABLE public.email_verification_tokens;
DROP TABLE public.creator_sessions;
DROP VIEW public.creator_dashboard_stats;
DROP MATERIALIZED VIEW public.creator_analytics_summary;
DROP TABLE public.workshops;
DROP TABLE public.orders;
DROP TABLE public.enrollments;
DROP TABLE public.creators;
DROP TABLE public.categories;
DROP VIEW public.analytics_summary;
DROP TABLE public.analytics_clicks;
DROP FUNCTION public.update_workshop_run_attendance();
DROP FUNCTION public.update_updated_at_column();
DROP FUNCTION public.update_session_status();
DROP FUNCTION public.refresh_analytics_summary();
DROP FUNCTION public.get_top_workshops(creator_uuid uuid, limit_count integer);
DROP FUNCTION public.get_session_date_display(session_dates jsonb, lang text);
DROP FUNCTION public.get_analytics_for_period(creator_uuid uuid, start_date date, end_date date);
DROP FUNCTION public.generate_order_number();
DROP FUNCTION public.ensure_workshop_run();
DROP FUNCTION public.clone_workshop_sessions(p_workshop_id uuid, p_new_start_date date, p_run_name character varying);
DROP FUNCTION public.calculate_conversion_rate(creator_uuid uuid, days_back integer);
DROP EXTENSION "uuid-ossp";
--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: calculate_conversion_rate(uuid, integer); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.calculate_conversion_rate(creator_uuid uuid, days_back integer DEFAULT 30) RETURNS numeric
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.calculate_conversion_rate(creator_uuid uuid, days_back integer) OWNER TO most3mr;

--
-- Name: clone_workshop_sessions(uuid, date, character varying); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.clone_workshop_sessions(p_workshop_id uuid, p_new_start_date date, p_run_name character varying DEFAULT NULL::character varying) RETURNS TABLE(new_run_id uuid, sessions_created integer)
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_run_id UUID;
    v_date_offset INTERVAL;
    v_sessions_count INTEGER;
    v_first_session_date DATE;
BEGIN
    -- Generate new run ID
    v_run_id := gen_random_uuid();
    
    -- Get the date offset from the original sessions
    SELECT MIN(session_date) INTO v_first_session_date
    FROM workshop_sessions
    WHERE workshop_id = p_workshop_id
    AND parent_run_id IS NULL OR parent_run_id = (
        SELECT id FROM workshop_runs WHERE workshop_id = p_workshop_id ORDER BY created_at DESC LIMIT 1
    );
    
    v_date_offset := p_new_start_date - v_first_session_date;
    
    -- Create workshop run record
    INSERT INTO workshop_runs (id, workshop_id, run_name, start_date, end_date)
    SELECT 
        v_run_id,
        p_workshop_id,
        COALESCE(p_run_name, 'Session ' || TO_CHAR(p_new_start_date, 'Month YYYY')),
        MIN(session_date + v_date_offset),
        MAX(session_date + v_date_offset)
    FROM workshop_sessions
    WHERE workshop_id = p_workshop_id;
    
    -- Clone sessions with new dates
    INSERT INTO workshop_sessions (
        workshop_id, session_date, start_time, end_time, duration,
        timezone, location, location_ar, max_attendees, session_number,
        parent_run_id, metadata
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
        session_number,
        v_run_id,
        metadata
    FROM workshop_sessions
    WHERE workshop_id = p_workshop_id
    AND (parent_run_id IS NULL OR parent_run_id = (
        SELECT id FROM workshop_runs WHERE workshop_id = p_workshop_id ORDER BY created_at DESC LIMIT 1
    ))
    ORDER BY session_date, start_time;
    
    GET DIAGNOSTICS v_sessions_count = ROW_COUNT;
    
    RETURN QUERY SELECT v_run_id, v_sessions_count;
END;
$$;


ALTER FUNCTION public.clone_workshop_sessions(p_workshop_id uuid, p_new_start_date date, p_run_name character varying) OWNER TO most3mr;

--
-- Name: FUNCTION clone_workshop_sessions(p_workshop_id uuid, p_new_start_date date, p_run_name character varying); Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON FUNCTION public.clone_workshop_sessions(p_workshop_id uuid, p_new_start_date date, p_run_name character varying) IS 'Creates a new run of a workshop by cloning existing sessions with new dates';


--
-- Name: ensure_workshop_run(); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.ensure_workshop_run() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    run_exists BOOLEAN;
    new_run_id UUID;
    workshop_name TEXT;
    workshop_name_ar TEXT;
BEGIN
    -- Check if the run_id exists and belongs to the correct workshop
    IF NEW.run_id IS NOT NULL THEN
        SELECT EXISTS(
            SELECT 1 FROM workshop_runs 
            WHERE id = NEW.run_id AND workshop_id = NEW.workshop_id
        ) INTO run_exists;
        
        -- If run_id is valid, keep it
        IF run_exists THEN
            RETURN NEW;
        END IF;
    END IF;
    
    -- Try to find an existing default run for this workshop
    SELECT id INTO new_run_id
    FROM workshop_runs 
    WHERE workshop_id = NEW.workshop_id 
    AND run_name LIKE 'Default Run -%'
    LIMIT 1;
    
    -- If no default run exists, create one
    IF new_run_id IS NULL THEN
        -- Get workshop details
        SELECT name, COALESCE(title_ar, name) 
        INTO workshop_name, workshop_name_ar
        FROM workshops 
        WHERE id = NEW.workshop_id;
        
        -- Create new run
        new_run_id := gen_random_uuid();
        INSERT INTO workshop_runs (id, workshop_id, run_name, run_name_ar, start_date, end_date, status)
        VALUES (
            new_run_id,
            NEW.workshop_id,
            'Default Run - ' || workshop_name,
            'الدورة الافتراضية - ' || workshop_name_ar,
            NEW.session_date,
            NEW.session_date,
            'upcoming'
        );
    END IF;
    
    -- Set the valid run_id
    NEW.run_id := new_run_id;
    
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.ensure_workshop_run() OWNER TO most3mr;

--
-- Name: generate_order_number(); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.generate_order_number() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.order_number IS NULL THEN
        NEW.order_number := 'WQ' || TO_CHAR(NEW.created_at, 'YYYYMMDD') || '-' || LPAD(CAST((SELECT COALESCE(MAX(CAST(RIGHT(order_number, 4) AS INTEGER)), 0) + 1 FROM orders WHERE order_number LIKE 'WQ' || TO_CHAR(NEW.created_at, 'YYYYMMDD') || '-%') AS TEXT), 4, '0');
    END IF;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.generate_order_number() OWNER TO most3mr;

--
-- Name: get_analytics_for_period(uuid, date, date); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.get_analytics_for_period(creator_uuid uuid, start_date date, end_date date) RETURNS TABLE(total_clicks bigint, unique_visitors bigint, mobile_percentage numeric, top_country character varying, top_platform character varying, conversion_rate numeric)
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.get_analytics_for_period(creator_uuid uuid, start_date date, end_date date) OWNER TO most3mr;

--
-- Name: get_session_date_display(jsonb, text); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.get_session_date_display(session_dates jsonb, lang text DEFAULT 'en'::text) RETURNS text
    LANGUAGE plpgsql
    AS $$
DECLARE
    dates_array DATE[];
    first_date DATE;
    last_date DATE;
    total_count INTEGER;
BEGIN
    -- Convert JSONB array to PostgreSQL array
    SELECT array_agg(value::date ORDER BY value::date)
    INTO dates_array
    FROM jsonb_array_elements_text(session_dates);
    
    total_count := array_length(dates_array, 1);
    
    IF total_count IS NULL OR total_count = 0 THEN
        RETURN '';
    ELSIF total_count = 1 THEN
        -- Single day
        IF lang = 'ar' THEN
            RETURN to_char(dates_array[1], 'DD Mon YYYY');
        ELSE
            RETURN to_char(dates_array[1], 'Mon DD, YYYY');
        END IF;
    ELSE
        -- Multiple days
        first_date := dates_array[1];
        last_date := dates_array[total_count];
        
        IF lang = 'ar' THEN
            RETURN format('%s أيام: %s - %s', 
                total_count,
                to_char(first_date, 'DD Mon'),
                to_char(last_date, 'DD Mon YYYY')
            );
        ELSE
            RETURN format('%s days: %s - %s', 
                total_count,
                to_char(first_date, 'Mon DD'),
                to_char(last_date, 'Mon DD, YYYY')
            );
        END IF;
    END IF;
END;
$$;


ALTER FUNCTION public.get_session_date_display(session_dates jsonb, lang text) OWNER TO most3mr;

--
-- Name: get_top_workshops(uuid, integer); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.get_top_workshops(creator_uuid uuid, limit_count integer DEFAULT 5) RETURNS TABLE(workshop_id uuid, title character varying, title_ar character varying, enrollment_count bigint, revenue numeric)
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.get_top_workshops(creator_uuid uuid, limit_count integer) OWNER TO most3mr;

--
-- Name: refresh_analytics_summary(); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.refresh_analytics_summary() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW creator_analytics_summary;
END;
$$;


ALTER FUNCTION public.refresh_analytics_summary() OWNER TO most3mr;

--
-- Name: update_session_status(); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.update_session_status() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Update status to 'full' when max capacity is reached
    IF NEW.current_attendees >= NEW.max_attendees AND NEW.max_attendees > 0 THEN
        NEW.status = 'full';
        NEW.status_ar = 'ممتلئ';
    END IF;
    
    -- Update updated_at timestamp
    NEW.updated_at = NOW();
    
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_session_status() OWNER TO most3mr;

--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO most3mr;

--
-- Name: update_workshop_run_attendance(); Type: FUNCTION; Schema: public; Owner: most3mr
--

CREATE FUNCTION public.update_workshop_run_attendance() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Update the workshop run's total attendance
    UPDATE workshop_runs
    SET current_attendees = (
        SELECT SUM(current_attendees) 
        FROM workshop_sessions 
        WHERE parent_run_id = NEW.parent_run_id
    )
    WHERE id = NEW.parent_run_id;
    
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_workshop_run_attendance() OWNER TO most3mr;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: analytics_clicks; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.analytics_clicks (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    ip_address inet,
    user_agent text,
    referrer text,
    country character varying(100),
    country_ar character varying(100),
    city character varying(100),
    city_ar character varying(100),
    device character varying(50),
    device_ar character varying(50),
    os character varying(50),
    os_ar character varying(50),
    browser character varying(50),
    browser_ar character varying(50),
    platform character varying(50),
    platform_ar character varying(50),
    clicked_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.analytics_clicks OWNER TO most3mr;

--
-- Name: TABLE analytics_clicks; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON TABLE public.analytics_clicks IS 'Track page visits and click analytics for creators';


--
-- Name: analytics_summary; Type: VIEW; Schema: public; Owner: most3mr
--

CREATE VIEW public.analytics_summary AS
 SELECT analytics_clicks.creator_id,
    date_trunc('day'::text, analytics_clicks.clicked_at) AS date,
    count(*) AS total_clicks,
    count(DISTINCT analytics_clicks.ip_address) AS unique_visitors,
    count(
        CASE
            WHEN ((analytics_clicks.device)::text = 'Mobile'::text) THEN 1
            ELSE NULL::integer
        END) AS mobile_clicks,
    count(
        CASE
            WHEN ((analytics_clicks.device)::text = 'Desktop'::text) THEN 1
            ELSE NULL::integer
        END) AS desktop_clicks,
    count(
        CASE
            WHEN ((analytics_clicks.platform)::text = 'Instagram'::text) THEN 1
            ELSE NULL::integer
        END) AS instagram_clicks,
    count(
        CASE
            WHEN ((analytics_clicks.platform)::text = 'WhatsApp'::text) THEN 1
            ELSE NULL::integer
        END) AS whatsapp_clicks,
    count(
        CASE
            WHEN ((analytics_clicks.platform)::text = 'Direct'::text) THEN 1
            ELSE NULL::integer
        END) AS direct_clicks
   FROM public.analytics_clicks
  GROUP BY analytics_clicks.creator_id, (date_trunc('day'::text, analytics_clicks.clicked_at))
  ORDER BY analytics_clicks.creator_id, (date_trunc('day'::text, analytics_clicks.clicked_at)) DESC;


ALTER TABLE public.analytics_summary OWNER TO most3mr;

--
-- Name: VIEW analytics_summary; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON VIEW public.analytics_summary IS 'Daily analytics summary with device and platform breakdown';


--
-- Name: categories; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.categories (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    name_ar character varying(255),
    description text,
    description_ar text,
    color character varying(7) DEFAULT '#2DD4BF'::character varying,
    sort_order integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.categories OWNER TO most3mr;

--
-- Name: creators; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.creators (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    name_ar character varying(255),
    username character varying(50) NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    avatar text,
    plan character varying(50) DEFAULT 'free'::character varying,
    plan_ar character varying(50) DEFAULT 'مجاني'::character varying,
    is_active boolean DEFAULT true,
    email_verified boolean DEFAULT false,
    email_verified_at timestamp without time zone,
    last_login_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT creators_plan_check CHECK (((plan)::text = ANY ((ARRAY['free'::character varying, 'pro'::character varying])::text[])))
);


ALTER TABLE public.creators OWNER TO most3mr;

--
-- Name: TABLE creators; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON TABLE public.creators IS 'Main creators/users who create workshops and manage their stores';


--
-- Name: enrollments; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.enrollments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    workshop_id uuid NOT NULL,
    order_id uuid,
    session_id uuid,
    student_name character varying(255) NOT NULL,
    student_email character varying(255),
    student_phone character varying(20),
    total_price numeric(10,3) NOT NULL,
    status character varying(20) DEFAULT 'successful'::character varying,
    status_ar character varying(20) DEFAULT 'مكتمل'::character varying,
    enrollment_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    completion_status character varying(20) DEFAULT 'enrolled'::character varying,
    notes text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT enrollments_completion_status_check CHECK (((completion_status)::text = ANY ((ARRAY['enrolled'::character varying, 'in_progress'::character varying, 'completed'::character varying, 'dropped'::character varying])::text[]))),
    CONSTRAINT enrollments_status_check CHECK (((status)::text = ANY ((ARRAY['successful'::character varying, 'pending'::character varying, 'rejected'::character varying, 'cancelled'::character varying])::text[])))
);


ALTER TABLE public.enrollments OWNER TO most3mr;

--
-- Name: TABLE enrollments; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON TABLE public.enrollments IS 'Successful enrollments/registrations for workshops';


--
-- Name: orders; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.orders (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    order_number character varying(50),
    customer_name character varying(255) NOT NULL,
    customer_phone character varying(20) NOT NULL,
    customer_email character varying(255),
    total_amount numeric(10,3) NOT NULL,
    currency character varying(3) DEFAULT 'KWD'::character varying,
    status character varying(20) DEFAULT 'pending'::character varying,
    status_ar character varying(20) DEFAULT 'قيد الانتظار'::character varying,
    payment_method character varying(50),
    payment_reference character varying(255),
    order_source character varying(20) DEFAULT 'whatsapp'::character varying,
    notes text,
    is_viewed boolean DEFAULT false,
    viewed_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT orders_order_source_check CHECK (((order_source)::text = ANY ((ARRAY['whatsapp'::character varying, 'direct'::character varying, 'qr'::character varying])::text[]))),
    CONSTRAINT orders_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'paid'::character varying, 'cancelled'::character varying, 'refunded'::character varying])::text[]))),
    CONSTRAINT orders_total_check CHECK ((total_amount >= (0)::numeric))
);


ALTER TABLE public.orders OWNER TO most3mr;

--
-- Name: TABLE orders; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON TABLE public.orders IS 'Customer orders placed through WhatsApp or direct booking';


--
-- Name: workshops; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.workshops (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    category_id uuid,
    name character varying(255) NOT NULL,
    title character varying(255) NOT NULL,
    title_ar character varying(255),
    description text,
    description_ar text,
    price numeric(10,3) DEFAULT 0.000,
    currency character varying(3) DEFAULT 'KWD'::character varying,
    duration integer,
    max_students integer DEFAULT 0,
    status character varying(20) DEFAULT 'draft'::character varying,
    is_active boolean DEFAULT false,
    is_free boolean DEFAULT false,
    is_recurring boolean DEFAULT false,
    recurrence_type character varying(20),
    sort_order integer DEFAULT 0,
    view_count integer DEFAULT 0,
    enrollment_count integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    workshop_type character varying(20) DEFAULT 'single'::character varying,
    CONSTRAINT workshops_price_check CHECK ((price >= (0)::numeric)),
    CONSTRAINT workshops_recurrence_type_check CHECK (((recurrence_type IS NULL) OR ((recurrence_type)::text = ANY ((ARRAY['weekly'::character varying, 'monthly'::character varying, 'yearly'::character varying])::text[])))),
    CONSTRAINT workshops_status_check CHECK (((status)::text = ANY ((ARRAY['draft'::character varying, 'published'::character varying, 'archived'::character varying])::text[]))),
    CONSTRAINT workshops_workshop_type_check CHECK (((workshop_type)::text = ANY ((ARRAY['single'::character varying, 'consecutive'::character varying, 'spread'::character varying, 'custom'::character varying])::text[])))
);


ALTER TABLE public.workshops OWNER TO most3mr;

--
-- Name: TABLE workshops; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON TABLE public.workshops IS 'Individual workshops or courses offered by creators';


--
-- Name: creator_analytics_summary; Type: MATERIALIZED VIEW; Schema: public; Owner: most3mr
--

CREATE MATERIALIZED VIEW public.creator_analytics_summary AS
 SELECT c.id AS creator_id,
    c.username,
    count(DISTINCT ac.id) AS total_clicks,
    count(DISTINCT date(ac.clicked_at)) AS active_days,
    count(DISTINCT ac.ip_address) AS unique_visitors,
    count(DISTINCT w.id) AS total_workshops,
    count(DISTINCT
        CASE
            WHEN w.is_active THEN w.id
            ELSE NULL::uuid
        END) AS active_workshops,
    count(DISTINCT e.id) AS total_enrollments,
    COALESCE(sum(
        CASE
            WHEN ((e.status)::text = 'successful'::text) THEN e.total_price
            ELSE NULL::numeric
        END), (0)::numeric) AS total_revenue,
    count(DISTINCT o.id) AS total_orders,
    count(DISTINCT
        CASE
            WHEN ((o.status)::text = 'pending'::text) THEN o.id
            ELSE NULL::uuid
        END) AS pending_orders
   FROM ((((public.creators c
     LEFT JOIN public.analytics_clicks ac ON ((c.id = ac.creator_id)))
     LEFT JOIN public.workshops w ON ((c.id = w.creator_id)))
     LEFT JOIN public.enrollments e ON ((w.id = e.workshop_id)))
     LEFT JOIN public.orders o ON ((c.id = o.creator_id)))
  GROUP BY c.id, c.username
  WITH NO DATA;


ALTER TABLE public.creator_analytics_summary OWNER TO most3mr;

--
-- Name: MATERIALIZED VIEW creator_analytics_summary; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON MATERIALIZED VIEW public.creator_analytics_summary IS 'High-level creator analytics (refresh periodically for performance)';


--
-- Name: creator_dashboard_stats; Type: VIEW; Schema: public; Owner: most3mr
--

CREATE VIEW public.creator_dashboard_stats AS
 SELECT c.id AS creator_id,
    count(DISTINCT w.id) AS total_workshops,
    count(DISTINCT
        CASE
            WHEN (w.is_active = true) THEN w.id
            ELSE NULL::uuid
        END) AS active_workshops,
    count(DISTINCT e.id) AS total_enrollments,
    COALESCE(sum(
        CASE
            WHEN (((e.status)::text = 'successful'::text) AND (e.enrollment_date >= date_trunc('month'::text, (CURRENT_DATE)::timestamp with time zone))) THEN e.total_price
            ELSE NULL::numeric
        END), (0)::numeric) AS monthly_revenue,
    (COALESCE(sum(
        CASE
            WHEN ((w.is_active = true) AND (w.max_students > 0)) THEN (w.price * (w.max_students)::numeric)
            ELSE (0)::numeric
        END), (0)::numeric) * 0.7) AS projected_sales,
    (COALESCE(sum(
        CASE
            WHEN (w.is_active = true) THEN w.max_students
            ELSE NULL::integer
        END), (0)::bigint) - count(DISTINCT e.id)) AS remaining_seats
   FROM ((public.creators c
     LEFT JOIN public.workshops w ON ((c.id = w.creator_id)))
     LEFT JOIN public.enrollments e ON ((w.id = e.workshop_id)))
  GROUP BY c.id;


ALTER TABLE public.creator_dashboard_stats OWNER TO most3mr;

--
-- Name: VIEW creator_dashboard_stats; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON VIEW public.creator_dashboard_stats IS 'Aggregated statistics for creator dashboard display';


--
-- Name: creator_sessions; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.creator_sessions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    session_token character varying(255) NOT NULL,
    device_info text,
    ip_address inet,
    expires_at timestamp without time zone NOT NULL,
    last_activity timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.creator_sessions OWNER TO most3mr;

--
-- Name: email_verification_tokens; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.email_verification_tokens (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    token character varying(255) NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    used_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.email_verification_tokens OWNER TO most3mr;

--
-- Name: enrollment_tracking; Type: VIEW; Schema: public; Owner: most3mr
--

CREATE VIEW public.enrollment_tracking AS
 SELECT e.id,
    e.workshop_id,
    w.title AS workshop_name,
    w.title_ar AS workshop_name_ar,
    e.student_name,
    e.student_email,
    e.student_phone,
    e.total_price,
    e.status,
    e.status_ar,
    e.enrollment_date,
    e.completion_status,
    w.creator_id
   FROM (public.enrollments e
     JOIN public.workshops w ON ((e.workshop_id = w.id)))
  ORDER BY e.enrollment_date DESC;


ALTER TABLE public.enrollment_tracking OWNER TO most3mr;

--
-- Name: VIEW enrollment_tracking; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON VIEW public.enrollment_tracking IS 'Detailed enrollment tracking with workshop information';


--
-- Name: monthly_revenue_trends; Type: VIEW; Schema: public; Owner: most3mr
--

CREATE VIEW public.monthly_revenue_trends AS
 SELECT w.creator_id,
    date_trunc('month'::text, e.enrollment_date) AS month,
    count(e.id) AS enrollments_count,
    sum(e.total_price) AS total_revenue,
    avg(e.total_price) AS avg_order_value
   FROM (public.enrollments e
     JOIN public.workshops w ON ((e.workshop_id = w.id)))
  WHERE ((e.status)::text = 'successful'::text)
  GROUP BY w.creator_id, (date_trunc('month'::text, e.enrollment_date))
  ORDER BY w.creator_id, (date_trunc('month'::text, e.enrollment_date));


ALTER TABLE public.monthly_revenue_trends OWNER TO most3mr;

--
-- Name: VIEW monthly_revenue_trends; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON VIEW public.monthly_revenue_trends IS 'Monthly revenue and enrollment trends for reporting';


--
-- Name: notifications; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.notifications (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    type character varying(50) NOT NULL,
    title character varying(255) NOT NULL,
    title_ar character varying(255),
    message text NOT NULL,
    message_ar text,
    data jsonb,
    is_read boolean DEFAULT false,
    read_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.notifications OWNER TO most3mr;

--
-- Name: order_items; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.order_items (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    order_id uuid NOT NULL,
    workshop_id uuid NOT NULL,
    workshop_name character varying(255) NOT NULL,
    workshop_name_ar character varying(255),
    price numeric(10,3) NOT NULL,
    quantity integer DEFAULT 1,
    subtotal numeric(10,3) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    session_id uuid,
    run_id uuid,
    workshop_run_id uuid
);


ALTER TABLE public.order_items OWNER TO most3mr;

--
-- Name: order_management; Type: VIEW; Schema: public; Owner: most3mr
--

CREATE VIEW public.order_management AS
 SELECT o.id,
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
    count(oi.id) AS item_count,
    string_agg((oi.workshop_name)::text, ', '::text) AS workshop_names,
    string_agg((oi.workshop_name_ar)::text, '، '::text) AS workshop_names_ar
   FROM (public.orders o
     LEFT JOIN public.order_items oi ON ((o.id = oi.order_id)))
  GROUP BY o.id, o.order_number, o.creator_id, o.customer_name, o.customer_phone, o.customer_email, o.total_amount, o.status, o.status_ar, o.order_source, o.created_at, o.is_viewed
  ORDER BY o.created_at DESC;


ALTER TABLE public.order_management OWNER TO most3mr;

--
-- Name: VIEW order_management; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON VIEW public.order_management IS 'Comprehensive order details for management interface';


--
-- Name: password_reset_tokens; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.password_reset_tokens (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    token character varying(255) NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    used_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.password_reset_tokens OWNER TO most3mr;

--
-- Name: popular_workshops; Type: VIEW; Schema: public; Owner: most3mr
--

CREATE VIEW public.popular_workshops AS
 SELECT w.id,
    w.creator_id,
    w.title,
    w.title_ar,
    w.price,
    count(e.id) AS enrollment_count,
    avg(w.price) AS avg_price,
    sum(e.total_price) AS total_revenue
   FROM (public.workshops w
     LEFT JOIN public.enrollments e ON (((w.id = e.workshop_id) AND ((e.status)::text = 'successful'::text))))
  WHERE (w.is_active = true)
  GROUP BY w.id, w.creator_id, w.title, w.title_ar, w.price
  ORDER BY (count(e.id)) DESC;


ALTER TABLE public.popular_workshops OWNER TO most3mr;

--
-- Name: VIEW popular_workshops; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON VIEW public.popular_workshops IS 'Workshop popularity rankings by enrollment and revenue';


--
-- Name: promo_code_uses; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.promo_code_uses (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    promo_code_id uuid NOT NULL,
    order_id uuid NOT NULL,
    discount_amount numeric(10,3) NOT NULL,
    used_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.promo_code_uses OWNER TO most3mr;

--
-- Name: promo_codes; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.promo_codes (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    code character varying(50) NOT NULL,
    type character varying(20) DEFAULT 'percentage'::character varying,
    value numeric(10,3) NOT NULL,
    min_order_amount numeric(10,3) DEFAULT 0,
    max_uses integer DEFAULT 0,
    current_uses integer DEFAULT 0,
    starts_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamp without time zone,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT promo_codes_type_check CHECK (((type)::text = ANY ((ARRAY['percentage'::character varying, 'fixed'::character varying])::text[])))
);


ALTER TABLE public.promo_codes OWNER TO most3mr;

--
-- Name: workshop_runs; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.workshop_runs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    workshop_id uuid NOT NULL,
    run_name character varying(255),
    run_name_ar character varying(255),
    start_date date NOT NULL,
    end_date date NOT NULL,
    status character varying(20) DEFAULT 'upcoming'::character varying,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT workshop_runs_status_check CHECK (((status)::text = ANY ((ARRAY['upcoming'::character varying, 'active'::character varying, 'full'::character varying, 'completed'::character varying, 'cancelled'::character varying])::text[])))
);


ALTER TABLE public.workshop_runs OWNER TO most3mr;

--
-- Name: TABLE workshop_runs; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON TABLE public.workshop_runs IS 'Groups multiple sessions of the same workshop run together';


--
-- Name: workshop_sessions; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.workshop_sessions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    workshop_id uuid NOT NULL,
    session_date date NOT NULL,
    start_time time without time zone NOT NULL,
    end_time time without time zone,
    duration numeric(4,2),
    timezone character varying(50) DEFAULT 'Asia/Kuwait'::character varying,
    location text,
    location_ar text,
    max_attendees integer,
    current_attendees integer DEFAULT 0,
    is_completed boolean DEFAULT false,
    notes text,
    notes_ar text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    status character varying(20) DEFAULT 'upcoming'::character varying,
    status_ar character varying(50) DEFAULT 'قادم'::character varying,
    session_number integer DEFAULT 1,
    run_id uuid,
    metadata jsonb DEFAULT '{}'::jsonb,
    parent_run_id uuid,
    end_date date,
    day_count integer DEFAULT 1,
    session_dates jsonb,
    total_days integer DEFAULT 1,
    CONSTRAINT workshop_sessions_status_check CHECK (((status)::text = ANY ((ARRAY['upcoming'::character varying, 'active'::character varying, 'full'::character varying, 'completed'::character varying, 'cancelled'::character varying])::text[])))
);


ALTER TABLE public.workshop_sessions OWNER TO most3mr;

--
-- Name: COLUMN workshop_sessions.status; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON COLUMN public.workshop_sessions.status IS 'Current status of the session';


--
-- Name: COLUMN workshop_sessions.session_number; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON COLUMN public.workshop_sessions.session_number IS 'Sequential number for multi-day workshops (Day 1, 2, etc.)';


--
-- Name: COLUMN workshop_sessions.run_id; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON COLUMN public.workshop_sessions.run_id IS 'Links session to a specific run/batch of the workshop';


--
-- Name: session_availability; Type: VIEW; Schema: public; Owner: most3mr
--

CREATE VIEW public.session_availability AS
 SELECT ws.id AS session_id,
    ws.workshop_id,
    w.name AS workshop_name,
    COALESCE(w.title_ar, w.name) AS workshop_name_ar,
    ws.session_date,
    ws.session_dates,
    ws.total_days,
    ws.start_time,
    ws.end_time,
    ws.max_attendees,
    ws.current_attendees,
    GREATEST(0, (ws.max_attendees - ws.current_attendees)) AS available_seats,
        CASE
            WHEN ((ws.status)::text = 'cancelled'::text) THEN 'cancelled'::character varying
            WHEN (ws.session_date < CURRENT_DATE) THEN 'completed'::character varying
            WHEN ((ws.current_attendees >= ws.max_attendees) AND (ws.max_attendees > 0)) THEN 'full'::character varying
            ELSE COALESCE(ws.status, 'upcoming'::character varying)
        END AS calculated_status,
    ws.run_id,
    wr.run_name,
    c.name AS creator_name,
    c.username AS creator_username
   FROM (((public.workshop_sessions ws
     JOIN public.workshops w ON ((ws.workshop_id = w.id)))
     JOIN public.creators c ON ((w.creator_id = c.id)))
     LEFT JOIN public.workshop_runs wr ON ((ws.run_id = wr.id)))
  WHERE (w.is_active = true);


ALTER TABLE public.session_availability OWNER TO most3mr;

--
-- Name: shop_settings; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.shop_settings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    logo_url text,
    creator_name character varying(255),
    creator_name_ar character varying(255),
    sub_header text,
    sub_header_ar text,
    enrollment_whatsapp character varying(20),
    contact_whatsapp character varying(20),
    checkout_language character varying(10) DEFAULT 'both'::character varying,
    greeting_message text,
    greeting_message_ar text,
    currency_symbol character varying(10) DEFAULT 'KD'::character varying,
    currency_symbol_ar character varying(10) DEFAULT 'د.ك'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT shop_settings_checkout_language_check CHECK (((checkout_language)::text = ANY ((ARRAY['ar'::character varying, 'en'::character varying, 'both'::character varying])::text[])))
);


ALTER TABLE public.shop_settings OWNER TO most3mr;

--
-- Name: TABLE shop_settings; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON TABLE public.shop_settings IS 'Customizable branding and settings for each creator''s store';


--
-- Name: subscriptions; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.subscriptions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    plan character varying(50) NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying,
    current_period_start timestamp without time zone NOT NULL,
    current_period_end timestamp without time zone NOT NULL,
    cancel_at_period_end boolean DEFAULT false,
    stripe_subscription_id character varying(255),
    stripe_customer_id character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT subscriptions_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'cancelled'::character varying, 'past_due'::character varying, 'trialing'::character varying])::text[])))
);


ALTER TABLE public.subscriptions OWNER TO most3mr;

--
-- Name: url_settings; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.url_settings (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    creator_id uuid NOT NULL,
    username character varying(50) NOT NULL,
    changes_used integer DEFAULT 0,
    max_changes integer DEFAULT 5,
    last_changed timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT url_settings_changes_check CHECK ((changes_used <= max_changes))
);


ALTER TABLE public.url_settings OWNER TO most3mr;

--
-- Name: TABLE url_settings; Type: COMMENT; Schema: public; Owner: most3mr
--

COMMENT ON TABLE public.url_settings IS 'Manage custom usernames and URL changes for creators';


--
-- Name: workshop_images; Type: TABLE; Schema: public; Owner: most3mr
--

CREATE TABLE public.workshop_images (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    workshop_id uuid NOT NULL,
    image_url text NOT NULL,
    is_cover boolean DEFAULT false,
    sort_order integer DEFAULT 0,
    alt_text character varying(255),
    alt_text_ar character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.workshop_images OWNER TO most3mr;

--
-- Name: workshop_run_summary; Type: VIEW; Schema: public; Owner: most3mr
--

CREATE VIEW public.workshop_run_summary AS
SELECT
    NULL::uuid AS run_id,
    NULL::uuid AS workshop_id,
    NULL::character varying(255) AS workshop_name,
    NULL::character varying(255) AS workshop_name_ar,
    NULL::character varying(255) AS run_name,
    NULL::date AS start_date,
    NULL::date AS end_date,
    NULL::bigint AS total_sessions,
    NULL::bigint AS total_capacity,
    NULL::bigint AS total_enrolled,
    NULL::bigint AS full_sessions,
    NULL::character varying(20) AS status,
    NULL::uuid AS creator_id;


ALTER TABLE public.workshop_run_summary OWNER TO most3mr;

--
-- Name: analytics_clicks analytics_clicks_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.analytics_clicks
    ADD CONSTRAINT analytics_clicks_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: creator_sessions creator_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.creator_sessions
    ADD CONSTRAINT creator_sessions_pkey PRIMARY KEY (id);


--
-- Name: creator_sessions creator_sessions_session_token_key; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.creator_sessions
    ADD CONSTRAINT creator_sessions_session_token_key UNIQUE (session_token);


--
-- Name: creators creators_email_key; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.creators
    ADD CONSTRAINT creators_email_key UNIQUE (email);


--
-- Name: creators creators_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.creators
    ADD CONSTRAINT creators_pkey PRIMARY KEY (id);


--
-- Name: creators creators_username_key; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.creators
    ADD CONSTRAINT creators_username_key UNIQUE (username);


--
-- Name: email_verification_tokens email_verification_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.email_verification_tokens
    ADD CONSTRAINT email_verification_tokens_pkey PRIMARY KEY (id);


--
-- Name: email_verification_tokens email_verification_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.email_verification_tokens
    ADD CONSTRAINT email_verification_tokens_token_key UNIQUE (token);


--
-- Name: enrollments enrollments_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_pkey PRIMARY KEY (id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: orders orders_order_number_key; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_order_number_key UNIQUE (order_number);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: password_reset_tokens password_reset_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_pkey PRIMARY KEY (id);


--
-- Name: password_reset_tokens password_reset_tokens_token_key; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_token_key UNIQUE (token);


--
-- Name: promo_code_uses promo_code_uses_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.promo_code_uses
    ADD CONSTRAINT promo_code_uses_pkey PRIMARY KEY (id);


--
-- Name: promo_codes promo_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.promo_codes
    ADD CONSTRAINT promo_codes_pkey PRIMARY KEY (id);


--
-- Name: shop_settings shop_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.shop_settings
    ADD CONSTRAINT shop_settings_pkey PRIMARY KEY (id);


--
-- Name: subscriptions subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_pkey PRIMARY KEY (id);


--
-- Name: url_settings url_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.url_settings
    ADD CONSTRAINT url_settings_pkey PRIMARY KEY (id);


--
-- Name: workshop_images workshop_images_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.workshop_images
    ADD CONSTRAINT workshop_images_pkey PRIMARY KEY (id);


--
-- Name: workshop_runs workshop_runs_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.workshop_runs
    ADD CONSTRAINT workshop_runs_pkey PRIMARY KEY (id);


--
-- Name: workshop_sessions workshop_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.workshop_sessions
    ADD CONSTRAINT workshop_sessions_pkey PRIMARY KEY (id);


--
-- Name: workshops workshops_pkey; Type: CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.workshops
    ADD CONSTRAINT workshops_pkey PRIMARY KEY (id);


--
-- Name: idx_analytics_clicks_country; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_analytics_clicks_country ON public.analytics_clicks USING btree (creator_id, country);


--
-- Name: idx_analytics_clicks_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_analytics_clicks_creator ON public.analytics_clicks USING btree (creator_id);


--
-- Name: idx_analytics_clicks_date; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_analytics_clicks_date ON public.analytics_clicks USING btree (clicked_at);


--
-- Name: idx_analytics_clicks_device; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_analytics_clicks_device ON public.analytics_clicks USING btree (creator_id, device);


--
-- Name: idx_analytics_clicks_platform; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_analytics_clicks_platform ON public.analytics_clicks USING btree (creator_id, platform);


--
-- Name: idx_analytics_creator_date; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_analytics_creator_date ON public.analytics_clicks USING btree (creator_id, clicked_at DESC);


--
-- Name: idx_categories_active; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_categories_active ON public.categories USING btree (creator_id, is_active);


--
-- Name: idx_categories_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_categories_creator ON public.categories USING btree (creator_id);


--
-- Name: idx_creator_analytics_summary_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE UNIQUE INDEX idx_creator_analytics_summary_creator ON public.creator_analytics_summary USING btree (creator_id);


--
-- Name: idx_creator_sessions_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_creator_sessions_creator ON public.creator_sessions USING btree (creator_id);


--
-- Name: idx_creator_sessions_expires; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_creator_sessions_expires ON public.creator_sessions USING btree (expires_at);


--
-- Name: idx_creator_sessions_token; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_creator_sessions_token ON public.creator_sessions USING btree (session_token);


--
-- Name: idx_creators_active; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_creators_active ON public.creators USING btree (is_active);


--
-- Name: idx_creators_email; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_creators_email ON public.creators USING btree (email);


--
-- Name: idx_creators_username; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_creators_username ON public.creators USING btree (username);


--
-- Name: idx_email_tokens_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_email_tokens_creator ON public.email_verification_tokens USING btree (creator_id);


--
-- Name: idx_email_tokens_token; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_email_tokens_token ON public.email_verification_tokens USING btree (token);


--
-- Name: idx_enrollments_date; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_enrollments_date ON public.enrollments USING btree (enrollment_date);


--
-- Name: idx_enrollments_session; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_enrollments_session ON public.enrollments USING btree (session_id);


--
-- Name: idx_enrollments_status; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_enrollments_status ON public.enrollments USING btree (status);


--
-- Name: idx_enrollments_student_email; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_enrollments_student_email ON public.enrollments USING btree (student_email);


--
-- Name: idx_enrollments_successful; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_enrollments_successful ON public.enrollments USING btree (workshop_id, enrollment_date DESC) WHERE ((status)::text = 'successful'::text);


--
-- Name: idx_enrollments_workshop; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_enrollments_workshop ON public.enrollments USING btree (workshop_id);


--
-- Name: idx_enrollments_workshop_status_date; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_enrollments_workshop_status_date ON public.enrollments USING btree (workshop_id, status, enrollment_date DESC);


--
-- Name: idx_notifications_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_notifications_creator ON public.notifications USING btree (creator_id);


--
-- Name: idx_notifications_type; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_notifications_type ON public.notifications USING btree (creator_id, type);


--
-- Name: idx_notifications_unread; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_notifications_unread ON public.notifications USING btree (creator_id, is_read);


--
-- Name: idx_order_items_order; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_order_items_order ON public.order_items USING btree (order_id);


--
-- Name: idx_order_items_session; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_order_items_session ON public.order_items USING btree (session_id);


--
-- Name: idx_order_items_workshop; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_order_items_workshop ON public.order_items USING btree (workshop_id);


--
-- Name: idx_orders_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_orders_creator ON public.orders USING btree (creator_id);


--
-- Name: idx_orders_creator_status_date; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_orders_creator_status_date ON public.orders USING btree (creator_id, status, created_at DESC);


--
-- Name: idx_orders_customer_phone; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_orders_customer_phone ON public.orders USING btree (customer_phone);


--
-- Name: idx_orders_date; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_orders_date ON public.orders USING btree (created_at);


--
-- Name: idx_orders_number; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE UNIQUE INDEX idx_orders_number ON public.orders USING btree (order_number) WHERE (order_number IS NOT NULL);


--
-- Name: idx_orders_pending; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_orders_pending ON public.orders USING btree (creator_id, created_at DESC) WHERE ((status)::text = 'pending'::text);


--
-- Name: idx_orders_status; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_orders_status ON public.orders USING btree (creator_id, status);


--
-- Name: idx_password_tokens_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_password_tokens_creator ON public.password_reset_tokens USING btree (creator_id);


--
-- Name: idx_password_tokens_token; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_password_tokens_token ON public.password_reset_tokens USING btree (token);


--
-- Name: idx_promo_codes_active; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_promo_codes_active ON public.promo_codes USING btree (is_active, expires_at);


--
-- Name: idx_promo_codes_creator_code; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE UNIQUE INDEX idx_promo_codes_creator_code ON public.promo_codes USING btree (creator_id, code);


--
-- Name: idx_promo_uses_code; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_promo_uses_code ON public.promo_code_uses USING btree (promo_code_id);


--
-- Name: idx_promo_uses_order; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_promo_uses_order ON public.promo_code_uses USING btree (order_id);


--
-- Name: idx_shop_settings_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE UNIQUE INDEX idx_shop_settings_creator ON public.shop_settings USING btree (creator_id);


--
-- Name: idx_subscriptions_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_subscriptions_creator ON public.subscriptions USING btree (creator_id);


--
-- Name: idx_subscriptions_status; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_subscriptions_status ON public.subscriptions USING btree (status);


--
-- Name: idx_subscriptions_stripe; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_subscriptions_stripe ON public.subscriptions USING btree (stripe_subscription_id);


--
-- Name: idx_url_settings_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE UNIQUE INDEX idx_url_settings_creator ON public.url_settings USING btree (creator_id);


--
-- Name: idx_workshop_images_cover; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_images_cover ON public.workshop_images USING btree (workshop_id, is_cover);


--
-- Name: idx_workshop_images_workshop; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_images_workshop ON public.workshop_images USING btree (workshop_id);


--
-- Name: idx_workshop_runs_dates; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_runs_dates ON public.workshop_runs USING btree (start_date, end_date);


--
-- Name: idx_workshop_runs_workshop; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_runs_workshop ON public.workshop_runs USING btree (workshop_id);


--
-- Name: idx_workshop_sessions_date; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_sessions_date ON public.workshop_sessions USING btree (session_date);


--
-- Name: idx_workshop_sessions_dates_gin; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_sessions_dates_gin ON public.workshop_sessions USING gin (session_dates);


--
-- Name: idx_workshop_sessions_parent_run; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_sessions_parent_run ON public.workshop_sessions USING btree (parent_run_id);


--
-- Name: idx_workshop_sessions_run_id; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_sessions_run_id ON public.workshop_sessions USING btree (run_id);


--
-- Name: idx_workshop_sessions_status; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_sessions_status ON public.workshop_sessions USING btree (status);


--
-- Name: idx_workshop_sessions_upcoming; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_sessions_upcoming ON public.workshop_sessions USING btree (workshop_id, session_date, is_completed);


--
-- Name: idx_workshop_sessions_workshop; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshop_sessions_workshop ON public.workshop_sessions USING btree (workshop_id);


--
-- Name: idx_workshops_active; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshops_active ON public.workshops USING btree (creator_id, is_active);


--
-- Name: idx_workshops_active_sorted_partial; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshops_active_sorted_partial ON public.workshops USING btree (creator_id, sort_order) WHERE (is_active = true);


--
-- Name: idx_workshops_category; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshops_category ON public.workshops USING btree (category_id);


--
-- Name: idx_workshops_creator; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshops_creator ON public.workshops USING btree (creator_id);


--
-- Name: idx_workshops_creator_active_sort; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshops_creator_active_sort ON public.workshops USING btree (creator_id, is_active, sort_order);


--
-- Name: idx_workshops_status; Type: INDEX; Schema: public; Owner: most3mr
--

CREATE INDEX idx_workshops_status ON public.workshops USING btree (creator_id, status);


--
-- Name: workshop_run_summary _RETURN; Type: RULE; Schema: public; Owner: most3mr
--

CREATE OR REPLACE VIEW public.workshop_run_summary AS
 SELECT wr.id AS run_id,
    wr.workshop_id,
    w.name AS workshop_name,
    COALESCE(w.title_ar, w.name) AS workshop_name_ar,
    wr.run_name,
    wr.start_date,
    wr.end_date,
    count(ws.id) AS total_sessions,
    sum(ws.max_attendees) AS total_capacity,
    sum(ws.current_attendees) AS total_enrolled,
    sum(
        CASE
            WHEN ((ws.status)::text = 'full'::text) THEN 1
            ELSE 0
        END) AS full_sessions,
    wr.status,
    w.creator_id
   FROM ((public.workshop_runs wr
     JOIN public.workshops w ON ((wr.workshop_id = w.id)))
     LEFT JOIN public.workshop_sessions ws ON ((ws.run_id = wr.id)))
  GROUP BY wr.id, w.id;


--
-- Name: workshop_sessions ensure_workshop_run_trigger; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER ensure_workshop_run_trigger BEFORE INSERT OR UPDATE ON public.workshop_sessions FOR EACH ROW EXECUTE FUNCTION public.ensure_workshop_run();


--
-- Name: orders generate_order_number_trigger; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER generate_order_number_trigger BEFORE INSERT ON public.orders FOR EACH ROW EXECUTE FUNCTION public.generate_order_number();


--
-- Name: workshop_sessions trigger_update_run_attendance; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER trigger_update_run_attendance AFTER UPDATE OF current_attendees ON public.workshop_sessions FOR EACH ROW WHEN ((new.parent_run_id IS NOT NULL)) EXECUTE FUNCTION public.update_workshop_run_attendance();


--
-- Name: workshop_sessions trigger_update_session_status; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER trigger_update_session_status BEFORE UPDATE OF current_attendees ON public.workshop_sessions FOR EACH ROW EXECUTE FUNCTION public.update_session_status();


--
-- Name: categories update_categories_updated_at; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: creators update_creators_updated_at; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER update_creators_updated_at BEFORE UPDATE ON public.creators FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: enrollments update_enrollments_updated_at; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER update_enrollments_updated_at BEFORE UPDATE ON public.enrollments FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: orders update_orders_updated_at; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON public.orders FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: promo_codes update_promo_codes_updated_at; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER update_promo_codes_updated_at BEFORE UPDATE ON public.promo_codes FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: shop_settings update_shop_settings_updated_at; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER update_shop_settings_updated_at BEFORE UPDATE ON public.shop_settings FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: subscriptions update_subscriptions_updated_at; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON public.subscriptions FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: url_settings update_url_settings_updated_at; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER update_url_settings_updated_at BEFORE UPDATE ON public.url_settings FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: workshop_sessions update_workshop_sessions_updated_at; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER update_workshop_sessions_updated_at BEFORE UPDATE ON public.workshop_sessions FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: workshops update_workshops_updated_at; Type: TRIGGER; Schema: public; Owner: most3mr
--

CREATE TRIGGER update_workshops_updated_at BEFORE UPDATE ON public.workshops FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: analytics_clicks analytics_clicks_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.analytics_clicks
    ADD CONSTRAINT analytics_clicks_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: categories categories_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: creator_sessions creator_sessions_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.creator_sessions
    ADD CONSTRAINT creator_sessions_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: email_verification_tokens email_verification_tokens_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.email_verification_tokens
    ADD CONSTRAINT email_verification_tokens_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: enrollments enrollments_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE SET NULL;


--
-- Name: enrollments enrollments_session_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_session_id_fkey FOREIGN KEY (session_id) REFERENCES public.workshop_sessions(id) ON DELETE SET NULL;


--
-- Name: enrollments enrollments_workshop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT enrollments_workshop_id_fkey FOREIGN KEY (workshop_id) REFERENCES public.workshops(id) ON DELETE RESTRICT;


--
-- Name: enrollments fk_enrollment_session; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.enrollments
    ADD CONSTRAINT fk_enrollment_session FOREIGN KEY (session_id) REFERENCES public.workshop_sessions(id);


--
-- Name: notifications notifications_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_run_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_run_id_fkey FOREIGN KEY (run_id) REFERENCES public.workshop_runs(id);


--
-- Name: order_items order_items_session_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_session_id_fkey FOREIGN KEY (session_id) REFERENCES public.workshop_sessions(id);


--
-- Name: order_items order_items_workshop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_workshop_id_fkey FOREIGN KEY (workshop_id) REFERENCES public.workshops(id) ON DELETE RESTRICT;


--
-- Name: order_items order_items_workshop_run_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_workshop_run_id_fkey FOREIGN KEY (workshop_run_id) REFERENCES public.workshop_runs(id);


--
-- Name: orders orders_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: password_reset_tokens password_reset_tokens_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.password_reset_tokens
    ADD CONSTRAINT password_reset_tokens_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: promo_code_uses promo_code_uses_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.promo_code_uses
    ADD CONSTRAINT promo_code_uses_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: promo_code_uses promo_code_uses_promo_code_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.promo_code_uses
    ADD CONSTRAINT promo_code_uses_promo_code_id_fkey FOREIGN KEY (promo_code_id) REFERENCES public.promo_codes(id) ON DELETE CASCADE;


--
-- Name: promo_codes promo_codes_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.promo_codes
    ADD CONSTRAINT promo_codes_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: shop_settings shop_settings_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.shop_settings
    ADD CONSTRAINT shop_settings_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: subscriptions subscriptions_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: url_settings url_settings_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.url_settings
    ADD CONSTRAINT url_settings_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: workshop_images workshop_images_workshop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.workshop_images
    ADD CONSTRAINT workshop_images_workshop_id_fkey FOREIGN KEY (workshop_id) REFERENCES public.workshops(id) ON DELETE CASCADE;


--
-- Name: workshop_runs workshop_runs_workshop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.workshop_runs
    ADD CONSTRAINT workshop_runs_workshop_id_fkey FOREIGN KEY (workshop_id) REFERENCES public.workshops(id) ON DELETE CASCADE;


--
-- Name: workshop_sessions workshop_sessions_workshop_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.workshop_sessions
    ADD CONSTRAINT workshop_sessions_workshop_id_fkey FOREIGN KEY (workshop_id) REFERENCES public.workshops(id) ON DELETE CASCADE;


--
-- Name: workshops workshops_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.workshops
    ADD CONSTRAINT workshops_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE SET NULL;


--
-- Name: workshops workshops_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: most3mr
--

ALTER TABLE ONLY public.workshops
    ADD CONSTRAINT workshops_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.creators(id) ON DELETE CASCADE;


--
-- Name: analytics_clicks; Type: ROW SECURITY; Schema: public; Owner: most3mr
--

ALTER TABLE public.analytics_clicks ENABLE ROW LEVEL SECURITY;

--
-- Name: creators; Type: ROW SECURITY; Schema: public; Owner: most3mr
--

ALTER TABLE public.creators ENABLE ROW LEVEL SECURITY;

--
-- Name: enrollments; Type: ROW SECURITY; Schema: public; Owner: most3mr
--

ALTER TABLE public.enrollments ENABLE ROW LEVEL SECURITY;

--
-- Name: orders; Type: ROW SECURITY; Schema: public; Owner: most3mr
--

ALTER TABLE public.orders ENABLE ROW LEVEL SECURITY;

--
-- Name: shop_settings; Type: ROW SECURITY; Schema: public; Owner: most3mr
--

ALTER TABLE public.shop_settings ENABLE ROW LEVEL SECURITY;

--
-- Name: url_settings; Type: ROW SECURITY; Schema: public; Owner: most3mr
--

ALTER TABLE public.url_settings ENABLE ROW LEVEL SECURITY;

--
-- Name: workshops; Type: ROW SECURITY; Schema: public; Owner: most3mr
--

ALTER TABLE public.workshops ENABLE ROW LEVEL SECURITY;

--
-- PostgreSQL database dump complete
--


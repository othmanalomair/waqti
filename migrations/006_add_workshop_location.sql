-- Add location fields to workshops table
ALTER TABLE workshops ADD COLUMN location_name VARCHAR(255);
ALTER TABLE workshops ADD COLUMN location_link TEXT;
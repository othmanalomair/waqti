-- Migration to add 'private' to the workshop_type constraint
-- This allows private workshops to be created

-- Drop the existing constraint
ALTER TABLE workshops DROP CONSTRAINT workshops_workshop_type_check;

-- Add the new constraint with 'private' included
ALTER TABLE workshops ADD CONSTRAINT workshops_workshop_type_check 
    CHECK (workshop_type::text = ANY (ARRAY['single'::character varying, 'consecutive'::character varying, 'spread'::character varying, 'custom'::character varying, 'private'::character varying]::text[]));
-- Add store_layout field to shop_settings table
ALTER TABLE shop_settings
ADD COLUMN store_layout VARCHAR(10) DEFAULT 'grid' CHECK (store_layout IN ('grid', 'row'));

-- Update existing records to have default 'grid' layout
UPDATE shop_settings
SET
    store_layout = 'grid'
WHERE
    store_layout IS NULL;

-- Add comment for the new field
COMMENT ON COLUMN shop_settings.store_layout IS 'Layout preference for displaying workshops in store (grid or row)';

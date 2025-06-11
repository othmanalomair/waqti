-- Migration to add deleted field to workshops table for soft delete functionality
-- This prevents workshops with existing orders from being accidentally deleted

-- Add deleted field to workshops table
ALTER TABLE workshops ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT false;

-- Add index for deleted field to optimize queries
CREATE INDEX idx_workshops_deleted ON workshops (deleted);

-- Add composite index for creator queries filtering out deleted workshops
CREATE INDEX idx_workshops_creator_active_not_deleted ON workshops (creator_id, is_active, deleted) WHERE deleted = false;

-- Update any workshop views if needed (optional)
-- No existing views need modification as they will automatically filter by deleted=false when we update queries
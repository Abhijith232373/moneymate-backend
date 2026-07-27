-- Add logo_url column to stores table for Business Profile logo management
ALTER TABLE stores ADD COLUMN IF NOT EXISTS logo_url TEXT;

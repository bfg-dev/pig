-- +pig Name: Add cost to goods table
-- +pig Requiremets: Create goods
-- +pig Up
ALTER TABLE goods ADD COLUMN cost INTEGER;

-- +pig Down
ALTER TABLE goods DROP COLUMN cost;
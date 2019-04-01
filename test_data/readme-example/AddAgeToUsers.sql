-- +pig Name: Add age to users table
-- +pig Requirements: Create users
-- +pig Up
ALTER TABLE users ADD COLUMN age INTEGER;

-- +pig Down
ALTER TABLE users DROP COLUMN age;
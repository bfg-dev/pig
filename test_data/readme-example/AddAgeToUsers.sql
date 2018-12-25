-- +pig Name: Add age to users table
-- +pig Requiremets: Create users
-- +pig Up
ALTER TABLE users ADD COLUMN age INTEGER;

-- +pig Down
ALTER TABLE users DROP COLUMN age;
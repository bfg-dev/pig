-- +pig Name: Create purchases
-- +pig Requirements: Create users, Create goods
-- +pig Up
CREATE TABLE purchases (
    "id" serial,
    "created" timestamp DEFAULT current_timestamp,
    "userId" INTEGER,
    "goodId" INTEGER
);

-- +pig Down
DROP TABLE purchases;
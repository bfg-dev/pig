-- +pig Name: Create goods
-- +pig Up
CREATE TABLE goods (
    "id" serial,
    "created" timestamp DEFAULT current_timestamp,
    "name" varchar(40)
);

-- +pig Down
DROP TABLE goods;
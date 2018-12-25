-- +pig Name: Create users
-- +pig Up
CREATE TABLE users (
    "id" serial,
    "created" timestamp DEFAULT current_timestamp,
    "fname" varchar(40),
    "sname" varchar(40),
    "email" varchar(40),
    "password" varchar(40)
);

-- +pig Down
DROP TABLE users;
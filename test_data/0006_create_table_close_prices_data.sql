-- +pig Name: Create close_proces_data table
-- +pig Up
CREATE TABLE close_prices_data (
  "id"           SERIAL,
  "created"      TIMESTAMP DEFAULT current_timestamp,
  "last_updated" TIMESTAMP DEFAULT current_timestamp,
  "from"         VARCHAR(20) NOT NULL,
  "to"           VARCHAR(20) NOT NULL,
  "rate"         JSONB       NOT NULL
);
-- new line

-- +pig Down
DROP TABLE close_prices_data;

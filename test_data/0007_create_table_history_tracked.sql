-- +pig Name: Create history_tracked table
-- +pig Up
CREATE TABLE history_tracked (
  "id"         SERIAL,
  "created"    TIMESTAMP            DEFAULT current_timestamp,
  "from"       VARCHAR(20) NOT NULL,
  "to"         VARCHAR(20) NOT NULL,
  "name"       VARCHAR(20) NOT NULL,
  "native"     VARCHAR(20) NOT NULL,
  "is_deleted" BOOLEAN     NOT NULL DEFAULT FALSE,
  "data"       JSONB       NOT NULL
);

-- +pig Down
DROP TABLE history_tracked;

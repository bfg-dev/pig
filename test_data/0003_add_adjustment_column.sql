-- +pig Name: Add adjustment column
-- +pig Requiremets: init
-- +pig Up
ALTER TABLE "tracked_rates" add column "adjustment" numeric(5,4) NOT NULL default 1;

-- +pig Down
ALTER TABLE "tracked_rates" DROP column "adjustment";
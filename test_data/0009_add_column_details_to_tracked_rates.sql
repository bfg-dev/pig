-- +pig Name: Add details column
-- +pig Requiremets: init
-- +pig Up
ALTER TABLE "tracked_rates"
  ADD COLUMN "details" JSONB DEFAULT NULL;

-- +pig Down
ALTER TABLE "tracked_rates"
  DROP COLUMN "details";

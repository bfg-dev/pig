-- +pig Name: Add source column
-- +pig Requiremets: init
-- +pig Up
ALTER TABLE "tracked_rates" ADD PRIMARY KEY (id);
ALTER TABLE "tracked_rates" add column "source" varchar(50);

-- +pig Down
ALTER TABLE "tracked_rates" DROP CONSTRAINT "tracked_rates_pkey";
ALTER TABLE "tracked_rates" DROP column "source";
-- +pig Name: Add asset rates table
-- +pig Requiremets: init
-- +pig Up
alter table "tracked_rates" add column "last_val" numeric(12,4) not null default 0;
alter table "tracked_rates" add column "last_time" timestamp;

-- +pig Down
alter table "tracked_rates" drop column "last_val";
alter table "tracked_rates" drop column "last_time";

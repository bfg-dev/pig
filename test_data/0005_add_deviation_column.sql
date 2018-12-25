-- +pig Name: Add deviation column
-- +pig Requiremets: init
-- +pig Up
alter table "tracked_rates" add column "deviation" smallint not null default 0;

-- +pig Down
alter table "tracked_rates" drop column "deviation";

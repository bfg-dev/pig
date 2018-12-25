-- +pig Name: init
-- +pig Up
create table tracked_rates(
    "id" serial,
    "created" timestamp DEFAULT current_timestamp,
    "from" varchar(20) NOT NULL,
    "to" varchar(20) NOT NULL,
    "is_deleted" boolean default false
);

create table currencies_rates(
    "id" serial,
    "created" timestamp DEFAULT current_timestamp,
    "last_updated" timestamp DEFAULT current_timestamp,
    "from" varchar(20) NOT NULL,
    "to" varchar(20) NOT NULL,
    "rate" numeric(11,3) NOT NULL
);

create unique index "currencies_rates_txid_block_logindex" on currencies_rates ("from", "to", "last_updated");

insert into tracked_rates("from", "to") values ('btc', 'usd'), ('btc', 'eur'), ('ltc', 'usd'), ('ltc', 'eur'), ('eth', 'usd'), ('eth', 'eur');

-- +pig Down
drop table tracked_rates;
drop table currencies_rates;
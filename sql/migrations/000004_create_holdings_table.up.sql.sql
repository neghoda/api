BEGIN;

create table holdings (
    fund_ticker varchar    NOT NULL,
    name        varchar    NOT NULL,
	share       varchar    NOT NULL,
    weight      varchar    NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;
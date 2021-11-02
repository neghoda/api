BEGIN;

create table countries (
    fund_ticker varchar   NOT NULL,
    name        varchar   NOT NULL,
    weight      varchar   NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;

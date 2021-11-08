BEGIN;

CREATE TABLE sectors (
    fund_ticker VARCHAR   NOT NULL,
    name        VARCHAR   NOT NULL,
    weight      VARCHAR   NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;

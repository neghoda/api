BEGIN;

CREATE TABLE funds (
    id          UUID      NOT NULL,
    name        VARCHAR   NOT NULL,
    ticker      VARCHAR   NOT NULL,
    description TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;

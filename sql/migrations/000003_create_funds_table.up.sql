BEGIN;

create table funds (
    id          uuid      NOT NULL,
    name        varchar   NOT NULL,
	  ticker      varchar   NOT NULL,
    description text      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;

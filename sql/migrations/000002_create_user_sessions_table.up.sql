BEGIN;

CREATE TABLE user_sessions (
    id            UUID NOT NULL,
    token_id      UUID NOT NULL,
    user_id       UUID NOT NULL,
    refresh_token TEXT NOT NULL,
    created_at    TIMESTAMP NOT NULL,
    updated_at    TIMESTAMP NOT NULL,
    expired_at    TIMESTAMP NOT NULL,


    CONSTRAINT user_sessions_pk PRIMARY KEY (id),
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users ON DELETE restrict
);

COMMIT;

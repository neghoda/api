BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION trigger_set_timestamp() RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


create table user_sessions (
    id uuid not null,
    token_id  uuid not null,
    user_id uuid not null,
    refresh_token text not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    expired_at timestamp not null,


    constraint user_sessions_pk primary key (id),
    constraint user_id_fk foreign key (user_id) references users on delete restrict
);

COMMIT;

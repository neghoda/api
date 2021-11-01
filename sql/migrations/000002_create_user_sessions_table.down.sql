BEGIN;

DROP FUNCTION IF EXISTS trigger_set_timestamp;
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP TABLE IF EXISTS user_sessions;

COMMIT;

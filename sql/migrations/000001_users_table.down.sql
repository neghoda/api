BEGIN;

DROP TABLE IF EXISTS users;
DROP TRIGGER IF EXISTS set_users_timestamp ON users;
DROP FUNCTION IF EXISTS trigger_set_timestamp();

COMMIT;



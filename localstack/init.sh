#!/bin/bash
set -e
for database in api-db api-db-test;
 do
psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$database'" | grep -q 1 || \
psql -U postgres <<-EOSQL
CREATE DATABASE "$database" WITH owner=postgres;
EOSQL
done
psql -U postgres -tc "CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS plpgsql;"

#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER sample;
    CREATE DATABASE sample_state;
    CREATE DATABASE sample_metadata;
	GRANT ALL PRIVILEGES ON DATABASE sample_state TO sample;
    GRANT ALL PRIVILEGES ON DATABASE sample_metadata TO sample;
EOSQL
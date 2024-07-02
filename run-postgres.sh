#! /usr/bin/env bash
set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo 'Starting database...'

PASSSWORD=daprrulz

docker run \
  --name sample-postgres \
  -e POSTGRES_PASSWORD=$PASSSWORD \
  -p 5432:5432 \
  -v "$SCRIPT_DIR/db/init-db.sh:/docker-entrypoint-initdb.d/init-db.sh" \
  --rm \
  -d \
  postgres

echo "Database started successfully."
echo "Connect using: psql -h localhost -U sample -P $PASSSWORD -d sample_state"
echo "Connection string: host=localhost user=postgres password=$PASSSWORD port=5432 connect_timeout=10 database=sample_state"
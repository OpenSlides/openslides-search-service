#!/bin/bash

ENTER_PSQL=$1

# This is called within the docker-compose.dev models container to setup meta/models in the postgres container
until pg_isready -h "$DATABASE_HOST" -p "$DATABASE_PORT" -U "$DATABASE_USER"; do
  echo "Waiting for Postgres server '$DATABASE_HOST' to become available..."
  sleep 3
done

psql -1 -v ON_ERROR_STOP=1 -h postgres -p 5432 -U openslides -d openslides -f scheme_relational.sql
psql -1 -v ON_ERROR_STOP=1 -h postgres -p 5432 -U openslides -d openslides -f base_data.sql
psql -1 -v ON_ERROR_STOP=1 -h postgres -p 5432 -U openslides -d openslides -f test_data.sql
psql -1 -v ON_ERROR_STOP=1 -h postgres -p 5432 -U openslides -d openslides -f mock_data.sql

if [ -n "$ENTER_PSQL" ]; then psql -h postgres -U openslides -p 5432 -d openslides; fi

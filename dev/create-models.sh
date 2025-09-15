#!/bin/bash

# This is called within the docker-compose.dev models container to setup meta/models in the postgres container
echo "Port: $DATABASE_PORT , User $DATABASE_USER"
until pg_isready -h "$DATABASE_HOST" -p "$DATABASE_PORT" -U "$DATABASE_USER"; do
  echo "Waiting for Postgres server '$DATABASE_HOST' to become available..."
  sleep 3
done

bash scripts/apply_db_schema.sh
bash scripts/apply_data.sh base_data.sql
bash scripts/apply_data.sh test_data.sql

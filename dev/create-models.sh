#!/bin/bash

# This is called within the docker-compose.dev models container to setup meta/models in the postgres container
until pg_isready -h "$DATABASE_HOST" -p "$DATABASE_PORT" -U "$DATABASE_USER"; do
  echo "Waiting for Postgres server '$DATABASE_HOST' to become available..."
  sleep 3
done

psql -1 -v ON_ERROR_STOP=1 -h postgres -p 5432 -U openslides -d openslides -f mock_data.sql

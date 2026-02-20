#!/bin/sh

if [ ! $ANONYMOUS_ONLY -a $DATABASE_HOST -a $DATABASE_PORT ]; then
    while ! nc -z "$DATABASE_HOST" "$DATABASE_PORT"; do
        echo "waiting for $DATABASE_HOST:$DATABASE_PORT"
        sleep 1
    done

    echo "$DATABASE_HOST:$DATABASE_PORT is available"
fi

exec "$@"

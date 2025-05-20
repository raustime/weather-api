#!/bin/sh

: "${PGHOST:=db}"
: "${PGPORT:=5432}"
: "${PGUSER:=postgres}"

echo "Waiting for postgres at $PGHOST:$PGPORT as user $PGUSER..."

until pg_isready -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" > /dev/null 2>&1; do
  echo "Waiting for postgres..."
  sleep 2
done

echo "Postgres is up - executing command"
exec "$@"
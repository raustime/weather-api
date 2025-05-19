#!/bin/sh

echo "Waiting for postgres..."

until pg_isready -h db -p 5432 -U postgres > /dev/null 2>&1; do
  echo "Waiting for postgres..."
  sleep 2
done

echo "Postgres is up - executing command"
exec "$@"
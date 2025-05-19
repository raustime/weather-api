#!/bin/bash
set -e

until nc -z db 5432; do
  echo "Waiting for postgres..."
  sleep 2
done

echo "Postgres is up - running tests"
exec "$@"
#!/bin/sh
set -e

echo "Esperando Postgres (${DB_HOST}:${DB_PORT}) iniciar..."

until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "$DB_HOST:$DB_PORT - no response"
  sleep 2
done

echo "Postgres está pronto! Iniciando aplicação..."
exec "$@"
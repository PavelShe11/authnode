#!/bin/bash
set -e

psql_command() {
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -c "$1"
}

echo "Configuring pg_cron for database $POSTGRES_DB..."
psql_command "ALTER SYSTEM SET cron.database_name TO '$POSTGRES_DB';"

pg_ctl restart -D "$PGDATA" -m fast -w

echo "Creating pg_cron extension..."
psql_command "CREATE EXTENSION IF NOT EXISTS pg_cron;"

if [ "$POSTGRES_USER" != "postgres" ]; then
  echo "Granting pg_cron usage to $POSTGRES_USER..."
  psql_command "GRANT USAGE ON SCHEMA cron TO $POSTGRES_USER;"
fi

echo "pg_cron extension has been installed and configured."
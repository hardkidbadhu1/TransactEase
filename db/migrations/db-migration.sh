#!/bin/sh

echo "Running migration for transact-api database...."

DB_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable"

echo "DB URL: $DB_URL"
migrate -path=/migrate/db/migrations -database=$DB_URL -verbose up

echo "Migration completed successfully."

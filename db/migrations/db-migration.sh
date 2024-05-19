#!/bin/sh

echo "Running migration for TransactEase database...."

DB_URL="mysql://$MYSQL_USER:$MYSQL_PASSWORD@$DB_HOST:$DB_PORT/$MYSQL_DATABASE"

echo "DB URL: $DB_URL"
migrate -path=. -database=$DB_URL -verbose up

echo "Migration completed successfully."

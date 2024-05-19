#!/bin/bash

set -e

host="$1"
shift
cmd="$@"

timeout=300  # 5 minutes
interval=1
elapsed=0

until pg_isready -h "$POSTGRES_HOST" -U "$POSTGRES_USER"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep $interval
  elapsed=$((elapsed + interval))
  if [ $elapsed -ge $timeout ]; then
    >&2 echo "Timeout reached: Postgres is still unavailable after $timeout seconds"
    exit 1
  fi
done

>&2 echo "Postgres is up - executing command"
exec $cmd

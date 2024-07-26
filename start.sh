#!/bin/sh

set -e

echo "DB_SOURCE is set to: $DB_SOURCE"
# echo "Contents of /app directory:"
# ls -la /app

# echo "Contents of /app/db directory:"
# ls -la /app/db

# echo "Contents of /app/db/migration directory:"
# ls -la /app/db/migration

echo "run db migration"
/app/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "start app"
exec "$@"
#!/bin/sh

echo "Waiting for ClickHouse to become available..."
for i in {1..30}; do
  clickhouse-client --query "SELECT 1" && break
  echo "ClickHouse is not ready yet... attempt $i"
  sleep 1
done

echo "Running ClickHouse migration from init.sql..."
clickhouse-client --queries-file=/docker-entrypoint-initdb.d/init.sql || echo "⚠️ Migration failed"

echo "Starting ClickHouse server..."
exec /entrypoint.sh

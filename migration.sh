#!/bin/bash
source .env

export DSN="host=pg port=5432 dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD sslmode=$POSTGRES_SSL_MODE"

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${DSN}" up -v
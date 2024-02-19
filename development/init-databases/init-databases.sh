#!/bin/bash

set -e
set -u

function create_user_and_database() {
  local database=$1
  echo "  Creating user and database '$database'"
  psql -U $POSTGRES_USER --c "CREATE DATABASE IF NOT EXISTS $database"
}

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
  echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
  for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
    create_user_and_database $db
  done
  echo "Multiple databases created"
fi

#!/bin/bash

set -e
set -u

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
  echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
  for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
    createdb -U $POSTGRES_USER $db
  done
  echo "Multiple databases created"
fi

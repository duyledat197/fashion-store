#!/bin/bash

svcs=("user-management" "product-management" "coupon-management")

install_migrate() {
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
}

migrate_all() {
  host=$1
  port=$2
  username=$3
  password=$4
  for svc in $svc; do
    migrate -source /migration -database postgres://${username}:${password}@${host}:${port}/${svc} up
  done
}

install_migrate
migrate_all

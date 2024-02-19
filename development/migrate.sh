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
  for svc in ${svcs[@]}; do
    echo "migrating $svc service database"
    echo "attach database with host=${host}"
    echo "attach database with port=${port}"
    echo "attach database with username=${username}"
    echo "attach database with password=${password}"
    migrate -path /migration/${svc} -database postgres://${username}:${password}@${host}:${port}/${svc}?sslmode=disable up
  done
}

install_migrate
migrate_all $1 $2 $3 $4

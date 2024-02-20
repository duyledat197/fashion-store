PROJECT_NAME := test
PKG := github.com/$(PROJECT_NAME)
MOD := $(GOPATH)/pkg/mod
COMPOSE_FILE := ./development/docker-compose.yml

gen-proto:
	docker compose -f ${COMPOSE_FILE} up generate_pb_go

gen-mock:
	docker compose -f ${COMPOSE_FILE} up generate_mock

migrate-all:
	docker compose -f ${COMPOSE_FILE} up migrate

start-postgres:
	docker compose -f ${COMPOSE_FILE} up postgres -d

start-adminer:
	docker compose -f ${COMPOSE_FILE} up adminer -d

start-user:
	SERVICE=user-management ENV=dev go run main.go userManagement

start-product:
	SERVICE=product-management ENV=dev go run main.go productManagement

run:
	./developments/start.sh
test:
	go test ./...

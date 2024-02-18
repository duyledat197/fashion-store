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

run:
	./developments/start.sh
test:
	go test ./...

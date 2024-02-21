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

build:
	CGO_ENABLED=0 && go build -o app-exe
	
docker-build:
	docker build -t trintech/review -f ./development/Dockerfile .

start-postgres:
	docker compose -f ${COMPOSE_FILE} up postgres -d

start-adminer:
	docker compose -f ${COMPOSE_FILE} up adminer -d

start-user-dev:
	SERVICE=user-management ENV=dev go run main.go userManagement

start-user-stg:
	SERVICE=user-management ENV=stg go run main.go userManagement

start-product-dev:
	SERVICE=product-management ENV=dev go run main.go productManagement

start-gateway-dev:
	SERVICE=gateway ENV=dev go run main.go gateway

start-coupon-dev:
	SERVICE=coupon-management ENV=dev go run main.go couponManagement

run:
	./developments/start.sh
test:
	go test ./...

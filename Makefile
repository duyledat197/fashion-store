PROJECT_NAME := test
PKG := github.com/$(PROJECT_NAME)
MOD := $(GOPATH)/pkg/mod
COMPOSE_FILE := ./development/docker-compose.yml

# build:
# 	@go build -i -v $(PKG)/cmd/server
run:
	./developments/start.sh
test:
	go test ./...
gen-proto:
	docker compose -f ${COMPOSE_FILE} up generate_pb_go
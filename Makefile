ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_USERNAME ?= hpcadmin
POSTGRES_PASSWORD ?= superfancytestpasswordthatnobodyknows&
POSTGRES_DATABASE ?= hpcadmin_test

all: build

migrate:
	@migrate -path database/migration/ -database "postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable" -verbose up

server:
	@go build -o bin/hpcadmin-server cmd/hpcadmin-server/main.go

cli:
	@go build -o bin/hpcadmin-cli cmd/hpcadmin-cli/main.go

build: server cli

run-server: server
	@./bin/hpcadmin-server

run-cli: cli
	@./bin/hpcadmin-cli

compose-up:
	@docker compose up --build

compose-down:
	@docker compose down -v

docs: build
	@./bin/hpcadmin-server -docs=markdown

test: migrate
	@go test ./...

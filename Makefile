ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_USERNAME ?= postgres
POSTGRES_PASSWORD ?= postgres
POSTGRES_DATABASE ?= hpcadmin_test

migrate:
	@migrate -path database/migration/ -database "postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable" -verbose up

build:
	@go build -o bin/hpcadmin-server cmd/hpcadmin-server/main.go

run: build
	@./bin/hpcadmin-server

docs: build
	@./bin/hpcadmin-server -docs=markdown

test: migrate
	@go test ./...

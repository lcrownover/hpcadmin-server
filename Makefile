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

install:
	cp ./bin/hpcadmin-server /usr/local/bin/hpcadmin-server
	mkdir -p /etc/hpcadmin-server
	cp ./extras/config.yaml.template /etc/hpcadmin-server/config.yaml
	
clean:
	rm -rf ./bin
	rm -rf /etc/hpcadmin-server
	rm -f /usr/local/bin/hpcadmin-server

migrate:
	migrate -path database/migration/ -database "postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable" -verbose up

server:
	go build -o bin/hpcadmin-server cmd/hpcadmin-server/main.go

tidy:
	go mod tidy

build: server

run-server: server
	./bin/hpcadmin-server

compose-up:
	docker compose up --build

compose-down:
	docker compose down -v

docs: build
	./bin/hpcadmin-server -docs=markdown

test: migrate
	go test ./...

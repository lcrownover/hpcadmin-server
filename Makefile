ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

HPCADMIN_TEST_DATABASE_HOST ?= localhost
HPCADMIN_TEST_DATABASE_PORT ?= 5432
HPCADMIN_TEST_DATABASE_USERNAME ?= hpcadmin
HPCADMIN_TEST_DATABASE_PASSWORD ?= superfancytestpasswordthatnobodyknows&
HPCADMIN_TEST_DATABASE_NAME ?= hpcadmin_test

export HPCADMIN_TEST_DATABASE_HOST
export HPCADMIN_TEST_DATABASE_PORT
export HPCADMIN_TEST_DATABASE_USERNAME
export HPCADMIN_TEST_DATABASE_PASSWORD
export HPCADMIN_TEST_DATABASE_NAME

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

migrate_quiet:
	migrate -path database/migration/ -database "postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable" up -quiet

build:
	@go build -o bin/hpcadmin-server cmd/hpcadmin-server/main.go

tidy:
	go mod tidy

docs: build
	./bin/hpcadmin-server -docs=markdown

testdb_setup:
	bash ./test/scripts/testDatabaseSetup.sh
	bash ./test/scripts/testBootstrap.sh

testdb_teardown:
	bash ./test/scripts/testDatabaseTeardown.sh

test: build
	@-bash ./test/scripts/testDatabaseSetup.sh
	@-bash ./test/scripts/testBootstrap.sh
	@-bash ./test/scripts/testStartServer.sh
	@-go test ./...
	@-bash ./test/scripts/testStopServer.sh
	@-bash ./test/scripts/testDatabaseTeardown.sh


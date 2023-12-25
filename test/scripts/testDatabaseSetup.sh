#!/bin/bash

POSTGRES_USERNAME="hpcadmin"
POSTGRES_PASSWORD="superfancytestpasswordthatnobodyknows&"
POSTGRES_HOST="localhost"
POSTGRES_PORT="5432"
POSTGRES_DATABASE="hpcadmin_test"

# Spin up postgres container
docker run \
	--name hpcadmin_test \
	-e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
	-e POSTGRES_USER="$POSTGRES_USERNAME" \
	-e POSTGRES_DB="$POSTGRES_DATABASE" \
	-p $POSTGRES_PORT:$POSTGRES_PORT \
	-d \
	postgres:latest

sleep 2

migrate -path database/migration/ -database "postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable" -verbose up

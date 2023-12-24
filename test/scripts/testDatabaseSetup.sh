#!/bin/bash

# Spin up postgres container
docker run \
	--name hpcadmin_test \
	-e POSTGRES_PASSWORD="superfancytestpasswordthatnobodyknows&" \
	-e POSTGRES_USER="hpcadmin" \
	-e POSTGRES_DB="hpcadmin_test" \
	-p 5432:5432 \
	-d \
	postgres:latest

#!/bin/bash

migrate -path database/migration/ -database "postgresql://postgres:postgres@localhost/hpcadmin?sslmode=disable" -verbose up

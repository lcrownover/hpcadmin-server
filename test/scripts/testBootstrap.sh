#!/bin/bash

POSTGRES_USERNAME="hpcadmin"
POSTGRES_PASSWORD="superfancytestpasswordthatnobodyknows&"
POSTGRES_HOST="localhost"
POSTGRES_PORT="5432"
POSTGRES_DATABASE="hpcadmin_test"

TEST_USERNAME="timmyt"
TEST_EMAIL="timmyt@example.org"
TEST_FIRSTNAME="Timmy"
TEST_LASTNAME="Test"
TEST_APIKEY="testkey1"
TEST_ROLE="admin"

export PGPASSWORD=$POSTGRES_PASSWORD

# Create the test user
psql \
	--host=$POSTGRES_HOST \
	--username=$POSTGRES_USERNAME \
	-p $POSTGRES_PORT \
	-d $POSTGRES_DATABASE \
	-c "INSERT INTO users (username, email, firstname, lastname) VALUES ('$TEST_USERNAME', '$TEST_EMAIL', '$TEST_FIRSTNAME', '$TEST_LASTNAME');" >/dev/null

# Get the ID for the test user
USERID=$(psql -t \
	--host=$POSTGRES_HOST \
	--username=$POSTGRES_USERNAME \
	-p $POSTGRES_PORT \
	-d $POSTGRES_DATABASE \
	-c "SELECT id FROM users WHERE username = 'timmyt';" |
	tr -d '[:space:]')

# Create an API key for the test user
psql \
	--host=$POSTGRES_HOST \
	--username=$POSTGRES_USERNAME \
	-p $POSTGRES_PORT \
	-d $POSTGRES_DATABASE \
	-c "INSERT INTO api_keys (key, role, user_id) VALUES ('$TEST_APIKEY', '$TEST_ROLE', $USERID);" >/dev/null

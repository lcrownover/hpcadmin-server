#!/bin/bash

# Stop and remove postgres container
docker stop hpcadmin_test >/dev/null
docker rm hpcadmin_test >/dev/null

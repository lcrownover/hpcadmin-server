#!/bin/bash

./bin/hpcadmin-server -config ./test/data/testconfig.yaml >/dev/null &
echo $! >/tmp/hpcadmin-server.pid

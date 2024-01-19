#!/bin/bash

xargs kill < /tmp/hpcadmin-server.pid 2> /dev/null
rm -f /tmp/hpcadmin-server.pid

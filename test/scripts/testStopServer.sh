#!/bin/bash

xargs kill < /tmp/hpcadmin-server.pid
rm -f /tmp/hpcadmin-server.pid

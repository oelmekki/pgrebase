#!/usr/bin/env bash

cd "$(dirname "$0")"
. ./common.sh

nohup postgres -D ./pg -k . -h 127.0.0.1 -p $PORT -F > ./pg/server_logs &

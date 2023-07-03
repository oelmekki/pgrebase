#!/usr/bin/env bash

cd "$(dirname "$0")"
. ./common.sh

psql -h 127.0.0.1 -p $PORT -U postgres -d pgrebase

#!/usr/bin/env bash

cd "$(dirname "$0")"
. ./common.sh

kill $(cat pg/postmaster.pid | head -n 1)

#!/usr/bin/env bash

cd "$(dirname "$0")"
. ./common.sh

./stop_db.sh

if [[ -d ./pg ]]; then
  rm -rf pg
fi

initdb -D ./pg -U postgres
./start_db.sh
sleep 3
createdb -h 127.0.0.1 -p $PORT -U postgres pgrebase
psql -h 127.0.0.1 -p $PORT -U postgres -d pgrebase -c "CREATE TABLE users(id SERIAL, name varchar(255), active boolean NOT NULL DEFAULT false, bio text)"

echo "Done."

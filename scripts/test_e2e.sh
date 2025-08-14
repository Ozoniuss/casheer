#!/bin/bash

# Note that e2e tests assume a DB with name casheer.e2e.db

export BUILD_DOCKER=

function cleanup {
  echo "Removing db..."
  rm $DBNAME
}

trap cleanup EXIT

if [[ "$1" == "--build" ]]; then
    BUILD_DOCKER=true
fi



DBNAME=casheer.e2e.db

touch "$DBNAME"
chmod 666 "$DBNAME"

# Run migrations before using the end to end test DB in the container.
sqlite3 "$DBNAME" < scripts/sqlite/002_update_entry_value.down.sql && 
sqlite3 "$DBNAME" < scripts/sqlite/001_tables.down.sql && 
sqlite3 "$DBNAME" < scripts/sqlite/001_tables.up.sql &&
sqlite3 "$DBNAME" < scripts/sqlite/002_update_entry_value.up.sql &&

ls -l "$DBNAME"


echo "Using database $DBNAME."



if [[ ! -z "$BUILD_DOCKER" ]]; then
    echo "building docker image..."
    docker compose -f docker-compose.e2e.yml up -d --build 
else
    docker compose -f docker-compose.e2e.yml up -d
fi

sleep 1

go test -v ./e2e/... ;

docker compose -f docker-compose.e2e.yml down

echo 'done';

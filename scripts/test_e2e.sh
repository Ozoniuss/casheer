#!/bin/bash

echo "ce ma? $(pwd)"

# Note that DBNAME is used in other places as well. This exact name is needed.
export DBNAME='casheer.e2e.db'
export BUILD_DOCKER=

function cleanup {
  echo "Removing db..."
  rm $DBNAME
}

trap cleanup EXIT

# I think a custom db file name is completely unnecessary. But, this code 
# makes it easy for me to keep the file for debugging if something fails.
if [[ -n $1 ]]; then
    if [[ "$1" == "--build" ]]; then
        BUILD_DOCKER=true
    else
        DBNAME="$1"
    fi
fi

# create a database
if [[ ! -f $DBNAME ]]; then
    touch "$DBNAME"
fi


echo "Using database $DBNAME."

if [[ ! -z "$BUILD_DOCKER" ]]; then
    echo "building docker image..."
    docker compose -f docker-compose.e2e.yml up -d --build 
else
    docker compose -f docker-compose.e2e.yml up -d
fi

# Run migrations before using the end to end test DB in the container.
sqlite3 "$DBNAME" < scripts/sqlite/002_update_entry_value.down.sql && 
sqlite3 "$DBNAME" < scripts/sqlite/001_tables.down.sql && 
sqlite3 "$DBNAME" < scripts/sqlite/001_tables.up.sql &&
sqlite3 "$DBNAME" < scripts/sqlite/002_update_entry_value.up.sql &&

sleep 1

go test -v ./e2e/... ;

docker compose -f docker-compose.e2e.yml down

echo 'done';

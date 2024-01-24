export DBNAME='casheer.e2e.db'
export BUILD_DOCKER=

# I think a custom db file name is completely unnecessary. But, this code 
# makes it easy for me to keep the file for debugging if something fails.
if [[ -n $1 ]]; then
    if [[ "$1" == "--build" ]]; then
        BUILD_DOCKER=true
    else
        DBNAME="$1"
    fi
fi

if [[ ! -f $DBNAME ]]; then
    touch "$DBNAME"
fi

echo "Using database $DBNAME."

# should use DBNAME
if [[ ! -z "$BUILD_DOCKER" ]]; then
    echo "building docker image..."
    docker compose -f docker-compose.e2e.yml up -d --build 
else
    docker compose -f docker-compose.e2e.yml up -d
fi

# Perform migrations. The image already comes with a db, but if it were to be
# cached, we would have to re-write the db file every time. So instead, just 
# create a file here and overwrite it between runs.
sqlite3 "$DBNAME" < scripts/sqlite/002_update_entry_value.down.sql && 
sqlite3 "$DBNAME" < scripts/sqlite/001_tables.down.sql && 
sqlite3 "$DBNAME" < scripts/sqlite/001_tables.up.sql &&
sqlite3 "$DBNAME" < scripts/sqlite/002_update_entry_value.up.sql &&

sleep 1

go test -v ./e2e/... ;

docker compose -f docker-compose.e2e.yml down

rm $DBNAME

echo 'done';

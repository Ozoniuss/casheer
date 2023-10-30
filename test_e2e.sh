DBNAME='casheer.e2e.db'

if [[ -n $1 ]]; then
    DBNAME="$1"
fi 

echo "Using database $DBNAME"

sqlite3 "$DBNAME" < scripts/sqlite/002_update_entry_value.down.sql && 
sqlite3 "$DBNAME" < scripts/sqlite/001_tables.down.sql && 
sqlite3 "$DBNAME" < scripts/sqlite/001_tables.up.sql &&
sqlite3 "$DBNAME" < scripts/sqlite/002_update_entry_value.up.sql &&
go test -v ./e2e/entries_test.go ; 
echo 'done';

#!/bin/bash
# Export SQLite and MySQL databases for comparison
# Usage: ./scripts/export-databases.sh

set -e

# Load environment variables from docker/.env.dev
if [ -f "./docker/.env.dev" ]; then
    export $(grep -v '^#' ./docker/.env.dev | xargs)
    echo "Loaded environment variables from ./docker/.env.dev"
    echo $MARIADB_USER
fi

# Configuration
SQLITE_DB="./backend/testdata/prepared_test_data.db"
OUTPUT_DIR="./backend/testdata/exports"
TIMESTAMP=""        #"_"$(date +%Y%m%d_%H%M%S)

# Docker container name
MARIADB_CONTAINER="bamort-mariadb-dev"

# MySQL credentials (from docker/.env.dev)
#MARIADB_USER="${MARIADB_USER:-bamort}"
#MARIADB_PASSWORD="${MARIADB_PASSWORD:-secure_user_password}"
#MARIADB_DATABASE="${MARIADB_DATABASE:-bamort}"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "📦 Exporting databases to $OUTPUT_DIR"

# Export SQLite to SQL dump
if [ -f "$SQLITE_DB" ]; then
    echo "📤 Exporting SQLite: $SQLITE_DB"
    sqlite3 "$SQLITE_DB" .dump > "$OUTPUT_DIR/sqlite_dump$TIMESTAMP.sql"
    
    # Export tables as CSV
    echo "📊 Exporting SQLite tables to CSV..."
    sqlite3 "$SQLITE_DB" "SELECT name FROM sqlite_master WHERE type='table';" | while read table; do
        sqlite3 -header -csv "$SQLITE_DB" "SELECT * FROM $table;" > "$OUTPUT_DIR/sqlite_${table}$TIMESTAMP.csv"
    done
    
    echo "✅ SQLite export complete"
else
    echo "❌ SQLite database not found: $SQLITE_DB"
    exit 1
fi

# Export MySQL from Docker container
if docker ps | grep -q "$MARIADB_CONTAINER"; then
    echo "📤 Exporting MySQL from container: $MARIADB_CONTAINER"
    
    # SQL dump
    docker exec "$MARIADB_CONTAINER" mariadb-dump -u"$MARIADB_USER" -p"$MARIADB_PASSWORD" "$MARIADB_DATABASE" \
        > "$OUTPUT_DIR/mysql_dump$TIMESTAMP.sql"
    
    # Export tables as CSV
    echo "📊 Exporting MySQL tables to CSV..."
    docker exec "$MARIADB_CONTAINER" mariadb -u"$MARIADB_USER" -p"$MARIADB_PASSWORD" -D"$MARIADB_DATABASE" \
        -e "SHOW TABLES;" -sN | while read table; do
        docker exec "$MARIADB_CONTAINER" mariadb -u"$MARIADB_USER" -p"$MARIADB_PASSWORD" -D"$MARIADB_DATABASE" \
            --batch -e "SELECT * FROM \`$table\`;" | sed 's/\t/,/g' > "$OUTPUT_DIR/mysql_${table}$TIMESTAMP.csv"
    done
    
    echo "✅ MySQL export complete"
else
    echo "⚠️  MySQL container not running: $MARIADB_CONTAINER"
    echo "   Start with: docker-compose -f docker/docker-compose.dev.yml up -d"
    exit 1
fi

echo ""
echo "✅ Export completed!"
echo "📁 Output directory: $OUTPUT_DIR"
echo ""
echo "Files created:"
ls -lh "$OUTPUT_DIR"/*"$TIMESTAMP"*
echo ""
echo "💡 Compare dumps with: diff $OUTPUT_DIR/sqlite_dump$TIMESTAMP.sql $OUTPUT_DIR/mysql_dump$TIMESTAMP.sql"

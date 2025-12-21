# Data Transfer: SQLite to MariaDB

This document describes how to transfer data from the SQLite test database to the MariaDB container using the maintenance endpoint.

## Overview

The data transfer functionality allows you to migrate all data from your SQLite test database (`./testdata/prepared_test_data.db`) to the MariaDB container. This is useful when:

- Setting up a new MariaDB database with existing test data
- Migrating from SQLite to MariaDB for production use
- Synchronizing test data across different database systems

## Prerequisites

1. **MariaDB Container Running**: Ensure the MariaDB container is up and healthy:
   ```bash
   docker-compose -f docker-compose.dev.yml up mariadb
   ```

2. **Backend Server Running**: The backend must be running to access the maintenance endpoint:
   ```bash
   docker-compose -f docker-compose.dev.yml up backend-dev
   ```

3. **SQLite Source Database**: The file `./testdata/prepared_test_data.db` must exist.

## Usage Methods

### Method 1: Using the Shell Script (Recommended)

```bash
# Basic transfer (preserves existing data)
./transfer_sqlite_to_mariadb.sh

# Transfer with data clearing (removes all existing data first)
./transfer_sqlite_to_mariadb.sh clear
```

### Method 2: Direct API Call

```bash
# Basic transfer
curl -X POST http://localhost:8180/api/maintenance/transfer-sqlite-to-mariadb

# Transfer with data clearing
curl -X POST "http://localhost:8180/api/maintenance/transfer-sqlite-to-mariadb?clear=true"
```

### Method 3: Using a HTTP Client (Postman, Insomnia, etc.)

- **Method**: POST
- **URL**: `http://localhost:8180/api/maintenance/transfer-sqlite-to-mariadb`
- **Query Parameter** (optional): `clear=true`

## Transfer Process

The endpoint performs the following steps:

1. **Validation**: Checks if the SQLite source file exists
2. **Database Connections**: Establishes connections to both SQLite and MariaDB
3. **Schema Migration**: Ensures all table structures exist in MariaDB
4. **Data Clearing** (optional): Removes existing data if `clear=true` parameter is provided
5. **Data Copy**: Transfers data in batches for all tables in the correct order
6. **Statistics**: Returns transfer statistics and confirmation

## Supported Tables

The transfer includes all application tables:

### Core Tables
- Users (`user.User`)
- Character data (`models.Char`)
- Skills, Spells, Equipment
- Character properties and relationships

### Learning System Tables
- Sources, Character Classes, Skill Categories
- Learning costs and improvement data
- Spell schools and level costs

### Game Master Data
- Skills, Weapons, Spells, Equipment
- Containers, Transportation, Beliefs
- Character equipment and weapon assignments

## Response Format

### Success Response (HTTP 200)
```json
{
  "message": "Data transfer from SQLite to MariaDB completed successfully",
  "source_file": "/app/testdata/prepared_test_data.db",
  "target": "mariadb:3306/bamort",
  "statistics": {
    "users": 5,
    "chars": 12,
    "skills": 45,
    "spells": 78,
    "equipment": 234
  }
}
```

### Error Response (HTTP 4xx/5xx)
```json
{
  "error": "SQLite source file not found: /app/testdata/prepared_test_data.db"
}
```

## Configuration Switch

After successful transfer, you may want to switch your application to use MariaDB:

### Update `.env` file:
```bash
# From SQLite configuration
DATABASE_TYPE=sqlite
DATABASE_URL=./testdata/prepared_test_data.db

# To MariaDB configuration
DATABASE_TYPE=mysql
DATABASE_URL=bamort:bG4)efozrc@tcp(mariadb:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local
```

### Restart the backend:
```bash
docker-compose -f docker-compose.dev.yml restart backend-dev
```

## Important Notes

### ⚠️ Data Safety
- **Backup First**: Always backup your MariaDB data before running with `clear=true`
- **Test Environment**: Consider testing the transfer in a separate environment first
- **Incremental Transfer**: Without `clear=true`, the transfer uses upsert operations (INSERT ... ON DUPLICATE KEY UPDATE)

### 🔧 Performance
- Data is transferred in batches of 100 records for optimal performance
- Large datasets may take several minutes to transfer
- Monitor the backend logs for detailed progress information

### 🐛 Troubleshooting

#### "SQLite source file not found"
- Ensure `./testdata/prepared_test_data.db` exists
- Run from the backend directory
- Check file permissions

#### "Failed to connect to MariaDB target"
- Verify MariaDB container is running and healthy
- Check network connectivity between containers
- Validate MariaDB credentials

#### "Failed to migrate structures"
- Ensure MariaDB user has sufficient privileges
- Check for database schema conflicts
- Review backend logs for specific migration errors

#### Transfer Fails Partway Through
- Tables with foreign key constraints may fail if referenced data is missing
- Use `clear=true` to ensure clean transfer
- Check for data type incompatibilities between SQLite and MariaDB

## Logging

The transfer process generates detailed logs. Monitor them with:

```bash
# View backend logs
docker-compose -f docker-compose.dev.yml logs -f backend-dev

# View MariaDB logs
docker-compose -f docker-compose.dev.yml logs -f mariadb
```

## Related Endpoints

- `GET /api/maintenance/setupcheck` - Verify system status
- `GET /api/maintenance/mktestdata` - Create test data from live database
- `GET /api/maintenance/reconndb` - Reconnect to database
- `GET /api/maintenance/reloadenv` - Reload environment configuration

# BaMoRT Deployment CLI Tool

Command-line interface for BaMoRT database deployment and maintenance operations.

## Building

```bash
cd backend
go build -o deploy cmd/deploy/main.go
```

## Commands

### Fresh Installation

Create a new database installation from scratch:

```bash
./deploy install
```

This will:
1. Create all database tables using GORM AutoMigrate
2. Initialize version tracking
3. Import master data from files
4. Optionally create an admin user

**Interactive prompts:**
- Master data directory (default: `./masterdata`)
- Create admin user? (y/N)
- Admin username and password (if yes)

### Apply Migrations

Apply pending database migrations:

```bash
./deploy migrate
```

**Dry-run mode** (preview without applying):
```bash
./deploy migrate --dry-run
```

### Check Status

Display current database version and compatibility status:

```bash
./deploy status
```

Shows:
- Current database version
- Backend version
- Required DB version
- Compatibility status
- Pending migrations (if any)

### Create Backup

Create a JSON backup of the database:

```bash
./deploy backup
```

Automatically:
- Creates timestamped backup file
- Shows backup metadata (version, size, table count)
- Cleans up backups older than 30 days

### Sync Master Data

Import/update master data from files:

```bash
./deploy sync-masterdata [directory]
```

Default directory: `./masterdata`

**Dry-run mode:**
```bash
./deploy sync-masterdata --dry-run
```

### Rollback Migration

Rollback the last applied migration:

```bash
./deploy rollback
```

⚠️ **Warning:** Only use if migration caused issues. Requires rollback SQL to be defined.

### Version Information

Display version information:

```bash
./deploy version
```

## Usage Examples

### New Production Setup

```bash
# 1. Build the tool
go build -o deploy cmd/deploy/main.go

# 2. Fresh installation
./deploy install
# Enter master data path: /opt/bamort/masterdata
# Create admin user? y
# Admin username: admin
# Admin password: ********

# 3. Verify status
./deploy status

# 4. Create initial backup
./deploy backup
```

### Upgrading Existing Database

```bash
# 1. Check current status
./deploy status

# 2. Create backup before migration
./deploy backup

# 3. Preview migrations (dry-run)
./deploy migrate --dry-run

# 4. Apply migrations
./deploy migrate

# 5. Verify status
./deploy status
```

### Updating Master Data

```bash
# 1. Preview sync
./deploy sync-masterdata --dry-run

# 2. Apply sync
./deploy sync-masterdata

# Or specify custom directory
./deploy sync-masterdata /path/to/masterdata
```

## Environment Variables

The tool uses the same configuration as the backend application:

- `DATABASE_TYPE`: `mysql` or `sqlite` (default: `mysql`)
- `DATABASE_HOST`: Database host (default: `localhost`)
- `DATABASE_PORT`: Database port (default: `3306`)
- `DATABASE_NAME`: Database name (default: `bamort`)
- `DATABASE_USER`: Database user
- `DATABASE_PASSWORD`: Database password

## Exit Codes

- `0`: Success
- `1`: Error occurred (check stderr output)

## Color Output

The tool uses ANSI colors for better readability:

- 🟢 Green: Success messages
- 🔴 Red: Error messages
- 🟡 Yellow: Warnings and prompts
- 🔵 Cyan: Informational headers

## Safety Features

### Confirmation Prompts

Destructive operations require confirmation:
- Fresh installation (overwrites database)
- Applying migrations
- Rolling back migrations

### Dry-Run Mode

Test operations without making changes:
- `--dry-run` or `-n` flag
- Available for: `migrate`, `sync-masterdata`

### Automatic Backups

Recommended workflow:
1. `deploy backup` before any migration
2. `deploy migrate --dry-run` to preview
3. `deploy migrate` to apply

## Troubleshooting

### "Database not initialized"

If `deploy status` shows this error:

```bash
# Fresh installation required
./deploy install
```

### "Backend Too Old"

Database version is newer than backend:

```bash
# Upgrade backend application to latest version
# DO NOT rollback database
```

### "Migration Required"

Database is behind backend version:

```bash
# Create backup first
./deploy backup

# Apply migrations
./deploy migrate
```

### Connection Errors

Check environment variables and database accessibility:

```bash
# Test database connection
mysql -h $DATABASE_HOST -u $DATABASE_USER -p $DATABASE_NAME

# Or for SQLite
sqlite3 $DATABASE_PATH
```

## Integration with CI/CD

### Example GitLab CI

```yaml
deploy:production:
  stage: deploy
  script:
    - cd backend
    - go build -o deploy cmd/deploy/main.go
    - ./deploy status
    - ./deploy backup
    - ./deploy migrate
  only:
    - main
  environment:
    name: production
```

### Example Docker

```dockerfile
FROM golang:1.21

WORKDIR /app
COPY backend/ .

RUN go build -o deploy cmd/deploy/main.go

ENTRYPOINT ["./deploy"]
CMD ["status"]
```

Usage:
```bash
docker run bamort-deploy status
docker run -it bamort-deploy install
docker run bamort-deploy migrate
```

## Development

Run without building:

```bash
go run cmd/deploy/main.go status
go run cmd/deploy/main.go migrate --dry-run
```

## Security Notes

1. **Password Input**: Uses terminal mode for secure password entry (no echo)
2. **Backups**: Contain sensitive data - store securely
3. **Admin User**: Create strong passwords in production
4. **Environment Variables**: Never commit credentials to version control

# MariaDB Database Setup

This directory contains initialization scripts for the MariaDB database container used in the Bamort application development environment.

## Configuration

The MariaDB container is configured with the following settings:

- **Image**: mariadb:11.4
- **Database**: bamort
- **Username**: bamort  
- **Password**: bG4)efozrc
- **Root Password**: root_password_dev
- **Port**: 3306 (exposed to host)
- **Charset**: utf8mb4
- **Collation**: utf8mb4_unicode_ci

## Volume Mounts

- `mariadb_data`: Persistent data volume for database files
- `./init-db`: Initialization scripts directory

## Initialization Scripts

Scripts in this directory (`init-db/`) are executed automatically when the MariaDB container starts for the first time:

1. `01-init.sql`: Basic database setup and user privileges

## Health Check

The container includes a health check that verifies:
- Database connectivity
- InnoDB engine initialization

## Usage

The MariaDB container starts automatically when running:

```bash
docker-compose -f docker-compose.dev.yml up
```

The backend service will wait for the database to be healthy before starting, thanks to the `depends_on` configuration with `condition: service_healthy`.

## Connection Details

- **Host**: `mariadb` (within Docker network) or `localhost` (from host)
- **Port**: 3306
- **Database**: bamort
- **Username**: bamort
- **Password**: bG4)efozrc

## Data Persistence

Database data is stored in the `mariadb_data` Docker volume and persists across container restarts.

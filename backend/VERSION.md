# Version Management

## Current Version: 0.1.30

The backend version is managed in `/backend/config/version.go`.

## Updating the Version

To update the application version:

1. Edit `/backend/config/version.go`
2. Change the `Version` constant:
   ```go
   const Version = "0.1.31"  // Update this
   ```

## Git Commit Information

The git commit hash is automatically detected at runtime from:

1. `GIT_COMMIT` environment variable (preferred for Docker)
2. Git command output (works in development)
3. Falls back to "unknown" if neither is available

## Setting Git Commit in Docker

### Development Environment

Add to `docker/.env`:
```bash
GIT_COMMIT=$(git rev-parse --short HEAD)
```

Then restart the container:
```bash
cd docker
docker-compose -f docker-compose.dev.yml up -d backend-dev
```

### Production Build

Use build-time variable:
```bash
docker build --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) .
```

Or set in docker-compose.yml:
```yaml
environment:
  - GIT_COMMIT=${GIT_COMMIT:-unknown}
```

## API Endpoints

- **Public**: `GET /api/public/version` (no authentication)
- **Protected**: `GET /api/version` (requires JWT token)

Both return:
```json
{
  "version": "0.1.30",
  "gitCommit": "d0c177b"
}
```

## Frontend Integration

The landing page automatically fetches version information from `/api/public/version` and displays it.

To see the version in the frontend, visit: http://localhost:5173

# Bamort Services Quick Reference

## Development Environment URLs

| Service | URL | Credentials | Description |
|---------|-----|-------------|-------------|
| **Frontend** | http://localhost:5173 | - | Vue.js Application |
| **Backend API** | http://localhost:8180 | - | Go REST API |
| **phpMyAdmin** | http://localhost:8082 | root/root_password_dev | Database Management |
| **MariaDB** | localhost:3306 | bamort/bG4)efozrc | Direct Database Connection |

## Production Environment URLs

| Service | URL | Credentials | Description |
|---------|-----|-------------|-------------|
| **Frontend** | http://localhost:5173 | - | Vue.js Application |
| **Backend API** | http://localhost:8180 | - | Go REST API |
| **MariaDB** | localhost:3306 | bamort/[ENV_VAR] | Direct Database Connection |

> **Note**: phpMyAdmin is disabled in production by default. Uncomment the service in `docker-compose.yml` if needed (Port 8081).

## Container Names

| Service | Development | Production |
|---------|-------------|------------|
| MariaDB | bamort-mariadb-dev | bamort-mariadb |
| Backend | bamort-backend-dev | backend |
| Frontend | bamort-frontend-dev | frontend |
| phpMyAdmin | bamort-phpmyadmin-dev | bamort-phpmyadmin* |

*\*phpMyAdmin is disabled in production by default*

## Quick Commands

### Start Development Environment
```bash
cd /data/dev/bamort/docker

# Option 1: Use default configuration
docker-compose -f docker-compose.dev.yml up -d

# Option 2: Use custom environment variables
cp .env.example .env.dev
# Edit .env.dev with your settings
docker-compose -f docker-compose.dev.yml --env-file .env.dev up -d
```

### Start Production Environment
```bash
cd /data/dev/bamort/docker

# Create .env file with your configuration
cp .env.example .env
# Edit .env with your settings

docker-compose up -d
```

### View Logs
```bash
# All services
docker-compose -f docker-compose.dev.yml logs -f

# Specific service
docker-compose -f docker-compose.dev.yml logs -f mariadb
docker-compose -f docker-compose.dev.yml logs -f backend-dev
docker-compose -f docker-compose.dev.yml logs -f frontend-dev
docker-compose -f docker-compose.dev.yml logs -f phpmyadmin
```

### Database Access
```bash
# Via phpMyAdmin (Browser)
open http://localhost:8082

# Via Command Line
docker exec -it bamort-mariadb-dev mysql -u bamort -p bamort
```

## Port Summary

### Development
- **3306**: MariaDB
- **5173**: Vue.js Frontend (Vite dev server)
- **8080**: Backend API (accessible to frontend)
- **8081**: phpMyAdmin
- **8180**: Go Backend API (container port)

### Production  
- **3306**: MariaDB
- **5173**: Vue.js Frontend
- **8081**: phpMyAdmin (disabled by default)
- **8180**: Go Backend API

## Enabling phpMyAdmin in Production

phpMyAdmin is commented out in production for security reasons. To enable it:

1. **Edit `docker-compose.yml`**: Uncomment the phpMyAdmin service section
2. **Set Environment Variables**: Ensure `MARIADB_ROOT_PASSWORD` is set in your `.env` file
3. **Restart Services**: Run `docker-compose up -d phpmyadmin`

Access at: http://localhost:8081

```bash
# Enable phpMyAdmin in production
sed -i 's/^  # phpmyadmin:/  phpmyadmin:/' docker-compose.yml
sed -i 's/^  #   /    /' docker-compose.yml
docker-compose up -d phpmyadmin
```

# Environment variables for BaMoRT development environment

# API Configuration
# API_URL=http://localhost:8180
ENVIRONMENT=development
# Database Configuration (for development)
DATABASE_TYPE=mysql
DATABASE_URL="bamort:bG4)efozrc@tcp(192.168.0.36:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local"

# MariaDB Configuration (development)
MARIADB_ROOT_PASSWORD=root_password_dev
MARIADB_PASSWORD="bG4)efozrc"
MARIADB_DATABASE=bamort
MARIADB_USER=bamort

# Frontend Configuration
API_URL="http://localhost:8180"
VITE_API_URL="http://localhost:8180"
API_PORT=8180
BASE_URL="http://localhost:5173"
TEMPLATES_DIR=./templates
EXPORT_TEMP_DIR=./export_temp
LOG_LEVEL=debug
COMPOSE_PROJECT_NAME=bamort
CHROME_BIN="/usr/bin/chromium"

#./server
#/usr/local/go/bin/go run ./cmd/main.go
echo $DATABASE_URL
/home/de31a2/.local/bin/go/bin/go run ./cmd/main.go
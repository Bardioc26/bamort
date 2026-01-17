# How to setup
Setting up is quite easy. I both cases I would suggest docker

## Production Environment

* fetch certificate for frontend and backend 
  certbot certonly -d frontend.domain.de -d backend.domain.de --standalone
* Setup proxy
  Apache2:
    <IfModule mod_ssl.c>
    <VirtualHost *:443>
      ServerName frontend.domain.de

      SSLEngine on
      SSLCertificateFile /etc/letsencrypt/live/frontend.domain.de/fullchain.pem
      SSLCertificateKeyFile /etc/letsencrypt/live/frontend.domain.de/privkey.pem
      Include /etc/letsencrypt/options-ssl-apache.conf

      #-# Request header rules
      RequestHeader set X-Forwarded-Proto "https"

      ProxyPreserveHost On
      ProxyRequests Off

      ProxyPass / http://localhost:8181/
      ProxyPassReverse / http://localhost:8181/ timeout=120

      LogLevel warn
      ErrorLog ${APACHE_LOG_DIR}/frontend.domain.de.error.log

    </VirtualHost>
    <VirtualHost *:443>
      ServerName backend.domain.de
      
      SSLEngine on
      SSLCertificateFile /etc/letsencrypt/live/bamort.trokan.de/fullchain.pem
      SSLCertificateKeyFile /etc/letsencrypt/live/bamort.trokan.de/privkey.pem
      Include /etc/letsencrypt/options-ssl-apache.conf

      #-# Request header rules
      RequestHeader set X-Forwarded-Proto "https"

      ProxyPreserveHost On
      ProxyRequests Off

      ProxyPass / http://localhost:8182/
      ProxyPassReverse / http://localhost:8182/ timeout=120

      LogLevel warn
      ErrorLog ${APACHE_LOG_DIR}/backend.domain.de.error.log
      TransferLog ${APACHE_LOG_DIR}/backend.domain.de.transfer.log
      CustomLog ${APACHE_LOG_DIR}/backend.domain.de.access.log combined
    </VirtualHost>
    </IfModule>

* Set URLS  to configure frontend and backend
  edit ./docker/.env
    #- Environment variables for Bamort production environment

    #- API Configuration
    API_URL=https://backend.domain.de
    VITE_API_URL=https://backend.domain.de

    #- Database Configuration Backend
    DATABASE_TYPE=mysql
    #DATABASE_URL=bamort:your_secure_user_password@tcp(mariadb:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local

    #- MariaDB Configuration 
    MARIADB_ROOT_PASSWORD=your_secure_root_password
    MARIADB_PASSWORD=your_secure_user_password
    MARIADB_DATABASE=bamort
    MARIADB_USER=bamort

    API_PORT=8180
    BASE_URL=https://frontend.domain.de
    TEMPLATES_DIR=./templates
    EXPORT_TEMP_DIR=./export_temp
  
  !!! Do not configure DATABASE_URL for production in .env !!!

* edit ./frontend/src/utils/api.js because the env variable does not work reliably
    const API = axios.create({
      baseURL: import.meta.env.VITE_API_URL || 'https://backend.domain.de', // Use env variable with fallback
    })

* edit ./docker/docker-compose.yml
    Set VITE_API_URL and VITE_BASE_URL to https://backend.domain.de and https://frontend.domain.de

* run ./docker/start-prd.sh
* test https://backend.domain.de/api/public/version
  should responde like this: {"version":"0.1.30"}


## Development Environment
* Edit ./docker/.env
  set 
    API_URL=http://localhost:8180
    VITE_API_URL=http://localhost:8180
    API_PORT=8180
    BASE_URL=http://localhost:5173

* run ./docker/start-dev.sh

Frontend is not reachable at http://localhost:5173 and backend at http://localhost:8180.
Both containers are reacting on code changes and reload or rebuild automatically
Database is automatically created with env variables in MARIADB_* set your DATABASE_URL accordingly
Database is filled at first creation with values from ./docker/initdb/*.sql
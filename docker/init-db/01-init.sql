-- Initialize MariaDB database for Bamort application
-- This script runs automatically when the MariaDB container starts for the first time

-- Ensure UTF-8 charset and collation
ALTER DATABASE bamort CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Grant additional privileges to the bamort user if needed
GRANT ALL PRIVILEGES ON bamort.* TO 'bamort'@'%';
FLUSH PRIVILEGES;

-- Log initialization
SELECT 'MariaDB initialization completed for Bamort application' AS message;

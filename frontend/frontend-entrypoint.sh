#!/bin/sh

# Generate runtime config if VITE_API_URL is set
if [ -n "$VITE_API_URL" ]; then
  cat > /usr/share/nginx/html/config.js <<EOF
window.runtimeConfig = {
  apiUrl: "$VITE_API_URL"
}
EOF
else
  echo "Warning: VITE_API_URL not set, using build-time configuration"
fi

# Start nginx
exec nginx -g "daemon off;"

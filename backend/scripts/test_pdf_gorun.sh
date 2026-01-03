#!/bin/bash
# Test PDF export via running server

echo "Testing PDF export with go run..."
echo ""

# Start the server in background
cd /data/dev/bamort/backend
/usr/local/go/bin/go run ./cmd/main.go &
SERVER_PID=$!

echo "Server started with PID: $SERVER_PID"
echo "Waiting for server to start..."
sleep 5

# Check if server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "✗ Server failed to start"
    exit 1
fi

echo "Server is running, testing PDF export..."

# Get auth token (you'll need to implement this based on your auth)
# For now, just try to access the endpoint
curl -s "http://localhost:8180/api/pdf/templates" || echo "Server not responding yet..."

echo ""
echo "Press Ctrl+C to stop the server (PID: $SERVER_PID)"
echo "Then manually test: curl -H 'Authorization: Bearer YOUR_TOKEN' 'http://localhost:8180/api/pdf/export/18?template=Default_A4_Quer'"

# Wait for user interrupt
wait $SERVER_PID

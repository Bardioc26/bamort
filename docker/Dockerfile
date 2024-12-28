# =========== 1) Build stage ===========
FROM golang:1.20-alpine AS builder

# Create and set working directory
WORKDIR /app

# Copy go.mod and go.sum first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the backend code
COPY . .

# Build the Go binary
RUN go build -o server main.go

# =========== 2) Runtime stage ===========
FROM alpine:3.18

WORKDIR /app

# Copy the compiled binary from builder stage
COPY --from=builder /app/server /app

# Expose port 8080 (or your backend port)
EXPOSE 8080

# Run the Go server
CMD ["./server"]

# Stage 1: Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Ensure clean dependencies
RUN go mod tidy

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Stage 2: Lightweight runtime stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS requests (Vault API calls)
RUN apk --no-cache add ca-certificates

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create non-root user for security
RUN adduser -D -g '' appuser
USER appuser

# Expose service port
EXPOSE 5003

CMD ["./main"]
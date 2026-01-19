# Stage 1: Build stage
FROM golang:1.23 AS builder

# Build arguments for Vault configuration
ARG VAULT_ADDR
ARG VAULT_TOKEN
ARG SERVICE_NAME

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Create config file from build arguments
RUN echo "VAULT_ADDR=${VAULT_ADDR}" > .env && \
    echo "VAULT_TOKEN=${VAULT_TOKEN}" >> .env && \
    echo "SERVICE_NAME=${SERVICE_NAME}" >> .env

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Stage 2: Lightweight runtime stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Environment variables will be loaded from .env file at runtime
CMD ["./main"]














# # Stage 1: Build stage
# FROM golang:1.23 AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .
# RUN go mod tidy

# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# # Stage 2: Lightweight runtime stage
# FROM alpine:latest

# WORKDIR /app

# # Install CA certificates for HTTPS requests
# RUN apk --no-cache add ca-certificates

# # Copy the binary from builder stage
# COPY --from=builder /app/main .

# # Create non-root user
# RUN adduser -D -g '' appuser
# USER appuser

# CMD ["./main"]
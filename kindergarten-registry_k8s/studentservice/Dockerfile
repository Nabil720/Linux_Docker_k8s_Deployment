# Stage 1: Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies and tidy them
RUN go mod download && go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Stage 2: Runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

# Copy binary
COPY --from=builder /app/main .

# Copy .env file if it exists
COPY --from=builder /app/.env .env 2>/dev/null || :

RUN adduser -D -g '' appuser
USER appuser

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
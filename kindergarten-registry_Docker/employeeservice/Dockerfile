FROM golang:1.23

WORKDIR /app

# Copy go mod files first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Run go mod tidy to ensure all dependencies are resolved
RUN go mod tidy

RUN go build -o main .

CMD ["./main"]
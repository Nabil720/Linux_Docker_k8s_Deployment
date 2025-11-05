FROM golang:1.22

WORKDIR /app

# Copy go mod files first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

RUN go build -o main .

CMD ["./main"]
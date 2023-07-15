# Stage 1: Build the Go application
FROM golang:1.20.5 AS builder

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the Go source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -o app .

# Stage 2: Create a minimal production image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]

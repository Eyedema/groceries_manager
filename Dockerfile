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
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .


EXPOSE 8080

# Run the Go application
CMD ["./app"]

# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies - uses cached layers if possible
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out ./cmd/api

# Stage 2: Production
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/out .
# COPY --from=builder /app/.env .

# Expose the port your application listens on
EXPOSE 3000

# Command to run your application
CMD ["./out"]

# Step 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy the Go source code
COPY . .

# Download Go modules
RUN go mod tidy

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o order-service .

# Step 2: Create a minimal image to run the application
FROM alpine:latest

WORKDIR /root/

# Install any necessary libraries
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/order-service .
COPY --from=builder /app/static ./static


# Expose the application port
EXPOSE 8081

# Command to run the application
CMD ["./order-service"]

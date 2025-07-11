# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .


# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

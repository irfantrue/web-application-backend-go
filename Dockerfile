# syntax=docker/dockerfile:1

# Stage 1: Builder
FROM golang:1.27-alpine AS builder

# Install build dependencies
RUN apk update && apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy module files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a fully static binary (CGO disabled) with stripped debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o main ./cmd/api/main.go

# Stage 2: Runner
FROM alpine:3.18

# Install CA certificates for outbound HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

# Set timezone
ENV TZ=Asia/Jakarta

# Create a non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/main .

# Copy .env if present (remove this line if you don't use a .env file)
COPY --from=builder /app/.env .

# Switch to non-root user
USER appuser

# Expose the Fiber HTTP port
EXPOSE 8080

# Run the application
CMD ["./main"]

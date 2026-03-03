# syntax=docker/dockerfile:1

FROM golang:1.22-alpine as builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the clean DDD entrypoint
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /remi-server ./cmd/server

# Create a minimal runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /remi-server ./

# Expose the port the app listens on
EXPOSE 8080

# Run the app
CMD ["/app/remi-server"]

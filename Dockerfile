# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the example server (disabling CGO for a static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./example/main.go

# Run Stage
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server .

# Copy the frontend HTML file
COPY --from=builder /app/example/index.html .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./server"]

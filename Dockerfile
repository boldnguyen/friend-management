FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Install sqlboiler
RUN go get github.com/volatiletech/sqlboiler/v4@latest

# Use the official PostgreSQL image for running the database
FROM postgres:11-alpine

# Set environment variables for PostgreSQL
ENV POSTGRES_USER friend-management
ENV POSTGRES_PASSWORD "1234" 
ENV POSTGRES_DB friend-management
EXPOSE 5000

# Copy the Go application binary from the builder stage into the PostgreSQL image
COPY --from=builder /app/main /app/main

# Run the Go application
CMD ["/app/main"]
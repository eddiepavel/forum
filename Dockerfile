FROM golang:1.23.4 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o forum-app

# Use a minimal base image for the final container
FROM debian:bookworm-slim

# Install SQLite3 runtime dependency
RUN apt-get update && apt-get install -y --no-install-recommends \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/forum-app /app/forum-app

# Copy the entire assets directory
COPY assets /app/assets

# Copy the migrations directory
COPY migrations /app/migrations

# Ensure the database directory exists
RUN mkdir -p /app/database/file

# Expose the port the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./forum-app"]
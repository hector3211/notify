FROM golang:1.22-alpine

# Install necessary pkgs for CGO (libc-dev)
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
ENV CGO_ENABLED=1
RUN go build -o server ./cmd/main.go

# Copy your SQLite database file
COPY mydb.db /data/mydb.db

# Set the path for the SQLite database (inside the volume)
# Assuming your app references this path for the SQLite database
ENV NOTIFY_DB_PATH=/data/mydb.db

# Expose the port on which your app will run
EXPOSE 8080

# Command to run your application
CMD ["./server"]

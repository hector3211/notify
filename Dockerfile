FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o server ./cmd/main.go

# Set the path for the SQLite database (inside fly.io volume)
# Assuming your app references this path for the SQLite database
ENV NOTIFY_DB_PATH=/data/mydb.db

# Expose the port on which your app will run
EXPOSE 8080

# Command to run your application
CMD ["./server"]

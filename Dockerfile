# Base image
FROM golang:1.22

# Set working directory
WORKDIR /app

# Copy dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY .env .

COPY . .

# Build the Go app
RUN go build -o main .

# Install MySQL client
RUN apt-get update && apt-get install -y default-mysql-client

# Install Redis client
RUN apt-get install -y redis-tools

# Command to run the executable
CMD ["./main"]
# # Base image
# FROM golang:1.22

# # Set working directory
# WORKDIR /app

# # Copy dependencies
# COPY go.mod .
# COPY go.sum .

# # Download dependencies
# RUN go mod download

# # Copy the rest of the application
# COPY .env .

# COPY . .

# # Build the Go app
# RUN go build -o main .

# # Install MySQL client
# RUN apt-get update && apt-get install -y default-mysql-client

# # Install Redis client
# RUN apt-get install -y redis-tools

# # Command to run the executable
# # CMD ["./main"]
# CMD ["go mod init", "./main"]

# # FROM golang:1.22

# # RUN mkdir /app
# # WORKDIR /app

# # RUN go get github.com/gin-gonic/gin
# # CMD ["go mod init", "./main"]

FROM golang:alpine

RUN mkdir /app

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD . .

RUN go install github.com/githubnemo/CompileDaemon

ADD /db/my.cnf/my.cnf /etc/mysql/mysql.conf.d/mysqld.cnf

EXPOSE 9090

ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
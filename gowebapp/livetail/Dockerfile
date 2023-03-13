# syntax=docker/dockerfile:1

FROM golang:1.19.4-alpine as builder

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

# # Setup folders
RUN mkdir /app
WORKDIR /app


# Copy the source from the current directory to the working Directory inside the container
# COPY ./app ./app
# COPY main.go .
COPY go.mod go.sum /
# COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Setup folders
# RUN mkdir /app
# WORKDIR /app

COPY . .
# RUN go build -o main .
RUN CGO_ENABLED=0 go build -o /my_main -ldflags="-w -s"



######## Start a new stage from scratch #######
FROM golang:1.19.4-alpine
# FROM alpine  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /my_main .
COPY --from=builder /app/logs_message.txt .
COPY --from=builder /app/dummy_logs.txt .

# Expose port 8080 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./my_main"] 
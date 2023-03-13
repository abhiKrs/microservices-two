
FROM golang:1.19.4-alpine as builder

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

# Setup folders
RUN mkdir /app
WORKDIR /app


COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

RUN go build -o main .


######## Start a new stage from scratch #######
FROM golang:1.19.4-alpine  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/logs_message.txt .

# Expose port 8080 to the outside world
EXPOSE 4000

# Command to run the executable
CMD ["./main"] 
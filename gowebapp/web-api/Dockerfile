FROM golang:1.19.4-alpine as builder

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base && apk add --no-cache gcc musl-dev


# Setup folders
RUN mkdir /app
WORKDIR /app


COPY go.mod go.sum /

RUN go mod download

RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add pkgconf git bash build-base sudo
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make && make install

COPY . .

RUN go build -tags musl --ldflags "-extldflags -static" -o /my_main
# RUN CGO_ENABLED=0 go build -o /my_main -ldflags="-w -s"
# RUN CGO_ENABLED=1 go build -o /my_main
# EXPOSE 8080

# CMD [ "/app/main" ]

######## Start a new stage from scratch #######
FROM golang:1.19.4-alpine  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /my_main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./my_main"]

# FROM golang:1.19.4-alpine as builder

# RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base
# RUN apt-get update && \
#     apt-get -y install gcc

# # Setup folders
# RUN mkdir /app
# WORKDIR /app


# # Copy the source from the current directory to the working Directory inside the container
# # COPY ./app ./app
# # COPY main.go .
# # COPY go.mod .
# COPY go.mod go.sum /

# RUN go mod download

# # # Download all the dependencies
# # RUN go get -d -v ./...

# # # Install the package
# # RUN go install -v ./...

# # COPY ./app ./app
# # COPY main.go .
# COPY . .

# # RUN go build -o main .
# RUN CGO_ENABLED=0 go build -o /my_main -ldflags="-w -s"
# # EXPOSE 8080

# # CMD [ "/app/main" ]

# ######## Start a new stage from scratch #######
# FROM golang:1.19.4-alpine  

# RUN apk --no-cache add ca-certificates

# WORKDIR /root/

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /my_main .

# # Expose port 8080 to the outside world
# EXPOSE 8080

# # Command to run the executable
# CMD ["./my_main"]
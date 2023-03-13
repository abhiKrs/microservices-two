# FROM golang:1.19.4-alpine as builder

FROM flink:1.16.1-scala_2.12-java11



RUN apt-get update
RUN apt install -y build-essential checkinstall software-properties-common
RUN add-apt-repository universe
RUN apt update
RUN apt-get install -y bash gcc musl-dev openssl golang-go 
RUN apt install ca-certificates wget && update-ca-certificates

# Install python3.8
RUN apt install sudo
RUN apt-get install -y software-properties-common && sudo add-apt-repository ppa:deadsnakes/ppa && apt-get update && apt-get install -y \
  python3.8 \
  python3-pip \
  && rm -rf /var/lib/apt/lists/*
RUN ln -s /usr/bin/python3.8 /usr/bin/python


# install golang 1.19.4
RUN wget -O go.tgz https://golang.org/dl/go1.19.4.src.tar.gz
RUN tar -C /usr/local -xzf go.tgz
RUN cd /usr/local/go/src/ && ./make.bash
ENV PATH "/usr/local/go/bin:$PATH"
ENV GOPATH "/opt/go/"
ENV PATH "$PATH:$GOPATH/bin"
RUN go version


# ENV GOPATH /root
ENV FLINK_HOME "/opt/flink"
# ENV GOPATH $FLINK_HOME
ENV PATH $GOPATH/bin:$PATH
# RUN echo "${PATH}"
# RUN echo "${GOPATH}"
RUN mkdir -p $GOPATH/src $GOPATH/bin && chmod -R 777 $GOPATH

# WORKDIR $GOPATH
WORKDIR $FLINK_HOME
# WORKDIR /root/
# =======================
# RUN mkdir /app
# WORKDIR /app


# Copy the source from the current directory to the working Directory inside the container
# COPY ./app ./app
# COPY main.go ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN echo $(ls -las /opt/flink)
# RUN rm ${FLINK_HOME}/go.mod
# Install the package
# RUN go get -d -v ./...
# RUN go install -v ./...
COPY ./app ./app
COPY main.go ./

RUN  CGO_ENABLED=0 go build -o main .
# RUN CGO_ENABLED=0 go build -o /main -ldflags="-w -s"
# ========================


COPY ./bin ./bin
COPY ./*.py ./
# COPY ./requirements.txt .

# Expose port 8080 to the outside world
EXPOSE 8005

# Command to run the executable
CMD ["./main"]
# CMD ["main"]

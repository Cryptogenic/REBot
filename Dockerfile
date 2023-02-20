FROM golang:1.20

WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update

RUN apt install -y build-essential cmake git keystone libcapstone-dev

RUN git clone https://github.com/keystone-engine/keystone.git

WORKDIR /app/keystone
RUN mkdir build && cd build && cmake .. && make && make install

WORKDIR /app

COPY go.mod ./ 
COPY go.sum ./
COPY *.go ./

RUN go mod download

RUN go build

RUN go test


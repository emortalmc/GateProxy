FROM golang:alpine AS build

RUN apk add --no-cache --update curl ca-certificates openssl git tar bash sqlite fontconfig \
    && adduser --disabled-password --home /home/container container

# Build
# Copy the go source
COPY command ./command
COPY game ./game
COPY nbs ./nbs
COPY redisdb ./redisdb
COPY proxy.go ./

WORKDIR /app/

COPY go.mod go.sum /app/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o /app/proxy proxy.go

USER container
ENV USER=container HOME=/home/container
WORKDIR /home/container



COPY /app/proxy ./proxy
COPY ./entrypoint.sh /entrypoint.sh
CMD [ "/bin/bash", "/entrypoint.sh" ]
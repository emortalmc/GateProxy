FROM golang:alpine

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
COPY . /home/container/

USER container
ENV USER=container HOME=/home/container
WORKDIR /home/container

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o ./proxy proxy.go

COPY ./proxy /home/container/proxy
COPY ./entrypoint.sh /entrypoint.sh
CMD [ "/bin/bash", "/entrypoint.sh" ]
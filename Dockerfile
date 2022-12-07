FROM golang:alpine AS build

WORKDIR /app/

# Build
# Copy the go source
COPY command ./command
COPY game ./game
COPY nbs ./nbs
COPY redisdb ./redisdb
COPY proxy.go ./

COPY go.mod go.sum /app/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o proxy proxy.go

USER container
ENV USER=container HOME=/home/container
WORKDIR /home/container

FROM golang:alpine as exp
RUN apk add --no-cache --update curl ca-certificates openssl git tar bash sqlite fontconfig \
    && adduser --disabled-password --home /home/container container

COPY --from=build /app/proxy /proxy
COPY ./entrypoint.sh /entrypoint.sh
CMD [ "/bin/bash", "/entrypoint.sh" ]
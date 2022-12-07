FROM golang:1.19.4-bullseye

RUN         apt-get update -y \
            && apt-get install -y --no-install-recommends curl ca-certificates openssl git tar sqlite3 fontconfig libfreetype6 libstdc++6 lsof build-essential tzdata iproute2 locales \
            && apt-get clean \
            && rm -rf /var/lib/apt/lists/* \
            && useradd -m -d /home/container container \
            && locale-gen en_US.UTF-8

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
COPY . /app/

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o proxy proxy.go

USER container
ENV USER=container HOME=/home/container
WORKDIR /home/container

COPY /app/proxy .

COPY ./entrypoint.sh /entrypoint.sh
CMD [ "/bin/bash", "/entrypoint.sh" ]
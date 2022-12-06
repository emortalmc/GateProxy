FROM golang:1.19 AS build

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download\
   && useradd -m -d /home/container container
#                               && locale-gen en_US.UTF-8

USER        container
#ENV         LC_ALL=en_US.UTF-8
#ENV         LANG=en_US.UTF-8
#ENV         LANGUAGE=en_US.UTF-8
ENV         USER=container HOME=/home/container
WORKDIR /home/container

# Copy the go source
COPY command ./command
COPY game ./game
COPY nbs ./nbs
COPY redisdb ./redisdb
COPY proxy.go ./


COPY        ./entrypoint.sh /entrypoint.sh
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o proxy proxy.go

# Move binary into final image
FROM gcr.io/distroless/static:nonroot AS app
COPY --from=build /home/container/proxy /
CMD         [ "/bin/bash", "/entrypoint.sh" ]
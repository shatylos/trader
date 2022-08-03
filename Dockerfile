# Compile stage
FROM golang:1.18 AS build-env

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

ADD . /trader
WORKDIR /trader

# RUN go build -o /server
RUN go build -gcflags "all=-N -l" -o /server

# Final stage
FROM debian:buster

RUN mkdir /app && \
    mkdir /app/cert && \
    apt-get update && \
    apt-get install -y openssl && \
    openssl genrsa -des3 -passout pass:x -out /app/cert/server.pass.key 2048 && \
    openssl rsa -passin pass:x -in /app/cert/server.pass.key -out /app/cert/server.key && \
    rm /app/cert/server.pass.key && \
    openssl req -new -key /app/cert/server.key -out /app/cert/server.csr \
        -subj "/C=UA/ST=Kyiv/L=Boryspil/O=Shatylo/OU=IT Department/CN=shatylo.pp.ua" && \
    openssl x509 -req -days 365 -in /app/cert/server.csr -signkey /app/cert/server.key -out /app/cert/server.crt

EXPOSE 8080 40000

WORKDIR /app/

COPY --from=build-env /go/bin/dlv /app/
COPY --from=build-env /server /app/

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server"]

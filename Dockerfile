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

EXPOSE 8080 40000

WORKDIR /app/

COPY --from=build-env /go/bin/dlv /app/
COPY --from=build-env /server /app/

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server"]

# Build stage
FROM golang:1.16-alpine AS build-env

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/app

COPY . .
RUN go mod download
RUN go build -tags emoji -o backend

# final stage
FROM alpine:3.9.5
WORKDIR /app
COPY --from=build-env /go/src/app/backend backend

ENTRYPOINT ./backend

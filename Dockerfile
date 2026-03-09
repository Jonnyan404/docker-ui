FROM golang:1.24-alpine AS build

WORKDIR /app

ENV GO111MODULE=on \
	CGO_ENABLED=0

RUN apk add --no-cache git

# Cache deps first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest and build
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /app/server .


FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=build /app/server /app/server
COPY ./log4go.xml /app/log4go.xml
COPY ./html /app/html

LABEL AUTHOR="jonnyan404"
LABEL LANGUAGE="golang"
LABEL PRODUCT="docker"
LABEL COPYRIGHT="jonnyan404"
LABEL DECLAIM="All right reserved by jonnyan404"

EXPOSE 8999

ENTRYPOINT  ["./server"]

FROM golang:1.15-alpine as builder

WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://proxy.golang.org,direct"

RUN apk add --no-cache git

COPY go.mod .
COPY go.sum .

RUN go mod download -x
COPY . .

RUN go build -a -tags 'netgo osusergo' -o /go/bin/fuck-off-as-a-service main.go

LABEL description=fuck-off-as-aservice
LABEL builder=true
LABEL maintainer='Santiago Leira'

FROM alpine
COPY --from=builder go/bin/fuck-off-as-a-service /usr/local/bin

WORKDIR usr/local/bin
ENTRYPOINT [ "fuck-off-as-a-service", "serve" ]
EXPOSE 8080

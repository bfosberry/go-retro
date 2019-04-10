FROM golang:1.12.1

COPY ./ /go/src/github.com/codeship/retro/
WORKDIR /go/src/github.com/codeship/retro

RUN go get ./...

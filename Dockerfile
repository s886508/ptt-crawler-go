FROM golang:1.16-alpine AS build-env

ADD go.mod src
ADD go.sum src
ADD cmd/ src/cmd
ADD pkg/ src/pkg

RUN cd src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/ptt-crawler cmd/ptt-crawler/main.go

FROM alpine:3.13
COPY --from=build-env /go/bin/ptt-crawler /

ENTRYPOINT ["/ptt-crawler"]

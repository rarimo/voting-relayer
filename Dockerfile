FROM golang:1.22-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/rarimo/voting-relayer
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/voting-relayer /go/src/github.com/rarimo/voting-relayer


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/voting-relayer /usr/local/bin/voting-relayer
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["voting-relayer"]

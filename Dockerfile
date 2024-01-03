FROM golang:1.20.11-alpine3.18 AS builder
WORKDIR /go/src/github.com/app
COPY . ./
RUN go mod download
RUN mkdir bin
RUN apk add build-base
RUN cd cmd/app/ && go build -o ../../bin/go-vendors-api

FROM golang:1.20.11-alpine3.18
WORKDIR /go
COPY --from=builder /go/src/github.com/app/bin/go-vendors-api ./
ENTRYPOINT ["./go-vendors-api"]
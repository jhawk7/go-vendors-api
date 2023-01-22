FROM golang:1.19-alpine3.17 AS builder
WORKDIR /go/src/github.com/app
COPY . ./
RUN go mod download
RUN mkdir bin
RUN apk add build-base
RUN cd cmd/main/ && go build -o ../../bin/go-vendors-api

FROM golang:1.19-alpine3.17
WORKDIR /go
COPY --from=builder /go/src/github.com/app/bin/go-vendors-api ./
ENTRYPOINT ["./go-vendors-api"]
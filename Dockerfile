FROM golang:1.19
WORKDIR /usr/src/app
EXPOSE 8080
COPY . ./
RUN go mod download
RUN mkdir bin
RUN cd cmd/main/ && go build -o ../../bin/go-vendors-api
ENTRYPOINT ["./bin/go-vendors-api"]
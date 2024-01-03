BINARY_NAME=go-vendor-api
MAIN_FILE=./cmd/app/main.go
GO=$(shell which go)

all: makedir build

makedir:
	@if [ ! -d ./bin ] ; then mkdir -p ./bin ; fi

build:
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME} ${MAIN_FILE}
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME} ${MAIN_FILE}

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm -rf ./bin

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all
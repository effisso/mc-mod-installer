BINARY_NAME=out/mcmods

all: build test

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test -v ./...
 
run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
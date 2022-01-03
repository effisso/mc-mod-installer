BINARY_NAME=out/mcmods

all: build test

build:
	go build -o ${BINARY_NAME} main.go

build-release: build-all-targets test

build-all-targets:
	env GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}-win.exe main.go
	env GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME}-linux-arm main.go
	env GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux-amd main.go
	env GOOS=darwin GOARCH=arm64 go build -o ${BINARY_NAME}-mac-arm main.go
	env GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}-mac-amd main.go

test:
	go test -v ./...
 
run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
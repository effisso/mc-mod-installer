OUT_FOLDER=out
BINARY_NAME=mcmods

all: build test

build:
	go build -o ${OUT_FOLDER}/${BINARY_NAME} main.go

build-test-release: build-all-targets test zip

build-all-targets:
	env GOOS=windows GOARCH=amd64 go build -o ${OUT_FOLDER}/win/${BINARY_NAME}.exe main.go
	@echo env GOOS=linux GOARCH=arm64 go build -o ${OUT_FOLDER}/linux-arm/${BINARY_NAME} main.go
	@echo env GOOS=linux GOARCH=amd64 go build -o ${OUT_FOLDER}/linux-amd/${BINARY_NAME} main.go
	@echo env GOOS=darwin GOARCH=arm64 go build -o ${OUT_FOLDER}/mac-arm/${BINARY_NAME} main.go
	@echo env GOOS=darwin GOARCH=amd64 go build -o ${OUT_FOLDER}/mac-amd/${BINARY_NAME} main.go

test:
	go test -v ./...

lint:
	golint -set_exit_status ./...

local-coverage:
	go test --coverprofile coverage.out --covermode count -v ./...
	go tool cover -html coverage.out

ci-coverage:
	go test --coverprofile coverage.out --covermode count -v ./...
	go tool cover -func coverage.out
 
run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}

get-ci-deps:
	sudo apt-get install -y zip
	go install golang.org/x/lint/golint@latest

zip:
	zip -r ${OUT_FOLDER}/${BINARY_NAME}-windows.zip ${OUT_FOLDER}/win/${BINARY_NAME}.exe
	@echo tar -zcvf ${OUT_FOLDER}/${BINARY_NAME}-linux-arm.tar.gz ${OUT_FOLDER}/linux-arm/${BINARY_NAME}
	@echo tar -zcvf ${OUT_FOLDER}/${BINARY_NAME}-linux-amd.tar.gz ${OUT_FOLDER}/linux-amd/${BINARY_NAME}
	@echo tar -zcvf ${OUT_FOLDER}/${BINARY_NAME}-mac-arm.tar.gz ${OUT_FOLDER}/mac-arm/${BINARY_NAME}
	@echo tar -zcvf ${OUT_FOLDER}/${BINARY_NAME}-mac-amd.tar.gz ${OUT_FOLDER}/mac-amd/${BINARY_NAME}
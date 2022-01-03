OUT_FOLDER=out
BINARY_NAME=mcmods

all: build test

build:
	go build -o ${BINARY_NAME} main.go

build-release:  get-ci-deps build-all-targets test zip

build-all-targets:
	env GOOS=windows GOARCH=amd64 go build -o ${OUT_FOLDER}/win/${BINARY_NAME}.exe main.go
	env GOOS=linux GOARCH=arm64 go build -o ${OUT_FOLDER}/linux-arm/${BINARY_NAME} main.go
	env GOOS=linux GOARCH=amd64 go build -o ${OUT_FOLDER}/linux-amd/${BINARY_NAME} main.go
	env GOOS=darwin GOARCH=arm64 go build -o ${OUT_FOLDER}/mac-arm/${BINARY_NAME} main.go
	env GOOS=darwin GOARCH=amd64 go build -o ${OUT_FOLDER}/mac-amd/${BINARY_NAME} main.go

test:
	go test -v ./...
 
run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}

get-ci-deps:
	sudo apt-get install -y zip

zip:
	zip -r ${OUT_FOLDER}/${BINARY_NAME}-windows.zip ${OUT_FOLDER}/win/${BINARY_NAME}.exe
	tar -zcvf ${OUT_FOLDER}/${BINARY_NAME}-linux-arm.tar.gz ${OUT_FOLDER}/linux-arm/${BINARY_NAME}
	tar -zcvf ${OUT_FOLDER}/${BINARY_NAME}-linux-amd.tar.gz ${OUT_FOLDER}/linux-amd/${BINARY_NAME}
	tar -zcvf ${OUT_FOLDER}/${BINARY_NAME}-mac-arm.tar.gz ${OUT_FOLDER}/mac-arm/${BINARY_NAME}
	tar -zcvf ${OUT_FOLDER}/${BINARY_NAME}-mac-amd.tar.gz ${OUT_FOLDER}/mac-amd/${BINARY_NAME}
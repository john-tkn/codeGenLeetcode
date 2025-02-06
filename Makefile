BINARY_NAME=ltgen

build:
	go build -o ${BINARY_NAME} gen.go
run: build
	./${BINARY_NAME}
clean:
	go clean
	rm ${BINARY_NAME}

.PHONY: all clean build run
BINARY_NAME=main

all: clean build

clean:
	rm -f ${BINARY_NAME}

build:
	go build -o ${BINARY_NAME} cmd/main.go

run: build
	docker-compose -f docker-compose.yml up
	./${BINARY_NAME}

dev:
	air

dev-sync:
	air -- -sync

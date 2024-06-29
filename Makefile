.PHONY: all clean build run db dev dev-sync test
BINARY_NAME=main

all: clean build

clean:
	rm -f ${BINARY_NAME}

build:
	go build -o ${BINARY_NAME} cmd/main.go

run: build
	./${BINARY_NAME}

db:
	docker-compose -f docker-compose.yml up -d db

dev: db
	@go run github.com/air-verse/air@v1.52.3 \
	--build.cmd "go build --tags dev -o tmp/bin/${BINARY_NAME} ./cmd/" --build.bin "tmp/bin/${BINARY_NAME}" --build.delay "100" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true \
	--screen.clear_on_rebuild true \
	--log.main_only true

dev-sync: db
	@go run github.com/air-verse/air@v1.52.3 \
	--build.cmd "go build --tags dev -o tmp/bin/${BINARY_NAME} ./cmd/" --build.bin "tmp/bin/${BINARY_NAME}" --build.delay "100" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true \
	--screen.clear_on_rebuild true \
	--log.main_only true \
 	-- -sync

test:
	go test -v ./... -count=1

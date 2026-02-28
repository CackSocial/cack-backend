.PHONY: run build test lint swagger docker-up docker-down clean

# Go parameters
BINARY_NAME=server
MAIN_PACKAGE=./cmd/server

run:
	go run $(MAIN_PACKAGE)

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run ./...

swagger:
	swag init -g cmd/server/main.go -o docs

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

clean:
	rm -rf bin/ docs/
	go clean

.DEFAULT_GOAL := run

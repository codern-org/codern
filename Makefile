BINARY_NAME = codern
VERSION := $(shell git describe --tags --abbrev=0)

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build \
	-ldflags="-s -w -X 'github.com/codern-org/codern/internal/constant.Version=$(VERSION)'" \
	-o dist/$(BINARY_NAME) main.go

clean:
	go clean
	rm -rf dist/$(BINARY_NAME)

dev:
	ENVIRONMENT=development go run .

lint:
	golangci-lint run

migrate-db:
	go run ./internal/cmd/mysql_migration.go

swagger:
	swag init -o ./other/swagger

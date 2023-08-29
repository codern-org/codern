BINARY_NAME=codern

build:
	GOARCH=amd64 GOOS=linux go build -o dist/${BINARY_NAME} main.go

clean:
	go clean
	rm -rf dist/${BINARY_NAME}

deps:
	go mod download

dev:
	ENVIRONMENT=development go run .

lint:
	golangci-lint run

migrate-db:
	go run ./internal/cmd/mysql_migration.go

swagger:
	swag init

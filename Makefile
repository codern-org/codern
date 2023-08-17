BINARY_NAME=codern

build:
	GOARCH=amd64 GOOS=linux go build -o dist/${BINARY_NAME} main.go

clean:
	go clean
	rm -rf dist/${BINARY_NAME}

migrate-db:
	go run ./internal/cmd/mysql_migration.go

deps:
	go mod download

lint:
	golangci-lint run

run:
	ENVIRONMENT=development go run .

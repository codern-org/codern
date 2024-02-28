APP 			:= codern
VERSION		:= $(shell git describe --tags --abbrev=0 || echo "0.0.0")
DATE			:= $(shell date +"%Y-%m-%dT%H:%M:%SZ")

VER_FLAGS = -X 'github.com/codern-org/codern/internal/constant.Version=$(VERSION)'

LD_FLAGS ?=	-s -w

.PHONY: build
build:
	$(info Build version $(VERSION) ($(DATE)))

	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build \
		-ldflags="$(LD_FLAGS) $(VER_FLAGS)" \
		-o dist/$(APP) main.go

	@echo Build sucessfully

.PHONY: clean
clean:
	go clean
	rm -rf dist/$(APP)

.PHONY: dev
dev:
	ENVIRONMENT=development go run .

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	go vet
	golangci-lint run

.PHONY: vulncheck
vulncheck:
	govulncheck -show verbose ./...

.PHONY: sec
sec:
	gosec ./...

.PHONY: migrate-db
migrate-db:
	go run ./internal/cmd/mysql_migration.go

.PHONY: swagger
swagger:
	swag init --parseDependency -o ./other/swagger

FROM mcr.microsoft.com/devcontainers/go:1.21

ENV TZ=Asia/Bangkok

# Install golangci-lint for linting
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2

# Install golang-migrate for database migration
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install swagger doc generator from Go annotations
RUN go install github.com/swaggo/swag/cmd/swag@latest

CMD ["sleep", "infinity"]

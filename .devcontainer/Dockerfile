FROM mcr.microsoft.com/devcontainers/go:0-1.20

# Install golangci-lint for linting
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3

# Install golang-migrate for database migration
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

CMD ["sleep", "infinity"]

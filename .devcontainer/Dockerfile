FROM mcr.microsoft.com/devcontainers/go:1.21

# No need to install tzdata for Debian
ENV TZ=Asia/Bangkok

# Install golangci-lint for linting
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install swagger doc generator from Go annotations
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Install govulncheck for vulnerability detection 
RUN go install golang.org/x/vuln/cmd/govulncheck@latest

# Install gosec for security checker
RUN go install github.com/securego/gosec/v2/cmd/gosec@latest

CMD ["sleep", "infinity"]

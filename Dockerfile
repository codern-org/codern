# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
RUN apk add make && apk add ca-certificates

COPY ./ ./
RUN go mod download
RUN make build

# Runner stage
FROM scratch AS runner
WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/dist/ ./

EXPOSE 3000
CMD ["/app/codern"]

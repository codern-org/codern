# Build stage
FROM golang:1.20 AS builder
WORKDIR /app

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

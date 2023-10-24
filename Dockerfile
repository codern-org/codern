# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
RUN apk add make

COPY ./ ./
RUN go mod download
RUN make build

# Runner stage
FROM scratch AS runner
WORKDIR /app

COPY --from=builder /app/dist/ ./

EXPOSE 3000
CMD ["/app/codern"]

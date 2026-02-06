# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-service ./cmd/api-service

# Run stage
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /api-service .
COPY cmd/api-service/VERSION .

EXPOSE 8080
ENV ENV=PRD \
    PORT=8080 \
    VERSION=VERSION

CMD ["./api-service"]

FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api ./cmd/api/main.go

FROM golang:1.24-alpine AS dev
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY . .
CMD ["air", "-c", ".air.toml"]

FROM alpine:3.20 AS prod
RUN apk add --no-cache curl
COPY --from=builder /app/api /api
EXPOSE 8080
CMD ["/api"]

FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# CMD ["go", "run", "./cmd/api/main.go"]

RUN go build -o api ./cmd/api/main.go

# Final image
FROM alpine:latest
COPY --from=builder app/api /api

EXPOSE 8080

CMD ["/api"]

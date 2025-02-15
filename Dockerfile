FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/app

FROM alpine

WORKDIR /app

COPY --from=builder /app/app .
COPY config.yml .
COPY migrations ./migrations

CMD ["./app"]

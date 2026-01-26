FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tracker-bot ./cmd/tracker-bot
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrator

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/tracker-bot /app/tracker-bot
COPY --from=builder /app/migrate /app/migrator
COPY --from=builder /app/.env /app/.env

COPY migrations /app/migrations

CMD ["/app/tracker-bot"]
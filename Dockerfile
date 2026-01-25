FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tracker-bot ./cmd/tracker-bot

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/tracker-bot /app/tracker-bot
COPY --from=builder /app/.env /app/.env

CMD ["/app/tracker-bot"]
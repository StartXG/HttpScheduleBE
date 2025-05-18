FROM golang:1.24.2 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o http_schedule_be ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/http_schedule_be .
COPY --from=builder /app/etc .

EXPOSE 8080

CMD ["./http_schedule_be"]


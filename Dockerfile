FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/config/dbconf.env /app/config/dbconf.env
EXPOSE 8080
CMD ["./main"]

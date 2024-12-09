FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine
RUN apk update && apk add ca-certificates
WORKDIR /root/

COPY --from=builder /app/app .
COPY --from=builder /app/config/dbconf.env /root/config/dbconf.env
CMD ["./app"]
EXPOSE 8080
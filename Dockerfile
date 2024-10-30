FROM golang:1.23.2 AS builder
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp cmd/app/main.go
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .
COPY --from=builder /app/config.yaml .
EXPOSE 8080
CMD ["./myapp"]

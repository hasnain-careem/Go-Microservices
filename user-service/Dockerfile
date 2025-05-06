FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/user-service

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o user-service .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/user-service/user-service .
COPY --from=builder /app/user-service/.env ./ 

CMD ["./user-service"]

EXPOSE 50051 9091

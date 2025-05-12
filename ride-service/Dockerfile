FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/ride-service

RUN go mod tidy
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o ride-service .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/ride-service/ride-service .
COPY --from=builder /app/ride-service/.env ./ 

CMD ["./ride-service"]

EXPOSE 50052 9092

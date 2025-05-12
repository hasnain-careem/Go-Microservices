FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy the entire project including common module
COPY . .

WORKDIR /app/booking-service

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o booking-service .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/booking-service/booking-service .
COPY --from=builder /app/booking-service/.env ./ 

CMD ["./booking-service"]

EXPOSE 50053 9093

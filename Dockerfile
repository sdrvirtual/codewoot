FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go build -o app ./cmd/codewoot

FROM golang:1.24
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]


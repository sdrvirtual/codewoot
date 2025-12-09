FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o app ./cmd/codewoot

FROM golang:1.24-alpine
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/internal/db/migrations /app/internal/db/migrations
RUN apk update && apk add --no-cache ffmpeg
CMD ["./app"]

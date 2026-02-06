# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git make
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/sachapel ./cmd/server

# Production stage
FROM alpine:latest AS production

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/sachapel .
COPY --from=builder /build/static ./static
COPY --from=builder /build/migrations ./migrations

RUN mkdir -p /app/storage/bulletins/morning \
             /app/storage/bulletins/evening \
             /app/storage/photos \
             /app/storage/documents \
             /app/storage/uploads

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 3000

CMD ["./sachapel"]

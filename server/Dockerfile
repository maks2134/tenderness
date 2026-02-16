FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/tenderness ./cmd/tenderness/
COPY internal ./internal/
COPY migrations ./migrations/

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /tenderness-app ./cmd/tenderness

FROM alpine:latest

RUN apk --no-cache add curl

WORKDIR /root/

COPY --from=builder /tenderness-app .
COPY --from=builder /app/migrations ./migrations/

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/health || exit 1

CMD ["./tenderness-app"]
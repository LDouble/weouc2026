FROM golang:1.25.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api-server ./cmd/api-server

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /out/api-server /usr/local/bin/api-server

ENV API_SERVER_ENV=docker \
    API_SERVER_HOST=0.0.0.0 \
    API_SERVER_PORT=8080

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
  CMD wget -q -O - http://127.0.0.1:8080/readyz | grep -q '"status":"ready"' || exit 1

ENTRYPOINT ["api-server"]

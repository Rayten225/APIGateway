# syntax=docker/dockerfile:1.4
FROM golang:1.23.4 AS builder

WORKDIR /app

# Копируем go.mod и vendor
COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY . .

ENV CGO_ENABLED=0
# Сборка с флагом -mod=vendor
RUN go build -mod=vendor -o app .

# Финальный минимальный образ
FROM gcr.io/distroless/static
COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]

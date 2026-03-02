
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/collector_bin ./collector/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/dashboard_bin ./dashboard/

FROM alpine:latest AS runtime_collector

WORKDIR /app

COPY --from=builder /app/collector_bin /app/collector
RUN chmod +x /app/collector

# COPY --from=builder /app/collector/sensors.json /app/sensors.json

COPY --from=builder /app/.env /app/.env

CMD ["./collector"]

FROM alpine:latest AS runtime_dashboard

WORKDIR /app/dashboard

COPY --from=builder /app/dashboard_bin ./dashboard_bin
RUN chmod +x ./dashboard_bin

COPY --from=builder /app/.env /app/dashboard/.env


CMD ["./dashboard_bin"]

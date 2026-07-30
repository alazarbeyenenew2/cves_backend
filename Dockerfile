FROM golang:1.25.6-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates

ENV GOPROXY=https://proxy.golang.org,direct

COPY . .

RUN go mod download
RUN go mod tidy

RUN go build -o cves ./cmd

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates


COPY --from=builder /app/pulsechannels_channel /app/pulsechannels_channel

COPY --from=builder /app/config/config.yml /app/config/config.yaml

RUN mkdir -p logs
CMD ["/app/cves"]

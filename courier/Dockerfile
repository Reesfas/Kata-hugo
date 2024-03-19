FROM golang:1.19.1-alpine3.15 AS builder
COPY geotask /go/src/geotask
WORKDIR /go/src/geotask
# Create slimest possible image
RUN go build -ldflags="-w -s" -o /go/bin/server /go/src/geotask/cmd/api

FROM alpine:3.15
# Copy binary from builder
COPY --from=builder /go/bin/server /go/bin/server
COPY public /app/public
COPY docs /app/docs
COPY .env /app/.env

WORKDIR /app
# Set entrypoint
ENTRYPOINT ["/go/bin/server"]

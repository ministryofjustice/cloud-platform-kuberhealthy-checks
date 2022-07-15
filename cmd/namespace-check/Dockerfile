FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.* /app/
RUN go mod download

COPY . /app

WORKDIR /app/cmd/namespace-check

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

RUN addgroup -S kuberhealthy && adduser -S kuberhealthy -G kuberhealthy

# Build a small image
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
# https://github.com/aws/aws-sdk-go/issues/2322
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/cmd/namespace-check/namespace-check /app/namespace-check

USER kuberhealthy

CMD ["/app/namespace-check"]

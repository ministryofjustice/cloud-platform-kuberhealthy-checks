FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.* /app/
RUN go mod download

COPY . /app

WORKDIR /app/cmd/namespace-check

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

RUN addgroup -S kuberhealthy && adduser -S kuberhealthy -G kuberhealthy
RUN mkdir -p /home/kuberhealthy/.kube && chown -R kuberhealthy:kuberhealthy /home/kuberhealthy/.kube
RUN touch /home/kuberhealthy/.kube/config && chown -R kuberhealthy:kuberhealthy /home/kuberhealthy/.kube/config

# Build a small image
FROM scratch
ENV \
  USER=kuberhealthy \
  GROUP=kuberhealthy \
  HOME=/home/kuberhealthy

COPY --from=builder /etc/passwd /etc/passwd
# https://github.com/aws/aws-sdk-go/issues/2322
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /app/cmd/namespace-check/namespace-check /app/namespace-check
COPY --chown=user:group --from=builder /app/cmd/namespace-check/namespace-check /app/namespace-check
COPY --chown=user:group --from=builder /home/kuberhealthy /home/kuberhealthy

USER kuberhealthy

CMD ["/app/namespace-check"]

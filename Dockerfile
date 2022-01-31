FROM golang:latest as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /app
COPY observer ./
RUN make

##
## Deploy
##
FROM scratch as autobot
COPY --from=builder /app/bin/autobot /ingoda/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/ingoda/autobot"]

FROM scratch as nats-sniffer
COPY --from=builder /app/bin/nats-sniffer /ingoda/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/ingoda/nats-sniffer"]

FROM ubuntu:latest as selfcheck-syslog
COPY --from=builder /app/bin/selfcheck-syslog /ingoda/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/ingoda/selfcheck-syslog"]

FROM ubuntu:latest as listmanager
COPY --from=builder /app/bin/listmanager /ingoda/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/ingoda/listmanager"]

FROM ubuntu:latest as receiver
COPY --from=builder /app/bin/receiver /ingoda/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/ingoda/receiver"]

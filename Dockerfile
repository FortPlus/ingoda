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
ENV FP_LOG_CONF=/ingoda/config.json 
COPY --from=builder /app/bin/autobot /ingoda/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/ingoda/autobot"]


FROM scratch as nats-sniffer
ENV FP_LOG_CONF=/ingoda/config.json
COPY --from=builder /app/bin/nats-sniffer /ingoda/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/ingoda/nats-sniffer"]



FROM ubuntu:latest as selfcheck-syslog
ENV FP_LOG_CONF=/ingoda/config.json
COPY --from=builder /app/bin/selfcheck-syslog /ingoda/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/ingoda/selfcheck-syslog"]



FROM ubuntu:latest as ban
ENV FP_LOG_CONF=/ingoda/config.json
COPY --from=builder /app/bin/ban /ingoda/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/ingoda/ban"]



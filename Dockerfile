FROM golang:latest as builder

WORKDIR /app

COPY ./observer ./
RUN make


##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /bin/ /ingoda

EXPOSE 8030

USER nonroot:nonroot

ENTRYPOINT ["/ingoda"]


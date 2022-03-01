FROM golang:1.17 as builder
WORKDIR /go/src/product-review
COPY . .
RUN make build


FROM debian:latest
WORKDIR /root/
COPY --from=builder /go/src/product-review/build/product-review /usr/bin/product-review

# Copy CA certificates to prevent x509: certificate signed by unknown authority errors
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 80
CMD ["product-review", "serve"]
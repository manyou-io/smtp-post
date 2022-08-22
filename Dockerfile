FROM golang:1.19.0-alpine3.16 as builder

WORKDIR /usr/src/app
COPY . ./

RUN set -ex; \
    apk update; \
    apk add upx git; \
    go mod download; \
    go mod verify; \
    go build -ldflags='-s -w' -o smtp-post; \
    upx --brute smtp-post

FROM scratch

ENV SMTP_POST_DOMAIN="smtp-post" \
    SMTP_POST_BIND=":587" \
    SMTP_POST_READ_TIMEOUT="10" \
    SMTP_POST_WRITE_TIMEOUT="10" \
    SMTP_POST_MAX_RCPT="14" \
    SMTP_POST_MAX_SIZE="5242880" \
    SMTP_POST_API_KEY="smtp-post" \
    SMTP_POST_USERNAME="smtp-post" \
    SMTP_POST_PASSWORD="smtp-post" \
    SMTP_POST_ENDPOINT="" \
    SMTP_POST_TLS_CERT="" \
    SMTP_POST_TLS_KEY=""

COPY --from=builder /usr/src/app/smtp-post /usr/bin/

CMD ["/usr/bin/smtp-post", "--domain", "$SMTP_POST_DOMAIN", "--bind", "$SMTP_POST_BIND", "--read-timeout", "$SMTP_POST_READ_TIMEOUT", "--write-timeout", "$SMTP_POST_WRITE_TIMEOUT", "--max-rcpt", "$SMTP_POST_MAX_RCPT", "--max-size", "$SMTP_POST_MAX_SIZE", "--api-key", "$SMTP_POST_API_KEY", "--username", "$SMTP_POST_USERNAME", "--password", "$SMTP_POST_PASSWORD", "--endpoint", "$SMTP_POST_ENDPOINT", "--tls-cert", "$SMTP_POST_TLS_CERT", "--tls-key", "$SMTP_POST_TLS_KEY"]

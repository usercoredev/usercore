FROM golang:1.22.1-alpine AS builder
LABEL org.opencontainers.image.authors="Talut TASGIRAN <talut@tasgiran.com>"
WORKDIR /go/src/github.com/usercoredev/usercore
COPY . ./
ENV CGO_ENABLED=0 GOOS=linux
RUN apk update && apk add --no-cache git && \
    go mod download && \
    go build -ldflags="-s -w" -o usercore && \
    rm -rf /var/cache/apk/*

FROM alpine:latest
RUN apk --no-cache update && apk add --no-cache \
    ca-certificates \
    tzdata && \
    update-ca-certificates && \
    adduser -D -u 1000 usercore
COPY --from=builder --chown=usercore:usercore /go/src/github.com/usercoredev/usercore/usercore /app/usercore
USER usercore
CMD ["/app/usercore"]


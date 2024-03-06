FROM golang:alpine AS builder
LABEL maintainer="Talut TASGIRAN <talut@tasgiran.com>"
RUN apk update && apk add --no-cache git && rm -rf /var/cache/apk/*
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ENV CGO_ENABLED=0
ENV GOOS=linux
COPY . ./
RUN go build -o /usercore

FROM alpine:latest
RUN apk --no-cache update && apk add --no-cache \
    ca-certificates \
    tzdata \
    && update-ca-certificates
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
WORKDIR /
COPY --from=builder /usercore .
RUN chown appuser:appuser ./usercore && chmod +x ./usercore
USER appuser
CMD ["/usercore"]

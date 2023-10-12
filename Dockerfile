# Builder
FROM golang:1.20.6-alpine AS builder

WORKDIR /usr/local/src

RUN apk update && apk upgrade && \
    apk --update --no-cache add  make bash build-base

COPY . .

RUN go mod download
RUN make build

# Distribution
FROM alpine:latest

COPY --from=builder /usr/local/src/app /
COPY --from=builder /usr/local/src/configs/server.json configs/server.json
COPY --from=builder /usr/local/src/migrations/create.sql migrations/create.sql

CMD ["/app"]
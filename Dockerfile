FROM golang:1.24-alpine AS builder

RUN apk update && apk add --no-cache \
    git \
    make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata gettext

WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/migrate . 

COPY internal/config/config.yml.template ./internal/config/config.yml.template

COPY entrypoint.sh .

RUN chmod +x ./entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
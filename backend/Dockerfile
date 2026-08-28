FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Миграции вшиты в бинарник через embed — отдельно копировать не нужно.
RUN CGO_ENABLED=0 go build -o server ./cmd/api
RUN CGO_ENABLED=0 go build -o createadmin ./cmd/createadmin
RUN CGO_ENABLED=0 go build -o seedcourse ./cmd/seedcourse

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates \
    && adduser -D -u 10001 app

COPY --from=builder /app/server /app/createadmin /app/seedcourse ./

USER app

EXPOSE 8090

CMD ["./server"]

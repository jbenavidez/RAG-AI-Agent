FROM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP_CMD=./cmd/web
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/app ${APP_CMD}


FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/app /app/app

COPY static /app/static
COPY templates /app/templates

EXPOSE 8080

CMD ["/app/app"]
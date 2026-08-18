# Etapa 1: compilar la API utilizando Go.
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod .
COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /api1 .

# Etapa 2: imagen final mínima.
FROM scratch

COPY --from=builder /api1 /api1

EXPOSE 8080

ENTRYPOINT ["/api1"]
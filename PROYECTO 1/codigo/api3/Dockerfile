FROM docker.io/library/golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod .
COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /api3 .

FROM scratch

COPY --from=builder /api3 /api3

EXPOSE 8080

ENTRYPOINT ["/api3"]
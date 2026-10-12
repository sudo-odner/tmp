FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /bin/server /app/server

EXPOSE 5008/udp

CMD ["/app/server"]

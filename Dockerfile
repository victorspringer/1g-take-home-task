FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server

FROM alpine:3.20.3

ENV HTTP_PORT=8080
ENV GIN_MODE=release
ENV POSTGRES_CONN=postgres://user:pass@localhost:5432/challenge?sslmode=disable

COPY --from=builder /app/app /app

EXPOSE 8080

CMD ["/app"]

## step #1
FROM golang:1.19 AS builder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /app

ADD go.mod .
ADD go.sum .

# COPY . .
COPY /cmd/gophermart/main.go ./cmd/gophermart/main.go
COPY /internal/. ./internal/.
COPY /pkg/. ./pkg/.

RUN go mod download
RUN go build -o ./cmd/gophermart/gophermart ./cmd/gophermart/

## step #2
FROM alpine AS product

WORKDIR /app

COPY --from=builder /app/cmd/gophermart/gophermart ./cmd/gophermart/gophermart
COPY /migrations/. ./migrations/.

EXPOSE 8080

CMD ["./cmd/gophermart/gophermart"]


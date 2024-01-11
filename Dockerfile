# Builder image
FROM golang:alpine AS builder

WORKDIR /build
COPY . /build

RUN go build -o go-serverless cmd/main.go

# Final image
FROM alpine

WORKDIR /app
COPY --from=builder /build/go-serverless /app/go-serverless

ENTRYPOINT /app/go-serverless

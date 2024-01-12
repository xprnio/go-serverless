# --- Builder image
FROM golang:alpine AS builder
WORKDIR /build

COPY go.mod /build
RUN go mod download

COPY . /build
RUN go build -o go-serverless cmd/main.go

# --- Final image
FROM scratch
WORKDIR /

# Letting the application know we're running in docker
ENV OS_ENV=docker

# Application data path
ENV DATA_PATH=/data

# Information about the 
ENV CONTEXT_NAME=serverless_context
ENV CONTEXT_PATH=/tmp/context

COPY --from=builder /build/go-serverless /go-serverless
ENTRYPOINT ["/go-serverless"]

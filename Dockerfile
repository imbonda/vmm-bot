# First stage: Build the Go application
FROM --platform=${BUILDPLATFORM} golang:1.22-alpine3.19 AS builder

ARG BUILDPLATFORM
ARG TARGETARCH
ARG TARGETOS
ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux

WORKDIR /go/src/app/cmd

# Copy the necessary stuff only
COPY ./pkg ../pkg
COPY ./internal ../internal

COPY go.mod go.sum ../

# Download dependencies
RUN go mod download

# Copy the service code
COPY ./cmd/config config
COPY ./cmd/interfaces interfaces
COPY ./cmd/service service
COPY ./cmd/main.go .

# Build the Go application
RUN GOARCH=${TARGETARCH} GOOS=${TARGETOS} go build -o trader.${TARGETARCH} .

# Second stage: Create the final Alpine based image
FROM alpine:3
ARG TARGETARCH
RUN apk add --no-cache bash
# Copy the binary from the builder stage
COPY --from=builder /go/src/app/cmd/trader.${TARGETARCH} /app/trader

# Expose the port your application will run on
EXPOSE 8000
#
WORKDIR /app/
# Run the application
ENTRYPOINT ["/app/trader"]

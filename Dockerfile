# syntax=docker/dockerfile:1.4

#──────────────────────────────────────────────────────────────
# Stage 1: Builder
#──────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /api

# Copy dependency manifests first for optimal layer caching
COPY go.mod go.sum ./

# Download dependencies (cached separately from source code)
RUN go mod download

# Copy full source code
COPY . .

# Compile the binary
RUN mkdir -p ./tmp && go build -o ./tmp/main ./app/...

#──────────────────────────────────────────────────────────────
# Stage 2: Runtime
#──────────────────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /api

# Copy only the compiled binary from the builder stage
COPY --from=builder /api/tmp/main ./tmp/main

EXPOSE 10000

CMD ["./tmp/main"]

# syntax=docker/dockerfile:1.4

#──────────────────────────────────────────────────────────────
# Stage 1: Builder
#──────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency manifests first for optimal layer caching
COPY go.mod go.sum ./

# Download dependencies (cached separately from source code)
RUN go mod download

# Copy full source code
COPY . .

# Compile the binary
RUN go build -o helptrix-api ./app/...

#──────────────────────────────────────────────────────────────
# Stage 2: Runtime
#──────────────────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/helptrix-api .

EXPOSE 8080

CMD ["./helptrix-api"]

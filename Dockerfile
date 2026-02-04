# Build Stage
FROM golang:1.25-alpine AS builder

# Install system tools
RUN apk add --no-cache git

# 1. Install Main Buffalo CLI
RUN go install github.com/gobuffalo/cli/cmd/buffalo@latest

# ---> FIX: Install Soda (The direct database tool) <---
RUN go install github.com/gobuffalo/pop/v6/soda@latest

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the app
RUN go build -o bin/app .

# Final Stage
FROM alpine:latest
WORKDIR /root/

RUN apk add --no-cache ca-certificates

# Copy app binary and config
COPY --from=builder /app/bin/app .
COPY --from=builder /app/public ./public
COPY --from=builder /app/database.yml .

# Copy the tools
COPY --from=builder /go/bin/buffalo /usr/local/bin/buffalo
# ---> Copy Soda <---
COPY --from=builder /go/bin/soda /usr/local/bin/soda

EXPOSE 3000
CMD ["./app"]
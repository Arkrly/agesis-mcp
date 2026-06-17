# Step 1: Build Go Backend & Mock MCP
FROM golang:1.25-alpine AS go-builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o aegis-mcp ./cmd/aegis-mcp/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o mock-mcp ./test/mock-mcp/main.go

# Step 2: Build React Frontend
FROM node:20-alpine AS node-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN echo "VITE_API_URL=" > .env
RUN npm run build

# Step 3: Final Production Image
FROM alpine:3.21
WORKDIR /app

# Copy certificates and timezone data from builder to avoid network dependency in final stage
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy Go binaries
COPY --from=go-builder /app/aegis-mcp .
COPY --from=go-builder /app/mock-mcp .

# Copy Static Frontend
COPY --from=node-builder /app/dist ./ui

# Copy Default Configs
COPY config/ ./config/

# Create an entrypoint script to run both services
RUN echo '#!/bin/sh' > entrypoint.sh && \
    echo './mock-mcp & # Start mock server in background' >> entrypoint.sh && \
    echo './aegis-mcp # Start main gateway in foreground' >> entrypoint.sh && \
    chmod +x entrypoint.sh

# Environment Defaults
ENV AEGIS_LISTEN_ADDR=:8080
ENV AEGIS_UPSTREAM_URL=http://localhost:9090/mcp
ENV AEGIS_POLICY_FILE=config/policy.rego

EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]

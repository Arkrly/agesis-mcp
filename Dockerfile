# Step 1: Build Go Backend & Mock MCP
FROM golang:1.24-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o aegis-mcp ./cmd/aegis-mcp
RUN CGO_ENABLED=0 GOOS=linux go build -o mock-mcp ./test/mock-mcp/main.go

# Step 2: Build React Frontend
FROM node:20-alpine AS node-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
# Set API URL to empty so it uses the same host in production
RUN echo "VITE_API_URL=" > .env
RUN npm run build

# Step 3: Final Production Image
FROM alpine:latest
WORKDIR /app

# Install dependencies for the entrypoint script
RUN apk add --no-cache ca-certificates

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

# Environment Defaults for Railway
ENV AEGIS_LISTEN_ADDR=:8080
ENV AEGIS_UPSTREAM_URL=http://localhost:9090/mcp
ENV AEGIS_JWT_SHARED_SECRET=hackathon-demo-secret-key-12345
ENV AEGIS_POLICY_FILE=config/policy.rego

EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]

# Multi-stage Dockerfile for PageMail
# Stage 1: Build the Go backend
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

# Install dependencies for building
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o pagemail ./cmd/pagemail

# Stage 2: Build the Next.js frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Copy package files
COPY frontend/package*.json ./
RUN npm ci

# Copy source code and build
COPY frontend/ .
RUN npm run build

# Stage 3: Final runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates chromium wget

# Create app user
RUN addgroup -g 1001 -S pagemail && \
    adduser -S pagemail -u 1001

# Create directories
RUN mkdir -p /app/files /app/frontend

# Copy backend binary
COPY --from=backend-builder /app/pagemail /app/pagemail
RUN chmod +x /app/pagemail

# Copy frontend static export
COPY --from=frontend-builder /app/dist /app/frontend/dist

# Set ownership
RUN chown -R pagemail:pagemail /app

# Switch to non-root user
USER pagemail

WORKDIR /app

# Environment variables
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./pagemail"]
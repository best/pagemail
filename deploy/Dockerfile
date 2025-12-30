# ==============================================================================
# Stage 1: Build Frontend
# ==============================================================================
FROM node:20-alpine AS frontend-builder

# Use China mirrors
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN npm config set registry https://registry.npmmirror.com

WORKDIR /app/web

COPY web/package*.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

# ==============================================================================
# Stage 2: Build Backend
# ==============================================================================
FROM golang:1.23-alpine AS backend-builder

# Use China mirrors
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache git ca-certificates tzdata

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG BUILD_TIME

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -s -w" \
    -o /app/bin/pagemail ./cmd/pagemail

# ==============================================================================
# Stage 3: Production Image
# ==============================================================================
FROM alpine:3.20

# Use China mirrors
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    nginx \
    supervisor \
    chromium \
    nss \
    freetype \
    harfbuzz \
    fontconfig \
    font-noto \
    font-noto-cjk \
    font-noto-emoji

# Create non-root user
RUN addgroup -g 1000 pagemail && \
    adduser -u 1000 -G pagemail -s /bin/sh -D pagemail

WORKDIR /app

# Copy built artifacts
COPY --from=backend-builder /app/bin/pagemail /app/pagemail
COPY --from=frontend-builder /app/web/dist /app/web/dist

# Copy configuration files
COPY deploy/nginx.conf /etc/nginx/nginx.conf
COPY deploy/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Create necessary directories including crashpad
RUN mkdir -p /app/storage /var/log/nginx /var/run/nginx /var/log/supervisor /tmp/crashpad && \
    chown -R pagemail:pagemail /app /var/log/nginx /var/run/nginx /var/log/supervisor /tmp/crashpad

# Set Chrome path and disable crashpad for rod
ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROME_CRASHPAD_PIPE_NAME=
ENV BREAKPAD_DUMP_LOCATION=/tmp/crashpad

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/v1/health || exit 1

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]

# ==============================================
# IAM Service - Optimized Multi-stage Dockerfile
# ==============================================

# ==============================================
# Stage 1: Dependencies and cache optimization
# ==============================================
FROM golang:1.24-alpine AS deps
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ==============================================
# Stage 2: Build stage
# ==============================================
FROM deps AS builder

# Copy source code
COPY . .

# Build optimized binary with security hardening
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -trimpath \
    -o iam-service ./src/main.go

# Verify binary (skip on ARM64 Mac)
# RUN file iam-service && \
#     ldd iam-service 2>&1 | grep -q "not a dynamic executable" || exit 1

# ==============================================
# Stage 3: Development stage
# ==============================================
FROM alpine:3.19 AS development

# Security: Create non-root user first
RUN addgroup -g 1001 -S appgroup && \
    adduser -S -D -h /app -s /bin/sh -G appgroup -u 1001 appuser

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    postgresql-client \
    && cp /usr/share/zoneinfo/UTC /etc/localtime \
    && echo "UTC" > /etc/timezone \
    && apk del tzdata

WORKDIR /app

# Copy binary and resources
COPY --from=builder --chown=appuser:appgroup /app/iam-service ./
COPY --from=builder --chown=appuser:appgroup /app/migrations ./migrations/
COPY --from=builder --chown=appuser:appgroup /app/scripts ./scripts/

# Make binary executable
RUN chmod +x ./iam-service

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

EXPOSE 8080

CMD ["./iam-service"]

# ==============================================
# Stage 4: Production stage (Distroless)
# ==============================================
FROM gcr.io/distroless/static-debian12:nonroot AS production

# Metadata
LABEL org.opencontainers.image.title="IAM Service" \
      org.opencontainers.image.description="Multi-tenant Identity & Access Management" \
      org.opencontainers.image.source="https://github.com/saas-mt/iam-service" \
      org.opencontainers.image.vendor="SaaS MT Team" \
      org.opencontainers.image.licenses="MIT"

WORKDIR /app

# Copy binary and essential files only
COPY --from=builder --chown=nonroot:nonroot /app/iam-service ./
COPY --from=builder --chown=nonroot:nonroot /app/migrations ./migrations/

# Use distroless nonroot user (uid=65532)
USER nonroot

EXPOSE 8080

ENTRYPOINT ["./iam-service"]

# ==============================================
# Default stage: Development
# ==============================================
FROM development

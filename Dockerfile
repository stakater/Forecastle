# Stage 1: Build frontend
FROM node:24 AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/yarn.lock ./
RUN yarn install --frozen-lockfile --network-timeout 1000000
COPY frontend/ ./
RUN yarn build

# Stage 2: Build Go binary
FROM --platform=${BUILDPLATFORM} golang:1.25 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy frontend build to internal/web/build for embedding
COPY --from=frontend-builder /app/frontend/build ./internal/web/build/

# Build the binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /forecastle ./cmd/forecastle

# Stage 3: Runtime image
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /forecastle .
USER nonroot:nonroot

ENTRYPOINT ["/forecastle"]

# Build stage
FROM golang:1.22-bookworm AS builder

# Install dependencies for Fyne
RUN apt-get update && apt-get install -y \
    gcc \
    libgl1-mesa-dev \
    xorg-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy everything
COPY . .

# Tidy, download and build
RUN go mod tidy && go build -o bitwarden-html-converter -ldflags="-s -w" .

# Runtime stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libgl1 \
    libx11-6 \
    libxcursor1 \
    libxrandr2 \
    libxinerama1 \
    libxi6 \
    libxxf86vm1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/bitwarden-html-converter .

ENTRYPOINT ["/app/bitwarden-html-converter"]

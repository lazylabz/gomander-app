FROM ubuntu:22.04

# Avoid interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Set working directory
WORKDIR /app

# Update package list and install basic dependencies
RUN apt-get update
RUN apt-get install -y curl wget git ca-certificates software-properties-common
RUN apt-get install -y build-essential pkg-config gcc g++ libc6-dev
RUN apt-get install -y libgtk-3-dev libwebkit2gtk-4.1-dev
RUN apt-get install -y libnss3-dev libxss1 libasound2-dev xvfb
RUN rm -rf /var/lib/apt/lists/*

# Install Node.js 18.x (LTS)
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs

# Detect architecture and install appropriate Go version
RUN ARCH=$(dpkg --print-architecture) && \
    case ${ARCH} in \
        amd64) GOARCH=amd64 ;; \
        arm64) GOARCH=arm64 ;; \
        armhf) GOARCH=armv6l ;; \
        *) echo "Unsupported architecture: ${ARCH}" && exit 1 ;; \
    esac && \
    wget https://go.dev/dl/go1.24.5.linux-${GOARCH}.tar.gz && \
    tar -C /usr/local -xzf go1.24.5.linux-${GOARCH}.tar.gz && \
    rm go1.24.5.linux-${GOARCH}.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOBIN="/go/bin"
ENV PATH="${GOBIN}:${PATH}"
ENV CGO_ENABLED=1
ENV GOOS=linux

# Create Go workspace
RUN mkdir -p /go/src /go/bin /go/pkg

# Install Wails CLI
RUN CGO_ENABLED=1 go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Install additional npm packages commonly used with Wails
RUN npm install -g pnpm@latest

# Verify installations and run wails doctor
RUN echo "=== Verifying Go installation ===" && \
    go version && \
    echo "\n=== Verifying Node.js installation ===" && \
    node --version && \
    npm --version && \
    echo "\n=== Verifying Wails installation ===" && \
    wails version && \
    echo "\n=== Running Wails Doctor ===" && \
    wails doctor || true

# Copy all files
COPY . .

COPY docker/docker-build-entrypoint.sh /usr/local/bin/docker-build-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-build-entrypoint.sh

ENTRYPOINT ["docker-build-entrypoint.sh"]

# Labels for documentation
LABEL description="Ubuntu-based Docker image for Wails application development"
LABEL version="1.0"
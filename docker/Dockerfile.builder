FROM golang:latest AS builder
WORKDIR /app

# Copy go.mod & go.sum first to cache dependencies
COPY go.mod go.sum ./

# Set up git + go env configs
RUN git config --global http.sslVerify false && \
    go env -w GOINSECURE="github.com,go.googlesource.com,golang.org,go.uber.org,google.golang.org,*.org" && \
    go env -w GOSUMDB=off && \
    go env -w GOPROXY=direct

# GODEBUG chỉ được phép set bằng ENV
ENV GODEBUG=x509ignoreCN=1

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Cache dependencies using Docker BuildKit
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy source code
COPY . .

CMD ["tail", "-f", "/dev/null"]


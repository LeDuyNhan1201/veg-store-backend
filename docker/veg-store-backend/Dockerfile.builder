FROM golang:latest AS builder
WORKDIR /app

# Step 1: copy go.mod (cache deps)
COPY go.mod go.sum ./

# Step 2: copy code
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    INSIDE_DOCKER=1 make force-download && go mod download

# Step 3: run make
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    INSIDE_DOCKER=1 make force-download && make prepare

CMD ["tail", "-f", "/dev/null"]


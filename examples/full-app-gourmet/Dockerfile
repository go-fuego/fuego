FROM alpine:latest AS css-builder
WORKDIR /app
RUN apk add --no-cache curl libstdc++ libgcc
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then TAILWIND_ARCH="x64"; \
    elif [ "$ARCH" = "aarch64" ]; then TAILWIND_ARCH="arm64"; \
    else echo "Unsupported architecture: $ARCH" && exit 1; fi && \
    curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-${TAILWIND_ARCH}-musl && \
    chmod +x tailwindcss-linux-${TAILWIND_ARCH}-musl && \
    mv tailwindcss-linux-${TAILWIND_ARCH}-musl /usr/local/bin/tailwindcss
COPY examples/full-app-gourmet/ .
RUN tailwindcss -i ./tailwind.css -o ./static/tailwind.min.css --minify

FROM golang:latest AS builder
WORKDIR /app
COPY examples/full-app-gourmet/go.mod examples/full-app-gourmet/go.sum ./
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/examples/full-app-gourmet
COPY --from=css-builder /app/static/tailwind.min.css /app/examples/full-app-gourmet/static/tailwind.min.css
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o gourmet .

FROM alpine:latest
WORKDIR /app
RUN mkdir -p data
COPY --from=builder /app/examples/full-app-gourmet/gourmet .
ENTRYPOINT [ "/app/gourmet" ]

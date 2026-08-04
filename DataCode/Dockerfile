# syntax=docker/dockerfile:1
FROM golang:1.25-bookworm AS build
WORKDIR /src
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates git && rm -rf /var/lib/apt/lists/*
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/parser ./cmd/parser

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates curl libreoffice-writer-nogui poppler-utils \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /out/parser /app/parser
COPY data /app/data
VOLUME ["/app/result"]
ENTRYPOINT ["/app/parser"]

# syntax=docker/dockerfile:1

FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY DataCode /src/DataCode
COPY platform/api /src/platform/api
WORKDIR /src/platform/api
RUN go mod download && CGO_ENABLED=0 go build -o /out/api ./cmd/api

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    poppler-utils \
    tesseract-ocr \
    tesseract-ocr-rus \
    libreoffice-writer-nogui \
    libreoffice-calc-nogui \
    unar \
    fonts-dejavu-core \
    fonts-liberation \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /out/api /usr/local/bin/api
ENV HTTP_ADDR=:8080
EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=5s --retries=5 CMD curl -fsS http://127.0.0.1:8080/api/v1/health || exit 1
CMD ["api"]

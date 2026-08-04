#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"
docker compose up -d --build "$@"
echo "UI  http://localhost:3000"
echo "API http://localhost:8080/api/v1/health"
echo "DB  postgres://zakupki:zakupki@localhost:5432/zakupki"

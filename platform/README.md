# Zakupki Platform

Каталог закупок: CSV-ingest → парсинг ЕИС → PostgreSQL → UI.  
Опционально: AI-анализ через [analizator_zakupok](https://github.com/rinat1313/analizator_zakupok) + LM Studio.

## Быстрый старт

```bash
./up.sh
# UI  http://localhost:3000
# API http://localhost:8080/api/v1/health
```

С AI (нужен соседний клон `analizator_zakupok` и LM Studio на `:1234`):

```bash
export LM_STUDIO_MODEL=<id>
docker compose --profile ai up -d --build
```

## Документация

- [API.md](docs/API.md) — эндпоинты
- [HANDOFF.md](docs/HANDOFF.md) — анализ работы агентов и бэклог

## Разработка API

```bash
cd api
export DATABASE_URL=postgres://zakupki:zakupki@127.0.0.1:5432/zakupki?sslmode=disable
export ANALIZATOR_URL=http://127.0.0.1:8088   # опционально
go test ./...
go run ./cmd/api
```

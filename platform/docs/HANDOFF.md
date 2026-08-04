# Продолжение работы агентов (handoff)

Дата: 2026-08-04  
Контекст: продолжение после `bc-07297476-…` (транскрипт из этой среды недоступен) и артефактов `bc-95d3d0cc-…`.

## Что уже сделано предыдущими агентами

### 1. `analizator_zakupok` (агент `bc-95d3d0cc`)
- Go-микросервис: chunking + чек-листы YAML + LM Studio (`/v1/chat/completions`).
- Отдельный репозиторий: https://github.com/rinat1313/analizator_zakupok
- Интеграция с файловым парсером: [ZakupkiParser PR #1](https://github.com/rinat1313/ZakupkiParser/pull/1) (`-analyze-url`, корневой compose).
- Документация интегратора: `docs/INTEGRATOR.md`.

### 2. Platform MVP (локальные изменения → cloud PR #2)
Ветка `cursor/cloud-agent-1785816804000-d75mp` / [PR #2](https://github.com/rinat1313/ZakupkiParser/pull/2):
- `platform/` — Postgres + Go API + Bootstrap UI.
- CSV-ingest → фоновый worker → карточка ЕИС через новый `DataCode/pkg/collect`.
- Каталог тендеров, статусы `analysis_status`, ручная оценка (`tender_assessments`).
- Адаптеры ЭТП (tektorg и др.) — **заглушки**.
- Courts / РНП — **заглушки**.
- Связи с `analizator_zakupok` **не было**.

### 3. Stub-репозитории под будущий сплит
Пустые README: `zakupki-core`, `zakupki-customer`, `zakupki-gateway`, `zakupki-parser`, `zakupki-platform`.  
Код пока монорепо-прототип в `ZakupkiParser`.

### 4. Full-stack launcher в `analizator_zakupok`
Кратко добавляли скрипты `up/parse/down` и `docs/STACK.md`, затем **откатили**: полный стек остаётся в PR #1 парсера.

## Что сделано в этом продолжении

1. HTTP-мост `ANALIZATOR_URL` → `POST /api/v1/tenders/{id}/analyze`.
2. Сбор корпуса из `object_name` + `documents.text_content` → вызов analizator с `text`.
3. Результат пишется в `tender_assessments` (`details.source=analizator_zakupok`), статус → `analyzed`.
4. Health API показывает `analizator: ok|unavailable|disabled`.
5. UI: кнопка **AI-анализ** в карточке тендера.
6. Compose profile `ai` для соседнего клона `analizator_zakupok`.
7. Юнит-тесты клиента моста.

## Как поднять с AI

```bash
# рядом:
#   ZakupkiParser/
#   analizator_zakupok/

cd ZakupkiParser/platform
export LM_STUDIO_MODEL=<id-из-LM-Studio>
docker compose --profile ai up -d --build
# UI http://localhost:3000
# В карточке тендера → «AI-анализ»
```

Без Docker AI: `export ANALIZATOR_URL=http://127.0.0.1:8088` при запуске API.

## Следующие шаги (бэклог)

1. Реальный Adapter хотя бы для одного ЭТП (tektorg).
2. Авто-анализ после успешного ingest (флаг категории / job).
3. Согласовать merge PR #1 и PR #2 (не потерять `-analyze-url` и `pkg/collect`).
4. Тесты ingest/upsert; пагинация list; auth.
5. Вынос в stub-репы, когда границы стабильны.
6. Courts / РНП — внешние источники.

## Карта потоков

```
CSV → platform ingest → EIS collect → Postgres documents.text_content
                                         ↓
                              POST …/tenders/{id}/analyze
                                         ↓
                              analizator_zakupok + LM Studio
                                         ↓
                              tender_assessments + analysis_status=analyzed
```

Параллельный (файловый) поток PR #1:

```
CLI parser → result/{reg}/valid_* → analizator (reg_number) → analysis/analysis.json
```

# Patch: auto-AI freeze (zakupki-core)

## Symptom

При включённом `auto_ai` анализатор обрабатывает несколько карточек и «замирает»,
хотя в каталоге ещё есть готовые к анализу (`analysis_status=none`, есть `text_content`).

## Root cause (zakupki-core)

Воркер `internal/autoai` занимает слот через `Control.BeginAnalyze` **до** вызова
`AnalyzeTender`, а `defer release()` внутри `AnalyzeTender` регистрировался **после**
`ListDocuments` / проверки корпуса. Early-return → слот не освобождается →
`FreeAnalyzeSlots() == 0` навсегда → новые карточки не берутся.

Дополнительно: ошибки LLM раньше ставили `analysis_status=other`, и auto-очередь
(`NextTenderReadyForAI` смотрит только `none`) больше их не видела — при том что UI
считает `other` eligible.

## Apply

Репозиторий `zakupki-core` (на `main`):

```bash
cd /path/to/zakupki-core
git apply /path/to/this/0001-Fix-auto-AI-freeze-from-analyze-slot-leak.patch
go test ./internal/autoai/ ./internal/control/ ./internal/store/
```

Либо скопировать файлы из этой папки:

| Файл здесь | Куда в zakupki-core |
|---|---|
| `worker.go` | `internal/autoai/worker.go` |
| `worker_test.go` | `internal/autoai/worker_test.go` |
| `control_test.go` | `internal/control/control_test.go` |
| `tenders_requeue.go.fragment` | вставить/заменить блок requeue в `internal/store/tenders.go` |

## What changes

1. `defer release()` сразу при pre-held слоте (до любых early-return).
2. Early/LLM ошибки → `none` (карточка снова в auto-очереди); пустой корпус → `other`.
3. Раз в ~1 мин: `analyzing` старше 15м → `none`; `other` старше 15м с текстом → `none`.

## Manual check (after deploy core)

1. Включить `auto_ai` на категории с несколькими готовыми карточками — в логах core
   должны идти подряд `auto-ai: start …` / `auto-ai: done …` до опустошения очереди.
2. `GET /api/v1/workers`: `analyze_running` не залипает на `analyze_capacity`, пока
   есть ready-карточки.
3. Симулировать краткий сбой LLM / рестарт analizator — карточки уходят в `none` и
   продолжают обрабатываться **без** рестарта core.

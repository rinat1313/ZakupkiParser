# Код разнесён по сервисам

Монолит `platform/` + `DataCode/` перенесён в:

| Было | Стало |
|------|--------|
| `platform/` compose/UI/API | `zakupki-platform` + `zakupki-gateway` + `zakupki-core` |
| `DataCode/` EIS parser | `zakupki-parser` |
| courts/rnp stubs | `zakupki-customer` |
| AI | `analizator_zakupok` |

Запуск: https://github.com/rinat1313/zakupki-platform

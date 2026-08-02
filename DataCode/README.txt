DataCode — заготовка под разработку парсера на Go
==================================================
Создано по TASK.txt (2026-08-02).
Полные руководства по парсингу — в ../result/
Образцы HTML — в ../html/ (не изменять).


СТРУКТУРА
---------
DataCode/
  README.txt              — этот файл
  AI_AGENT.txt            — краткая инструкция для AI-агента при написании Go-кода
  go.mod                  — модуль (заготовка)
  models/                 — структуры данных для парсинга / БД
    notice44.go
    notice223.go
    organization.go
    document.go
    common.go
  html_map/               — где в HTML лежат нужные поля
    map_44.txt
    map_223.txt
  docs/
    structs_usage.txt     — как работать со структурами
  internal/parser/        — сюда класть реализацию парсеров (пока пусто)
  pkg/eis/                — общий HTTP/клиент ЕИС (пока пусто)
  cmd/parser/             — entrypoint CLI (пока пусто)
  scripts/eis_curl/       — тестовые curl к ЕИС (запуск из IDE)


С ЧЕГО НАЧАТЬ РАЗРАБОТКУ
------------------------
0. Проверить доступ к ЕИС:
     bash DataCode/scripts/eis_curl/run_all.sh
   или один запрос: bash DataCode/scripts/eis_curl/02_notice44.sh
1. Прочитать AI_AGENT.txt и result/13_44_vs_223_razlichiya.txt
2. Реализовать internal/parser/parser44 и parser223 отдельно
3. Маппить HTML → models/* → СУБД
4. Фикстуры брать из html/44 фз/ и html/223 фз/


НЕ ДЕЛАТЬ
---------
- Один парсер на 44 и 223 с общими CSS-селекторами
- Менять файлы в html/ и over/
- Хардкодить только ea20 — брать href из поиска

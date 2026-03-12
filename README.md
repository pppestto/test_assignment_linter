# test_assignment_linter

Линтер для проверки лог-сообщений в **log/slog** и **go.uber.org/zap**: строчная первая буква, только ASCII (английский), без спецсимволов/эмодзи (только буквы, цифры, пробелы), без утечки чувствительных данных через конкатенацию с переменной.

## Сборка и запуск

```bash
go build -o loglint ./cmd/loglint
./loglint ./...

go test ./... -count=1
```

С конфигом (бонусы): `./loglint -config=loglint.example.json ./...`  
С автофиксом: `./loglint -fix ./...`

## Примеры (соответствие ТЗ)

| Правило | Неправильно | Правильно |
|--------|-------------|-----------|
| Строчная буква | `slog.Error("Failed...")` | `slog.Error("failed...")` |
| Английский | `slog.Info("запуск...")` | `slog.Info("starting...")` |
| Без спецсимволов | `"failed!!!"`, `"bad!"` | `"connection failed"`, буквы/пробелы только |
| Чувствительные | `"token: " + token` | `"token validated"` |

Интеграция с **golangci-lint** — см. ниже и `docs/GOLANGCI-LINT.md`.

---

# Интеграция с golangci-lint (этап 4)

Линтер реализован как `go/analysis.Analyzer` и подключается к golangci-lint через **Module Plugin System** (рекомендуемый способ для приватных линтеров).

## Что нужно

- Установленный `golangci-lint` (команда `golangci-lint custom` входит в стандартную поставку).
- Git, Go.

## Установка (automatic way)

### 1. Файл `.custom-gcl.yml`

В корне **этого** репозитория уже лежит `.custom-gcl.yml`: в нём указан локальный модуль и импорт пакета `linters`, где в `init()` вызывается `register.Plugin("loglint", ...)`.

Если клонируешь линтер в другое место — поправь `path:` на каталог с `go.mod` линтера.

Версия `version:` в `.custom-gcl.yml` должна совпадать с версией golangci-lint, которой собираешь custom-бинарник. Проверить свою версию:

```bash
golangci-lint version
```

При необходимости замени `version: v2.1.0` в `.custom-gcl.yml` на свою.

### 2. Сборка custom-бинарника

Из корня репозитория линтера:

```bash
golangci-lint custom
```

По умолчанию появится `./custom-gcl` (Windows: `custom-gcl.exe`). Дальше в CI и локально вызывай **этот** бинарник вместо обычного `golangci-lint`.

Имя и путь выхода можно задать в `.custom-gcl.yml` (`name`, `destination`).

### 3. Конфиг проекта

В **проверяемом** проекте положи `.golangci.yml` (или дополни существующий). Минимальный пример — см. `.golangci.example.yml` в репозитории линтера:

- `linters.enable` содержит `loglint`
- `linters.settings.custom.loglint.type: module`

Без `default: none` custom-линтеры тоже участвуют в прогоне; если везде отключено всё кроме списка — добавь `loglint` в `enable`.

### 4. Запуск

```bash
./custom-gcl run ./...
```

Или с явным включением:

```bash
./custom-gcl run -E loglint ./...
```

Флаг `-config` анализатора (JSON) при использовании через golangci-lint задаётся через `settings` custom-секции — см. доку golangci-lint по custom linters; при необходимости можно расширить плагин и читать `register.DecodeSettings` в `New()`.

## Альтернатива: только бинарник loglint

Без custom-сборки можно по-прежнему использовать:

```bash
go run github.com/pppestto/test_assignment_linter/cmd/loglint@latest ./...
# или
go build -o loglint ./cmd/loglint && ./loglint ./...
```

Плагин golangci-lint нужен, если хочется один прогон со всеми линтерами и единый `.golangci.yml`.

## Структура плагина

| Файл / пакет | Назначение |
|--------------|------------|
| `linters/plugin.go` | `register.Plugin("loglint", New)`, `BuildAnalyzers` → `addcheck.Analyzer`, `LoadModeTypesInfo` |
| `addcheck/` | Сам анализатор и правила |

---




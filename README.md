# test_assignment_linter

Линтер для проверки **первого строкового аргумента** в вызовах **log/slog** и **go.uber.org/zap** (Info, Error, Debug, Warn и т.д.):

| Правило | Суть |
|--------|------|
| **Lowercase** | Сообщение с маленькой буквы |
| **English** | Только ASCII (латиница, цифры, пунктуация ASCII — но см. следующее) |
| **Emoji** | По ТЗ допускаются только **буквы, цифры и пробелы**; всё остальное (эмодзи, `!`, `…`, кириллица и т.д.) — нарушение |
| **Sensitive** | Нет утечки через конкатенацию: `"token: " + token` режется на переменной; статическая строка `"token validated"` — ок |

Конфиг JSON (`-config`), чувствительные ключевые слова и regexp-паттерны — см. `loglint.example.json`. Поддерживается `-fix` (правки в духе lowercase и вычищения «лишних» символов по правилам ТЗ).

**Go:** в CI прогоняется на **1.22+**; локально нужен Go не ниже `go` из `go.mod`.

---

## Сборка и запуск

```bash
go build -o loglint ./cmd/loglint
./loglint ./...

go test ./... -count=1
```

- С конфигом: `./loglint -config=loglint.example.json ./...`
- С автофиксом: `./loglint -fix ./...`

Без сборки:

```bash
go run ./cmd/loglint ./...
```

---

## Примеры (соответствие правилам)

| Правило | Неправильно | Правильно |
|--------|-------------|-----------|
| Строчная буква | `slog.Error("Failed to connect")` | `slog.Error("failed to connect")` |
| Только ASCII / латиница | `slog.Info("запуск сервиса")` | `slog.Info("service starting")` |
| Только буквы, цифры, пробелы | `"failed!!!"`, `"ok 👍"`, `"cost: $10"` | `"connection failed"`, `"ok"`, `"cost 10 dollars"` |
| Чувствительные данные | `"token: " + token` | `"token validated"` или литерал без конкатенации с секретом |

---

## Интеграция с golangci-lint

Линтер — `go/analysis.Analyzer` с именем `loglint`, подключается через **Module Plugin System** (`golangci-lint custom`).

### Требования

- Установленный `golangci-lint` (команда `golangci-lint custom` в стандартной поставке).
- Git, Go.

### 1. `.custom-gcl.yml`

В корне репозитория уже есть `.custom-gcl.yml`: указан модуль и импорт пакета `linters`, в `init()` вызывается `register.Plugin("loglint", ...)`.

Если клонируешь в другое место — поправь `path:` на каталог с `go.mod` этого модуля.

**Важно:** поле `version:` должно **совпадать** с версией твоего `golangci-lint` (иначе сборка custom-бинарника может не совпасть по API).

```bash
golangci-lint version
```

Подставь выводимую версию в `version:` в `.custom-gcl.yml` (в репозитории сейчас задано значение-пример — замени на своё).

### 2. Сборка custom-бинарника

Из корня репозитория линтера:

```bash
golangci-lint custom
```

По умолчанию появится бинарник в текущей директории (часто `custom-gcl` / на Windows `custom-gcl.exe`). Имя и путь можно задать в `.custom-gcl.yml` (`name`, `destination`).

### 3. Конфиг проверяемого проекта

В **целевом** проекте в `.golangci.yml` (или аналоге) нужно включить кастомный линтер. Минимальный пример — **`.golangci.example.yml`** в этом репозитории:

- `linters.enable` содержит `loglint`
- `linters.settings.custom.loglint.type: module`

Если используешь `default: none` — не забудь добавить `loglint` в `enable`.

### 4. Запуск

```bash
./custom-gcl run ./...
# или явно:
./custom-gcl run -E loglint ./...
```

Передача JSON-конфига анализатора (`-config`) при запуске через golangci-lint настраивается в секции `settings` для custom-линтера; при необходимости плагин можно расширить (например, `register.DecodeSettings` в `New()`).

---

## Только бинарник, без custom-сборки

```bash
go run github.com/pppestto/test_assignment_linter/cmd/loglint@latest ./...
# или после go build:
./loglint ./...
```

Custom-плагин нужен, если хочешь один прогон со всеми линтерами и единый `.golangci.yml`.

---

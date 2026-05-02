# mcp-obsidian

MCP-сервер на Go для взаимодействия с [Obsidian](https://obsidian.md) через плагин [Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api).

## Требования

- **Go 1.21+**
- **Obsidian** с установленным и настроенным плагином [Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api)

## Установка плагина в Obsidian

1. Откройте Obsidian → **Настройки** → **Сторонние плагины** → **Обзор**
2. Найдите **Local REST API** и установите его
3. Включите плагин и скопируйте **API Key** из настроек плагина

## Сборка

```bash
git clone <repo-url> mcp-obsidian
cd mcp-obsidian
go mod tidy
go build -o mcp-obsidian ./...
```

Или напрямую:

```bash
go install github.com/user/mcp-obsidian@latest
```

## Запуск

### Через переменные окружения (рекомендуется)

```bash
export OBSIDIAN_API_KEY="ваш-api-ключ"
export OBSIDIAN_URL="http://127.0.0.1:27123"   # опционально, это значение по умолчанию
./mcp-obsidian
```

### Через флаги командной строки

```bash
./mcp-obsidian -key "ваш-api-ключ" -url "http://127.0.0.1:27123"
```

## Настройка в Claude Desktop

Добавьте в `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "obsidian": {
      "command": "/path/to/mcp-obsidian",
      "env": {
        "OBSIDIAN_API_KEY": "ваш-api-ключ",
        "OBSIDIAN_URL": "http://127.0.0.1:27123"
      }
    }
  }
}
```

## Настройка в VS Code (Copilot / Continue)

В файле `.vscode/mcp.json`:

```json
{
  "servers": {
    "obsidian": {
      "type": "stdio",
      "command": "/path/to/mcp-obsidian",
      "env": {
        "OBSIDIAN_API_KEY": "ваш-api-ключ"
      }
    }
  }
}
```

## Доступные инструменты

| Инструмент | Описание |
|---|---|
| `list_files` | Список файлов и папок в хранилище (или поддиректории) |
| `read_note` | Чтение содержимого заметки |
| `create_note` | Создание новой заметки |
| `update_note` | Перезапись существующей заметки |
| `append_to_note` | Добавление текста в конец заметки |
| `delete_note` | Удаление заметки |
| `search_notes` | Полнотекстовый поиск по хранилищу |
| `get_active_note` | Получение заметки, открытой в Obsidian прямо сейчас |
| `open_note` | Открытие заметки в Obsidian |
| `get_periodic_note` | Получение периодической заметки (дневная/недельная/месячная/годовая) |
| `server_info` | Информация о сервере и хранилище |

## Примеры использования

После подключения к Claude или другому AI-ассистенту вы можете:

- *"Покажи все заметки в папке Projects"*
- *"Прочитай заметку Daily/2024-01-15.md"*
- *"Создай новую заметку Ideas/МояИдея.md с содержимым..."*
- *"Найди все заметки, где упоминается 'Golang'"*
- *"Добавь в сегодняшнюю дневную заметку: встреча прошла хорошо"*
- *"Что сейчас открыто в Obsidian?"*

## Переменные окружения

| Переменная | По умолчанию | Описание |
|---|---|---|
| `OBSIDIAN_API_KEY` | *(обязательно)* | API-ключ из плагина Local REST API |
| `OBSIDIAN_URL` | `http://127.0.0.1:27123` | Базовый URL сервера Obsidian |

## Архитектура

```
mcp-obsidian/
├── main.go                      # Точка входа, инициализация сервера
├── go.mod                       # Модуль Go
├── internal/
│   ├── obsidian/
│   │   └── client.go            # HTTP-клиент для Obsidian Local REST API
│   └── tools/
│       └── tools.go             # Регистрация всех MCP-инструментов
└── README.md
```

Сервер использует транспорт **stdio** (стандартный ввод/вывод) согласно спецификации MCP, что обеспечивает совместимость со всеми MCP-клиентами (Claude Desktop, VS Code, Cursor и др.).

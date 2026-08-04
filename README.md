# Workspace Monorepo

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Node.js](https://img.shields.io/badge/Node.js-24-green.svg)](https://nodejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9-blue.svg)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.23-00ADD8.svg)](https://golang.org/)

> Полнофункциональный monorepo проект с TypeScript/Express API сервером, React UI и MCP сервером для Obsidian на Go.

## 📋 Содержание

- [Обзор](#-обзор)
- [Структура проекта](#-структура-проекта)
- [Технологический стек](#-технологический-стек)
- [Требования](#-требования)
- [Установка](#-установка)
- [Команды](#-команды)
- [Пакеты](#-пакеты)
  - [API Spec](#api-spec)
  - [API Zod](#api-zod)
  - [API Client React](#api-client-react)
  - [Database](#database)
  - [API Server](#api-server)
  - [Mockup Sandbox](#mockup-sandbox)
  - [Scripts](#scripts)
  - [MCP Obsidian (Go)](#mcp-obsidian-go)
- [Разработка](#-разработка)
- [Безопасность](#-безопасность)
- [Лицензия](#-лицензия)

## 🎯 Обзор

Это **pnpm workspace monorepo**, использующий TypeScript для типизации. Каждый пакет управляет своими зависимостями независимо. Проект включает в себя:

- **Backend**: Express 5 API сервер с PostgreSQL + Drizzle ORM
- **Frontend**: React + Vite sandbox для разработки UI
- **API Codegen**: Автоматическая генерация hooks и схем из OpenAPI спецификации
- **MCP Server**: Go сервер для интеграции с Obsidian через Local REST API

## 📁 Структура проекта

```
workspace/
├── artifacts/                    # Приложения (сервера, UI)
│   ├── api-server/              # Express API сервер
│   └── mockup-sandbox/          # React UI sandbox
├── lib/                         # Общие библиотеки
│   ├── api-client-react/        # React Query hooks для API
│   ├── api-spec/                # OpenAPI спецификация
│   ├── api-zod/                 # Zod схемы для валидации
│   └── db/                      # Database слой (Drizzle ORM)
├── mcp-obsidian/                # MCP сервер на Go
├── scripts/                     # Утилиты и скрипты
├── package.json                 # Корневой package.json
├── pnpm-workspace.yaml          # Конфигурация workspace
└── tsconfig.base.json           # Базовая TS конфигурация
```

## 🛠 Технологический стек

### Backend
| Компонент | Версия | Описание |
|-----------|--------|----------|
| **Node.js** | 24 | Runtime окружение |
| **TypeScript** | 5.9 | Статическая типизация |
| **Express** | 5 | Web фреймворк |
| **PostgreSQL** | - | База данных |
| **Drizzle ORM** | 0.45+ | TypeScript ORM |
| **Zod** | 3.25+ | Валидация схем |
| **Pino** | 9+ | Логирование |

### Frontend
| Компонент | Версия | Описание |
|-----------|--------|----------|
| **React** | 19.1.0 | UI библиотека |
| **Vite** | 7.3+ | Build tool |
| **TanStack Query** | 5.90+ | Data fetching |
| **TailwindCSS** | 4.1+ | Стилизация |
| **Radix UI** | - | Headless компоненты |
| **Framer Motion** | 12.2+ | Анимации |

### DevOps & Tools
| Компонент | Версия | Описание |
|-----------|--------|----------|
| **pnpm** | - | Менеджер пакетов |
| **Orval** | 8.5+ | API codegen |
| **esbuild** | 0.27+ | Сборщик |
| **Go** | 1.23+ | Язык для MCP сервера |
| **MCP SDK** | 0.32+ | Model Context Protocol |

## 📦 Требования

Перед началом работы убедитесь, что установлены:

- **Node.js** версии 24 или выше
- **pnpm** (рекомендуется последняя версия)
- **Go** 1.23+ (для MCP сервера)
- **PostgreSQL** (для локальной разработки API)

```bash
# Проверка версий
node --version      # v24.x.x
pnpm --version      # 9.x.x
go version          # go1.23.x
```

## 🚀 Установка

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd workspace
```

### 2. Установка зависимостей

```bash
# Установить все зависимости workspace
pnpm install
```

### 3. Настройка переменных окружения

Создайте файл `.env` в корне проекта или в соответствующих пакетах:

```bash
# Database (обязательно для @workspace/db и @workspace/api-server)
DATABASE_URL=postgresql://user:password@localhost:5432/dbname

# API Server
PORT=3000
NODE_ENV=development
LOG_LEVEL=info

# MCP Obsidian (опционально)
OBSIDIAN_API_KEY=your-api-key
OBSIDIAN_URL=http://127.0.0.1:27123
```

## ⌨️ Команды

### Глобальные команды (из корня проекта)

```bash
# Типизация всех пакетов
pnpm run typecheck

# Сборка всех пакетов (с предварительной проверкой типов)
pnpm run build

# Запустить typecheck только для библиотек
pnpm run typecheck:libs
```

### Работа с отдельными пакетами

```bash
# Генерация API клиентов и Zod схем из OpenAPI
pnpm --filter @workspace/api-spec run codegen

# Миграции базы данных (dev only)
pnpm --filter @workspace/db run push

# Принудительный push схемы БД
pnpm --filter @workspace/db run push-force

# Запуск API сервера в режиме разработки
pnpm --filter @workspace/api-server run dev

# Сборка API сервера
pnpm --filter @workspace/api-server run build

# Запуск Mockup Sandbox (React UI)
pnpm --filter @workspace/mockup-sandbox run dev

# Сборка Mockup Sandbox
pnpm --filter @workspace/mockup-sandbox run build

# Запуск скрипта hello
pnpm --filter @workspace/scripts run hello
```

## 📦 Пакеты

### @workspace/api-spec

**Путь:** `lib/api-spec/`

Содержит OpenAPI спецификацию и конфигурацию для генерации кода.

**Зависимости:**
- `orval` ^8.5.2

**Команды:**
```bash
pnpm --filter @workspace/api-spec run codegen
```

**Файлы:**
- `openapi.yaml` — OpenAPI 3.1.0 спецификация
- `orval.config.ts` — Конфигурация генератора

---

### @workspace/api-zod

**Путь:** `lib/api-zod/`

Zod схемы для валидации API запросов/ответов, сгенерированные из OpenAPI spec.

**Зависимости:**
- `zod` ^3.25.76

**Экспорты:**
```typescript
import { HealthCheckResponse } from '@workspace/api-zod';
```

---

### @workspace/api-client-react

**Путь:** `lib/api-client-react/`

React Query hooks для работы с API, сгенерированные автоматически.

**Зависимости:**
- `@tanstack/react-query` ^5.90.21
- `react` >=18 (peer)

**Использование:**
```typescript
import { useHealthCheck, setBaseUrl, setAuthTokenGetter } from '@workspace/api-client-react';

// Настройка базового URL и токена
setBaseUrl('/api');
setAuthTokenGetter(() => 'your-token');

// Использование хука
const { data } = useHealthCheck();
```

---

### @workspace/db

**Путь:** `lib/db/`

Слой работы с базой данных на основе Drizzle ORM.

**Зависимости:**
- `drizzle-orm` ^0.45.2
- `drizzle-zod` ^0.8.3
- `pg` ^8.20.0
- `zod` ^3.25.76

**Dev Dependencies:**
- `drizzle-kit` ^0.31.9
- `@types/pg` ^8.18.0

**Использование:**
```typescript
import { db, pool } from '@workspace/db';
import { usersTable, insertUserSchema } from '@workspace/db/schema';

// Типы
type InsertUser = z.infer<typeof insertUserSchema>;
type User = typeof usersTable.$inferSelect;

// Запрос
const users = await db.select().from(usersTable);
```

**Миграции:**
```bash
# Применить миграции (dev)
pnpm --filter @workspace/db run push

# Принудительно (осторожно!)
pnpm --filter @workspace/db run push-force
```

---

### @workspace/api-server

**Путь:** `artifacts/api-server/`

Express 5 API сервер с логгированием, CORS и валидацией.

**Зависимости:**
- `express` ^5
- `@workspace/api-zod` workspace:*
- `@workspace/db` workspace:*
- `cors` ^2
- `pino` ^9
- `pino-http` ^10
- `cookie-parser` ^1.4.7
- `drizzle-orm` ^0.45.2

**Dev Dependencies:**
- `esbuild` ^0.27.3
- `esbuild-plugin-pino` ^2.3.3
- `pino-pretty` ^13
- `@types/express` ^5.0.6

**Запуск:**
```bash
# Development режим
pnpm --filter @workspace/api-server run dev

# Production сборка
pnpm --filter @workspace/api-server run build
pnpm --filter @workspace/api-server run start
```

**Структура:**
```
api-server/src/
├── index.ts           # Точка входа
├── app.ts             # Настройка Express приложения
├── lib/
│   └── logger.ts      # Pino logger
├── middlewares/       # Кастомные middleware
└── routes/
    ├── index.ts       # Роутер
    └── health.ts      # Health check endpoint
```

**Endpoints:**
- `GET /api/healthz` — Проверка здоровья сервера

---

### @workspace/mockup-sandbox

**Путь:** `artifacts/mockup-sandbox/`

React + Vite песочница для разработки UI компонентов с полной поддержкой Radix UI и TailwindCSS.

**Ключевые зависимости:**
- `react` 19.1.0
- `react-dom` 19.1.0
- `vite` ^7.3.2
- `@tailwindcss/vite` ^4.1.14
- `@tanstack/react-query` ^5.90.21
- `framer-motion` ^12.2.24
- `lucide-react` ^0.545.0
- `wouter` ^3.3.5
- `zod` ^3.25.76

**UI библиотеки:**
- Полный набор Radix UI компонентов
- shadcn/ui совместимые компоненты
- Recharts для графиков
- Sonner для уведомлений
- Vaul для drawer компонентов

**Запуск:**
```bash
# Development сервер
pnpm --filter @workspace/mockup-sandbox run dev

# Production сборка
pnpm --filter @workspace/mockup-sandbox run build

# Preview production сборки
pnpm --filter @workspace/mockup-sandbox run preview
```

**Особенности:**
- Алиас `@` для `src/` директории
- Cartographer плагин для отладки
- Runtime error modal в development
- TailwindCSS v4 с JIT компиляцией

---

### @workspace/scripts

**Путь:** `scripts/`

Утилиты и вспомогательные скрипты.

**Команды:**
```bash
# Пример скрипта
pnpm --filter @workspace/scripts run hello

# Typecheck
pnpm --filter @workspace/scripts run typecheck
```

**Git hooks:**
```bash
# Post-merge hook (автоматически устанавливает зависимости и обновляет БД)
./scripts/post-merge.sh
```

---

### mcp-obsidian (Go)

**Путь:** `mcp-obsidian/`

MCP (Model Context Protocol) сервер на Go для взаимодействия с Obsidian через плагин [Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api).

**Требования:**
- Go 1.23+
- Obsidian с плагином Local REST API

**Зависимости:**
- `github.com/mark3labs/mcp-go` v0.32.0
- `github.com/google/uuid` v1.6.0
- `github.com/spf13/cast` v1.7.1

**Сборка:**
```bash
cd mcp-obsidian
go mod tidy
go build -o mcp-obsidian .
```

**Запуск:**
```bash
# Через переменные окружения
export OBSIDIAN_API_KEY="your-api-key"
export OBSIDIAN_URL="http://127.0.0.1:27123"  # опционально
./mcp-obsidian

# Через флаги
./mcp-obsidian -key "your-api-key" -url "http://127.0.0.1:27123"
```

**Настройка в Claude Desktop:**
```json
{
  "mcpServers": {
    "obsidian": {
      "command": "/path/to/mcp-obsidian",
      "env": {
        "OBSIDIAN_API_KEY": "your-api-key",
        "OBSIDIAN_URL": "http://127.0.0.1:27123"
      }
    }
  }
}
```

**Доступные инструменты:**

| Инструмент | Описание |
|------------|----------|
| `list_files` | Список файлов и папок в хранилище |
| `read_note` | Чтение содержимого заметки |
| `create_note` | Создание новой заметки |
| `update_note` | Перезапись существующей заметки |
| `append_to_note` | Добавление текста в конец заметки |
| `delete_note` | Удаление заметки |
| `search_notes` | Полнотекстовый поиск по хранилищу |
| `get_active_note` | Получение текущей открытой заметки |
| `open_note` | Открытие заметки в Obsidian |
| `get_periodic_note` | Получение периодической заметки (daily/weekly/monthly) |
| `server_info` | Информация о сервере и хранилище |

**Архитектура:**
```
mcp-obsidian/
├── main.go              # Точка входа, инициализация MCP сервера
├── go.mod               # Go модуль
├── go.sum               # Зависимости
├── internal/
│   ├── obsidian/
│   │   └── client.go    # HTTP клиент для Obsidian API
│   └── tools/
│       └── tools.go     # Регистрация MCP инструментов
└── README.md            # Документация
```

Подробнее см. [mcp-obsidian/README.md](./mcp-obsidian/README.md)

## 👨‍💻 Разработка

### Добавление нового пакета

1. Создайте директорию в `lib/` или `artifacts/`
2. Инициализируйте `package.json`:
```json
{
  "name": "@workspace/new-package",
  "version": "0.0.0",
  "private": true,
  "type": "module"
}
```
3. Добавьте зависимости из каталога `pnpm-workspace.yaml`
4. Настройте `tsconfig.json` если нужно

### Работа с API

1. Измените `lib/api-spec/openapi.yaml`
2. Запустите генерацию:
```bash
pnpm --filter @workspace/api-spec run codegen
```
3. Используйте сгенерированные хуки в React или схемы в сервере

### Работа с базой данных

1. Добавьте модель в `lib/db/src/schema/`
2. Примените миграции:
```bash
pnpm --filter @workspace/db run push
```

**Пример модели:**
```typescript
// lib/db/src/schema/users.ts
import { pgTable, text, serial, timestamp } from "drizzle-orm/pg-core";
import { createInsertSchema } from "drizzle-zod";
import { z } from "zod";

export const usersTable = pgTable("users", {
  id: serial("id").primaryKey(),
  email: text("email").notNull().unique(),
  name: text("name"),
  createdAt: timestamp("created_at").defaultNow().notNull(),
});

export const insertUserSchema = createInsertSchema(usersTable).omit({
  id: true,
  createdAt: true
});

export type InsertUser = z.infer<typeof insertUserSchema>;
export type User = typeof usersTable.$inferSelect;
```

### Сборка проекта

```bash
# Полная сборка всех пакетов
pnpm run build

# Только проверка типов
pnpm run typecheck
```

## 🔒 Безопасность

### Supply-chain защита

В проекте включена защита от supply-chain атак через `minimumReleaseAge` в `pnpm-workspace.yaml`:

```yaml
minimumReleaseAge: 1440  # 1 день
```

Это предотвращает установку пакетов, опубликованных менее 24 часов назад.

**Исключения:**
- `@replit/*` — пакеты Replit
- `stripe-replit-sync` — доверенный пакет

⚠️ **Не отключайте эту настройку!** Отключение оставляет проект уязвимым для атак через вредоносные npm пакеты.

### Sensible defaults

- Логирование автоматически redact'ит чувствительные данные (authorization, cookies)
- CORS настроен с безопасными значениями по умолчанию
- Строгие проверки типов TypeScript включены

## 📄 Лицензия

MIT License — см. [LICENSE](LICENSE) файл для деталей.

---

## 🔗 Полезные ссылки

- [Документация pnpm workspaces](https://pnpm.io/workspaces)
- [Drizzle ORM документация](https://orm.drizzle.team/)
- [Orval — API codegen](https://orval.dev/)
- [MCP Protocol](https://modelcontextprotocol.io/)
- [Obsidian Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api)

---

**Made with ❤️ using TypeScript, Go, and modern web technologies**

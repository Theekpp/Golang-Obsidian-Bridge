# Workspace

## Overview

pnpm workspace monorepo using TypeScript. Each package manages its own dependencies.

## Stack

- **Monorepo tool**: pnpm workspaces
- **Node.js version**: 24
- **Package manager**: pnpm
- **TypeScript version**: 5.9
- **API framework**: Express 5
- **Database**: PostgreSQL + Drizzle ORM
- **Validation**: Zod (`zod/v4`), `drizzle-zod`
- **API codegen**: Orval (from OpenAPI spec)
- **Build**: esbuild (CJS bundle)

## Key Commands

- `pnpm run typecheck` — full typecheck across all packages
- `pnpm run build` — typecheck + build all packages
- `pnpm --filter @workspace/api-spec run codegen` — regenerate API hooks and Zod schemas from OpenAPI spec
- `pnpm --filter @workspace/db run push` — push DB schema changes (dev only)
- `pnpm --filter @workspace/api-server run dev` — run API server locally

See the `pnpm-workspace` skill for workspace structure, TypeScript setup, and package details.

## MCP Obsidian Server (Go)

A standalone MCP server written in Go for interacting with Obsidian via the Local REST API plugin.

- **Location**: `mcp-obsidian/`
- **Language**: Go 1.25
- **MCP library**: `github.com/mark3labs/mcp-go v0.32.0`
- **Transport**: stdio (JSON-RPC 2.0)

### Build

```bash
cd mcp-obsidian
go build -o mcp-obsidian .
```

### Run

```bash
export OBSIDIAN_API_KEY="your-key"
./mcp-obsidian
```

### Tools provided

`list_files`, `read_note`, `create_note`, `update_note`, `append_to_note`, `delete_note`, `search_notes`, `get_active_note`, `open_note`, `get_periodic_note`, `server_info`

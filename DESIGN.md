# Contributor notes

User-facing docs: [README.md](README.md). Product API and naming live in [hypermesh-docs](https://github.com/FyberLabs/hypermesh-docs) — this file only records how *this* binary wires to them.

## Surfaces

| Command | HTTP |
|---|---|
| `catalog` / `catalog show` | `GET /api/v1/hypermesh/catalog` |
| `classes` | `GET /api/v1/hypermesh/classes` |
| `hosts` | `GET /api/v1/hypermesh/renter/hosts` |
| `checkout` | `POST /api/v1/hypermesh/leases` |
| `lease list` / `show` / `complete` | `GET` / `POST` under `/api/v1/hypermesh/leases` |
| `chat` / `prompt` / `completions create` | `POST {HYPERMESH_CHAT_BASE}/v1/chat/completions` |
| `mcp …` | files under `HYPERMESH_CONFIG_DIR` |

Chat uses the chat base, not the control-plane renter chat path.

## Checkout body

```json
{
  "kind": "p2_loaded_model",
  "renter_user_id": "<uuid>",
  "catalog_id": "<from catalog>",
  "success_url": "…",
  "cancel_url": "…",
  "reserved_hours": 1,
  "purpose": "renter",
  "device_id": "<from hosts>"
}
```

`device_id` is the UUID from `hosts`. `--catalog-id` / `--model` default to whatever `DefaultCatalogID` is in code today; that constant is a convenience default, not a catalog inventory.

## Script mode

`--script`: assistant text on stdout, diagnostics on stderr, exit `1` on failure. `scripts/hypermesh-prompt.ps1` execs the same binary.

## Local MCP files

| File | Shape |
|---|---|
| `mcp.json` | `{ "mcpServers": { … } }` |
| `mcp-profiles.json` | `{ "active", "profiles" }` |
| `mcp-bindings.json` | `{ "bindings": [ … ] }` |

The visor owns attach, prefer-MCP, tool audit, and vault `source: "mcp"`.

# Hypermesh CLI — design notes

Contributor lock for this repo. Public how-to lives in [README.md](README.md). Do not invent a second API.

## Names

- Binary: `hypermesh` (alias `hm`).
- Product: **Hypermesh** / **Hyperme.sh**. Never bare Hyperme.

## Auth

One org API key from api-keys (`purpose: renter`) for both control plane and router.

| Header | Value |
|---|---|
| `X-Api-Key` | renter org key |
| `X-Tenant-ID` | tenant id |

Reject `hm_dev_`, `hm_rtr_`, and `hm_site_` as renter identity.

No OIDC dance in the CLI. No SIWE. No WireGuard.

## Bases

| Env | Default |
|---|---|
| `HYPERMESH_API_BASE` | `https://api.test.hyperme.sh` |
| `HYPERMESH_CHAT_BASE` | `https://chat.test.hyperme.sh` |

Config: `~/.config/hypermesh/config.toml`. Credentials file `~/.config/hypermesh/credentials` mode `0600`. Override directory with `HYPERMESH_CONFIG_DIR`.

## Locked routes

| Command | HTTP |
|---|---|
| `catalog` / `catalog show` | `GET /api/v1/hypermesh/catalog` (public). Show is a client-side filter. |
| `classes` | `GET /api/v1/hypermesh/classes` (public) |
| `hosts` | `GET /api/v1/hypermesh/renter/hosts` (same auth as `POST /leases`). Optional `catalog_id` query only. |
| `checkout` | `POST /api/v1/hypermesh/leases` |
| `lease list` | `GET /api/v1/hypermesh/leases` |
| `lease show` | `GET /api/v1/hypermesh/leases/{id}` |
| `lease complete` | `POST /api/v1/hypermesh/leases/{id}/complete` |
| `chat` / `prompt` / `completions create` | `POST {HYPERMESH_CHAT_BASE}/v1/chat/completions` |

Do **not** call `POST /api/v1/hypermesh/renter/chat/completions` (always-409 stub).

## LeaseCreate

Default checkout body (flags override catalog id and hours):

```json
{
  "kind": "p2_loaded_model",
  "renter_user_id": "<uuid>",
  "catalog_id": "<catalog id>",
  "success_url": "…",
  "cancel_url": "…",
  "reserved_hours": 1,
  "purpose": "renter",
  "device_id": "<uuid>"
}
```

`DefaultCatalogID` in code is the current API default (`llama-3.1-8b-q4`). It is a default flag value, not a product ceiling — callers pass `--catalog-id` for other catalog entries.

`device_id` is the plane UUID from `hosts`, not `public_label`. Required when `purpose` is `renter`. The client rejects a missing or non-UUID `device_id` before POST. The CLI does not auto-pick a host.

`hosts` rows: `device_id`, `public_label` (display only), `class_id`, `certified`, `online`, `sell_state`. No serial, secrets, or host console fields.

Lease fields the CLI prints: `id`, `status`, `checkout_url`. Stripe Checkout only.

Status motion: `offered` → `paid` → `starting` → `active` → `ended` | `failed` | `refunded`.

`checkout` prints `id` and `checkout_url` immediately, opens `checkout_url` unless `--no-open`, and `--wait` polls `GET /leases/{id}` until `active` or `failed`.

## Router chat

Paid lease ticket plus the same renter key.

- `lease_id` in the JSON body
- and `X-Hypermesh-Lease-Id` / `X-Lease-Id`
- OpenAI-shaped `model` + `messages`

Never log prompt bodies.

## Shell

`prompt --script`, `chat --script`, and `completions create --script` are the non-interactive mode for bash.

- stdout is only the assistant text (one trailing newline if the model text has none)
- diagnostics go to stderr and are not part of the model text
- exit status is `0` on success and `1` on any failure (the HTTP status is not the process status)
- `--json` stays the raw completion and cannot be combined with `--script`
- empty assistant text is a failure; the CLI does not dump the completion JSON onto stdout

`scripts/hypermesh-prompt.ps1` execs this same binary (`prompt --script`). It does not open a second HTTP client. Chat remains `POST {HYPERMESH_CHAT_BASE}/v1/chat/completions`.

## Local MCP

Config under `HYPERMESH_CONFIG_DIR` (default `~/.config/hypermesh`):

| File | Shape |
|---|---|
| `mcp.json` | Cursor-compatible `{ "mcpServers": { … } }` |
| `mcp-profiles.json` | `{ "active", "profiles": { name: { servers, config } } }` |
| `mcp-bindings.json` | `{ "bindings": [ { "server", "wm_class"?, "app_id"?, "executable"? } ] }` |

Commands: `mcp catalog ls`, `mcp list`, `mcp profile ls|create|use|add`, `mcp profile config set`, `mcp import cursor|docker`, `mcp doctor`, `mcp bindings ls|add`. Doctor checks PATH (stdio) or TCP reachability (http/sse). It does not call tools. The visor owns MCP client lifecycle, prefer-MCP detection, tool-invoke audit, and vault `source: "mcp"` (locator `server/tool` against `mcp.json` — role B, separate from session tool attach).

## Principles

Thin client over the REST lock in [FyberLabs/hypermesh-docs](https://github.com/FyberLabs/hypermesh-docs):

- [customer-interfaces.md](https://github.com/FyberLabs/hypermesh-docs/blob/main/customer-interfaces.md) — one API, CLI is a key, not a second login
- [full-model.md](https://github.com/FyberLabs/hypermesh-docs/blob/main/full-model.md) — product shape for loaded models
- [payments.md](https://github.com/FyberLabs/hypermesh-docs/blob/main/payments.md) — Stripe Checkout, Fyber Labs merchant of record
- [router.md](https://github.com/FyberLabs/hypermesh-docs/blob/main/router.md) — customers scale on HTTPS; host WireGuard is not a renter path
- [public-sites.md](https://github.com/FyberLabs/hypermesh-docs/blob/main/public-sites.md) — do not dump internal locks onto customer copy

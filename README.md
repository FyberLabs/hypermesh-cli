# hypermesh-cli

Thin **Hypermesh** CLI for [Hyperme.sh](https://hyperme.sh): browse the catalog, check out a lease, chat on the router, and manage local MCP config for the desktop visor.

Binaries: `hypermesh` and `hm`.

This is a client over the Hypermesh HTTP API. It is not a second control plane.

Design notes for contributors: [DESIGN.md](DESIGN.md). Product docs: [FyberLabs/hypermesh-docs](https://github.com/FyberLabs/hypermesh-docs).

## Install

Go 1.22+ on the `PATH`.

```bash
git clone https://github.com/FyberLabs/hypermesh-cli.git
cd hypermesh-cli
make build          # bin/hypermesh and bin/hm
make vet
make test
make install        # PREFIX=/usr/local (override as needed)
```

CI (`.github/workflows/test.yml`) runs `go vet`, `make test`, and `make build` on every push and pull request.

Or:

```bash
go install github.com/FyberLabs/hypermesh-cli/cmd/hypermesh@latest
go install github.com/FyberLabs/hypermesh-cli/cmd/hm@latest
```

## Auth

Use a renter org API key (`purpose: renter`) for both the control plane and the chat router. Headers: `X-Api-Key` and `X-Tenant-ID`.

Do not use `hm_dev_`, `hm_rtr_`, or `hm_site_` keys — those are host, router, and site credentials.

```bash
hypermesh auth login \
  --api-key "$HYPERMESH_API_KEY" \
  --tenant-id "$HYPERMESH_TENANT_ID" \
  --renter-user-id "$HYPERMESH_RENTER_USER_ID"

hypermesh auth whoami
hypermesh auth logout
```

Config lives under `~/.config/hypermesh` (override with `HYPERMESH_CONFIG_DIR`):

| File | Purpose |
|---|---|
| `config.toml` | bases and preferences |
| `credentials` | API key (mode `0600`) |
| `mcp.json` | Cursor-shaped MCP servers |
| `mcp-profiles.json` | named MCP profiles |
| `mcp-bindings.json` | focused-app → server bindings |

## Quick start

1. Browse the public catalog and hardware classes (no key required).
2. List renter-safe hosts and copy a `device_id` (plane UUID). `public_label` is display only.
3. Check out a lease for a catalog model on that device, then pay Stripe Checkout.
4. When the lease is `active`, chat with the same key plus the lease id.

```bash
hypermesh catalog
hypermesh catalog show "$CATALOG_ID"
hypermesh classes

hypermesh hosts
hypermesh hosts --catalog-id "$CATALOG_ID"

hypermesh checkout \
  --device-id "$DEVICE_ID" \
  --catalog-id "$CATALOG_ID" \
  --renter-user-id "$HYPERMESH_RENTER_USER_ID" \
  --success-url "https://hyperme.sh/ok" \
  --cancel-url "https://hyperme.sh/cancel" \
  --no-open \
  --wait

hypermesh lease list
hypermesh lease show "$LEASE_ID"

hypermesh prompt --script --lease-id "$LEASE_ID" "hello"
hypermesh chat --script --lease-id "$LEASE_ID" --message "hello"
hypermesh completions create --script --lease-id "$LEASE_ID" --model "$CATALOG_ID" --message "hello"

# Bash: stdout is only the assistant text. Exit status is 1 on failure.
text=$(hypermesh prompt --script --lease-id "$LEASE_ID" "hello")

# PowerShell calls the same binary (not a second HTTP client).
./scripts/hypermesh-prompt.ps1 -LeaseId "$LEASE_ID" "hello"

hypermesh lease complete "$LEASE_ID"
```

Checkout defaults `--catalog-id` to the current API default (`llama-3.1-8b-q4`). Pass another catalog id when you want a different model. Missing `--device-id` fails before POST — the CLI does not pick a host.

`--json` works on every command except together with `--script`. Failures exit `1` (not the HTTP status). Prompt bodies are not logged.

`--script` is for shells: stdout is only the assistant text; errors stay on stderr. An empty assistant message is a failure.

Chat is `POST $HYPERMESH_CHAT_BASE/v1/chat/completions` with `lease_id` in the body and `X-Hypermesh-Lease-Id` / `X-Lease-Id`.

Lease status: `offered` → `paid` → `starting` → `active` → `ended` | `failed` | `refunded`. `--wait` polls until `active` or `failed`.

## Local MCP

Configure local MCP servers for the [hypermesh-visor](https://github.com/FyberLabs/hypermesh-visor) companion. The visor attaches the active profile when a session opens. When a binding matches the focused app and that server is healthy, the companion prefers `mcp:<id>` over mouse/type.

```bash
hypermesh mcp catalog ls
hypermesh mcp import cursor                 # or: hypermesh mcp import cursor /path/to/mcp.json
hypermesh mcp import docker                 # adds docker mcp gateway run
hypermesh mcp profile ls
hypermesh mcp profile create frontend
hypermesh mcp profile use frontend
hypermesh mcp profile add filesystem
hypermesh mcp profile config set filesystem.cwd=/tmp
hypermesh mcp bindings add fixture --wm-class FixtureApp
hypermesh mcp bindings ls
hypermesh mcp doctor
hypermesh mcp list
```

## Environment

| Env | Default |
|---|---|
| `HYPERMESH_API_BASE` | `https://api.test.hyperme.sh` |
| `HYPERMESH_CHAT_BASE` | `https://chat.test.hyperme.sh` |

Also: `HYPERMESH_API_KEY`, `HYPERMESH_TENANT_ID`, `HYPERMESH_RENTER_USER_ID`, `HYPERMESH_LEASE_ID`, `HYPERMESH_SUCCESS_URL`, `HYPERMESH_CANCEL_URL`, `HYPERMESH_CONFIG_DIR`.

`whoami` reads local config only.

## Commands

| Command | Route |
|---|---|
| `auth login\|whoami\|logout` | local config |
| `catalog` / `catalog show` | `GET /api/v1/hypermesh/catalog` |
| `classes` | `GET /api/v1/hypermesh/classes` |
| `hosts` | `GET /api/v1/hypermesh/renter/hosts` |
| `checkout` | `POST /api/v1/hypermesh/leases` |
| `lease list\|show\|complete` | `GET/POST /api/v1/hypermesh/leases…` |
| `chat` / `prompt` / `completions create` | `POST {chat base}/v1/chat/completions` |
| `mcp …` | local MCP config under the Hypermesh config dir |

## Not this CLI

Host enroll, device secrets, WireGuard, MHS clustering, and a second login (OIDC / SIWE) live elsewhere. This binary is the renter client.

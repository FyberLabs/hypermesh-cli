# hypermesh-cli

Thin **Hypermesh** CLI for [Hyperme.sh](https://hyperme.sh): browse the catalog, check out a lease, chat on the router, and manage local MCP config for the desktop visor.

Binaries: `hypermesh` and `hm`.

```bash
git clone https://github.com/FyberLabs/hypermesh-cli.git
cd hypermesh-cli
make build    # bin/hypermesh and bin/hm
make test
make install  # PREFIX=/usr/local
```

Or: `go install github.com/FyberLabs/hypermesh-cli/cmd/hypermesh@latest` (and the same for `cmd/hm`). Needs Go 1.22+.

## Auth

Use a renter org API key and tenant id (same key for API and chat).

```bash
hypermesh auth login \
  --api-key "$HYPERMESH_API_KEY" \
  --tenant-id "$HYPERMESH_TENANT_ID" \
  --renter-user-id "$HYPERMESH_RENTER_USER_ID"

hypermesh auth whoami
hypermesh auth logout
```

Config directory: `~/.config/hypermesh` (or `HYPERMESH_CONFIG_DIR`).

| File | Purpose |
|---|---|
| `config.toml` | API / chat bases and preferences |
| `credentials` | API key (mode `0600`) |
| `mcp.json` | MCP servers |
| `mcp-profiles.json` | named MCP profiles |
| `mcp-bindings.json` | focused-app → server bindings |

## Usage

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

text=$(hypermesh prompt --script --lease-id "$LEASE_ID" "hello")
./scripts/hypermesh-prompt.ps1 -LeaseId "$LEASE_ID" "hello"

hypermesh lease complete "$LEASE_ID"
```

Pick `CATALOG_ID` from `hypermesh catalog` and `DEVICE_ID` from `hypermesh hosts`. Checkout opens Stripe Checkout unless you pass `--no-open`; `--wait` polls until the lease is ready.

With `--script`, stdout is the assistant text and failures exit `1`. Use `--json` when you want the raw response (not with `--script`).

## Local MCP

For the [hypermesh-visor](https://github.com/FyberLabs/hypermesh-visor) companion:

```bash
hypermesh mcp catalog ls
hypermesh mcp import cursor
hypermesh mcp import docker
hypermesh mcp profile ls
hypermesh mcp profile create frontend
hypermesh mcp profile use frontend
hypermesh mcp profile add filesystem
hypermesh mcp bindings add filesystem --wm-class Code
hypermesh mcp bindings ls
hypermesh mcp doctor
hypermesh mcp list
```

## Environment

| Env | Default |
|---|---|
| `HYPERMESH_API_BASE` | `https://api.test.hyperme.sh` |
| `HYPERMESH_CHAT_BASE` | `https://chat.test.hyperme.sh` |

Also: `HYPERMESH_API_KEY`, `HYPERMESH_TENANT_ID`, `HYPERMESH_RENTER_USER_ID`, `HYPERMESH_LEASE_ID`, `HYPERMESH_SUCCESS_URL`, `HYPERMESH_CANCEL_URL`, `HYPERMESH_CONFIG_DIR`, `HYPERMESH_VISOR_URL`.

## Commands

| Command | What it hits |
|---|---|
| `auth …` | local config |
| `catalog` / `classes` | public catalog API |
| `hosts` | renter hosts |
| `checkout` / `lease …` | leases |
| `chat` / `prompt` / `completions create` | chat completions |
| `mcp …` | local MCP config |

More product context: [FyberLabs/hypermesh-docs](https://github.com/FyberLabs/hypermesh-docs). Contributor notes: [DESIGN.md](DESIGN.md).

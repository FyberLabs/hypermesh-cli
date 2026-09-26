# Contributor notes

Everything a reader needs to *use* this CLI is in [README.md](README.md). This file is only for people changing the binary.

## Layout

| Path | Role |
|---|---|
| `cmd/hypermesh`, `cmd/hm` | entrypoints |
| `internal/cli` | cobra commands |
| `internal/api` | HTTP client, URLs, lease body |
| `internal/config` | config + credentials |
| `internal/mcp` | local MCP files |
| `scripts/hypermesh-prompt.ps1` | PowerShell wrapper around `prompt --script` |

## Defaults in code

`DefaultCatalogID` / `DefaultAPIBase` / `DefaultChatBase` live in `internal/api/urls.go`. They seed flags when omitted — they are not a catalog inventory. Prefer reading ids from `GET /catalog` in examples and tests that care about a specific model.

## Chat path

Router chat is `{HYPERMESH_CHAT_BASE}/v1/chat/completions`. The control-plane path `/api/v1/hypermesh/renter/chat/completions` is unused by this client (see `PathRenterChatStub`).

## Tests

```bash
make vet
make test
make build
```

GitHub Actions (`.github/workflows/test.yml`) runs the same on push and pull request.

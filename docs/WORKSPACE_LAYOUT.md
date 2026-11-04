# Thoughts workspace layout

| Location | Purpose |
| --- | --- |
| `cmd/`, `internal/` | Go CLI, local service, Markdown/MCP operations and embedded frontend |
| `web/` | React/TypeScript/Tailwind/shadcn source and frontend tests |
| `adapters/` | Claude Desktop manifest and Claude Code/Codex stdio examples |
| `branding/`, `legal/`, `licenses/`, `release/` | Artwork, policy templates, dependency attribution and publisher/support metadata |
| `scripts/` | Windows builds, checks, packaging and inspect-first installation helpers |
| `docs/` | Current reports and documentation; `assets/` contains the reviewed isolated screenshot |
| `plugins/` | Preserved local Codex desktop-preview integration, outside the proposed public stdio distribution |
| `deferred/macos/` | Deferred implementation/history; not a supported platform or an active Windows build input |
| `Windows测试包/` | Local, ignored review copies: one clearly labelled current candidate and separate old versions |
| `dist/windows/<version>/` | Immutable build outputs; preserve paths while test services use them |
| `历史归档/` | Local-only original requirements, historical reports/screenshots, early executable and unrelated media |
| `.artifacts/` | Local test libraries, sessions, reports and build caches; may contain private runtime state |
| `.tools/` | Local development tools, not runtime dependencies for end users |

The Chinese start-here guide and installed-program shortcut are workstation conveniences, not package installation or automatic migration. No source, note library or installed extension is moved merely by opening a guide.

## Thoughts naming and compatibility

Current source uses `cmd/thoughts`, `internal/thoughts`, Go module `thoughts.local/thoughts`, and `thoughts.exe`. Frontend helpers and CSS use Thoughts names. Framework APIs such as React `useMemo` and Inno Setup `UpdateReadyMemo` retain their required upstream names.

Legacy on-disk Trash/runtime paths, browser preference keys, local protocol identifiers and explicit old configuration overrides remain unchanged so existing data, running integrations and software rollback remain readable. Existing historical executable versions are never renamed or deleted. New release packages and client examples use only `thoughts.exe`.

## Local archive policy

The 2026-11-01 organization pass preserves every moved file, with before/after SHA256 and a local move receipt. No permanent deletion, Git history rewriting, commit, push, security-policy change or library migration is part of that pass. Earlier screenshots and reports can contain workstation/private-library information; the entire historical archive stays outside the public source allowlist and is ignored by Git.

The original private Git history still contains the separately documented findings. Moving a working-copy report to an archive does not sanitize history. See [SECURITY_AUDIT.md](SECURITY_AUDIT.md) before any future publication.

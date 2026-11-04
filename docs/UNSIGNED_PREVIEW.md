# Thoughts unsigned experimental preview

Version: `0.3.0-preview.21`, Windows x64. **UNSIGNED — NOT A STABLE RELEASE.**

The maintainer explicitly authorized publishing an unsigned prerelease while preparing a SignPath Foundation application. This authorization supersedes the earlier source-only publication restriction for this preview; it does not establish signing trust, approve a stable release or authorize changes to Windows security settings. The product, executable, Go module and current source directories use Thoughts names. Existing data and historical software versions are preserved.

## Download and status

Release: [v0.3.0-preview.21](https://github.com/xuli82221-pixel/hello-Thoughts-/releases/tag/v0.3.0-preview.21). Published on 2026-11-04 at 11:00 Asia/Shanghai, after the manually triggered build, validation and maintainer upload completed:

- `thoughts-0.3.0-preview.21-windows-amd64-setup.exe`: offline current-user installer.
- `thoughts-0.3.0-preview.21-windows-amd64.zip`: portable package containing `thoughts.exe` and inspected client configuration examples.
- `thoughts-claude-desktop-0.3.0-preview.21-windows-amd64.mcpb`: Claude Desktop extension with the same `thoughts.exe` bytes.
- `SHA256SUMS.txt`, `BUILDINFO.json`, `UNSIGNED_NOTICE.txt`, `ACCEPTANCE.json`: checksums, source/build identity, unsigned notice and bounded test summaries.

Do not interpret GitHub's automatically generated source ZIP as the application. See [INSTALL.md](../INSTALL.md) for checksum and installation steps. These files have no trusted publisher signature. If Windows SAC, SmartScreen or organization policy blocks them, **stop and keep protection enabled**. Do not change execution policy, remove download provenance or use a self-signed certificate to simulate trust. SHA256 checks integrity, not publisher identity.

Build source: [`0b155f8a4d0381315f742bb6661e2e0eb6c738c4`](https://github.com/xuli82221-pixel/hello-Thoughts-/commit/0b155f8a4d0381315f742bb6661e2e0eb6c738c4). [Successful build and test run](https://github.com/xuli82221-pixel/hello-Thoughts-/actions/runs/41762640507). Core SHA256: `6f045833ae29839955edee44fed9f36b0b76b84566951b47c3615ce5f4c9a021`.

All seven uploaded asset digests were compared with the downloaded CI outputs before publication. ZIP and MCPB contain the same core bytes. Local file inspection confirmed `NotSigned` for the core, setup and registration helper. The GitHub release is public and explicitly marked prerelease; the tag points to the exact build commit above. No private history or local user data was uploaded.

## Build and verification boundaries

The read-only, manually dispatched [Windows preview workflow](../.github/workflows/unsigned-preview.yml) builds the exact public Git commit, uses locked frontend dependencies and a pinned Go toolchain, and verifies the official Inno Setup compiler's hash and publisher signature. No signing certificate or additional GitHub Secrets is used. It does not create or publish GitHub Releases; the authorized maintainer upload is separate.

The workflow runs frontend and native backend tests, package structure/checksum validation, actual stdio MCP operations and an isolated installer install/reinstall/uninstall check using synthetic data. The final `ACCEPTANCE.json` must report their actual results. GitHub-hosted Windows tests do not establish clean Smart App Control acceptance, real Claude Desktop/Claude Code/Codex host acceptance, or behavior on every Windows configuration. The earlier preview.20 installer was blocked on the maintainer's SAC-enabled machine; that history is retained and is not reclassified as a pass.

Only the explicit package files and sanitized summaries are uploaded. Raw QA homes, local sessions, tokens, logs, user paths and private libraries are excluded. `BUILDINFO.json` records the public source commit and Actions run. The release tag must point to that exact commit. Checksums are verified again after download and upload.

## Signing and privacy

**Pending — SignPath Foundation application not yet approved.** The application has not been submitted by this work. If approved, free code signing provided by SignPath.io, certificate by SignPath Foundation. See [Code signing policy](CODE_SIGNING_POLICY.md) and the [application draft](SIGNPATH_APPLICATION.md). Future signing requires provider review and a trusted SignPath integration; an unsigned CI build is not SignPath verification.

Thoughts stores notes and images locally, does not require a Thoughts account and does not upload local notes automatically. When authorized online AI clients read notes, their tool results may reach those providers under their own policies. Files are not encrypted by Thoughts; keep independent backups. See [PRIVACY.md](../PRIVACY.md).

## Naming and existing data

New programs and client configurations use `thoughts.exe`; source uses `cmd/thoughts`, `internal/thoughts` and `thoughts.local/thoughts`. This is a naming change, not a new note-management feature. No personal library, installed client, old executable or backup is renamed or deleted. Existing Trash/runtime storage and browser preference identifiers remain compatible. Save drafts and follow the existing reviewed update procedure if you deliberately upgrade a library; earlier local-development MCP launchers must be re-registered to the new `thoughts.exe`, because an old launcher may reject the newly named pin. A new release does not update your installation automatically.

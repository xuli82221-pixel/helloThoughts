# Historical shared-library and upgrade checks — preview.7–8

This is a version-scoped historical test report. Current published preview.21 status and remaining acceptance requirements are in [release readiness](RELEASE_READINESS_20261104.md). Its local-only statements, test counts, hashes and earlier language counts apply to the named earlier versions.

<!-- historical-records:start -->

## Historical records / 历史记录

Date: 2026-11-01. Historical candidate: **0.3.0-preview.8**, Windows x64. Local changes only; no commit, push or public release. **M0 has not passed.**

## Default sharing follow-up — preview.8

The default is now the Windows account profile's `Thoughts` home, obtained from the Windows process owner's profile rather than client `HOME`, `APPDATA` or working directory. Unset/empty, `auto` and the old literal `${HOME}/Thoughts` resolve there. An explicit custom home remains available for intentional separation. Explicit legacy `MEMO_CONFIG_DIR` is retained; a known old user-level Memo configuration stops automatic first initialization for migration rather than silently creating another empty library. Existing Store-app-specific/custom legacy locations are not searched or merged automatically.

The MCPB advanced setting is optional and defaults to `auto`. The registration helper also defaults to the same home and refuses conflicting existing registrations. `paths --json` inspects the selected home/config/library without creating files or starting a service. CLI initialization now uses a shared initialization lock as well as preserving the existing selection. Different Windows accounts or PCs are separate; no cloud sync was introduced.

| preview.8 check | Actual result |
| --- | --- |
| Frontend, Go, TypeScript/Vite and Windows x64 build | 17 frontend tests, 53 top-level Go tests and build passed; existing bundle advisory remains |
| Codex-first and Claude-first initialization | Passed in isolated Windows-account profile fixtures; each later host reads the real file saved by the earlier host, with shared service metadata |
| Different host HOME/APPDATA and current directory | Actual subprocess tests against the configured shared home passed for unset, auto and the older literal placeholder; same main library and service PID |
| Previously installed Codex/Claude launchers | Both direct status checks use preview.8 and the same main service PID/library after managed upgrade |
| Repeat init / explicit custom home / legacy guard / initialization lock | Go tests passed; path inspection creates no files |
| Actual two-MCP/web write, conflict protection, upgrade and rollback | Isolated runtime smoke passed on preview.8 |
| User-data integrity after actual preview.7 → preview.8 upgrade | All 36 original/backup files and 20 copied main-library baseline files unchanged by SHA256 |
| Portable ZIP and MCPB | CRC/path/per-file checks passed and identical core hash; optional auto manifest verified by package validator |
| Registration helper | Inspect-only default chooses the same account home; no client configuration changed by this check |
| Native-client limits | Existing Claude MCPB remains preview.7 with its correct absolute home and delegates to preview.8. The new optional-auto preview.8 MCPB was not imported in the GUI. Claude Code CLI is absent from PATH and standard checked install locations; native dispatch untested. Current Codex host reload/dispatch remains unverified |

preview.8 exe SHA256: `2c241f69717e9d6717492ec679c59c0cd5cca909d5f1933b9c235d079a61ba1f`; Authenticode remains NotSigned. Package hashes are in the immutable preview.8 `SHA256SUMS.txt`. No note body was sent to an AI during the read-only main-library status checks. Detailed local results remain in ignored artifacts. The following sections retain the preview.7 acceptance record; they do not imply clean-Windows or native-client acceptance for preview.8.

## Implemented behavior

All clients must set the same absolute `THOUGHTS_HOME`. The home holds the selected library configuration, managed core pin, sessions and reading metadata. A matching home lets separate MCP launchers reuse an identity-verified service. A vault-level exclusive Windows file lock also rejects a second service configured through another home. Unknown identities are not killed or trusted based on a PID or port alone.

Repeat initialization preserves the selected library. A conflicting new selection is rejected before creating that new destination. Migration remains deliberate: back up, copy, verify and configure; no implicit merging of other libraries occurs.

The new `runtime install`, `upgrade`, `status` and `rollback` commands accept only an explicitly supplied local executable and SHA256. Managed immutable versions live under the selected home. Each launcher verifies the selected pin, path and file hash before delegation; missing or modified files fail closed. Candidate capability and health checks are bounded. Downloaded-file security metadata is retained. No network updater, signature substitution or security-setting change was added.

A compatible service refuses upgrade shutdown while it has an active busy-editor lease. The lease contains an editor ID, sequence and busy status, never draft text. Software rollback preserves Markdown, images and configuration; it does not restore old note contents. Keep an independent library backup. Old launchers referenced by clients must remain in place.

Reader state stores validated note paths, tabs, selected page and scroll positions in the home, outside browser-origin storage. It contains no draft text, session tokens or note bodies. A new localhost port can therefore restore saved reading state after authentication. Concurrent browsers use the latest saved reading metadata; this is not independent per-device session sync. Browser closure or Windows restart can still lose unsaved text.

## Test environment and commands

Development host: Windows 11 x64; Go 1.27.3; existing local Node/npm and bundled Python. Test libraries were newly created isolated local directories. No acceptance note was written into the user's main library.

```powershell
pwsh -File scripts/build-desktop-adapters.ps1 -Python '<installed-python>'
python scripts/validate-desktop-adapters.py dist/windows/0.3.0-preview.7
python scripts/runtime-smoke.py --exe '<reviewed preview.7 exe>' --candidate '<reviewed compatible test exe>' --root '<new isolated directory>'
```

The upgrade candidate in the smoke test is an explicitly supplied build of this project labelled `0.3.0-preview.5-upgrade-test.1`. It is isolated test software, not a release or signing workaround. Its protocol compatibility tests do not certify arbitrary older versions.

| Check | Expected | Actual |
| --- | --- | --- |
| Build and unit tests, preview.7 | Frontend, TypeScript/Vite, Go and Windows x64 packages pass | 17 frontend tests and 49 top-level Go tests passed; TypeScript/Vite and build passed. Existing bundle-size advisory remains |
| Archive validation, preview.7 | CRC, path, per-file hashes and identical cores | Both packages passed; same executable SHA256 |
| Repeat init | Existing selection and file bytes retained | Passed, including rejection of another destination |
| Actual two MCP processes plus web API | Same PID/library; writes readable by other clients | Passed on preview.7 using real isolated Markdown |
| Stale write | Reject old SHA256 without replacing content | Passed |
| Managed install and repeat install | Reuse verified core | Passed |
| Wrong hash or busy editor | No active-pin change | Passed |
| Service upgrade | Health passes; compatible old MCP processes still work | Passed; new launcher delegates to the selected core |
| Software rollback | Previous core starts; notes/configuration/reader metadata unchanged | Passed |
| Physical vault guard | Second configuration cannot acquire the same vault lock | Passed in Windows Go test |
| Browser restart on a different port | Existing tabs and exact scroll position restored in preview | Passed in actual embedded browser with preview.5 and repeated after upgrade to preview.7 on another port: Quick notes, Other, Resume; scroll 1182.67 CSS px; 0 console errors. No personal note material |
| Authorized main-library migration | Preserve original libraries and verify copied bytes | 36 original/backup files verified; 20 main-library baseline files copied and unchanged after actual preview.5 → preview.7 upgrade |
| Installed launcher delegation | Old managed launcher and direct new exe resolve same home/service | Passed after actual preview.7 upgrade, including direct new exe with no home environment override |
| Codex registration | Shared home; preserve unrelated settings | Official stdio registration verified; old local plugin disabled and retained. Current host must restart before native tool dispatch can be tested |
| Claude Desktop | Official MCPB import and matching Shared Thoughts home | Official update completed in Claude Desktop 4.4212.0; preview.7 installed binary hash matches; shared home saved. Actual host lists all 14 tools with no connection error. Direct installed-entry status checks match Codex entry's core PID/library. Native conversation tool dispatch remains untested |
| Claude Code | Native client MCP dispatch | Not tested: client CLI is not installed on this host |
| Clean Windows / SAC On and Off | Real isolated environments and native clients | Not available; not passed |

The preview.7 isolated upgrade/rollback smoke command completed in approximately 10 seconds. Its first run failed while decoding UTF-8 service metadata through the Windows locale default in the **test harness**. Explicit UTF-8 decoding fixed the harness; the final full run passed. Failure records contain error types and function locations, not note contents or tokens. Tests did not modify Windows security settings.

The source build and immutable package precede this expanded report; package instructions can be older than this working-copy document. Detailed commands, timestamps, private backup paths and local browser screenshots are retained in ignored local artifacts. They must not be included in a public source upload.

Claude Desktop follow-up: the owner explicitly confirmed the exact local package and shared directory before the Windows Install action. The extension was updated through the official GUI, not by editing its installation cache. The default `${HOME}/Thoughts` was replaced with the exact selected absolute home and saved through Claude's confirmation dialog. The initial connection error cleared after the settings took effect; all 14 tools appeared. Existing tool approval settings were not broadened. No prompt was sent to Claude, and no user note body was read, written or transmitted during this configuration check. Original extension metadata/settings remain backed up.

## Integrity and limits

preview.7 executable SHA256: `33863fc467e91955ef5d1c51cd469e8a79e9f015a815988a37bbd9d5d6ce631c`. The executable is **NotSigned**. MCPB SHA256: `91ac24132ac4a5fbdaee9ca609beb08749a96f9c780b9755f86f2e0ece0a165c`. Checksum verification is not publisher authentication or clean-Windows acceptance.

CLI/stdio integration tests do not prove native host dispatch. Browser restoration was observed on preview.5 and repeated on preview.7 after another service/port change. No acceptance claim is extended to other computers, unsigned download trust, host exit, arbitrary versions, crash-proof draft recovery or automatic note history.

For normal shared changes use MCP/web writes with revisions. Other editors can still write directly to Markdown, and their raw filesystem writes cannot participate in the same revision handshake. Do not delete locks to recover from a mismatch; diagnose the service and home first.

The latest preview.8 working-copy source allowlist contains 139 files and had zero configured heuristic findings. This does not certify the entire private working directory or history: the existing private commit still has two files flagged for personal paths/library references. Private archives and local test/configuration material remain excluded. No history was rewritten and no public export or push was performed.

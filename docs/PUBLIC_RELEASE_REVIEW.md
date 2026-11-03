# Public release preparation — 2026-11-01

**Current status — 2026-11-04: SOURCE-ONLY PUBLIC PREVIEW; NO BINARY RELEASE.** The owner authorized importing the reviewed current-source snapshot as a fresh public history. The original private repository/history remains preserved. Current locally installed core is the alignment/eight-language `0.3.0-preview.20`, executable SHA256 `df959a3a34d051d8587f586f9418d2f9b6c13df17a2c1af28b3b9d791678bfda`. Preview.19 remains preserved. Older sections below record earlier preparation states, not a prohibition on this authorized source import.

## Current Windows testing and end-user installation

The Windows development Go executable remains blocked on this machine, but the Linux-container Go toolchain successfully builds Windows executables. The reviewed current executable has been activated locally with SAC still enabled. Frontend checks passed 25 tests, Go checks in Linux passed 49 with ten explicit Windows-only skips, and two notice-generation tests cover eight languages. The previously skipped ten tests have now all passed natively on Windows, with zero skips; see [the test list and scope](WINDOWS_NATIVE_TESTS_20261104.md). This is not a claim that the complete native suite or all three real AI clients have been accepted.

Release-preparation follow-up: all 59 backend tests passed natively on Windows, with no failures/skips. An offline current-user installer is implemented. An earlier installer passed isolated install/reinstall/data-retention/hash/busy-editor/upgrade/rollback/uninstall tests, but **SAC blocked the final unsigned installer with Windows error 4551 and Code Integrity events 3033/3077**. Final installer/shortcut acceptance is incomplete; earlier hashes do not certify it. Current ZIP/MCPB/setup are review artifacts only. See [release readiness](RELEASE_READINESS_20261104.md).

This repository publishes the allowlisted current source, with no Windows binary Release assets. The [installation quick start](../INSTALL.md#ordinary-windows-users) distinguishes future Setup.exe/portable downloads from GitHub source ZIPs. Trusted signing, real publisher-page/contact verification and native-client/clean-machine acceptance remain separate work. The source import does not certify or distribute the blocked unsigned installer.

## Historical preview.20 source work — Windows build was blocked, not yet installed

Follow-up: only real note tabs remain; utilities open without tabs. Five complete local privacy/use/support notices and per-language settings links are generated from one reviewed text source. 24 frontend and 2 notice-generation tests plus TypeScript/Vite passed. No public privacy URL is invented, no release metadata is represented as verified, and no newly built Windows core is claimed.

Restore independent left-side note trees; remove the Quick notes/project overview UI and migrate legacy overview tabs to selected-note reading. Preserve sidebar Trash actions and add menus to search/archive result rows. Frontend tests and TypeScript/Vite checks are recorded separately. Windows Code Integrity events 3033/3077 block the development Go compiler, so no preview.20 executable/package, Go-test pass, installed upgrade or full UI acceptance is claimed. The installed/local packaged core remains preview.19. No security settings or original Git history changed. See [workspace report](PROJECT_WORKSPACE_POLISH.md).

## Collection width and vector web logo — 0.3.0-preview.19

Collections use the full available workspace width with responsive gutters; document reading width remains bounded. Web branding/favicon now use the existing SVG source, synchronized by the build scripts. Native host assets retain PNG compatibility. Browser visual acceptance remains pending. See [workspace report](PROJECT_WORKSPACE_POLISH.md). No public operation or history change.

## Workspace responsive layout — 0.3.0-preview.18

Remove redundant collection paper icons. Consolidate document/list spacing into a shared reading measure and container queries that respond to workspace width, including sidebar collapse. Narrow rows retain dates on their own line; controls wrap and long content remains bounded. Browser viewport screenshots remain unverified under the existing policy block. Details and verification limits: [workspace report](PROJECT_WORKSPACE_POLISH.md). No publication or Git history operations.

## Parent navigation and duplicate heading cleanup — 0.3.0-preview.17

Remove generated document/collection heading repetition, add save-guarded parent collections and persist validated project tabs/scroll. Quick notes overview no longer duplicates its list in the sidebar; homonymous tabs show location. Frontend 21/21, Go 58 top-level tests, builds and package checks passed. Actual isolated service restart preserved parent-view metadata and Markdown hashes; main activation preserved 33 baseline data files. Browser/native-host acceptance remains pending. See [workspace report](PROJECT_WORKSPACE_POLISH.md). No publication or Git history operations.

## Sidebar scroll follow-up — 0.3.0-preview.16

Group branding/search into one opaque, non-scrolling header and constrain navigation to a separate scroll viewport, preventing group shrink and keeping header/footer space. Replace the repeated stack-of-books project icons with one folder icon beside the Projects section heading. New browser scroll/visual acceptance remains pending under the existing conversation policy block. See [workspace report](PROJECT_WORKSPACE_POLISH.md). No public operations or history changes.

## Sidebar search and pill tabs — 0.3.0-preview.15

Reference-based monochrome SVG navigation icons, content-sized pill tabs, sidebar search and an upper-right content toolbar implemented. Document ellipsis actions moved to sidebar note/context menus, including Quick notes. Existing save/conflict guards preserved. Frontend 19/19, Go 57 top-level tests, build, package integrity and isolated two-MCP/web/upgrade/rollback passed. Main core activation preserved all 33 baseline data-file hashes; installed Codex/Claude launchers resolve the same service. Real browser acceptance remains incomplete under the existing conversation policy block. No new screenshot, native-host acceptance or public operation is claimed. See [workspace report](PROJECT_WORKSPACE_POLISH.md).

## Sidebar overlap fix — 0.3.0-preview.14

Owner screenshot identified a project-label overlap caused by obsolete absolute SVG positioning conflicting with new padding. The old rules were removed and relevant row declarations consolidated. 19 frontend tests, 57 Go tests, TypeScript/Vite and package validation passed; main core updated with 33 data-file hashes unchanged. Post-fix real-browser visual acceptance remains incomplete under the current conversation tool policy. No publication or history changes. See [workspace report](PROJECT_WORKSPACE_POLISH.md).

## Project/workspace follow-up — 0.3.0-preview.12

Project archive and whole-directory recoverable Trash/restore, explicit search scope and compact headers/navigation implemented. 19 frontend and 57 Go tests, TypeScript/Vite, package checks and actual isolated CLI/HTTP/MCP lifecycle/upgrade tests passed. Main core activated with 33 baseline data files hash-verified unchanged. Current conversation browser policy denied the isolated QA page; real new-interface screenshots and responsive/interaction acceptance remain incomplete. No alternate access was attempted. Details: [project/workspace report](PROJECT_WORKSPACE_POLISH.md). No public operations; M0 remains not passed.

## Default shared library follow-up — 0.3.0-preview.8

All new adapters default to the same Windows account Thoughts home; optional auto configuration and read-only paths diagnostics added. Codex-first/Claude-first real-file fixtures, host-environment independence, explicit overrides, legacy guard and shared initialization lock passed. Frontend 17/17, Go 53 top-level tests, TypeScript/Vite, package checks and actual isolated MCP/web/upgrade/rollback passed. The actual shared core was upgraded to preview.8; old installed entries and default/auto entries report the same service/library, with all baseline file hashes unchanged. New preview.8 MCPB GUI import and native Claude Code/Codex dispatch remain unverified. Existing Claude preview.7 package continues using the correct shared absolute home. Details: [shared library acceptance](SHARED_LIBRARY_UPGRADE_ACCEPTANCE.md). No public operations; M0 remains not passed.

## Shared library and upgrade follow-up — 0.3.0-preview.7

Explicit shared home, repeat-init preservation, physical vault guard, reviewed local managed-core upgrades/rollback and server-persisted reading state implemented. Frontend 17/17, 49 top-level Go tests, TypeScript/Vite, archive integrity and actual isolated two-MCP/web/upgrade/rollback tests passed. An authorized user library was backed up, copied with hash verification and upgraded without editing its note contents. Codex stdio registration is configured but host restart/native dispatch remains pending. Claude Desktop official preview.7 MCPB update and shared-home configuration completed; actual host lists 14 tools without a connection error. Direct installed-entry status checks confirm the same core PID/library as the Codex entry; native conversation tool dispatch is still untested. Claude Code is not installed. Details and limits: [shared library acceptance](SHARED_LIBRARY_UPGRADE_ACCEPTANCE.md). Unsigned/clean-Windows gates remain open; M0 has not passed.

## Current UI follow-up — 0.3.0-preview.4

Theme settings reduced to import only; independent note tabs, per-tab close/rename, fixed top-right edit/preview, and bounded proportional image previews. Final frontend 15/15, Go suite, TypeScript/Vite and package checks passed. Real browser tests include immediate switch/save, restored preview tabs, inactive rename and seven tabs at 900px width with active-tab visibility after resize. Details and measured limits: [UI acceptance](UI_TABS_IMAGES_ACCEPTANCE.md). Final exe remains NotSigned; local Defender scan found no matching threat records. Clean Windows and real-client acceptance remain pending; M0 has not passed. New review folder: `Windows测试包/当前版本-0.3.0-preview.4/`. No installed client or actual user library was replaced.

## Earlier local follow-up — 0.3.0-preview.2

- White background/black drawing, same brand artwork, 512px package icon; familiar folder and square-pencil sidebar icons.
- Light appearance only, with legacy dark preference/cache fallback. No notes or language preferences deleted.
- Optional Base mainnet sponsorship address/QR configured by the owner; independent QR decoding and exact clipboard matching passed. No wallet connection, payment SDK, private keys, preset amount/token or paid features.
- Frontend 12/12, TypeScript/Vite, full Go suite (43 top-level tests), Windows x64 build and actual stdio/file integration passed. Browser preview/edit/save/disk consistency/reload passed; external change observed after 424 ms; Chinese settings and 900px window checked, no horizontal overflow or console errors.
- Two newly built archives pass integrity checks and use identical core bytes. Both exe copies and the actual running service are NotSigned; MCPB verification reports unsigned. Defender single-file scan completed without matching threat records.
- See [Windows security assessment](WINDOWS_SECURITY_COMPATIBILITY.md). Clean SAC On/Off, real-client import/tool dispatch and full host lifecycle remain unverified. **M0 has not passed.**
- New packages are under `dist/windows/0.3.0-preview.2/` and copied to `Windows测试包/旧版本/本地打磨-0.3.0-preview.2/`. No installed client package or actual user library was replaced. Prior public-source export is stale; regenerate and re-audit only when publication preparation resumes.

## Historical preparation — 0.3.0-preview.1

The remaining sections record the previous candidate. Its tests, screenshots and signing state do not certify a later candidate.

## Repository and scope

- Workspace confirmed as Thoughts (Go + React/TypeScript/Vite/Tailwind/shadcn/Lucide).
- Original branch `main`, one private commit `77410b8`, original `origin` preserved.
- New target: **https://github.com/xuli82221-pixel/hello-Thoughts-**; verified public and empty. Independent local remote `public`; no replacement of `origin` or force push.
- Windows x64; Claude Desktop MCPB, Claude Code/Codex stdio. Same core executable and explicit shared `THOUGHTS_HOME`. Legacy libraries/configuration are not silently moved.
- Paid-theme price/checkout code removed. No payment, subscription, cloud service, wallet connection or invented sponsorship destination.
- Private legacy plugin/deferred-platform files remain locally for preservation; they are outside the proposed public export and active packages.

## Prepared files

English README, MIT LICENSE, `.gitignore`, SECURITY.md, INSTALL.md for all three clients, CHANGELOG.md, PRIVACY.md, THIRD_PARTY_NOTICES.txt, dependency review, security/history review, source/export scripts and a reviewed isolated screenshot. Complete proposed file list: [PUBLIC_SOURCE_FILES.md](PUBLIC_SOURCE_FILES.md). An ignored review export includes a path/size/SHA256 upload manifest and source ZIP; no Git objects are included.

## Test results

| Check | Actual result |
| --- | --- |
| Frontend tests | 11 passed |
| TypeScript and Vite | Passed; existing chunk-size advisory remains |
| Go full suite | 42 top-level tests passed, no failed/skipped tests |
| Windows x64 build | Passed; actual standalone candidate executed |
| stdio MCP | Actual initialize/list 14 tools, no implicit initialization; passed |
| Real-file integration | Create/read/update/search, stale revisions, rename/move, local image upload/read, Trash/restore and path isolation passed |
| Reinitialization | Same configuration/files retained; conflicting library selection rejected |
| Lifecycle | Closing test MCP leaves service; identity-verified isolated service termination/reopen preserves Markdown |
| Browser | Actual standalone webpage, default preview/edit/save, image and Chinese text, refresh after reload; 0 console errors |
| External Markdown refresh | Observed 273 ms in the isolated browser test |
| Settings | No paid-theme/checkout/wallet UI displayed |
| Codex registration | Real CLI, isolated CODEX_HOME: inspect-only/no write, add, identical repeat, bad SHA rejection, conflicting entry rejection, unrelated setting retained |
| Claude Desktop format | Official `@anthropic-ai/mcpb@2.1.2` validation passed; optional 512px icon recommendation (current icon 256px) |
| Archives | Both passed CRC/safe-path/contents/license/per-file SHA checks and ship identical core bytes |
| Real Claude Desktop import/tool dispatch | **Pending for this candidate**; existing user extension/configuration not replaced |
| Real Claude Code tool dispatch | **Pending; Claude Code CLI not found on this machine** |
| Real Codex host tool dispatch | **Pending**; isolated registration and direct protocol tests do not replace this |
| Clean Windows / full host-exit lifecycle | **Pending** |
| Authenticode | **NotSigned** for executable and registration helper |

The first Go attempt overlapped frontend asset replacement and failed due to removed old embedded asset names. It was rerun after frontend completion. A restricted-sandbox test then could not perform normal NTFS/junction/temp-file operations; all tests passed when explicitly authorized with project-local temporary directories. No Windows security setting was disabled or bypassed. Detailed failure and success logs remain local under `.artifacts` and are not uploaded.

Core tests exercise actual files/HTTP in isolated directories. The standalone smoke test invokes the built executable through real stdio/CLI and reads its disk files. The screenshot contains test text/logo explicitly labelled as acceptance material, not fabricated user notes. No actual user library was edited or migrated.

## Local candidate packages

Generated in `dist/windows/0.3.0-preview.1/` with two packages and SHA256SUMS.txt. The portable ZIP serves Claude Code and Codex; the MCPB serves Claude Desktop. Read index.json for signing and pending-acceptance flags. Copies are in `Windows测试包/旧版本/开源候选-0.3.0-preview.1/` for local review. They are candidates, not a published release. No Node/Go/Python/Docker runtime is needed by these embedded-binary packages.

## Blocking history finding and required decision

The private history contains workstation/private-library references in two historical reports; see [SECURITY_AUDIT.md](SECURITY_AUDIT.md). Recognizable credentials were not found in Git, but public history review **did not pass**. Current redaction alone cannot fix this. All private history remains unchanged.

Recommended: keep the original private repository/history, and use the reviewed current-source export as a fresh first commit in the new public repository. No import has been performed. The owner must explicitly confirm this history strategy and the exact upload inventory before any public push. If the original ancestry must be published, separately authorize cleanup and re-audit; no force push is planned.

Signing, clean Windows and three real-client acceptance remain required before claiming the universal Windows install/use goal is achieved. Publishing source as an explicitly unfinished preview is a separate decision from releasing supported binaries.

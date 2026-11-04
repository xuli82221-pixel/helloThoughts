# Install Thoughts — Windows x64

An explicitly authorized unsigned experimental prerelease is available in [hello-Thoughts-](https://github.com/xuli82221-pixel/hello-Thoughts-). Precompiled packages need no Thoughts account, API key, Node, Go, Python or Docker. Install your AI client separately; its own account requirements still apply.

**Code signing policy:** [Pending — SignPath Foundation application not yet approved](docs/CODE_SIGNING_POLICY.md). All current Windows programs remain unsigned. The experimental preview offers unsigned Setup.exe/ZIP/MCPB assets; no signed stable package is available. The Foundation application has not been submitted by this preparation; a future core signature alone will not satisfy the installer/helper release gates.

## Ordinary Windows users

**Current availability — 2026-11-04:** `0.3.0-preview.21` is an explicitly authorized **unsigned experimental prerelease**. Download assets from its GitHub Release. The public Windows build and checks completed successfully. The executable is `thoughts.exe`. No trusted publisher signature, clean-SAC acceptance or complete native AI-client acceptance is claimed. Ordinary users who need those assurances should wait for a signed stable release. Stop if Windows blocks the files; do not disable protection or bypass warnings.

The preview includes an **offline installer**:

1. From **Releases**, download `thoughts-<version>-windows-amd64-setup.exe` and `SHA256SUMS.txt`. **Code → Download ZIP** is source code, not the application.
2. Complete [Verify first](#verify-first), including the stated **NotSigned** status. A valid publisher signature is a requirement for a future signed stable release, not a property of this preview. Double-click Setup.exe and click **Install**. No administrator prompt, terminal commands or developer tools are required. Setup chooses the Windows language automatically; its eight languages match the application languages.
3. Finish setup and open **Thoughts** from the desktop or Start menu. Setup initializes the shared library and the managed runtime automatically. It does not install or configure an AI client.

The installer is offline and per-user. Program/shortcut files live under `%LOCALAPPDATA%\Programs\Thoughts`; the shared home is `%USERPROFILE%\Thoughts`. Existing library selection, notes and images are retained. A reinstallation updates through the checked runtime mechanism; unsaved editors can block an upgrade until their drafts are saved. No service or autostart registration is added. Opening Thoughts starts the local service when needed.

Use **Windows Settings → Apps → Installed apps → Thoughts → Uninstall** to remove the launcher, shortcuts and installer registration. **Uninstall intentionally retains the shared home, configuration and managed executable versions**, because other MCP clients may still use them. It does not wipe notes/images or stop an unrelated/shared client service. Independent note backups remain the user's responsibility; software rollback is not a note backup.

The final local unsigned installer was actually blocked by SAC with Windows error 4551. It is a review artifact, not an accepted end-user package. Earlier local installer-flow passes do not certify the final file. Stop if Windows blocks it; keep protection enabled. This historical block remains relevant to unsigned previews. Signed stable distribution still requires trusted signing and clean-machine testing; this preview is separately authorized and does not claim to resolve SAC trust.

### Portable ZIP alternative

Download `thoughts-<version>-windows-amd64.zip` and verify its archive checksum. Extract it into a new stable folder such as `C:\Tools\Thoughts`. Open PowerShell in that folder and run the following. Each command stops on an error; no security setting is changed.

```powershell
$ThoughtsExe = (Resolve-Path -LiteralPath '.\thoughts.exe').Path
$ExecutableSha = (Get-Content -LiteralPath '.\checksums.json' -Raw | ConvertFrom-Json).files.'thoughts.exe'
& $ThoughtsExe init --json
if ($LASTEXITCODE) { throw 'Library initialization failed; keep the error and stop.' }
& $ThoughtsExe runtime install --source $ThoughtsExe --sha256 $ExecutableSha --json
if ($LASTEXITCODE) { throw 'Runtime installation failed; keep the error and stop.' }
& $ThoughtsExe open --browser system --json
if ($LASTEXITCODE) { throw 'Thoughts could not open; keep the error for diagnosis.' }
```

The browser opens the local application. The default library is `%USERPROFILE%\Thoughts\library`; repeat initialization preserves an existing configured library. On later launches, run `thoughts.exe open --browser system --json` from the same extracted folder. Keep this launcher folder in place. Windows x64 is the application target; using Linux to build the Windows executable does not make the application a Linux release.

An AI client is optional for reading and editing through the webpage. For AI access, use the matching route below:

| What the user wants | Package / next step |
| --- | --- |
| Local webpage without an AI client | Offline Setup.exe (recommended), or portable ZIP |
| Claude Desktop tools | Reviewed `.mcpb`; follow the Claude Desktop extension instructions below |
| Codex or Claude Code tools | Portable ZIP; register its managed executable using the client-specific instructions below |

All clients should reuse the same default home. Setup.exe does not register MCP automatically; select an AI integration separately using the instructions below. There is no automatic download/update service. Publishing source alone does not create release downloads or establish Windows signing trust.

## Verify first

Download Setup.exe (or the portable ZIP / Claude Desktop MCPB) and `SHA256SUMS.txt` from the same reviewed release. Review the repository owner, release notes and signing status. Do not pipe downloaded scripts into a shell.

For each package, use its exact downloaded filename:

```powershell
$ArchiveName = Read-Host 'Exact downloaded Setup.exe, ZIP or MCPB filename'
$Matches = @(Get-Content -LiteralPath '.\SHA256SUMS.txt' | Where-Object { $_ -match ('^[a-fA-F0-9]{64}  ' + [regex]::Escape($ArchiveName) + '$') })
if ($Matches.Count -ne 1) { throw 'Missing or ambiguous checksum' }
$Expected = $Matches[0].Substring(0,64).ToLowerInvariant()
if ((Get-FileHash -LiteralPath $ArchiveName -Algorithm SHA256).Hash.ToLowerInvariant() -ne $Expected) { throw 'SHA256 mismatch. Stop.' }
```

For Setup.exe, inspect **Properties → Digital Signatures** before running it, or use `Get-AuthenticodeSignature` on the exact file. This preview is expected to report `NotSigned` and has no trusted publisher to authenticate. Future signed releases must match their documented publisher. A hash is not publisher authentication. For the ZIP alternative, extract into a new stable directory and inspect the helper, configuration examples, security documentation and extracted `thoughts.exe` signature. Unsigned candidates may be blocked. Do not disable protection, change execution policy or bypass warnings.

## One shared home

From preview.8, all three adapters default to the **same Windows account profile's `Thoughts` directory**. Installation order, current project, `HOME` and host-specific `APPDATA` do not select another library. Claude Desktop's optional advanced setting defaults to `auto`; no folder selection is required. Empty/unset `THOUGHTS_HOME` and `auto` resolve to the same home. A retained literal `${HOME}/Thoughts` from an older extension is also recognized as that default.

```powershell
$ThoughtsHome = Join-Path ([Environment]::GetFolderPath('UserProfile')) 'Thoughts'
$ThoughtsExe = 'C:\Tools\Thoughts\thoughts.exe' # replace with the real extracted path
& $ThoughtsExe paths --json # read-only, before initialization too
```

The home holds `config/` and the default `library/`. Use automatic selection for all clients. An explicit absolute `THOUGHTS_HOME` remains an advanced override for deliberately using another library; use the same override in every client that should share it. Keep home settings on updates/reinstalls. Avoid temporary/cloud-synced folders. Different Windows accounts or PCs are separate; this is not cloud synchronization.

Older Memo users: inspect and back up the existing library. No automatic migration/merge happens. Do not point a new home at existing material casually; a different home means a separate library.

For a deliberately customized home, use the ordinary long-form local path. DOS
8.3 aliases and reparse/junction aliases are not accepted compatibility targets
for the managed core and may be refused by its strict path checks. Keep the
automatic shared home for ordinary installation.

When the common default is absent but the known old Windows-user `AppData/Roaming/Memo/config.json` exists, automatic initialization stops for migration instead of creating an empty replacement. Explicit `MEMO_CONFIG_DIR` remains available for legacy/testing use. Old Store-app-specific or custom locations still require deliberate review; they are not searched or merged automatically. Prefer one migrated shared home rather than leaving different clients on legacy overrides.

Repeat `init` preserves an existing configured library. Supplying a different `--vault` is rejected rather than replacing that selection. For a deliberate migration, stop writes, back up configuration and all library files, copy to the stable destination, verify file hashes, then configure all clients to use the same home. Keep the original until verification is complete. Updating an executable or reinstalling an extension must not remove this home.

## Migrating an earlier local development installation

This is the first public binary preview. Earlier local-development MCP launchers must be deliberately re-registered to the new package's `thoughts.exe` when migrating. An old launcher may reject the newly named runtime pin. Save drafts, back up the library and configuration, close the clients, and review the existing managed-update procedure before changing a shared installation. The release does not update any installed client or core automatically. Old runtime files and backups are preserved; reading legacy pins does not imply every old launcher understands a new filename.

## Managed core updates and rollback

This preview supports reviewed **local** updates; it does not download updates automatically. Use one managed core so extension upgrades and separately registered clients do not start competing core versions. Keep any registered launcher executable in place: it delegates to the checked current pin.

After verifying the extracted package and `checksums.json`, select the intended home and initialize it (repeat initialization preserves an existing selection), then install the managed core before registering it:

```powershell
$env:THOUGHTS_HOME = $ThoughtsHome
& $ThoughtsExe init --json
$ExeSha = (Get-Content -LiteralPath 'C:\Tools\Thoughts\checksums.json' -Raw | ConvertFrom-Json).files.'thoughts.exe'
& $ThoughtsExe runtime install --source $ThoughtsExe --sha256 $ExeSha --json
& $ThoughtsExe runtime status --json
```

Use the returned installed executable path as `$ThoughtsExe` for Claude Code/Codex registration below. MCPB launchers use the same pin when their Shared Thoughts home matches. Initialize only after the user has selected the intended home; runtime installation requires that existing library configuration.

For a later reviewed package extracted into a **new** directory, verify its archive and executable checksums first:

```powershell
$CandidateExe = 'C:\Tools\Thoughts-new\thoughts.exe'
$CandidateSha = (Get-Content -LiteralPath 'C:\Tools\Thoughts-new\checksums.json' -Raw | ConvertFrom-Json).files.'thoughts.exe'
& $ThoughtsExe runtime upgrade --source $CandidateExe --sha256 $CandidateSha --json
& $ThoughtsExe doctor --json
# Explicitly return to the previous software version when needed:
& $ThoughtsExe runtime rollback --json
```

The command verifies the candidate, requests shutdown only from the authenticated compatible service, switches the pin and waits for health. A reported busy editor rejects the switch; save or resolve the draft first. A startup failure attempts to restore the previous pin. Unknown or incompatible service identities are never forcibly killed. Retain the old executable, home and independent library backup. Rollback changes **software only**, not Markdown or image contents, and is not a note revision system.

Opening from every host must report the same physical `library` path and service PID. A second home aimed at an already served vault is refused; choose the same home rather than deleting locks. Reopen using `open` after a restart because the port and session may change. Saved tabs and scroll positions are loaded from the home even in a fresh browser origin. Unsaved drafts remain in their original tab and are not guaranteed after closing it.

An updated MCPB may still delegate to the previously pinned core. Check `runtime status` and `doctor`; importing an extension is not itself a core upgrade. Restart/reconnect a host when it keeps an old MCP process. Never delete the shared home as part of an upgrade or uninstall unless the user separately chooses to delete their data.

## Claude Desktop — official MCPB format

1. Verify the MCPB SHA256. Its ZIP-compatible contents include `manifest.json`, `server/thoughts.exe` and notices; inspect before installing.
2. Go to **Settings → Desktop app → Extensions → Manage**, or **Settings → Extensions** in older versions.
3. Use **Advanced settings → Install Extension…**, select the reviewed `.mcpb` and review local access.
4. Leave **Shared Thoughts home (advanced)** on `auto`. For an existing deliberately customized home, retain its absolute path instead. Check `thoughts_status` for the actual library.
5. Check local MCP connection/tools in Extensions or the conversation's Connectors menu. Refresh/restart the client only if required.
6. Run the real installation test below. An installed icon alone is not successful tool execution.

Connecting alone creates no notes or window. The binary stdio server needs no additional runtime. [Official installation instructions](https://support.claude.com/en/articles/10949351-getting-started-with-local-mcp-servers-on-claude-desktop), [MCPB specification](https://github.com/modelcontextprotocol/mcpb/blob/main/MANIFEST.md).

## Claude Code — local stdio

Inspect `claude mcp get thoughts` before adding if that name might exist. Do not replace an unrelated entry.

```powershell
claude mcp add --transport stdio --scope user thoughts --env "THOUGHTS_HOME=$ThoughtsHome" -- "$ThoughtsExe" mcp
claude mcp get thoughts
```

Check `/mcp` inside Claude Code, then run the test below. User scope is cross-project. The executable is invoked directly; no `cmd /c`, `npx` or remote server is needed. The equivalent configuration example is in `adapters/stdio/claude-code.json.example`; merge only this entry. [Official Claude Code documentation](https://code.claude.com/docs/en/mcp).

With a preview.8+ launcher, `--env "THOUGHTS_HOME=auto"` can replace the absolute home above. Ask Claude Code: “Use Thoughts MCP. Check thoughts_status first and confirm the library path; then read the note I name.” Do not create a per-project library or copy the notes into the code repository. If `claude` is not installed, the installation command is only a prepared instruction, not a completed client connection test. A VS Code AI extension must separately support local stdio MCP; the same executable/args/environment can be used in that extension's supported configuration.

## OpenAI Codex — local stdio

Inspect `codex mcp get thoughts --json` first if the name might exist.

```powershell
codex mcp add thoughts --env "THOUGHTS_HOME=$ThoughtsHome" -- "$ThoughtsExe" mcp
codex mcp get thoughts --json
codex mcp list
```

The official CLI registers a named entry; do not replace the entire `~/.codex/config.toml`. `adapters/stdio/codex.toml.example` shows the table-only form. Check `/mcp` and actually call tools. This is not a ChatGPT web plugin. [Official OpenAI documentation](https://developers.openai.com/codex/mcp/).

## Inspect-first helper

After verifying/extracting the archive, read its embedded executable checksum:

```powershell
$ExeSha = (Get-Content -LiteralPath '.\checksums.json' -Raw | ConvertFrom-Json).files.'thoughts.exe'
.\register-mcp.ps1 -Client codex -Executable $ThoughtsExe -ExpectedSha256 $ExeSha -ThoughtsHome $ThoughtsHome
```

The default only displays the plan. Review it, then repeat with `-Apply`; choose `-Client claude-code` for Claude Code. It verifies SHA256 and refuses a conflicting named registration. It does not initialize/migrate files, download software or change security settings. It delegates to the official client CLI. A registration success message is not a handshake/tool-call test.

`-ThoughtsHome` is optional from preview.8; omitting it resolves the Windows account's shared default and registers that absolute path. Supplying another path is an explicit advanced override. The helper does not use the active project or a host-specific `HOME` to choose the library.

You may ask your assistant to download the exact reviewed assets, check SHA256, explain the commands and perform the approved registration while preserving unrelated settings. Never give it permission to weaken system protection.

## Real installation test

Ask the client: “Initialize Thoughts in the configured shared home. Create `Inbox/acceptance-can-delete.md` as a clearly labelled installation test, read it back, update using the current revision, then open the local webpage.” Verify the actual Markdown file and read the web-saved result back through the client.

Compare `thoughts_status`/library paths in all clients and read the same file. Test search, rename, move, images, recoverable Trash and restore using only designated acceptance files. Directly edit the file externally and check webpage refresh. Reject stale revisions; preserve drafts on conflicts.

`thoughts_open` with `browser: "none"` returns a one-use opening URL for an embedded browser; `browser: "system"` opens the default browser. Never publish the ticket or copy session cookies between browsers. Bare localhost URLs do not authorize a new browser.

## Open without an AI client

```powershell
$env:THOUGHTS_HOME = $ThoughtsHome
& $ThoughtsExe doctor --json
# First use only, when creating this selected local library is authorized:
& $ThoughtsExe init --json
& $ThoughtsExe open --browser system --json
```

Later use `open`, without reinitializing. Port/session can change after restart; reopen through the program. Test Windows restart and data persistence in the actual environment. Installation or a protocol test does not replace real-client and clean-Windows acceptance.

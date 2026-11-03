# Security policy

## Reporting

Email `hello@spxipo.com` privately with the affected version, relative source locations and a minimal reproduction using an isolated library. Never attach personal notes, opening links, runtime/session files, credentials or keys. Do not post private data or exploitable details in public issues. No response-time or bounty promise is made.

## Scope and boundaries

- Windows x64 preview only; no production-audit or signing claim. No macOS, Linux, mobile, remote MCP or ordinary ChatGPT web plugin support.
- The stdio adapter inherits its client's OS permissions. It exposes specific note-management tools, not arbitrary shell commands or a general filesystem reader.
- Public setup uses one absolute local `THOUGHTS_HOME` across clients, containing `config/` and a default `library/`. Legacy configurations remain compatible. No automatic migration or merging occurs.
- HTTP binds only IPv4 `127.0.0.1` on an assigned port. Host/Origin and cross-site checks, single-use expiring tickets, HttpOnly SameSite cookies and CSRF protect browser access.
- Opening tickets/private runtime records are credentials. Never share them. A bare localhost URL does not transfer another browser's session.
- `os.Root`, relative-path validation and regular-file checks confine note/image access, rejecting traversal and symlink/junction escapes. Images are decoded and bounded. Theme packages accept JSON tokens and validated previews, not scripts. Remote images and active raw HTML are not rendered.
- SHA256 revisions reject stale updates. Preserve drafts on conflict; never blindly retry using a fresh revision. Trash is recoverable, but automatic revision history is absent. Back up the library independently.
- Log code intentionally omits bodies and tokens, but logs can contain paths/timing/errors. Review before sharing. Files are not encrypted. Malware or another process with the same user's OS permissions is outside the browser-session threat model.

## Networking and AI

There is no note upload endpoint, analytics, cloud database or automatic updater. Runtime API calls are loopback; explicit opening launches the browser with a local URL. Users can follow links in notes or support/privacy links. An AI host receives contents requested through tools and may send them to its provider. Its behavior is outside Thoughts' local storage guarantee.

## Lifecycle and upgrades

The webpage service is separate from stdio and browser tabs. Reuse checks library, instance identity and exact version. Never delete a live lock or terminate an unknown process to force an upgrade. Save drafts before stopping a verified service. OS restart stops the process; authorized reopen starts it again. Files persist separately. Full client-exit/clean-Windows acceptance must be recorded independently.

## Publication gate

See the [Code signing policy](docs/CODE_SIGNING_POLICY.md). SignPath Foundation approval is pending; all current Windows programs remain unsigned. MFA for GitHub and SignPath is mandatory for future signing members but was not verified in this preparation. No trusted signing integration or approved binary release is claimed.

Review current source, full Git history, assets, archives and licenses. A clean current tree does not sanitize old commits. Public upload stops on sensitive findings; history is never rewritten without explicit owner authorization. SHA256 is integrity checking, not Authenticode trust. Never disable Windows protection or bypass security warnings for a preview.

# Historical source/history security audit

The original private history remains excluded from the public repository. Current [preview.21 release status](RELEASE_READINESS_20261104.md) is public unsigned prerelease; its source/run/checksums are documented in [unsigned preview evidence](UNSIGNED_PREVIEW.md). Nothing in the earlier audit below overrides that release or certifies it as signed, stable or independently security-audited.

<!-- historical-records:start -->

## Historical records — source import and earlier candidates

## Reviewed source publication

The target is `https://github.com/xuli82221-pixel/hello-Thoughts-.git`. On 2026-11-04 the owner authorized the reviewed current-source snapshot to be published with a fresh public history. The repository was verified empty before import. The original private repository and `origin` are preserved. That source-import operation excluded binaries, original private history, personal libraries and credentials; the later preview.21 binary release was a separate explicitly authorized operation.

## History findings

The latest read-only reachable/reflog scan examined **1 commit, 190 blobs**, with findings in the same two historical reports below. An earlier read-only `git fsck --full --no-reflogs --unreachable` check reported no unreachable objects. Heuristic credential checks did not find a recognizable API key, token, password marker or private key in reachable history. This is not proof that every possible password or mnemonic is absent.

The old commit contains workstation paths/private library references in:

| File | Historical lines | Risk |
| --- | --- | --- |
| `docs/DATA_RELIABILITY_AUDIT.md` | 10–11, 23 | Private library/project references and workstation paths |
| `docs/THOUGHTS_BRAND_RESULTS.json` | 57–59 | Workstation-specific absolute paths |

Matched values are intentionally omitted. These findings block publicly pushing the existing history. Current copies were redacted, and originals preserved in ignored local archives. **Deleting values in the working copy does not sanitize the old commit.** No history rewrite, force push or deletion of old Git objects was performed.

Authorized approach: import the separately reviewed current-source snapshot as the new public repository's first commit, while keeping all original history private. The export contained no `.git` directory; publication uses an independent Git directory. The owner approved this source-only import on 2026-11-04. Publishing identical original ancestry would still require separate history cleanup, authorization and re-audit.

## Working tree and artifacts

2026-11-04 follow-up: a fresh read-only history audit examined one commit and
295 reachable/reflog blobs and still found workstation paths in the same two
historical reports. Existing history remains unsuitable for public push. The
current-source allowlist is checked independently by `release-preflight.py` and
export-time scanning, including the new installer script and explicit CI file.
Raw machine-specific inventories and test logs remain local. A reviewed source
snapshot does not rewrite history or authorize publication.

The .2 working-tree scan, including ignored local artifacts, produced 156 location-only findings, including local runtime/session credentials and workstation paths. These files must remain excluded, not be copied wholesale into a public repository. A separate source-allowlist check is saved per local iteration; neither that check nor UI acceptance replaces a fresh complete pre-publication audit. At that historical audit stage, no new source export, commit or push was performed; source publication and later preview.21 release followed separate reviews and authorizations. Dependency caches and binary/video contents were inventoried rather than treated as reviewed source; unrelated media and original input archives are excluded. No other project was modified.

The proposed source is an explicit allowlist: core, frontend, three client adapters, build/test scripts, publisher metadata, licenses and selected public documentation. It excludes real libraries/images, original screenshots, chats, `.env`, logs, runtime records, caches, user configuration, compiled binaries, private legacy plugin packaging and historical local reports. The only public screenshot is newly captured from a clearly labelled isolated acceptance library, using the candidate's actual standalone webpage. Source-export scanning runs before copying and reports locations/rules only.

Both Windows candidate archives were checked for safe paths, duplicate entries, CRC, per-file/package SHA256, licensing notices and identical core executables. No library/session/config/log files are bundled. `-trimpath` is used; package inspection checks that the development account's profile path is absent from executable bytes. These checks cannot establish complete binary provenance or detect every possible encoded secret.

## Code review boundaries

- No note-upload, analytics, telemetry, cloud database or remote-listener implementation was found in the reviewed runtime. The HTTP server binds `127.0.0.1`; local clients reject redirects/unverified instances.
- MCP returns explicitly requested data to the AI host. This is intentionally documented and is not a promise that the host never forwards contents to its provider.
- Note/image paths use validated relative paths, `os.Root`, regular-file and size/format checks. Tests cover traversal, Windows junction/symlink escape, stale writes and authorization.
- Browser sessions use one-use tickets, origin/host checks and CSRF. Raw Markdown HTML and remote images are not active. Theme imports contain bounded validated data rather than scripts.
- Runtime process execution is limited to starting the same local executable and opening its locally constructed URL with the system browser. No arbitrary command tool was found. Development tests can run explicit local tooling; they are not exposed over MCP/HTTP.
- No wallet connection, key collection or payment processing is present. Paid-theme code/strings have been removed. An owner-configured Base mainnet sponsorship address and locally generated address-only QR are displayed optionally, with no preset amount/token. The generator validates the EIP-55 checksum and independently decodes the image; it makes no RPC or third-party QR request.

The review is a bounded source/package inspection, not an independent penetration test or formal security certification. Same-user malware, OS compromise and AI-host policy remain outside the local browser-authentication boundary. See SECURITY.md for reporting and limitations.

## Licensing

The existing MIT project license is retained. Production dependency/Go/font license texts and shadcn source attribution are included in THIRD_PARTY_NOTICES.txt. The full npm metadata inventory found MIT, OFL-1.1, ISC, Apache-2.0, CC-BY-4.0, MPL-2.0, BSD-3-Clause and 0BSD; MPL/CC entries are build-only. No unknown/UNLICENSED/GPL/AGPL metadata was found. Dependencies retain their licenses and are not relicensed by the project MIT file. See DEPENDENCIES.md; this is a technical review, not legal advice.

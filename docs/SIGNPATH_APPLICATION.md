# SignPath Foundation application preparation

Prepared: 2026-11-04. **Draft for maintainer review; not submitted.**

## Project facts

| Field | Draft answer / evidence |
| --- | --- |
| Project name | Thoughts |
| Public source / project homepage | https://github.com/xuli82221-pixel/hello-Thoughts- |
| Description | Free, local-first Markdown workspace for AI assistants on Windows. The Go core embeds a React webpage for local reading/editing; standard stdio MCP lets Claude Desktop, Claude Code and OpenAI Codex use the same local library. |
| Functionality | Markdown create/read/edit/search/rename/move, local images, project navigation/archive, recoverable Trash, SHA256 stale-write protection, local theme import, shared core updates/rollback and saved reading state. UI and local notices support eight languages. No built-in AI, Thoughts account, cloud note upload, subscriptions or paid features. See [English baseline](../README.md). |
| License | [MIT](../LICENSE) for project code. Dependencies/fonts and Inno Setup retain their own licenses; see [notices](../THIRD_PARTY_NOTICES.txt) and [dependency inventory](DEPENDENCIES.md). Whole-package eligibility still requires review. |
| Stack | Go (`go.mod`), React/TypeScript/Vite, Node.js 24 build tooling, Python packaging tools, Inno Setup 7 installer tooling. Precompiled runtime needs none of these tools. |
| Release target | Windows x64, Windows 10 21H2 or newer. Linux CI/cross-compilation is development infrastructure, not a Linux product release. |
| Expected signing file | First-party `thoughts.exe`, built from `cmd/thoughts`, reused in ZIP/MCPB. Initial application excludes third-party binaries, installer/uninstaller runtime and `register-mcp.ps1`; scope must be confirmed by the Foundation. See [Code signing policy](CODE_SIGNING_POLICY.md#what-we-sign). |
| Current version/status | `0.3.0-preview.21` unsigned experimental prerelease, separately authorized by the maintainer. The release build is from public source and records its exact commit/run. No signed stable package or Foundation approval is claimed. Earlier preview.20 unsigned setup was blocked by SAC (4551). See [unsigned preview](UNSIGNED_PREVIEW.md) for the current release and acceptance boundaries. |
| Download instructions | [INSTALL.md](../INSTALL.md) clearly separates current source access and future reviewed binary installation. There is no invented download URL or standalone product website. |
| Privacy | [PRIVACY.md](../PRIVACY.md). Local files are not encrypted; users must back them up. Authorized tool results go to the selected AI client and may reach its online provider. |
| Code signing policy | [CODE_SIGNING_POLICY.md](CODE_SIGNING_POLICY.md): pending, conditional provider wording, artifact scope, provenance, manual approval, roles, MFA and privacy boundaries. |
| Maintainer/contact | [xuli82221-pixel](https://github.com/xuli82221-pixel); configured public support contact: hello@spxipo.com. Contact ownership/deliverability has not been verified by sending mail. No email was sent. |
| Roles | The collaborator API returned only xuli82221-pixel (`admin`). The same maintainer has author, reviewer and manual approver responsibilities. No independent reviewer or configured SignPath team is claimed. |
| CI/CD | [Source checks](../.github/workflows/ci.yml): Windows/Linux matrix, pinned checkout/setup actions, read-only `contents`, npm lockfile install/tests/build, Go tests/vet, legal/preflight checks and npm audit. Local package entry point: `scripts/build-desktop-adapters.ps1`. No automated signing or release publication exists. A separate manually triggered unsigned-preview build uploads review artifacts without write access or Secrets. |

## Earlier public evidence snapshot (before unsigned-preview authorization)

Authenticated GitHub API inspection on 2026-11-04, before the subsequent unsigned-preview authorization, established:

- Public repository, MIT license metadata, one administrator; the authenticated maintainer has admin/write access.
- Repository created `2026-10-31T23:15:19Z`; 0 stars, 0 forks and 0 subscribers at inspection. These are a dated snapshot, not a forecast or a user count.
- No Releases and no version tags. Source ZIP download does not meet the same claim as an already distributed executable.
- Latest inspected `Source checks` run succeeded for public commit `dc32fa9ca7044901d9a6cf70d2493e8a4852a345`. Earlier failed runs exist. A successful source run is not a signed build or clean-machine/client acceptance.
- Description and nine Topics were updated to the requested values. Homepage remains empty because no separate project website was verified; the repository README serves as the project landing page.

Evidence can be rechecked at [repository metadata](https://api.github.com/repos/xuli82221-pixel/hello-Thoughts-), [Releases](https://github.com/xuli82221-pixel/hello-Thoughts-/releases), [tags](https://github.com/xuli82221-pixel/hello-Thoughts-/tags) and [Actions](https://github.com/xuli82221-pixel/hello-Thoughts-/actions). Permissions were checked through the authenticated collaborator endpoint; public contributor lists alone do not prove signing authority. Recheck live facts before submitting.

## Eligibility assessment

Read against the official [application page](https://signpath.org/apply.html) and [Foundation conditions](https://signpath.org/terms.html), accessed 2026-11-04. The application page embeds a hosted form; this document prepares answers but does not fill or submit it. The conditions page labels its code of conduct a draft; the provider's current decision controls.

| Requirement | Assessment | Remaining action |
| --- | --- | --- |
| Open-source license; no proprietary components or commercial dual licensing | Project code is MIT; no commercial dual-license offering identified. Mixed third-party/font/installer licenses are documented, not automatically Foundation-approved. | Review the full shipped dependency/license inventory, including OFL fonts and the Inno Setup runtime, with the provider. |
| Active maintenance | Recent public source and CI activity support current maintenance. | Continue maintenance; no long-term support guarantee is inferred. |
| Already released in the form to be signed | The earlier source-only snapshot did not meet this requirement. A separately authorized unsigned preview is now being built and published; its actual public assets must be supplied to the Foundation for review. | Ask whether a source-only preview is eligible or must defer. Do not publish an unsigned/unaccepted formal release to satisfy a checkbox. |
| Functionality documented on project/download page | README and INSTALL describe scope, local data boundaries and future package routes. | Publish this reviewed documentation manually; an approved binary download remains unavailable. |
| Verifiable reputation and project control | Control is supported by authenticated admin permission. **Reputation is not established:** very young repository and 0-star/0-fork snapshot; no usage evidence or user counts. | Provide only real independent evidence if it develops; Foundation decides sufficiency. No minimum star threshold is invented. |
| Own project / own binaries only | First-party `thoughts.exe` has a documented source path. | Restrict artifact filters to it; resolve bundled runtime/installer scope explicitly. Never sign upstream binaries as Thoughts. |
| No malware, unwanted software, hacking features or privacy compromise | Reviewed documented product scope is local note management, not vulnerability exploitation; no built-in AI or automatic uploads. Prior audit limits remain in [security audit](SECURITY_AUDIT.md). | Maintainer must review actual shipped source/dependencies; this is not a malware-free certification. |
| Announce system changes; provide uninstall | INSTALL describes current-user program files, shortcuts, shared home and data-retaining uninstall; installer source implements the intended route. | Final signed installer/shortcut and clean-machine acceptance are still pending. |
| MFA for all GitHub/SignPath team members | **Unverified.** Repository access does not establish MFA; no SignPath team enrollment was verified. | Maintainer confirms/enables both and configures roles before signing access. |
| Authors, reviewers and approvers assigned | Policy assigns the sole actual maintainer all three responsibilities. | Confirm the single-maintainer arrangement with the Foundation; configure corresponding SignPath roles and manual approval. |
| Verifiable automated source-to-binary build | **Incomplete.** Public source CI passes; local packaging is not a trusted SignPath build integration. The manual preview build records its public commit and run; its version tag identifies the same commit. | Implement reviewed automated release artifacts tied to exact public commit/tag, trusted integration and retained checksums. |
| Enforced product/version metadata | **Not demonstrated.** No core PE resource/metadata enforcement was found in the inspected source/build path. | Implement and verify Foundation product/version restrictions after design review; no product code changed here. |
| Manual approval for every signed release | Required by policy, **not configured** in CI/SignPath. | Configure protected human approval; keep release publication separately gated by acceptance. |
| Public code signing/privacy policy and attribution | Documents prepared, with conditional provider attribution and explicit unsigned status. | Maintainer must publish reviewed docs; provider approval remains pending. Do not use an approved-provider badge prematurely. |

**Conclusion:** Thoughts has a plausible open-source/local-first project basis, but Foundation acceptance of the released form, verifiable reputation, trusted SignPath provenance, metadata enforcement, MFA verification and full package signing/acceptance remain unresolved. This is not an eligibility determination or an approval prediction.

## Maintainer submission checklist

- [ ] Review all facts, licenses and contact ownership; update dated GitHub evidence.
- [ ] Review and manually commit/push only the intended public documentation/source snapshot; preserve private historical material. This preparation performs neither action.
- [ ] Decide whether to ask the Foundation about the source-only/reputation gaps now or defer the application. Do not claim an existing binary release.
- [ ] Confirm the first-party-only scope, installer/runtime questions and one-person role arrangement with SignPath.
- [ ] Confirm GitHub MFA, enroll in SignPath and enable MFA there; never share recovery codes.
- [ ] Manually submit at the official application page and retain the actual response privately.
- [ ] If approved, design/review trusted CI, tags, artifact metadata, manual approvals and least-privilege credentials before adding any Secrets.
- [ ] Complete signature/checksum, installer/uninstaller, clean-Windows and native-client acceptance, then obtain separate human release approval.

No application was submitted, no credentials were created, no signing Secrets were added, and no binaries were published by this preparation.

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
| Download instructions | [Published unsigned preview](https://github.com/xuli82221-pixel/hello-Thoughts-/releases/tag/v0.3.0-preview.21) provides Setup.exe/ZIP/MCPB and four evidence/checksum files. [INSTALL.md](../INSTALL.md) explains current installation, unsigned status and future signed-stable requirements. No separate product website is invented. |
| Privacy | [PRIVACY.md](../PRIVACY.md). Local files are not encrypted; users must back them up. Authorized tool results go to the selected AI client and may reach its online provider. |
| Code signing policy | [CODE_SIGNING_POLICY.md](CODE_SIGNING_POLICY.md): pending, conditional provider wording, artifact scope, provenance, manual approval, roles, MFA and privacy boundaries. |
| Maintainer/contact | [xuli82221-pixel](https://github.com/xuli82221-pixel); configured public support contact: hello@spxipo.com. Contact ownership/deliverability has not been verified by sending mail. No email was sent. |
| Roles | The collaborator API returned only xuli82221-pixel (`admin`). The same maintainer has author, reviewer and manual approver responsibilities. No independent reviewer or configured SignPath team is claimed. |
| CI/CD | [Source checks](../.github/workflows/ci.yml): Windows/Linux matrix, pinned checkout/setup actions, read-only `contents`, npm lockfile install/tests/build, Go tests/vet, legal/preflight checks and npm audit. Local package entry point: `scripts/build-desktop-adapters.ps1`. There is no automated signing. Neither source checks nor the manual unsigned-preview build publishes Releases automatically; preview.21 was uploaded and published separately with maintainer authorization. The unsigned-preview workflow uploads build artifacts without repository write access or added Secrets. |

## Verified public release evidence — 2026-11-04

- [v0.3.0-preview.21](https://github.com/xuli82221-pixel/hello-Thoughts-/releases/tag/v0.3.0-preview.21) is public, not a draft, and explicitly marked prerelease. It was published at 2026-11-04 11:00 Asia/Shanghai and contains seven assets: Setup.exe, portable ZIP, Claude Desktop MCPB, SHA256SUMS, BUILDINFO, ACCEPTANCE and UNSIGNED_NOTICE.
- The version tag points to public build commit `0b155f8a4d0381315f742bb6661e2e0eb6c738c4`. [Windows build and bounded acceptance run](https://github.com/xuli82221-pixel/hello-Thoughts-/actions/runs/41762640507) succeeded. Uploaded asset digests match the reviewed CI outputs. These facts establish a published unsigned build, not Foundation provenance verification or signing trust.
- GitHub metadata at this review showed a public MIT repository created `2026-10-31T23:15:19Z`, 1 star, 0 forks and 0 subscribers. These are a dated snapshot, not user counts or proof of sufficient reputation. The earlier 0-star/source-only inspection preceded this release and is not the current release status.
- The latest authenticated collaborator inspection returned one administrator, `xuli82221-pixel`; no independent reviewer or configured SignPath team is inferred. Recheck permissions and MFA before submission.
- Description and nine Topics are configured. Homepage remains empty because no separate project website was verified; the repository README serves as the project landing page.

Evidence can be rechecked at [repository metadata](https://api.github.com/repos/xuli82221-pixel/hello-Thoughts-), [Releases](https://github.com/xuli82221-pixel/hello-Thoughts-/releases), [tags](https://github.com/xuli82221-pixel/hello-Thoughts-/tags) and [Actions](https://github.com/xuli82221-pixel/hello-Thoughts-/actions). Recheck live facts before submitting.

## Eligibility assessment

Read against the official [application page](https://signpath.org/apply.html) and [Foundation conditions](https://signpath.org/terms.html), accessed 2026-11-04. The application page embeds a hosted form; this document prepares answers but does not fill or submit it. The conditions page labels its code of conduct a draft; the provider's current decision controls.

| Requirement | Assessment | Remaining action |
| --- | --- | --- |
| Open-source license; no proprietary components or commercial dual licensing | Project code is MIT; no commercial dual-license offering identified. Mixed third-party/font/installer licenses are documented, not automatically Foundation-approved. | Review the full shipped dependency/license inventory, including OFL fonts and the Inno Setup runtime, with the provider. |
| Active maintenance | Recent public source and CI activity support current maintenance. | Continue maintenance; no long-term support guarantee is inferred. |
| Already released in the form to be signed | The earlier source-only snapshot did not meet this requirement. A separately authorized [unsigned preview is now public](https://github.com/xuli82221-pixel/hello-Thoughts-/releases/tag/v0.3.0-preview.21), with Setup.exe/ZIP/MCPB and source-linked build evidence. The Foundation must decide whether this prerelease meets its released-form requirement. | Supply the actual prerelease URL, assets, tag and build evidence; ask the Foundation whether this released form is eligible. Do not claim a signed stable release. |
| Functionality documented on project/download page | Published README and INSTALL describe scope, local data boundaries and the actual unsigned package routes. | Keep the published download instructions accurate; trusted signatures and stable-release acceptance remain pending. |
| Verifiable reputation and project control | Control is supported by authenticated admin permission. **Reputation is not established:** very young repository and the dated 1-star/0-fork snapshot; no independent usage evidence or user counts. | Provide only real independent evidence if it develops; Foundation decides sufficiency. No minimum star threshold is invented. |
| Own project / own binaries only | First-party `thoughts.exe` has a documented source path. | Restrict artifact filters to it; resolve bundled runtime/installer scope explicitly. Never sign upstream binaries as Thoughts. |
| No malware, unwanted software, hacking features or privacy compromise | Reviewed documented product scope is local note management, not vulnerability exploitation; no built-in AI or automatic uploads. Prior audit limits remain in [security audit](SECURITY_AUDIT.md). | Maintainer must review actual shipped source/dependencies; this is not a malware-free certification. |
| Announce system changes; provide uninstall | INSTALL describes current-user program files, shortcuts, shared home and data-retaining uninstall; installer source implements the intended route. | Final signed installer/shortcut and clean-machine acceptance are still pending. |
| MFA for all GitHub/SignPath team members | **Unverified.** Repository access does not establish MFA; no SignPath team enrollment was verified. | Maintainer confirms/enables both and configures roles before signing access. |
| Authors, reviewers and approvers assigned | Policy assigns the sole actual maintainer all three responsibilities. | Confirm the single-maintainer arrangement with the Foundation; configure corresponding SignPath roles and manual approval. |
| Verifiable automated source-to-binary build | Public automated Windows build and commit/tag/run mapping are completed for preview.21. **SignPath trusted build integration is incomplete.** The public workflow and checksum records are evidence for review, not Foundation origin verification. | Present the existing public automated build and retained checksums; configure the Foundation-approved trusted build/signing integration after approval. |
| Enforced product/version metadata | **Not demonstrated.** No core PE resource/metadata enforcement was found in the inspected source/build path. | Implement and verify Foundation product/version restrictions after design review; no product code changed here. |
| Manual approval for every signed release | Required by policy, **not configured** in CI/SignPath. | Configure protected human approval; keep release publication separately gated by acceptance. |
| Public code signing/privacy policy and attribution | Documents are public, with conditional provider attribution and explicit unsigned status. | Provider approval remains pending; keep the public policies current. Do not use an approved-provider badge prematurely. |

**Conclusion:** Thoughts has a plausible open-source/local-first project basis, but Foundation acceptance of the released form, verifiable reputation, trusted SignPath provenance, metadata enforcement, MFA verification and full package signing/acceptance remain unresolved. This is not an eligibility determination or an approval prediction.

## Maintainer submission checklist

- [ ] Review all facts, licenses and contact ownership; update dated GitHub evidence.
- [x] Publish the reviewed source, documentation and explicitly authorized unsigned preview; preserve original private history and user material. Recheck the linked release evidence before applying.
- [ ] Decide whether to ask the Foundation to review the published unsigned prerelease and reputation evidence now or defer the application. Use the actual unsigned prerelease and its build evidence; do not claim a signed or stable release.
- [ ] Confirm the first-party-only scope, installer/runtime questions and one-person role arrangement with SignPath.
- [ ] Confirm GitHub MFA, enroll in SignPath and enable MFA there; never share recovery codes.
- [ ] Manually submit at the official application page and retain the actual response privately.
- [ ] If approved, design/review trusted CI, tags, artifact metadata, manual approvals and least-privilege credentials before adding any Secrets.
- [ ] Complete signature/checksum, installer/uninstaller, clean-Windows and native-client acceptance, then obtain separate human release approval.

The unsigned prerelease and documentation are public. No SignPath application was submitted by this work, no signing credentials or Secrets were created, and no trusted signature or stable release is claimed.

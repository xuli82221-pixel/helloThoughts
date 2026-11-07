# Code signing policy

Updated: 2026-11-04. Scope: Thoughts, Windows x64, unsigned preview `0.3.0-preview.21`.

## Signing status

**Pending — SignPath Foundation application not yet approved.**

The application has not been submitted by this preparation. All current Thoughts Windows programs remain unsigned. The maintainer has separately authorized an explicitly labeled unsigned experimental prerelease; there is no signed stable release, trusted publisher claim, SignPath certification or signed download. Source publication and passing tests do not confer signing trust. Do not disable SAC, bypass SmartScreen or use a self-signed certificate to imply trusted signing.

## Signing provider

**If approved, free code signing provided by SignPath.io, certificate by SignPath Foundation.**

This is a conditional future arrangement, not a statement that a subscription or certificate has been granted. Provider references: [SignPath.io](https://about.signpath.io/), [SignPath Foundation](https://signpath.org/), [conditions](https://signpath.org/terms.html), [application](https://signpath.org/apply.html). The Foundation decides eligibility and the approved artifact configuration.

## What we sign

The proposed initial Foundation scope is **`thoughts.exe`**, the first-party Windows x64 executable built from `cmd/thoughts` and the Thoughts source, including its embedded React UI. The same executable is included in the portable ZIP and Claude Desktop MCPB. New executable names and source paths use Thoughts; legacy stored identifiers remain for data/runtime compatibility.

Inspection of `scripts/build-desktop-adapters.ps1`, `scripts/pack-desktop-adapters.py` and `release/windows/Thoughts.iss` confirms these package outputs:

| Output | Proposed treatment |
| --- | --- |
| `thoughts.exe` | Proposed first-party executable; pending Foundation review and source verification |
| `thoughts-<version>-windows-amd64.zip` | Container with the core; retain package SHA256, not an Authenticode executable |
| Claude Desktop `.mcpb` | Container with the same core; core signing does not sign the extension/container |
| `thoughts-<version>-windows-amd64-setup.exe` and generated uninstaller | Use the upstream Inno Setup runtime; excluded from the initial first-party-only Foundation scope pending explicit provider review of package/runtime rules and a compliant signing design |
| `register-mcp.ps1` | First-party script, not an executable; excluded from this initial application scope |

Never sign another project, an upstream executable/DLL, build tools, AI clients or third-party software using Thoughts' signing identity. Dependencies retain their notices and licenses. Shipping a third-party runtime inside an installer does not make that runtime Thoughts-authored source.

The current local `-PublicRelease` build gate requires trusted signatures for the core, registration script and Setup.exe. Signing only `thoughts.exe` will **not** satisfy that gate. The generated uninstaller also needs release-time signature verification. The installer/helper signing path must be resolved before a signed stable distribution. A separate explicitly authorized `-UnsignedPreview` mode produced the public preview.21 experimental packages; it does not satisfy or disable the signed `-PublicRelease` gate. See [published unsigned preview](UNSIGNED_PREVIEW.md).

## Build provenance

The source of record is the [public GitHub repository](https://github.com/xuli82221-pixel/helloThoughts). For any future signing request:

- Build from the public, reviewed source in an isolated, verifiable automated CI run. Include the checked-in workflow, build scripts, dependency lockfiles and third-party notices in review. Never upload an arbitrary local binary for signing as if it had verified CI provenance.
- Record the exact Git commit, matching version tag, build/run identity and toolchain versions. All binaries in a build must carry the enforced product name `Thoughts` and consistent product version.
- Configure SignPath's trusted build integration, artifact filters and metadata restrictions. Do not give pull requests or untrusted code access to signing credentials or approval authority.
- A named approver must manually approve **every** signing request/release after reviewing source, provenance and acceptance evidence. No unattended signing or publishing.
- Keep SHA256 checksums for unsigned inputs, signed outputs and final packages, and preserve the commit/tag/run mapping and signature/timestamp verification results. Signing changes executable bytes; recompute downstream checksums after signing.
- Verify the final signed artifacts and packaged copies, then complete clean-Windows and native-client acceptance before human release approval.

**Current implementation:** `.github/workflows/ci.yml` performs read-only source checks on Windows and Linux. The source-check workflow does not upload release artifacts, request signing or publish releases. The manual unsigned-preview workflow uploads bounded build artifacts for maintainer review and has no release-publishing step. Local packaging/checksum scripts exist, but are not SignPath provenance verification. The unsigned prerelease is tied to its exact public build commit and version tag. Enforced PE product/version metadata, protected manual signing approval and the SignPath integration are pending. The separate manual unsigned-preview workflow records commit/run provenance; that is not Foundation provenance verification. No Secrets were added during this preparation.

## Team roles

The authenticated GitHub collaborator API on 2026-11-04 returned one account: [xuli82221-pixel](https://github.com/xuli82221-pixel), repository owner and `admin`. No other collaborators were returned. This records repository permissions; it does not establish SignPath membership or MFA.

| Role | Responsible maintainer | Responsibility |
| --- | --- | --- |
| Author / committer | xuli82221-pixel | Maintain source, dependencies, build scripts and workflows |
| Reviewer | xuli82221-pixel | Review outside contributions and changes affecting signing/security; inspect the release diff and CI evidence |
| Approver | xuli82221-pixel | Manually decide each signing request and release after acceptance |

These are the responsibilities assigned by this policy for the current single-maintainer project, not a claim of independent two-person review or already configured SignPath roles. Recheck access and update this table before onboarding another signer.

## MFA

Every signing team member **must enable MFA for both GitHub and SignPath** before receiving signing access. GitHub account MFA and SignPath enrollment/MFA were not verified in this review; neither is claimed complete. The maintainer must confirm both and configure the corresponding roles before signing is enabled. Do not publish recovery codes or authentication evidence containing credentials.

## Privacy

See [PRIVACY.md](../PRIVACY.md). Thoughts stores Markdown and local images in the user's local library. It has no Thoughts account, cloud note database, automatic note upload, telemetry, built-in AI, subscription or paid feature. The local reader binds to `127.0.0.1`.

When a user authorizes Claude Desktop, Claude Code, OpenAI Codex or another online AI client to read notes, tool results are delivered to that client and may be sent to its provider under its own policy and account settings. Thoughts' local storage guarantee does not cover those external copies. External links and optional wallet use have their own data boundaries. Files are not encrypted by Thoughts; independent backups remain the user's responsibility. Installation/system changes and data-retaining uninstall are described in [INSTALL.md](../INSTALL.md).

## Pending checklist

- [ ] Present the explicitly authorized unsigned prerelease to the Foundation for released-form review, and resolve the verifiable-reputation gap; no Foundation acceptance is implied.
- [ ] Review all bundled component/font/installer licenses against the Foundation's component requirements.
- [ ] Maintainer reviews [application materials](SIGNPATH_APPLICATION.md), verifies contact details and submits the official application manually.
- [ ] Receive Foundation approval and confirmation of artifact scope; approval is not guaranteed.
- [ ] Verify GitHub and SignPath MFA; configure real members, roles and manual approvals.
- [ ] Review the existing public commit/tag-based unsigned build and implement SignPath trusted provenance integration, PE metadata restrictions and least-privilege signing access.
- [ ] Resolve installer, uninstaller and registration-script signing outside the initial core-only scope; retain existing release gates.
- [ ] Validate signatures, timestamp chains, final package checksums, clean Windows/SAC behavior and all target clients.
- [ ] Obtain separate human approval for a signed stable release after its signing and acceptance checks; the explicitly authorized unsigned prerelease is already public.

The unsigned preview is public. This policy does not establish Foundation approval, certificate issuance, CI signing integration, signed-artifact acceptance or a signed stable release.

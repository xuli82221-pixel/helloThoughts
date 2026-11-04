# Historical preview.20 native Windows tests — 2026-11-04

This is a version-scoped historical test report. Current published preview.21 status and remaining acceptance requirements are in [release readiness](RELEASE_READINESS_20261104.md). Its local-only statements, test counts, hashes and earlier language counts apply to the named earlier versions.

<!-- historical-records:start -->

## Historical records / 历史记录

For preview.20, the ten tests previously skipped by the Linux build container were run using the compiled Windows AMD64 test executable on Windows. **All ten top-level tests passed, with zero failures, zero skips and process exit code 0.** The first test also passed both installation-order subtests.

The compiled Go source snapshot was compared with the current source before execution and matched. Test executable SHA256: `c7414e29a819eaa8c2b65c7d00bf3c4193d7515a50b6d7e213881feb8f370561`.

| Test | Behavior checked | Result |
| --- | --- | --- |
| `TestDefaultHomeIndependentOfHostAndInstallationOrder` | Codex-first and Claude-first setup resolve the same library despite host-specific HOME/APPDATA; later clients read the file saved by the first | PASS |
| `TestDefaultDoesNotSilentlyAbandonLegacyConfiguration` | Legacy configuration prevents silently creating an empty replacement; an explicit approved home is respected and legacy bytes remain intact | PASS |
| `TestPathsIsReadOnlyAndInitUsesOneSharedLock` | Path inspection creates no files, concurrent initialization is refused and successful initialization persists the shared configuration | PASS |
| `TestInitReinstallKeepsConfiguredLibraryAndRejectsSilentSwitch` | Reinstallation preserves the configured library and existing note; requesting a different library is rejected before creating it | PASS |
| `TestMCPSetupPreservesConfigurationAndExistingFiles` | MCP setup keeps existing configuration and note bytes, refuses switching libraries and does not replace corrupt configuration silently | PASS |
| `TestMCPUsesSameServiceAndRefusesStaleWrites` | MCP and the local HTTP service share the test store; stale revisions, escaping paths and malformed writes are rejected | PASS |
| `TestSharedHomeAcrossHostProfiles` | Different host-local configuration locations still use one explicitly selected home and retain the existing note | PASS |
| `TestServiceLock` | The Windows service lock excludes a second holder and is reusable after a clean close | PASS |
| `TestInitIdempotent` | Repeated initialization retains an existing Markdown file unchanged | PASS |
| `TestVaultLockRejectsASecondConfigHome` | A second store cannot acquire the same physical vault lock; the lock is available after the first holder closes | PASS |

The runner used a separate working directory and dedicated TEMP/TMP folders. Selected tests use `t.TempDir()` and the default-home tests substitute a temporary profile. The real library was not selected for test writes. SAC stayed enabled; no Windows security policy was changed. No installed runtime, preview.19 artifact or personal note/image was replaced by the test run.

Follow-up during release preparation: the complete compiled backend suite was also run natively on Windows: **59 top-level tests passed, zero failures and zero skips**. The full-suite harness supplied its required `testdata` fixtures and left host-selection overrides unset; tests select their own temporary profiles/configurations. An initial harness attempt used an unnecessary global home override and omitted fixtures; those harness errors were corrected before the passing run. No product assertion was removed or relaxed.

These results do not replace clean-PC installation, restart behavior, native Codex/Claude Desktop/Claude Code integration, signing or public-release acceptance. Local raw logs and machine-specific paths are kept outside the proposed public source. See [release readiness](RELEASE_READINESS_20261104.md).

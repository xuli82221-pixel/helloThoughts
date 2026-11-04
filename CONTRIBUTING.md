# Contributing to Thoughts

Thoughts targets Windows x64 (Windows 10 21H2 or newer). Go embeds the React build;
precompiled end-user downloads do not require development tools. Other runtime
platforms are not yet supported. Report security issues privately using SECURITY.md.

## Build and test

Use Go matching `go.mod`, Node.js 24, Python 3.12+ and, for Setup.exe, the digitally
signed official Inno Setup 7 compiler from <https://jrsoftware.org/isdl.php>.
Install the frontend's exact lockfile dependencies, then build the web assets
before Go tests or compilation, because `internal/webui` embeds those assets:

```text
cd web
npm ci
npm test
npm run build
cd ..
go test ./... -count=1
go vet ./...
python scripts/test_legal_notices.py
python -m unittest discover -s scripts -p "test_release_preflight.py"
python scripts/release-preflight.py
```

On Windows, the normal full package entry point is documented in README.md and
requires `-InnoCompiler`. It produces an offline Setup.exe, portable ZIP, MCPB,
checksums and an explicit pending-publication report in a new directory. Existing
outputs are never deleted. Use a trusted isolated toolchain if SAC rejects your
local compiler; never disable SAC or alter execution policy for a build.

Linux can test the platform-independent Go code and cross-compile with
`GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o thoughts.exe ./cmd/thoughts`.
Linux skips Windows-only assertions; native Windows testing remains required.
The Windows installer compiler and actual installation tests require Windows.

Windows canonical-path tests require a normal long-form temporary directory.
If TEMP/TMP contains a DOS 8.3 alias, select a dedicated existing long-form test
directory for those variables before running Go tests; CI uses `runner.temp`.
This changes test fixtures only. Do not relax path checks, change security policy
or treat a passing test as acceptance of arbitrary short-path/reparse aliases.

## Data and review boundaries

- Use new temporary homes and clearly labelled synthetic fixtures for testing.
- Do not commit actual libraries, screenshots of personal notes, tokens, logs,
  certificates, private keys, caches or machine-specific configuration.
- Preserve stale-write checks, path isolation, default-home compatibility and
  unsaved-editor protection. Never silently switch libraries or delete old cores.
- Uninstall only removes the program launcher, shortcuts and uninstall entry.
  The shared home and managed versions remain for user data and other MCP clients.
- Never publish the private repository's old history. A reviewed source export
  contains only the current allowlist; preparing it does not authorize a push.

CI checks Windows and Linux source changes with read-only GitHub permissions.
It does not publish releases. Signing, clean-machine acceptance and native-client
acceptance must be recorded independently; a green CI run is not release approval.

Follow the [Code signing policy](docs/CODE_SIGNING_POLICY.md) and consult the
[SignPath application draft](docs/SIGNPATH_APPLICATION.md) before release work.
The Foundation application is not approved; GitHub/SignPath MFA, trusted build
provenance, artifact metadata and manual signing approvals remain to be verified
or configured. Documentation and translations use README.md as the English
content baseline; preserve commands, paths, privacy boundaries and unsigned status.

For documentation changes, run `python scripts/check-docs.py` to check all eight
READMEs, relative links and anchors, and inclusion of linked files in the public
source export. Review translation meaning manually as well; the script cannot
judge linguistic accuracy or certify external services.

The separately authorized experimental preview uses the manual read-only
[unsigned build workflow](.github/workflows/unsigned-preview.yml). It records
commit/run provenance and unsigned status without signing or publishing a release.
See [preview boundaries](docs/UNSIGNED_PREVIEW.md); signed stable-release gates
remain unchanged.

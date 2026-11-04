# Dependency and license review

Reviewed local dependency metadata on 2026-11-01. This is a technical inventory, not a legal opinion.

2026-11-04 follow-up: the offline Setup.exe uses the unmodified Inno Setup 7.1.1 runtime.
Its original copyright/about information is retained; its license is included in
`licenses/InnoSetup-LICENSE.txt`, the installed license folder and third-party notices.
The official compiler download and compiler were checked for valid Pyrsys B.V. signatures.
Compiler tooling is a build dependency; end users do not install it.
Fresh npm audit reported zero known vulnerabilities; Go checks are recorded in the
[release readiness report](RELEASE_READINESS_20261104.md). No new npm/Go runtime dependency was added.

Clean-checkout CI follow-up: Node.js 24 type declarations are now explicit locked
development dependencies (`@types/node` and its MIT-licensed `undici-types` dependency).
They type-check the existing Node test files and are not bundled into the application.
CI uses the locally reviewed Go 1.27.3 patch rather than the minimum go.mod patch.

- Project code remains MIT licensed with the existing credited author.
- 337 npm lock entries reviewed; 193 production entries. Windows-inapplicable optional packages not installed locally were checked using lockfile metadata.
- Detected licenses: 0BSD, Apache-2.0, BSD-3-Clause, CC-BY-4.0, ISC, MIT, MPL-2.0, OFL-1.1.
- No GPL/AGPL, unknown or UNLICENSED metadata was found by the scanner. Metadata alone is not a comprehensive compliance proof.
- MPL-2.0 and CC-BY-4.0 entries are listed below. They are not relicensed under MIT. Build-only tooling is not bundled into the runtime; production notices retain shipped third-party terms.
- Embedded font files retain OFL-1.1 terms and reserved-name conditions. Fonts have not been renamed or sold separately.
- Go runtime/dependency license texts and production JavaScript/font notices are reproduced in THIRD_PARTY_NOTICES.txt and included in the portable/MCPB archives and installed notice files.
- The WebP test fixture retains its separate WEBP_LICENSE.
- Vendored shadcn/ui-derived component code must retain the upstream MIT attribution in THIRD_PARTY_NOTICES.txt.

| Dependency | License | Build only |
| --- | --- | --- |
| caniuse-lite 1.0.30001819 | CC-BY-4.0 | True |
| lightningcss 1.32.0 | MPL-2.0 | True |
| lightningcss-android-arm64 1.32.0 | MPL-2.0 | True |
| lightningcss-darwin-arm64 1.32.0 | MPL-2.0 | True |
| lightningcss-darwin-x64 1.32.0 | MPL-2.0 | True |
| lightningcss-freebsd-x64 1.32.0 | MPL-2.0 | True |
| lightningcss-linux-arm-gnueabihf 1.32.0 | MPL-2.0 | True |
| lightningcss-linux-arm64-gnu 1.32.0 | MPL-2.0 | True |
| lightningcss-linux-arm64-musl 1.32.0 | MPL-2.0 | True |
| lightningcss-linux-x64-gnu 1.32.0 | MPL-2.0 | True |
| lightningcss-linux-x64-musl 1.32.0 | MPL-2.0 | True |
| lightningcss-win32-arm64-msvc 1.32.0 | MPL-2.0 | True |
| lightningcss-win32-x64-msvc 1.32.0 | MPL-2.0 | True |

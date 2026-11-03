# Project management and workspace polish — 2026-11-03

Current local candidate/core: **0.3.0-preview.19**, Windows x64. Local only; no commit, push or public release. **M0 remains not passed.**

## Tree-first navigation correction — preview.20 source only, NOT installed

### Note-only tabs and five-language notices follow-up

Visible tabs now contain only real notes. Search, archive management and Trash open as utility views with a close control returning to the selected note or an empty pane; they do not add tabs. Restoring older metadata drops all non-note tabs while preserving note selection and utility page state. Both Quick notes and Projects toggle their left trees without routing to an overview. Backend metadata validation continues to recognize old values for compatibility, not as an exposed navigation feature.

All five UI languages now have full local privacy and use/license/support notices: English, Korean, Japanese, Simplified and Traditional Chinese. Each covers local files/preferences/sessions; absence of telemetry/account/cloud upload; AI-client/provider boundaries; external links; retention, Trash and backups; lack of encryption and other local-program access; voluntary support contact; notice changes; MIT/third-party licensing; client permissions; and optional Base donations. MIT's original English text is retained and localized summaries identify it as authoritative. No legal compliance certification, verified public policy URL or marketplace acceptance is claimed.

Reviewed text is in `legal/notices.json`, rendered by `scripts/legal_notices.py` and the HTML template. Release metadata generation emits five HTML pages and multilingual `PRIVACY.md`; settings links select the same language as the UI. Pages have static language navigation, local assets and no scripts or remote resources. Contact remains the supplied publisher/email; missing public policy/terms URLs are not fabricated. The optional-support network label is localized too.

Updated verification: **24 frontend tests passed**, including note-only restoration and per-language privacy routing; TypeScript/Vite passed. **2 Python notice checks passed**, covering language/section parity, local language links, original license text and escaping. All five pages are present in the built frontend. Go/package, installed-core and browser acceptance remain blocked/incomplete as below. Existing preview.19 binaries/packages and real notes were not modified. Evidence: tree20-notices-frontend-tests.log, tree20-notices-frontend-build.log and tree20-notices-results.json under ignored artifacts.

### Initial tree correction

The earlier parent-navigation iteration introduced overview pages and tied Quick notes children to the document page. That made the navigation disappear on an overview and was contrary to the intended left-tree/right-document interaction. This is corrected in source: Quick notes, Projects and each project expand independently; changing branches does not change the right pane or collapse unrelated branches. Note selection alone opens the document after the existing save guard. Quick notes has a sidebar creation action and every note retains its sidebar actions, including recoverable Trash.

Quick notes/project overview rendering and the return-to-overview button are removed. Old overview tabs are migrated out of the visible workspace, while real note tabs and the selected note are retained. With no selected note, the right pane stays empty. Search, archive management and Trash remain explicit utility views without tabs; search/archive rows retain the added action menu. Tree preferences are local UI metadata; note contents are not changed. The existing reader-state schema still accepts earlier metadata for compatibility.

Validation: 23 frontend tests passed, including legacy overview migration and independent expansion/reveal behavior. TypeScript and Vite builds run with the installed Node runtime. **Go tests, executable/package build, isolated deletion API smoke, installed-core upgrade and browser UI acceptance are NOT completed for preview.20.** Windows blocked `.tools/go/bin/go.exe`; a direct launch reports “应用程序控制策略已阻止此文件。” Code Integrity Operational recorded events 3033 and 3077 on 2026-11-03. These events establish an enforced application-control block but do not by themselves distinguish SAC from another policy origin. No existing Go compiler was found on PATH or in the two standard install paths checked. No protection, policy or certificate was changed, and no renamed compiler or other bypass was used.

The installed preview.19 core and current downloadable local packages remain unchanged. Source/package metadata preview.20 is not proof of installation. No real user note was deleted, moved, initialized or rewritten. Evidence is under ignored artifacts: preview20-build.log, preview20-build-retry.log, tree20-frontend-tests.log and tree20-frontend-build.log. No public operation was performed.

## Full-width collections and vector web branding — preview.19

Owner screenshots confirm that the shared reading-width cap from preview.18 wastes available space on collection pages. Collections now fill the workspace between the responsive gutters: headers, separators, hover rows and date columns extend together. Only document bodies and empty-state text keep the reading measure. Sidebar expansion/collapse and both narrow container breakpoints remain supported by the same rules.

Navigation/action icons were already SVG (Lucide or local currentColor SVG components). The web brand and favicon previously loaded PNG even though a matching vector source existed. Both now load the local branding SVG; release/branding scripts synchronize that source into web public assets. Host/plugin PNG assets remain supported. There is no new logo design, external asset or image-generation dependency.

These changes are based on the supplied screenshots and source review. New browser viewport screenshots remain pending under the existing conversation policy denial; static/build checks do not establish visual acceptance. Evidence is stored in ignored preview19-build.log and layout19-results.json. No public Git operation is authorized or performed.

preview.19: 21 frontend tests, 58 top-level Go tests, TypeScript/Vite, Windows x64 build and both package-integrity checks passed. The embedded SVG bytes match the canonical source. Main activation succeeded; all 33 baseline data files retained their SHA256 values, and both installed entry executables resolve the same preview.19 service/library. No additional broad runtime smoke was needed for these layout/static-asset changes; browser/native-host acceptance is not claimed.

## Workspace-based responsive layout — preview.18

Remove repeated paper icons from note collection/search rows; retain meaningful tab/navigation icons. Existing document/collection gutter rules were spread across overlapping viewport media queries and later overrides. Replace those layout declarations with one left-aligned content policy: text measure capped at 48rem plus side padding, rather than expanding with a collapsed sidebar. Both the reader and collection pages use the same gutter and maximum width.

The workspace is now an inline-size query container. Available content width at or below 700px uses 24px gutters; at or below 450px uses 16px gutters, with 32px as the wider default. Narrow note rows place dates on a second row instead of hiding them or squeezing the preview. Collection headers/media controls wrap; long parent labels shrink; table/code blocks retain local horizontal scrolling and image frames retain their ratio while shrinking. Theme packages cannot change these layout variables. These dimensions describe source rules, not browser measurements.

The supplied screenshots differ in visible note content and sidebar state; browser zoom can also change effective CSS width. Identical line breaks across differing zoom or available width are not promised. The target is consistent rules and readable text in both browsers. Browser screenshot/viewport acceptance is still incomplete under the existing conversation policy denial; no alternate access route is used. Existing test/build, package and managed-activation evidence is recorded in ignored preview18-build.log and layout18-results.json. This visual change adds no dependencies or data operations.

preview.18 verification: 21 frontend and 58 top-level Go tests, TypeScript/Vite, Windows x64 build and ZIP/MCPB integrity passed. Managed activation succeeded, all 33 baseline main-library data hashes remained unchanged, and both installed entry executables report the same preview.18 service. Existing protocol smoke was not repeated for these layout-only changes. No new visual, browser-click or native-host acceptance is claimed.

## One document identity and real parent navigation — preview.17

Generated document titles and the repeated project/date row are removed. The tab names the active note; the existing Markdown reading-body behavior and actual files are unchanged. Save status and edit remain at the upper-right of the content region. Modification time is available in the save-status tooltip. At the upper-left, a save-guarded parent button returns a quick note to Quick notes or a project note to its actual project collection. A project collection returns to Projects or Archived projects, according to its current archive state. This is hierarchical navigation, not browser-history navigation.

Projects is now a real collection of local project folders, reachable from its sidebar section label. Project collections show their real notes and allow creation. Empty UI-created projects land in their collection after the library refresh completes. Independent note tabs remain open when returning to a parent. Project tabs and their scroll positions persist through the existing local reader-state endpoint; project paths are validated, and removed project tabs are pruned without discarding conflicted note drafts. Earlier previews may not restore these new project-view metadata values after a software rollback; Markdown files are unaffected.

Collection pages no longer repeat a large title already shown by their tab; accessible headings remain. Quick notes overview no longer duplicates its central list in the sidebar. Sidebar note actions remain available while reading a quick note. Homonymous note tabs display their owning project/Quick notes context.

Validation: 21 frontend tests, 58 top-level Go tests, TypeScript/Vite, Windows x64 build and both package-integrity checks passed. An actual isolated library test created Markdown files, saved project collection tabs/page/scroll, stopped and reopened the service, verified restored state, rejected an escaping parent path and verified unchanged Markdown hashes. Main activation preserved all 33 baseline data-file hashes; both installed entry executables resolve the same preview.17 service/library. Evidence: preview17-final-build.log, parent17-acceptance/results.json and layout17-results.json in ignored artifacts. An initial build is preserved separately and was not installed. No public operation was performed.

New browser screenshots and click/scroll visual acceptance remain incomplete under this conversation's existing browser-policy denial. No alternate access route was attempted. API, persistence and source checks do not constitute visual or native-host acceptance.

## Sidebar scroll boundary and project icon simplification — preview.16

The owner's preview.15 screenshot showed scrolling navigation peeking immediately below the search field and an overly repetitive stack-of-books icon. Branding and search now share a non-scrolling, opaque sidebar header with 12px bottom padding and a separator. The navigation has its own bounded scroll viewport and painting context; groups do not shrink into each other. The footer also keeps its own space. This replaces the search field's former standalone sibling placement.

Only the Projects section heading has a simple folder-outline SVG, inheriting the existing color. Individual projects have their disclosure arrow, name and hover menu; the repeated stack icons are removed. Quick notes and settings retain their existing open-book and gear icons.

The previous conversation browser denial remains in effect. This change is based on the supplied screenshot and source layout review; new real-browser scroll/visual acceptance is pending. No alternate browser or port is used to bypass that denial. The checks and activation outcome are recorded separately in ignored local artifacts.

preview.16 passed the existing 19 frontend and 57 top-level Go tests, TypeScript/Vite, Windows x64 build and both package-integrity checks. Managed activation succeeded; all 33 baseline main-library data files retained identical SHA256 values. Direct calls through the installed Codex and Claude entry executables report the same preview.16 service/library. The prior preview.15 isolated protocol/upgrade smoke remains the last such run; it was not repeated for these layout/icon-only changes. No new browser or native-host acceptance is claimed. Evidence: preview16-build.log, layout16-results.json and private upgrade/hash records under ignored artifacts.

## Pill tabs, navigation icons and content controls — preview.15

The sidebar uses local currentColor SVG outlines based on the owner's references: stacked books for projects, an open book for Quick notes, and a gear for settings. These inherit existing theme colors. The Thoughts brand is unchanged.

Tabs fit their labels inside 32px pills with a 6px gap and a 200px maximum width (170/150px on narrower windows). Fixed tab widths, vertical separators and the active underline were removed. Separate note tabs, close actions, keyboard navigation, automatic active-tab visibility and the overflow menu remain.

Library search is a field directly below the sidebar brand, outside the scrolling project tree. Focusing it opens the search results page after the existing save guard. Input updates search results there; the existing archive-scope checkbox remains. Ctrl+K opens the sidebar and focuses/selects this same field, including after sidebar collapse. It does not add a second toolbar search control.

Save state and edit/preview are in a compact strip at the upper-right edge of the content region, directly below the tabs and outside the scrolling document. They are not in the title row or constrained to the document's text width. Normal saved state includes a check and localized text; pending, failure, conflict and disconnection remain explicit. Viewing remains the default.

The document ellipsis menu is removed. Sidebar note menus and right-click menus offer rename, move, copy content, copy absolute file path and recoverable Trash. Quick notes shows its note rows when active; the selected archived note is also available under Archived projects. Copying the active note retains its current draft; copying an inactive note reads its real content from the local service rather than the content-free library listing.

preview.15: 19 frontend tests, 57 top-level Go tests, TypeScript/Vite, Windows x64 build and both package-integrity checks passed. Actual isolated two-MCP/web/CAS/reinstall/upgrade/rollback tests passed. Main activation retained rollback and all 33 baseline data-file hashes. Direct status calls through both installed launchers resolve the same preview.15 service and shared library. No real note body was logged or edited by this operation. No additional dependency was added.

**Visual acceptance is incomplete.** This conversation's existing browser-policy denial has not been bypassed via another browser or port. No new screenshot or measured visual/interaction pass is claimed. Native-host dispatch, clean Windows and unsigned-download gates remain separate. Local evidence: preview15-build.log, preview15-runtime-smoke.log, runtime-upgrade-preview15/results.json and layout15-results.json under ignored artifacts.

## Sidebar label regression — preview.14

The owner screenshot of preview.12 showed project chevrons/folders overlapping the project names. This was a CSS regression: an older rule absolutely positioned the SVGs while later padding changes moved the text into the same area. It was not evidence of Markdown encoding corruption. The obsolete positioning rules were removed and the project toggle, child-row and tree-spacing declarations consolidated into single normal-flow flex rules.

preview.14 build, 19 frontend tests, 57 Go tests and package integrity passed. The installed core was updated with rollback retained; all 33 baseline data-file hashes remained unchanged. Both installed entries report the same preview.14 service. No encoding conversion, note edit or path migration was performed. Actual post-fix browser visual verification remains incomplete because this conversation's browser access is blocked. Details are in the ignored local layout14-fix-results.json.

## Resulting behavior

- Sidebar branding and the workspace tabs use a shared 48px header. Collapsed navigation shows only its expand button, with no duplicate Thoughts logo in the tab strip. The sidebar is 224px wide.
- Project rows are smaller. A hovered/focused menu provides add note, archive and move project to Trash. Opening a project focuses the expanded tree on that project. The reader's existing note tabs remain independent. An open-tabs menu helps reach overflowed tabs without a permanent scrollbar.
- The content header has one note title plus a compact project/time line. Save state and edit/preview occupy the content region's upper-right strip. Collection rows show one preview line and omit repeated filenames in Quick notes.
- New UI-created projects are empty by default; starter nodes are an explicit checkbox. The CLI/MCP creation defaults are retained for compatibility.
- Deleted external note tabs are removed, while a missing note with an unresolved draft keeps its recovery tab. A recovery copy falls back to Inbox when its project is gone, rather than recreating the deleted project automatically.

## Project data rules

Archive writes a small `.thoughts-project.json` marker within the project; it does not move or rewrite Markdown or images. Archived projects leave the main project tree and appear under Archived projects. Unarchive restores their sidebar visibility. The active project snapshot includes every regular file's content and relative path, not just Markdown, so attachment changes invalidate a stale project action.

Project deletion is recoverable. The complete directory moves to `.memo/project-trash/<random-id>`, with synced metadata written first. No permanent delete endpoint was added. Restore refuses an occupied original path and preserves both projects. Folder links/reparse points are refused; traversal, container paths, malformed metadata and unsupported files fail safely. Regular files, empty directories and relative image references are kept together. A project action is limited to 10,000 files and 1 GiB. External non-cooperating filesystem writers cannot participate in a fully atomic revision handshake.

Lifecycle endpoints use the same authenticated loopback service, Origin and CSRF controls as note writes. These project actions are exposed in the reader/HTTP API; the existing 14 MCP tools are unchanged in this iteration.

## Search scope

The sidebar field is labelled **Search library**. It searches case-insensitive text in note filenames, Markdown headings and body text in active projects and Quick notes. Trash is excluded. Archived projects are excluded by default and can be explicitly included with the checkbox. This is not a whole-computer search or an image/OCR search. Ctrl+F remains the browser's current-page find action.

## Earlier lifecycle validation (preview.12)

| Check | Result |
| --- | --- |
| Frontend tests | 19 passed, including archive reading-state restoration and missing-tab draft recovery |
| Go tests | 57 top-level tests passed, including four new project lifecycle/security tests |
| TypeScript/Vite and Windows x64 build | Passed; pre-existing bundle-size advisory remains |
| ZIP/MCPB integrity | Passed; same core bytes, CRC/path/per-file SHA256 verified |
| Actual isolated CLI/HTTP lifecycle smoke | Archive search filtering, external-attachment revision conflict, whole-project Trash/restore, collision refusal and Markdown/image hash preservation passed |
| Actual isolated two-MCP/web/upgrade/rollback smoke | Passed on preview.12; not a native host acceptance test |
| Authorized main-core activation | preview.8 → preview.12 succeeded, managed rollback retained. 33 pre-upgrade main-library data files unchanged by SHA256; runtime locks excluded from hashing |
| Installed Codex/Claude entry checks | Direct status calls from both installed launchers resolve the same preview.12 PID and shared library; no note bodies read |
| Real browser visual/interaction acceptance | **Not completed.** The preview.9 isolated QA page was denied by the current conversation's browser tool. No alternative browser/port was used to bypass that refusal; preview.12 was not retried through another origin. No new screenshot is claimed |

Detailed records are under ignored local artifacts: preview12-build.log, layout12-acceptance/results.json, runtime-upgrade-preview12/results.json, layout12-main-upgrade-integrity.json and layout12-main-upgrade-integrity.json. Intermediate preview.9 artifacts remain local. These records and private session files are not public source.

The main core is installed, but the Claude extension package still has its existing launcher version and delegates through the managed pin. Open the installed Thoughts entry again after the service upgrade to get a fresh browser session; refreshing an old service address cannot create one. The user reported successful opening in a new conversation; the exact old-conversation host permission/cache cause remains unconfirmed.

Native host lifecycle, new-interface screenshots, exact alignment/responsive measurements, unsigned download compatibility and clean SAC On/Off environments remain outstanding. Authenticode is still NotSigned. Compilation and protocol tests do not satisfy those gates.

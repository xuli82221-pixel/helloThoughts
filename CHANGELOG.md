# Changelog

## 0.3.0-preview.21 — unsigned experimental prerelease

- Standardize the executable, source directories, Go module, frontend helpers, CSS and client examples as Thoughts (`thoughts.exe`). Preserve legacy data, protocol identifiers and immutable runtime pins for compatibility and rollback.
- Add an explicitly authorized manual unsigned-preview build with commit/run provenance, checksums and bounded synthetic acceptance evidence. Signed stable-release gates remain unchanged.
- Publish only as a clearly labeled unsigned prerelease; no Foundation approval, clean-SAC or complete native-client acceptance is claimed. No note-management feature is added.


## Documentation preparation — 2026-11-04 (unreleased)

- Add complete eight-language READMEs with shared navigation and a dedicated Donate page using the existing maintainer configuration.
- Prepare a pending Code signing policy and factual SignPath application assessment; no submission, approval, signing or binary release.
- Clarify source-only availability and future installation routes. Product version remains `0.3.0-preview.20`; no product feature or core logic changes.
- Replace the README illustration with an actual preview.20 capture from a new synthetic library; preserve the earlier image and record screenshot provenance.

## 0.3.0-preview.20 — public source preview; no binary release

- Restore the navigation tree as the primary note browser. Quick notes, Projects and each project expand/collapse independently of the right pane; selecting a note opens its content without collapsing other branches.
- Remove Quick notes/project overview pages and their return-to-overview action. Migrate old overview tabs to the last selected note or an empty reading pane, preserving note tabs and utility views. Existing reader-state schemas remain readable for compatibility.
- Keep only note tabs in the visible workspace. Search, archive management and Trash remain closable utility views without creating tabs; legacy non-note tabs are removed on restoration.
- Align navigation disclosure/icon/text/action columns and content toolbars across the interface, including narrow widths and actual scrollbar space.
- Add German, French and Spanish; bundle privacy, use/license, support and optional-donation notices in all eight UI languages.
- Add a sidebar Quick notes creation action; keep rename/move/copy/Trash in every note's sidebar menu. Search/archive result rows also expose hover/focus actions with revision checks and the existing delete confirmation.
- Build the Windows core through the Linux Go container while leaving SAC enabled; activate the reviewed core locally and retain preview.19. All 59 native Windows backend tests pass.
- Add an offline per-user Setup.exe with eight installer languages, shared-home preservation, managed upgrades/rollback and desktop/Start-menu shortcuts. Uninstall retains shared data/configuration/managed versions for other clients.
- Add source preflight gates, read-only Windows/Linux CI and contributor instructions. Local unsigned candidates and a fresh source review snapshot do not authorize public upload.
- SAC blocked the final unsigned Setup.exe (error 4551); earlier installer-flow passes do not certify the final file. Keep it review-only pending trusted signing and final acceptance.

## 0.3.0-preview.19 — local only, unreleased

- Let note/project/search collections fill the available workspace width, retaining responsive gutters. Keep the reading measure for document bodies only.
- Use the existing SVG brand source for the web logo and favicon; synchronize it during branding generation and release builds. Preserve PNG assets for packaged host integrations.

## 0.3.0-preview.18 — local only, unreleased

- Remove redundant paper icons from Quick notes, project-note and search lists.
- Replace accumulated viewport-based document/list spacing rules with shared, left-aligned reading width and workspace container queries. Cap text measure at 48rem; adjust gutters using available content width, including sidebar collapse.
- Reflow narrow note rows with visible dates below the text; allow collection headers and media controls to wrap. Keep code/table overflow local and image previews proportional.
- Constrain long parent labels to preserve the top-right save/edit controls.

## 0.3.0-preview.17 — local only, unreleased

- Remove generated document titles and repeated location/date rows. Tabs carry document identity; modification time is available on the save-status tooltip.
- Add save-guarded parent navigation from a quick note to Quick notes, from a project note to its project list, and from that project to Projects or Archived projects. Parent collections have real-file lists and restorable tabs/scroll positions.
- Remove duplicate large collection headings and the sidebar's duplicate Quick notes list while its overview is open. Keep sidebar note actions while reading a note.
- Distinguish homonymous note tabs using their owning project. Validate persisted project locations and discard missing project tabs while preserving conflicted note drafts.

## 0.3.0-preview.16 — local only, unreleased

- Keep branding and search inside one solid, non-scrolling sidebar header with bottom spacing and a separator. Bound the project scroll viewport, isolate its painting, and prevent navigation groups from shrinking into one another.
- Replace the repeated book-stack project icons with one simple folder outline beside the Projects section title. Individual project rows retain only their disclosure arrow, name and hover menu.

## 0.3.0-preview.15 — local only, unreleased

- Use reference-based, currentColor SVG book-stack, open-book and gear icons for projects, Quick notes and settings.
- Fit independent note tabs to their labels in compact pills, retaining close buttons, keyboard navigation and overflow selection.
- Move library search into the sidebar below the brand. Ctrl+K expands navigation and focuses the same field; results and scope controls remain in the content pane.
- Put save status and edit/preview against the top-right edge of the content region, separate from the tab strip and document title.
- Remove the document ellipsis menu. Sidebar note menus/context menus now include rename, move, copy content/path and recoverable Trash, including Quick notes and the selected archived note.

## 0.3.0-preview.14 — local only, unreleased

- Remove obsolete absolute-positioned project icons that overlapped project labels after sidebar spacing changed. Use normal flex layout and align child-note text with project labels. No note encoding, contents or paths are changed.

## 0.3.0-preview.12 — local only, unreleased

- Align sidebar branding and note tabs on a shared 48px header. Collapsed navigation keeps only its expand button; no duplicate brand in tabs.
- Compact project rows, hover actions, note metadata, collection lists and save status; add an open-tabs menu without a permanent horizontal scrollbar.
- Add project archive/unarchive and whole-project recoverable Trash/restore, including attachments, tree SHA256 conflict checks and restore collision protection. Archive preserves file paths.
- Explain library search scope; exclude archived projects by default, allow explicit inclusion, and keep Trash excluded.
- Remove stale externally deleted note tabs while retaining conflicted draft tabs; recover copies into Inbox when their project no longer exists.
- Default new UI-created projects to empty, with starter nodes as an explicit option. Preserve existing data.

## 0.3.0-preview.8 — local only, unreleased

- Default all clients to the Windows account's Thoughts home, independent of host environment and installation order. Empty/unset and auto home selections share the same configuration and files.
- Make the Claude MCPB home setting optional with an auto default; retain explicit custom-home support and recognize the older literal default placeholder.
- Make MCP registration home selection optional; resolve the same Windows-user default and preserve conflicting registrations for manual review.
- Add read-only paths diagnostics, shared initialization locking and a migration stop for a known legacy configuration instead of silently creating a separate empty default library.
- Preserve explicit legacy configuration overrides for compatibility and isolated tests. No automatic library merging, cloud storage or new runtime dependencies.

## 0.3.0-preview.7 — local only, unreleased

- Lock the physical library as well as its configuration directory, preventing two differently configured homes from starting competing services for the same files.
- Preserve an existing library during repeat initialization; reject a conflicting new selection.
- Add explicit local runtime install/upgrade/status/rollback, SHA256 and capability checks, immutable executable versions and fail-closed managed pins. Preserve downloaded-file security metadata.
- Share a managed core across MCP launchers and allow only identity-verified compatible local protocol versions to reuse it. Installed executables derive their own managed home when launched directly.
- Refuse service replacement while a connected editor reports a pending save, composition or attachment. Keep prior core pins for software rollback.
- Save path-only reading state and scroll positions outside browser-origin storage; restore it after service restart with a different port. Unsaved text is not persisted by this mechanism.
- Add real-file, two-MCP-process, web API, upgrade/rollback and browser restart acceptance. Native client and clean-Windows acceptance remain separate gates. No public upload.

## 0.3.0-preview.4 — local only, unreleased

- Keep the active tab and its close button visible when resizing a crowded tab strip.

## 0.3.0-preview.3 — local only, unreleased

- Reduce theme settings to one import action and import errors.
- Open notes in separate reusable tabs, with per-tab close, rename/move updates and saved path-only tab restoration. Switching and closing active notes retain the existing save/conflict protection.
- Move edit/preview to the fixed top-right workspace toolbar.
- Use a consistent maximum 360 × 240 local-image preview frame, preserve proportions and original files, and keep full-image viewing.
- Keep this iteration local; no installed extension replacement or public upload.

## 0.3.0-preview.2 — local only, unreleased

- Preserve the brand drawing with a white background and black strokes; generate a 512px package icon.
- Use familiar open/closed folders for projects and a square-pencil icon for Quick notes.
- Keep light appearance only; legacy dark selections fall back to Classic without changing notes or language.
- Add optional Base mainnet sponsorship with an offline-generated address-only QR, independently decoded during generation. No wallet connection or payment SDK.
- Add a read-only Windows security evidence probe. Unsigned binaries and missing clean SAC On/Off acceptance remain release limitations; no security policy was changed.
- Keep all edits local. Earlier review exports and packages are historical snapshots, not the current upload inventory.

## 0.3.0-preview.1 — unreleased

- Prepare free, open-source Windows x64 distribution for Claude Desktop (MCPB), Claude Code and Codex (local stdio).
- Remove paid-theme pricing, checkout metadata and purchase UI. No subscriptions, payments or wallet integration.
- Add explicit shared `THOUGHTS_HOME`; preserve legacy paths when absent. No implicit library migration.
- Package one portable ZIP and one MCPB with the same core, instructions, notices and SHA256.
- Add checksum-verified, inspect-first MCP registration, refusing conflicting existing entries.
- Prepare English documentation, security policy, source/history audit and public upload inventory.

Historical preview.1 publication status: upload was paused for owner review. As of 2026-11-04, the reviewed preview.20 source is public with fresh history; binary distribution remains pending signing and acceptance. Real clients, clean Windows and private-history findings are reported separately.

## 0.2.0-preview.4 — local candidate

- Open/switch/refresh existing notes in preview; explicit editing, new blank notes start editing.
- Built locally; executable blocked by application control. Development-frontend checks did not count as standalone acceptance.

## 0.2.0-preview.3 — local installed preview

- Thoughts branding, favicon, opening Skill and local shortcut.
- Windows packaging, local images/themes and desktop adapter prototypes.
- No public marketplace release or clean-Windows acceptance.

## Earlier development

Local Markdown, authenticated webpage, projects, search, images, SHA256 conflicts and recoverable Trash predate the first private Git snapshot. Private history remains locally; it is not automatically suitable for public publication.

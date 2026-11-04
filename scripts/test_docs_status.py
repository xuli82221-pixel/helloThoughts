"""Regression cases for the published-preview status synchronization failure."""
import importlib.util
from pathlib import Path
import unittest

SPEC = importlib.util.spec_from_file_location("check_docs", Path(__file__).with_name("check-docs.py"))
DOCS = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(DOCS)


class ReleaseStatusTests(unittest.TestCase):
    def test_new_preface_does_not_hide_stale_body_status(self):
        text = ("preview.21 已公开发布，后面是旧记录。\n\n"
                "当前状态：未发布二进制发行版。\n本次不创建公开二进制 Release。")
        self.assertEqual(len(DOCS.release_status_issues("docs/example.md", text)), 2)

    def test_explicit_history_keeps_earlier_release_and_hash_records(self):
        text = ("Current status: 0.3.0-preview.21 unsigned prerelease is public.\n"
                + DOCS.HISTORY_MARKER + "\n## Historical records — preview.20\n"
                "Current candidate: 0.3.0-preview.20; NO BINARY RELEASE.\n"
                "Historical SHA256 evidence remains unchanged.\n")
        self.assertEqual(DOCS.release_status_issues("docs/example.md", text), [])

    def test_signed_stable_pending_and_read_only_builds_are_not_missing_release(self):
        text = ("Current public version: 0.3.0-preview.21 unsigned prerelease.\n"
                "No signed stable release is available. The workflow does not publish releases.\n"
                "SignPath approval and clean-SAC acceptance remain pending.\n")
        self.assertEqual(DOCS.release_status_issues("docs/example.md", text), [])

    def test_current_version_cannot_silently_revert_to_preview20(self):
        for label in ("Current candidate:", "Current version:", "当前核心为"):
            with self.subTest(label=label):
                issues = DOCS.release_status_issues("docs/example.md", label + " 0.3.0-preview.20")
                self.assertEqual(issues[0]["rule"], "staleCurrentVersion")

    def test_readiness_requires_actual_release_link_and_unsigned_boundary(self):
        issues = DOCS.release_status_issues("docs/RELEASE_READINESS_20261104.md", "Current version: 0.3.0-preview.21")
        self.assertEqual({item["rule"] for item in issues},
                         {"missingPublishedReleaseEvidence", "missingUnsignedBoundary"})


if __name__ == "__main__":
    unittest.main()

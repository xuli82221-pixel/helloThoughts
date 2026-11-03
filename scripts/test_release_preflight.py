import importlib.util
import json
from pathlib import Path
import tempfile
import unittest

SPEC = importlib.util.spec_from_file_location("preflight", Path(__file__).with_name("release-preflight.py"))
PREFLIGHT = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(PREFLIGHT)


class SourceBoundaryTests(unittest.TestCase):
    def test_private_artifacts_and_escaping_paths_are_rejected(self):
        for name in ("../outside.md", "/outside.md", "C:/private.md", "docs\\private.md", ".git/config",
                     "library/note.md", "web/node_modules/a/index.js", ".env", "dist/memo.exe", "signing.pfx"):
            with self.subTest(name=name):
                self.assertFalse(PREFLIGHT.safe_relative(name))
        self.assertTrue(PREFLIGHT.safe_relative("release/windows/Thoughts.iss"))
        self.assertTrue(PREFLIGHT.safe_relative(".github/workflows/ci.yml"))

    def test_mismatched_release_versions_are_rejected(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            for relative in ("web/package.json", "web/package-lock.json", "adapters/claude-desktop/manifest.json"):
                file = root / relative
                file.parent.mkdir(parents=True, exist_ok=True)
                file.write_text(json.dumps({"version": "1.0.0"}), encoding="utf-8")
            runtime = root / "internal/memo/runtime.go"
            runtime.parent.mkdir(parents=True)
            runtime.write_text('const Version = "1.0.1"', encoding="utf-8")
            self.assertEqual(PREFLIGHT.inspect_versions(root)[0]["rule"], "version_mismatch")
            runtime.write_text('const Version = "1.0.0"', encoding="utf-8")
            self.assertEqual(PREFLIGHT.inspect_versions(root), [])

    def test_local_document_links_must_exist(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            readme = root / "README.md"
            readme.write_text("[install](INSTALL.md#start) [site](https://example.org) [here](#top)", encoding="utf-8")
            self.assertEqual(len(PREFLIGHT.inspect_links(root, [Path("README.md")])), 1)
            (root / "INSTALL.md").write_text("# Install", encoding="utf-8")
            self.assertEqual(PREFLIGHT.inspect_links(root, [Path("README.md")]), [])


if __name__ == "__main__":
    unittest.main()

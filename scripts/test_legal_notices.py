import json
from html.parser import HTMLParser
from pathlib import Path
import re
import tempfile
import unittest
from legal_notices import LANGUAGES, render_notices

ROOT = Path(__file__).resolve().parent.parent


class Page(HTMLParser):
    def __init__(self):
        super().__init__(); self.links = []; self.ids = []; self.assets = []; self.lang = None; self.scripts = 0
    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if tag == 'html': self.lang = attrs.get('lang')
        if tag == 'script': self.scripts += 1
        if tag == 'a': self.links.append(attrs.get('href', ''))
        if tag in ('link', 'img', 'script'): self.assets.append(attrs.get('src') or attrs.get('href') or '')
        if 'id' in attrs: self.ids.append(attrs['id'])


class NoticesTest(unittest.TestCase):
    def test_all_languages_have_complete_sections_local_links_and_original_license(self):
        source = (ROOT / 'web/src/lib/locales.ts').read_text(encoding='utf-8')
        ui = re.search(r"export const locales = \[([^]]+)\]", source).group(1)
        self.assertEqual(tuple(re.findall(r"'([^']+)'", ui)), LANGUAGES)
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp)
            files = render_notices(ROOT, {'name': 'Test publisher', 'supportEmail': 'support@example.invalid'}, output)
            expected_ids = None
            for language, file in zip(LANGUAGES, files):
                page = (output / 'web/public' / file).read_text(encoding='utf-8')
                parsed = Page(); parsed.feed(page)
                self.assertEqual(parsed.lang, language); self.assertEqual(parsed.scripts, 0)
                ids = set(parsed.ids)
                self.assertEqual(len(ids), len(parsed.ids))
                if expected_ids is None: expected_ids = ids
                self.assertEqual(ids, expected_ids)
                self.assertIn('usage', ids); self.assertIn('license', ids)
                self.assertIn('privacy-agents', ids); self.assertIn('usage-donations', ids)
                self.assertEqual(set(files), {link for link in parsed.links if link.endswith('.html')})
                self.assertTrue(all(not link.startswith(('https:', 'http:')) for link in parsed.links))
                self.assertTrue(all(asset.startswith('/') and not asset.startswith('//') for asset in parsed.assets))
                self.assertNotIn('@@', page)
                self.assertIn('THE SOFTWARE IS PROVIDED', page)
            text = (output / 'PRIVACY.md').read_text(encoding='utf-8')
            for name in ('English', '한국어', '日本語', '简体中文', '繁體中文', 'Deutsch', 'Français', 'Español'):
                self.assertIn('## ' + name, text)

    def test_publisher_text_is_escaped(self):
        with tempfile.TemporaryDirectory() as temp:
            output = Path(temp)
            render_notices(ROOT, {'name': '<script>alert(1)</script>', 'supportEmail': 'a@example.invalid'}, output)
            page = (output / 'web/public/privacy.html').read_text(encoding='utf-8')
            parsed = Page(); parsed.feed(page)
            self.assertEqual(parsed.scripts, 0)
            self.assertIn('&lt;script&gt;', page)


if __name__ == '__main__': unittest.main()

"""Render the eight bundled notices from reviewed text. No network or user data."""
import html
import json
from pathlib import Path

LANGUAGES = ('en', 'ko', 'ja', 'zh-Hans', 'zh-Hant', 'de', 'fr', 'es')


def notice_filename(language):
    return 'privacy.html' if language == 'en' else 'privacy-' + language + '.html'


def render_notices(root, profile, output_root=None):
    output_root = output_root or root
    data = json.loads((root / 'legal/notices.json').read_text(encoding='utf-8'))
    assert data['schema'] == 1 and tuple(data['locales']) == LANGUAGES
    template = (root / 'legal/privacy.template.html').read_text(encoding='utf-8')
    license_text = (root / 'LICENSE').read_text(encoding='utf-8').strip()
    public = output_root / 'web/public'
    public.mkdir(parents=True, exist_ok=True)
    markdown = ['# Thoughts · Privacy and use / 隐私与使用说明', data['updated']]
    baseline = data['locales']['en']

    def text(value):
        return value.replace('@@PUBLISHER@@', profile['name']).replace('@@EMAIL@@', profile['supportEmail'])

    for language, entry in data['locales'].items():
        for group in ('privacy', 'usage'):
            assert [section['id'] for section in entry[group]] == [section['id'] for section in baseline[group]]
            assert all(section['heading'].strip() and section['paragraphs'] and all(p.strip() for p in section['paragraphs']) for section in entry[group])
        assert all(entry[key].strip() for key in ('name', 'title', 'languageLabel', 'updatedLabel', 'privacyTitle', 'intro', 'usageTitle', 'licenseOriginal', 'contactLabel'))
        links = ''.join('<a href="' + notice_filename(lang) + '" lang="' + lang + '" hreflang="' + lang + '"' + (' aria-current="page"' if lang == language else '') + '>' + html.escape(item['name']) + '</a>' for lang, item in data['locales'].items())
        content = ['<article><p class="eyebrow">' + html.escape(entry['updatedLabel'] + ' · ' + data['updated']) + '</p><h1>' + html.escape(entry['privacyTitle']) + '</h1><p class="intro">' + html.escape(text(entry['intro'])) + '</p>']
        markdown.extend(['## ' + entry['name'], entry['updatedLabel'] + ': ' + data['updated'], '### ' + entry['privacyTitle'], text(entry['intro'])])
        for group in ('privacy', 'usage'):
            if group == 'usage':
                content.append('<section id="usage"><h2>' + html.escape(entry['usageTitle']) + '</h2>')
                markdown.append('### ' + entry['usageTitle'])
            for section in entry[group]:
                heading_tag = 'h2' if group == 'privacy' else 'h3'
                content.append('<div id="' + group + '-' + section['id'] + '"><' + heading_tag + '>' + html.escape(section['heading']) + '</' + heading_tag + '>')
                markdown.append('#### ' + section['heading'])
                for paragraph in section['paragraphs']:
                    content.append('<p>' + html.escape(text(paragraph)) + '</p>')
                    markdown.append(text(paragraph))
                content.append('</div>')
            if group == 'usage':
                content.append('<details id="license"><summary>' + html.escape(entry['licenseOriginal']) + '</summary><pre>' + html.escape(license_text) + '</pre></details></section>')
        content.append('</article>')
        replacements = {'LANG': language, 'TITLE': entry['title'], 'LANGUAGE_LABEL': entry['languageLabel'], 'CONTACT_LABEL': entry['contactLabel'], 'EMAIL': profile['supportEmail'], 'PUBLISHER': profile['name']}
        page = template
        for key, value in replacements.items():
            page = page.replace('@@' + key + '@@', html.escape(value, quote=True))
        page = page.replace('@@LANGUAGE_LINKS@@', links).replace('@@CONTENT@@', ''.join(content))
        if '@@' in page:
            raise ValueError('Unresolved notice token: ' + language)
        (public / notice_filename(language)).write_text(page, encoding='utf-8')
    markdown.extend(['## MIT License — English original', '```text\n' + license_text + '\n```'])
    (output_root / 'PRIVACY.md').write_text('\n\n'.join(markdown) + '\n', encoding='utf-8')
    return [notice_filename(lang) for lang in LANGUAGES]

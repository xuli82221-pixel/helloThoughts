"""Read-only public documentation checks; does not build, sign or publish."""
from collections import Counter
import importlib.util
import json
from pathlib import Path
import re
import unicodedata
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parent.parent
READMES = ["README.md", "README.zh-CN.md", "README.zh-TW.md", "README.ja.md",
           "README.ko.md", "README.de.md", "README.fr.md", "README.es.md"]
LINK = re.compile(r'!?\[[^\]]*\]\(([^)]+)\)')
FENCE = re.compile(r'^```[^\n]*\n.*?^```\s*$', re.M | re.S)
HISTORY_MARKER = '<!-- historical-records:start -->'
PUBLISHED_VERSION = '0.3.0-preview.21'
PUBLISHED_URL = 'https://github.com/xuli82221-pixel/helloThoughts/releases/tag/v' + PUBLISHED_VERSION
STALE_RELEASE_CLAIMS = re.compile(
    r'未发布二进制发行版|尚未发布二进制发行版|本次不创建公开二进制\s*Release'
    r'|SOURCE-ONLY PUBLIC PREVIEW|NO (?:PUBLIC )?(?:WINDOWS )?BINARY RELEASE(?:S)?\b'
    r'|binary publication remains pending|Package archives are immutable local candidates'
    r'|an approved binary download remains unavailable|future reviewed binary installation'
    r'|before any formal binary distribution', re.I)


def release_status_issues(name, text):
    """A historical preface alone cannot hide contradictory current statements."""
    current = FENCE.sub('', text.split(HISTORY_MARKER, 1)[0])
    issues = []
    for number, line in enumerate(current.splitlines(), 1):
        if STALE_RELEASE_CLAIMS.search(line):
            issues.append({'file': name, 'line': number, 'rule': 'staleCurrentReleaseClaim'})
        match = re.search(r'(?:current (?:status|version|public version|candidate|core|local candidate/core)|当前(?:状态|核心|版本))'
                          r'[^\n]*?(0\.3\.0-preview\.\d+)', line, re.I)
        if match and match.group(1) != PUBLISHED_VERSION:
            issues.append({'file': name, 'line': number, 'rule': 'staleCurrentVersion'})
    if name in {'docs/RELEASE_READINESS_20261104.md', 'docs/PUBLIC_RELEASE_REVIEW.md',
                'docs/SIGNPATH_APPLICATION.md'}:
        if PUBLISHED_VERSION not in current or PUBLISHED_URL not in current:
            issues.append({'file': name, 'rule': 'missingPublishedReleaseEvidence'})
        if not re.search(r'unsigned|未签名', current, re.I):
            issues.append({'file': name, 'rule': 'missingUnsignedBoundary'})
    return issues


def load_export():
    spec = importlib.util.spec_from_file_location("public_export", ROOT / "scripts/prepare-public-source.py")
    export = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(export)
    return export.source_files()


def anchors(text):
    result = set(re.findall(r'<(?:a|[a-z]+)\b[^>]*(?:id|name)=["\']([^"\']+)', text))
    counts = Counter()
    for heading in re.findall(r'^#{1,6}\s+(.+?)\s*#*$', FENCE.sub('', text), re.M):
        heading = re.sub(r'\[([^\]]+)\]\([^)]+\)', r'\1', heading)
        slug = ''.join(c for c in heading.lower() if c in '-_ ' or unicodedata.category(c)[0] in 'LN')
        slug = slug.replace(' ', '-')
        result.add(slug if not counts[slug] else f'{slug}-{counts[slug]}')
        counts[slug] += 1
    return result


def check():
    files = load_export()
    exported = {p.as_posix() for p in files}
    markdown = [p for p in files if p.suffix == '.md']
    errors = []
    texts = {p.as_posix(): (ROOT / p).read_text(encoding='utf-8-sig') for p in markdown}
    baseline = texts['README.md']
    baseline_links = Counter(LINK.findall(baseline))
    baseline_blocks = FENCE.findall(baseline)
    baseline_sections = len(re.findall(r'^## ', baseline, re.M))
    language_results = []
    for name in READMES:
        text = texts.get(name, '')
        sections = len(re.findall(r'^## ', text, re.M))
        features = len(re.findall(r'^- ', text, re.M))
        checks = {
            'uniformNavigation': text.splitlines()[:1] == baseline.splitlines()[:1],
            'completeSections': sections == baseline_sections == 10,
            'tenFeatures': features == 10,
            'unchangedCodeExamples': FENCE.findall(text) == baseline_blocks,
            'matchingLinkTargets': Counter(LINK.findall(text)) == baseline_links,
            'versionAndDate': PUBLISHED_VERSION in text and '2026-11-04' in text,
            'signingPolicyAndDonate': 'docs/CODE_SIGNING_POLICY.md' in text and 'DONATE.md' in text,
        }
        language_results.append({'file': name, 'sections': sections, 'checks': checks})
        errors.extend({'file': name, 'rule': key} for key, ok in checks.items() if not ok)
    checked_links = 0
    for name, text in texts.items():
        errors.extend(release_status_issues(name, text))
        for target in LINK.findall(FENCE.sub('', text)):
            target = target.split(' "', 1)[0].strip('<>')
            parsed = urlsplit(target)
            if parsed.scheme or parsed.netloc:
                continue
            checked_links += 1
            path = (ROOT / Path(name).parent / unquote(parsed.path)).resolve() if parsed.path else ROOT / name
            if not path.is_file():
                errors.append({'file': name, 'rule': 'missingRelativeTarget', 'target': target})
                continue
            relative = path.relative_to(ROOT).as_posix()
            if relative not in exported:
                errors.append({'file': name, 'rule': 'targetMissingFromPublicExport', 'target': target})
            if parsed.fragment and path.suffix == '.md':
                if unquote(parsed.fragment) not in anchors(path.read_text(encoding='utf-8-sig')):
                    errors.append({'file': name, 'rule': 'missingAnchor', 'target': target})
    return {'status': 'PASS' if not errors else 'FAIL', 'markdownFiles': len(markdown),
            'relativeLinksChecked': checked_links, 'readmes': language_results, 'errors': errors,
            'scope': 'structure, links/anchors/export, identical examples, current release-status contradictions; translation meaning and historical evidence review are manual'}


if __name__ == '__main__':
    result = check()
    print(json.dumps(result, ensure_ascii=False, indent=2))
    raise SystemExit(result['status'] != 'PASS')

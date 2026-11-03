"""Real CLI/HTTP project lifecycle acceptance in a NEW isolated library only.

Keeps its isolated service running for browser review; never opens a browser or
touches a user's configured home. Opening tickets are saved privately, not printed.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import urllib.request
import urllib.error


def run(exe, root, image):
    root.mkdir(parents=True, exist_ok=False)
    home = root / 'home'
    env = dict(os.environ, THOUGHTS_HOME=str(home))
    env.pop('MEMO_CONFIG_DIR', None)
    env.pop('THOUGHTS_NO_DELEGATE', None)
    def cli(*args, content=None):
        p = subprocess.run([str(exe), *args, '--json'], env=env, input=content, text=True, encoding='utf-8', capture_output=True, check=True, timeout=25, creationflags=subprocess.CREATE_NO_WINDOW)
        return json.loads(p.stdout)
    cli('init')
    cli('projects', 'create', '--name', 'Atlas', '--language', 'en')
    cli('projects', 'create', '--name', 'Research', '--template=false')
    cli('projects', 'create', '--name', 'Empty', '--template=false')
    tasks = cli('notes', 'read', '--path', 'Projects/Atlas/Tasks.md')
    cli('notes', 'update', '--path', tasks['path'], '--revision', tasks['revision'], content='# Tasks\n\nThis is an isolated UI acceptance note.\n\n- [ ] Review compact header alignment\n- [x] Preserve real local Markdown\n\n中文排版验收：笔记、标签、归档。\n')
    cli('notes', 'create', '--path', 'Inbox/Release checklist.md', content='# Release checklist\n\nIsolated acceptance material. Search needle in Quick notes.\n')
    cli('notes', 'create', '--path', 'Projects/Research/Archived needle.md', content='# Archived needle\n\nOnly an explicitly enabled archive search should find this note.\n')
    opening = cli('open', '--browser', 'none')
    state = json.loads(next((home / 'config/instances').glob('*/runtime.json')).read_text(encoding='utf-8'))
    assert Path(state['vault']) == home / 'library'
    def http(method, route, body=None, code=200):
        data = None if body is None else json.dumps(body).encode()
        req = urllib.request.Request(state['url']+route, method=method, data=data, headers={'Authorization':'Bearer '+state['token'], 'Content-Type':'application/json'})
        try:
            with urllib.request.urlopen(req, timeout=15) as response:
                assert response.status == code
                return json.load(response)
        except urllib.error.HTTPError as e:
            assert e.code == code, 'unexpected API status'
    snap = http('GET', '/api/project?path=Projects/Research')
    archived = http('POST', '/api/projects/archive', {'path':snap['path'],'revision':snap['revision'],'archived':True})
    assert archived['archived']
    assert len(http('GET','/api/search?q=needle')) == 1
    assert len(http('GET','/api/search?q=needle&includeArchived=1')) == 2
    cli('projects','create','--name','Lifecycle','--template=false')
    cli('notes','create','--path','Projects/Lifecycle/Keep.md',content='# Keep\n\n![image](assets/image.png)\n')
    folder = home/'library/Projects/Lifecycle'; (folder/'assets').mkdir()
    (folder/'assets/image.png').write_bytes(image.read_bytes())
    before = {str(p.relative_to(folder)):hashlib.sha256(p.read_bytes()).hexdigest() for p in folder.rglob('*') if p.is_file()}
    snap = http('GET','/api/project?path=Projects/Lifecycle')
    (folder/'extra.txt').write_text('external file',encoding='utf-8')
    http('POST','/api/projects/trash',{'path':snap['path'],'revision':snap['revision']},409)
    snap = http('GET','/api/project?path=Projects/Lifecycle')
    deleted = http('POST','/api/projects/trash',{'path':snap['path'],'revision':snap['revision']})
    assert not folder.exists()
    cli('projects','create','--name','Lifecycle','--template=false')
    http('POST','/api/project-trash/restore',{'id':deleted['id'],'revision':deleted['revision']},409)
    empty = http('GET','/api/project?path=Projects/Lifecycle')
    http('POST','/api/projects/trash',{'path':empty['path'],'revision':empty['revision']})
    restored = http('POST','/api/project-trash/restore',{'id':deleted['id'],'revision':deleted['revision']})
    for rel, digest in before.items(): assert hashlib.sha256((folder/rel).read_bytes()).hexdigest() == digest
    assert (folder/'extra.txt').read_text(encoding='utf-8') == 'external file'
    # Leave this fixture in Trash for browser review, not in the user's library.
    http('POST','/api/projects/trash',{'path':restored['path'],'revision':restored['revision']})
    fresh = cli('open','--browser','none')
    (root/'browser-open-private.json').write_text(json.dumps(fresh),encoding='utf-8')
    result={'status':'PASS','version':opening['version'],'checks':['archive hides from default search; explicit archive search includes it','stale folder revision including external attachment refused','whole-project Trash and restore preserve Markdown and image hashes','restore collision preserves both folders','real isolated fixtures retained for browser review'],'userLibraryTouched':False,'browserAcceptance':'PENDING'}
    (root/'results.json').write_text(json.dumps(result,indent=2),encoding='utf-8')
    return result

if __name__ == '__main__':
    p=argparse.ArgumentParser();p.add_argument('--exe',type=Path,required=True);p.add_argument('--root',type=Path,required=True);p.add_argument('--image',type=Path,required=True)
    a=p.parse_args();print(json.dumps(run(a.exe.resolve(),a.root.resolve(),a.image.resolve())))

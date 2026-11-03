import fs from 'node:fs'
import path from 'node:path'
import {execFileSync} from 'node:child_process'
import { fileURLToPath } from 'node:url'
const root=path.dirname(path.dirname(fileURLToPath(import.meta.url)))
const target=process.argv[2],goRoot=process.argv[3]
if(!target||!goRoot)throw new Error('Usage: node notices.mjs PLUGIN_DIR GOROOT')
const lock=JSON.parse(fs.readFileSync(path.join(root,'web/package-lock.json'),'utf8'))
const parts=['Thoughts third-party notices\nBundled production dependencies. Build-only packages are excluded.\n', '\n=== shadcn/ui component source ===\n'+fs.readFileSync(path.join(root,'licenses/shadcn-ui-MIT.txt'),'utf8')]
for(const [relative,entry] of Object.entries(lock.packages)){
  if(!relative||entry.dev)continue
  const dir=path.join(root,'web',relative)
  const pkg=JSON.parse(fs.readFileSync(path.join(dir,'package.json'),'utf8'))
  parts.push(`\n=== ${pkg.name} ${pkg.version} (${pkg.license??'see notice'}) ===\n`)
  const files=fs.readdirSync(dir).filter(n=>/^(licen[sc]e|copying|notice)(\.|$)/i.test(n)&&fs.statSync(path.join(dir,n)).isFile())
  for(const name of files)parts.push(fs.readFileSync(path.join(dir,name),'utf8'))
  if(!files.length)parts.push(`License metadata: ${entry.license??pkg.license??'unspecified'}\n`)
}
parts.push('\n=== Go standard library and runtime ===\n'+fs.readFileSync(path.join(goRoot,'LICENSE'),'utf8'))
parts.push('\n=== Inno Setup installer runtime ===\nThe offline Windows setup includes the unmodified Inno Setup installer runtime.\n'+fs.readFileSync(path.join(root,'licenses/InnoSetup-LICENSE.txt'),'utf8'))
const goBinary=path.join(goRoot,'bin',process.platform==='win32'?'go.exe':'go')
const runtimeModules=new Set(execFileSync(goBinary,['list','-deps','-f','{{if .Module}}{{.Module.Path}}{{end}}','./cmd/memo'],{cwd:root,encoding:'utf8'}).trim().split(/\r?\n/).filter(Boolean))
const modules=execFileSync(goBinary,['list','-m','-json','all'],{cwd:root,encoding:'utf8'}).trim().split(/\r?\n(?=\{)/).map(value=>JSON.parse(value))
for(const module of modules.filter(value=>!value.Main&&runtimeModules.has(value.Path))){
 const licenses=fs.readdirSync(module.Dir).filter(name=>/^(LICENSE|COPYING|NOTICE)(\.|$)/i.test(name)&&fs.statSync(path.join(module.Dir,name)).isFile())
 if(!licenses.length)throw new Error('Missing Go dependency license: '+module.Path)
 parts.push('\n=== '+module.Path+' '+module.Version+' ===\n')
 for(const name of licenses)parts.push(fs.readFileSync(path.join(module.Dir,name),'utf8'))
}
fs.writeFileSync(path.join(target,'THIRD_PARTY_NOTICES.txt'),parts.join('\n'),'utf8')

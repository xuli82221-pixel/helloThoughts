import type { Note } from './thoughts-api'
import {localeTag,t} from './i18n'
export type Filter={kind:'all'|'recent'|'folder';path?:string}
export type FolderNode={name:string;path:string;count:number;children:FolderNode[]}
export function folderTree(notes:Note[]):FolderNode[]{
  const roots:FolderNode[]=[]
  for(const note of notes){let level=roots;let path='';for(const name of note.path.split('/').slice(0,-1)){path=path?path+'/'+name:name;let node=level.find(n=>n.path===path);if(!node){node={name,path,count:0,children:[]};level.push(node)}node.count++;level=node.children}}
  const sort=(nodes:FolderNode[])=>{nodes.sort((a,b)=>a.name.localeCompare(b.name,localeTag()));nodes.forEach(node=>sort(node.children))};sort(roots);return roots
}
export function folderPaths(nodes:FolderNode[]):string[]{return nodes.flatMap(n=>[n.path,...folderPaths(n.children)])}
export function parentPath(path:string){return path.includes('/')?path.slice(0,path.lastIndexOf('/')):''}
export function relativeDate(value:string){
  const date=new Date(value),now=new Date(),days=Math.floor((new Date(now.getFullYear(),now.getMonth(),now.getDate()).getTime()-new Date(date.getFullYear(),date.getMonth(),date.getDate()).getTime())/86400000)
  const clock=date.toLocaleTimeString(localeTag(),{hour:'2-digit',minute:'2-digit',hour12:false})
  if(days===0)return t('today',{time:clock})
  if(days===1)return t('yesterday',{time:clock})
  return date.toLocaleDateString(localeTag(),{month:'short',day:'numeric'})
}
export function fullDate(value:string){return new Date(value).toLocaleString(localeTag(),{year:'numeric',month:'long',day:'numeric',hour:'2-digit',minute:'2-digit',hour12:false})}
export function readingBody(content:string){return content.replace(/^\uFEFF?#[ \t]+[^\n]*(?:\r?\n|$)/,'').trimStart()}
export function titleFromDraft(draft:string,fallback:string){return draft.match(/^#\s+(.+)$/m)?.[1]?.trim()||fallback}
export function headingID(text:string){return 'section-'+text.trim().toLowerCase().replace(/\s+/g,'-')}

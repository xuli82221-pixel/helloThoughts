export const noteTab=(path:string)=>path?'note:'+path:'empty'
export const tabPath=(id:string)=>id.startsWith('note:')?id.slice(5):null
export type CollectionTab='inbox'|'trash'|'search'|'archive'|'projects'|`project:${string}`
const pages=new Set(['inbox','trash','search','archive','projects'])
export const projectTab=(path:string):`project:${string}`=>`project:${path}`
export function projectTabPath(id:string){
 if(!id.startsWith('project:'))return null
 const path=id.slice(8),parts=path.split('/'),bytes=new TextEncoder()
 if(bytes.encode(path+'/check.md').length>1000||parts.some(part=>!part||part.startsWith('.')||/[\\:<>"|?*\u0000-\u001f]/.test(part)||/[. ]$/.test(part)||bytes.encode(part).length>240||/^(CON|PRN|AUX|NUL|COM[0-9]|LPT[0-9])$/i.test(part.split('.')[0])))return null
 if(parts.length===2)return parts[0]==='Projects'||parts[0]==='项目'?path:null
 return parts.length===1&&!['inbox','assets'].includes(path.toLowerCase())&&path!=='Projects'&&path!=='项目'?path:null
}
export const isCollectionTab=(id:string):id is CollectionTab=>pages.has(id)||projectTabPath(id)!==null
function valid(id:unknown):id is string{
 if(typeof id!=='string')return false
 if(isCollectionTab(id))return true
 const path=tabPath(id)
 return !!path&&/\.md$/i.test(path)&&!/[\\:\u0000-\u001f]/.test(path)&&path.split('/').every(part=>part!==''&&part!=='.'&&part!=='..')
}
export function restoreTabs(value:unknown,selected:string){
 const tabs=Array.isArray(value)?[...new Set(value.filter(valid))]:[]
 const id=noteTab(selected)
 return selected&&valid(id)?addTab(tabs,id):tabs
}
export function addTab(tabs:string[],id:string){return valid(id)&&!tabs.includes(id)?[...tabs,id]:tabs}
export function removeTab(tabs:string[],active:string,id:string){
 const index=tabs.indexOf(id),next=tabs.filter(tab=>tab!==id)
 return {tabs:next,active:active===id?(next[Math.min(index,next.length-1)]??'empty'):active}
}
export function replaceNoteTab(tabs:string[],from:string,to:string){
 return [...new Set(tabs.map(id=>id===noteTab(from)?noteTab(to):id))]
}
export function pruneNoteTabs(tabs:string[],paths:Set<string>,draftPath?:string,projects?:Set<string>){
 return tabs.filter(id=>{const project=projectTabPath(id);if(project&&projects)return projects.has(project);const path=tabPath(id);return !path||paths.has(path)||path===draftPath})
}

export type TreePage='empty'|'document'|'search'|'trash'|'archive'
export const isUtilityPage=(id:string):id is 'search'|'trash'|'archive'=>id==='search'||id==='trash'||id==='archive'
export function restoreTreeNavigation(value:unknown,selected:string,savedPage:string){
 const tabs=restoreTabs(value,selected).filter(id=>tabPath(id)!==null)
 const page:TreePage=isUtilityPage(savedPage)?savedPage:tabs.includes(noteTab(selected))?'document':'empty'
 return {tabs,page}
}
export function treeBranchOpen(expanded:Record<string,boolean>,key:string){return key==='ui:inbox'||key==='ui:projects'?expanded[key]!==false:expanded[key]===true}
export function toggleTreeBranch(expanded:Record<string,boolean>,key:string){return {...expanded,[key]:!treeBranchOpen(expanded,key)}}
export function revealTreeNote(expanded:Record<string,boolean>,projectPath?:string){return projectPath?{...expanded,'ui:projects':true,[projectPath]:true}:{...expanded,'ui:inbox':true}}

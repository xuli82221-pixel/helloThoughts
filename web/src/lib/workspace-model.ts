import type {Library,Note,Project} from './memo-api'

export function noteName(note:Pick<Note,'path'> & {name?:string}){return note.name??note.path.split('/').at(-1)!.replace(/\.md$/i,'')}
export function projectForNote(library:Library,path:string){return library.projects.find(project=>project.nodes.some(note=>note.path===path))}
export function nodeLabel(project:Project,note:Note){const relative=note.path.startsWith(project.path+'/')?note.path.slice(project.path.length+1):note.path;return relative.replace(/\.md$/i,'').split('/').join(' / ')}
export function validName(value:string){
 if(!value||value!==value.trim()||value.startsWith('.')||/[\/\\:<>"|?*\u0000-\u001f]/.test(value)||/[. ]$/.test(value)||new TextEncoder().encode(value+'.md').length>240)return false
 const stem=value.split('.')[0].toUpperCase()
 return !/^(CON|PRN|AUX|NUL|COM[0-9]|LPT[0-9])$/.test(stem)
}
export function nextUntitledName(base:string,notes:Note[]){const used=new Set(notes.map(note=>noteName(note).toLocaleLowerCase()));let name=base,index=2;while(used.has(name.toLocaleLowerCase()))name=base+' '+index++;return name}

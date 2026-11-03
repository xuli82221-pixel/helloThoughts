import {t} from './i18n'
import {messages,type MessageKey} from './locales'
import {normalizeReaderState} from './reader-state'
export type Note = { path:string; name:string; title:string; content:string; revision:string; modified:string; preview?:string;snippet?:string }
export type Project = {path:string;name:string;nodes:Note[];archived?:boolean}
export type ProjectSnapshot={path:string;name:string;revision:string;noteCount:number;archived:boolean}
export type ProjectTrashEntry={id:string;originalPath:string;name:string;revision:string;noteCount:number;deletedAt:string}
export type Library = {projects:Project[];inbox:Note[];assetVersion?:string}
export type TrashEntry = {id:string;originalPath:string;name:string;deletedAt:string;revision:string;reason:string}
export type Session = { csrf:string; vault:string; vaultPath:string; version:string }
export class ApiError extends Error { constructor(public status:number,message:string,public code=''){super(message)} }
let csrf=''
export async function api<T>(url:string,method='GET',body?:unknown,signal?:AbortSignal):Promise<T>{
  const response=await fetch(url,{method,signal,credentials:'same-origin',headers:{'Content-Type':'application/json','X-Memo-CSRF':csrf},body:body===undefined?undefined:JSON.stringify(body)})
  const value=await response.json()
  if(!response.ok)throw new ApiError(response.status,value.error?.message??'',value.error?.code??'')
  return value
}
export async function uploadBinary<T>(url:string,file:Blob):Promise<T>{
 const response=await fetch(url,{method:'POST',credentials:'same-origin',headers:{'Content-Type':'application/octet-stream','X-Memo-CSRF':csrf},body:file})
 const value=await response.json()
 if(!response.ok)throw new ApiError(response.status,value.error?.message??'',value.error?.code??'')
 return value as T
}
let sessionPromise:Promise<Session>|undefined
let readerCache:Record<string,unknown>|null=null,readerTimer:ReturnType<typeof setTimeout>|undefined,readerSaving=false,readerChanged=false
const readerKeys=new Set(['selected','tabs','workspace-page','scrolls'])
function readerPayload(){return normalizeReaderState({selected:readerCache?.selected,tabs:readerCache?.tabs,page:readerCache?.['workspace-page'],scrolls:readerCache?.scrolls})}
async function saveReader(){
 if(!readerCache||readerSaving)return
 readerSaving=true;readerChanged=false
 try{await api('/api/reader-state','PUT',readerPayload())}catch{/* Keep the local cache; notes are saved separately. */}
 finally{readerSaving=false;if(readerChanged)readerTimer=setTimeout(()=>void saveReader(),300)}
}
if(typeof window!=='undefined')window.addEventListener('pagehide',()=>{if(readerCache){clearTimeout(readerTimer);void fetch('/api/reader-state',{method:'PUT',credentials:'same-origin',keepalive:true,headers:{'Content-Type':'application/json','X-Memo-CSRF':csrf},body:JSON.stringify(readerPayload())}).catch(()=>{})}})
if(typeof window!=='undefined')window.addEventListener('hashchange',()=>{if(new URLSearchParams(location.hash.slice(1)).has('ticket'))location.reload()})
export function bootstrap(){
  if(!sessionPromise){
    const ticket=new URLSearchParams(location.hash.slice(1)).get('ticket')
    history.replaceState(null,'',location.pathname+location.search)
    sessionPromise=(ticket?api<Session>('/api/session','POST',{ticket}):api<Session>('/api/session')).then(async session=>{
     csrf=session.csrf
     try{
      const saved=await api<{exists:boolean;state:unknown}>('/api/reader-state')
      const p=normalizeReaderState(saved.exists?saved.state:{tabs:readPreference('tabs',[]),selected:readPreference('selected',''),page:readPreference('workspace-page','empty'),scrolls:readPreference('scrolls',{})})
      readerCache={tabs:p.tabs,selected:p.selected,'workspace-page':p.page,scrolls:p.scrolls}
     }catch{/* A damaged reading-state file must not prevent access to notes. */}
     return session
    })
  }
  return sessionPromise
}
export function errorMessage(error:unknown){
  if(error instanceof ApiError){const key=('error_'+error.code) as MessageKey;if(key in messages)return t(key);if(error.status===403)return t('error_forbidden');return t('error_generic')}
  return t('loadError')
}
export function readPreference<T>(key:string,fallback:T):T{if(readerCache&&key in readerCache)return readerCache[key] as T;try{const value=localStorage.getItem('memo.ui.'+key);return value===null?fallback:JSON.parse(value)}catch{return fallback}}
export function writePreference(key:string,value:unknown){try{localStorage.setItem('memo.ui.'+key,JSON.stringify(value))}catch{/* Optional cache. */}if(readerCache&&readerKeys.has(key)){readerCache[key]=value;readerChanged=true;clearTimeout(readerTimer);readerTimer=setTimeout(()=>void saveReader(),300)}}

import {useCallback,useEffect,useRef,useState} from 'react'
import {api,ApiError,bootstrap,errorMessage,readPreference,writePreference,type Library,type Note,type Project,type ProjectTrashEntry,type Session,type TrashEntry} from '@/lib/memo-api'
import {applyLocale,isLocale,t} from '@/lib/i18n'
import {initializeThemes} from '@/lib/themes'

type Editor={base:Note|null;draft:string;conflict:boolean;saveError:boolean;loading:boolean}
const empty:Editor={base:null,draft:'',conflict:false,saveError:false,loading:false}
const emptyLibrary:Library={projects:[],inbox:[]}
export function useMemoNotes(){
 const [session,setSession]=useState<Session|null>(null),[library,setLibrary]=useState<Library>(emptyLibrary),[trash,setTrash]=useState<TrashEntry[]>([])
 const [projectTrash,setProjectTrash]=useState<ProjectTrashEntry[]>([])
 const [editor,setEditor]=useState<Editor>(empty),editorRef=useRef(editor)
 const [connected,setConnected]=useState(false),[ready,setReady]=useState(false),[saving,setSaving]=useState(false),[mutating,setMutating]=useState(false),[error,setError]=useState('')
 const [composing,setComposing]=useState(false),composingRef=useRef(false)
 const [attachmentBusy,setAttachmentBusyState]=useState(false),attachmentBusyRef=useRef(false)
 const setAttachmentBusy=useCallback((value:boolean)=>{attachmentBusyRef.current=value;setAttachmentBusyState(value)},[])
 const pathRef=useRef(''),loadSequence=useRef(0),refreshSequence=useRef(0),initialSelection=useRef(false),mutatingRef=useRef(false),writingRef=useRef(false)
 const savePromise=useRef<Promise<boolean>|null>(null)
 const updateEditor=useCallback((value:Editor|((previous:Editor)=>Editor))=>{const next=typeof value==='function'?value(editorRef.current):value;editorRef.current=next;setEditor(next)},[])
 const notes=[...library.inbox,...library.projects.flatMap(project=>project.nodes)]
 const dirty=!!editor.base&&editor.draft!==editor.base.content
 const editorID=useRef(crypto.randomUUID()),editorSequence=useRef(0)
 useEffect(()=>{
  if(!session)return
  const busy=dirty||attachmentBusy||composing||saving
  const report=()=>void api('/api/editor-state','POST',{id:editorID.current,busy,sequence:++editorSequence.current}).catch(()=>{})
  report();const heartbeat=busy?setInterval(report,15000):undefined
  return()=>clearInterval(heartbeat)
 },[session,dirty,attachmentBusy,composing,saving])
 const bind=useCallback((note:Note)=>{pathRef.current=note.path;updateEditor({...empty,base:note,draft:note.content});writePreference('selected',note.path)},[updateEditor])
 const openRaw=useCallback(async(path:string)=>{
  const sequence=++loadSequence.current;pathRef.current=path;updateEditor({...empty,loading:true});setError('')
  try{const note=await api<Note>('/api/notes?path='+encodeURIComponent(path));if(sequence!==loadSequence.current)return false;bind(note);return true}
  catch(error){if(sequence===loadSequence.current){pathRef.current='';updateEditor(empty);setError(errorMessage(error))}return false}
 },[bind,updateEditor])
 const refresh=useCallback(async()=>{
  const sequence=++refreshSequence.current
  try{
   const [value,deleted,deletedProjects]=await Promise.all([api<Library>('/api/library'),api<TrashEntry[]>('/api/trash'),api<ProjectTrashEntry[]>('/api/project-trash')])
   if(sequence!==refreshSequence.current)return
   setLibrary(value);setTrash(deleted);setProjectTrash(deletedProjects)
   if(mutatingRef.current||writingRef.current||attachmentBusyRef.current)return
   const rows=[...value.inbox,...value.projects.flatMap(project=>project.nodes)]
   if(!initialSelection.current){initialSelection.current=true;const cached=readPreference('selected','');if(rows.some(row=>row.path===cached))await openRaw(cached);return}
   const current=editorRef.current,path=pathRef.current
   if(!path||current.loading||!current.base)return
   const metadata=rows.find(row=>row.path===path)
   if(!metadata){
    if(current.draft!==current.base.content||composingRef.current){updateEditor(previous=>({...previous,conflict:true}));return}
    pathRef.current='';updateEditor(empty);writePreference('selected','');return
   }
   if(metadata.revision===current.base.revision)return
   const remote=await api<Note>('/api/notes?path='+encodeURIComponent(path))
   if(sequence!==refreshSequence.current||pathRef.current!==path||mutatingRef.current||writingRef.current||attachmentBusyRef.current)return
   const latest=editorRef.current
   if(latest.loading||!latest.base||latest.base.revision===remote.revision)return
   if(latest.draft!==latest.base.content||composingRef.current)updateEditor(previous=>({...previous,conflict:true}))
   else updateEditor({...empty,base:remote,draft:remote.content})
  }catch(error){if(sequence===refreshSequence.current){if(error instanceof ApiError&&error.status===404)updateEditor(previous=>({...previous,conflict:true}));setError(errorMessage(error))}}
 },[openRaw,updateEditor])
 useEffect(()=>{
  let active=true,events:EventSource|undefined
  bootstrap().then(value=>{
   if(!active)return;setSession(value)
   void initializeThemes()
   events=new EventSource('/api/events')
   events.addEventListener('connected',()=>{if(active){setConnected(true);void refresh()}})
   events.addEventListener('changed',()=>{if(active)void refresh()})
   events.addEventListener('preferences',()=>{if(active)void initializeThemes()})
   events.onerror=()=>{if(active)setConnected(false)}
   void refresh().finally(()=>{if(active)setReady(true)})
  }).catch(error=>{if(active){setReady(true);setError(errorMessage(error))}})
  return()=>{active=false;events?.close()}
 },[refresh])
 const save=useCallback(async():Promise<boolean>=>{
  if(savePromise.current)return savePromise.current
  const current=editorRef.current
  if(!current.base||current.draft===current.base.content)return true
  if(mutatingRef.current||attachmentBusyRef.current||current.conflict||composingRef.current)return false
  const {path,revision}=current.base,content=current.draft
  const pending=(async()=>{
   writingRef.current=true;setSaving(true);setError('')
   try{
    const saved=await api<Note>('/api/notes','PUT',{path,content,revision})
    if(pathRef.current===path)updateEditor(previous=>({...previous,base:saved,conflict:false,saveError:false}))
    return true
   }catch(error){if(pathRef.current===path)updateEditor(previous=>({...previous,saveError:true,conflict:previous.conflict||(error instanceof ApiError&&error.status===409)}));setError(errorMessage(error));return false}
   finally{writingRef.current=false;setSaving(false);savePromise.current=null;void refresh()}
  })();savePromise.current=pending;return pending
 },[refresh,updateEditor])
 const ensureSaved=useCallback(async()=>{
  if(mutatingRef.current||attachmentBusyRef.current){setError(t('operationBusy'));return false}
  if(!await save())return false
  const value=editorRef.current
  if(value.base&&value.draft!==value.base.content){setError(t('saveBeforeSwitch'));return false}
  return true
 },[save])
 const load=useCallback(async(path:string)=>{if(!await ensureSaved())return false;return openRaw(path)},[ensureSaved,openRaw])
 const clearSelection=useCallback(async()=>{if(!await ensureSaved())return false;loadSequence.current++;pathRef.current='';updateEditor(empty);writePreference('selected','');return true},[ensureSaved,updateEditor])
 useEffect(()=>{
  if(!dirty||saving||mutating||attachmentBusy||composing||editor.conflict||editor.saveError||!connected)return
  const timer=setTimeout(()=>void save(),900);return()=>clearTimeout(timer)
 },[dirty,saving,mutating,attachmentBusy,composing,editor.draft,editor.conflict,editor.saveError,connected,save])
 useEffect(()=>{const beforeUnload=(event:BeforeUnloadEvent)=>{const value=editorRef.current;if(attachmentBusyRef.current||value.base&&value.draft!==value.base.content)event.preventDefault()};window.addEventListener('beforeunload',beforeUnload);return()=>window.removeEventListener('beforeunload',beforeUnload)},[])
 const changeDraft=useCallback((draft:string)=>updateEditor(previous=>({...previous,draft,saveError:false})),[updateEditor])
 const composition=useCallback((value:boolean)=>{composingRef.current=value;setComposing(value)},[])
 const reload=useCallback(async()=>{if(!pathRef.current)return false;return openRaw(pathRef.current)},[openRaw])
 const readForMutation=useCallback(async(note:Note)=>{
  if(!await ensureSaved())return null
  try{return await api<Note>('/api/notes?path='+encodeURIComponent(note.path))}catch(error){setError(errorMessage(error));return null}
 },[ensureSaved])
 const create=useCallback(async(path:string,content='')=>{
  const value=editorRef.current
  if(!value.conflict&&!await ensureSaved())return null
  if(mutatingRef.current)return null
  mutatingRef.current=true;setMutating(true);setError('')
  try{const note=await api<Note>('/api/notes','POST',{path,content});bind(note);initialSelection.current=true;return note}
  catch(error){setError(errorMessage(error));return null}
  finally{mutatingRef.current=false;setMutating(false);void refresh()}
 },[bind,ensureSaved,refresh])
 const rename=useCallback(async(note:Note,name:string)=>{
  if(!await ensureSaved())return null
  mutatingRef.current=true;setMutating(true);setError('')
  try{const renamed=await api<Note>('/api/notes/rename','POST',{path:note.path,name,revision:note.revision});if(pathRef.current===note.path)bind(renamed);return renamed}
  catch(error){setError(errorMessage(error));return null}
  finally{mutatingRef.current=false;setMutating(false);void refresh()}
 },[bind,ensureSaved,refresh])
 const trashNote=useCallback(async(note:Note)=>{
  if(!await ensureSaved())return false
  mutatingRef.current=true;setMutating(true);setError('')
  try{await api<TrashEntry>('/api/notes/trash','POST',{path:note.path,revision:note.revision});if(pathRef.current===note.path){pathRef.current='';updateEditor(empty);writePreference('selected','')}return true}
  catch(error){setError(errorMessage(error));return false}
  finally{mutatingRef.current=false;setMutating(false);void refresh()}
 },[ensureSaved,refresh,updateEditor])
 const move=useCallback(async(note:Note,directory:string)=>{
  if(!await ensureSaved())return null
  mutatingRef.current=true;setMutating(true);setError('')
  try{const moved=await api<Note>('/api/notes/move','POST',{path:note.path,directory,revision:note.revision});if(pathRef.current===note.path)bind(moved);return moved}
  catch(error){setError(errorMessage(error));return null}
  finally{mutatingRef.current=false;setMutating(false);void refresh()}
 },[bind,ensureSaved,refresh])
 const restore=useCallback(async(entry:TrashEntry)=>{
  if(!await ensureSaved())return null
  mutatingRef.current=true;setMutating(true);setError('')
  try{const note=await api<Note>('/api/trash/restore','POST',{id:entry.id,revision:entry.revision});bind(note);return note}
  catch(error){setError(errorMessage(error));return null}
  finally{mutatingRef.current=false;setMutating(false);void refresh()}
 },[bind,ensureSaved,refresh])
 const createProject=useCallback(async(name:string,language:string,template:boolean)=>{
  if(!await ensureSaved())return null
  mutatingRef.current=true;setMutating(true);setError('')
  try{return await api<Project>('/api/projects','POST',{name,language,template})}
  catch(error){setError(errorMessage(error));return null}
  finally{mutatingRef.current=false;setMutating(false);await refresh()}
 },[ensureSaved,refresh])
 const projectMutation=useCallback(async<T,>(route:string,body:unknown,removedPath?:string):Promise<T|null>=>{
  if(mutatingRef.current||!await ensureSaved())return null
  mutatingRef.current=true;setMutating(true);setError('')
  try{const value=await api<T>(route,'POST',body);if(removedPath&&pathRef.current.startsWith(removedPath+'/')){pathRef.current='';updateEditor(empty);writePreference('selected','')}return value}
  catch(error){setError(errorMessage(error));return null}
  finally{mutatingRef.current=false;setMutating(false);await refresh()}
 },[ensureSaved,refresh,updateEditor])
 return {session,library,notes,trash,projectTrash,projectMutation,editor,connected,ready,saving,mutating,attachmentBusy,setAttachmentBusy,error,setError,dirty,load,refresh,save,ensureSaved,clearSelection,readForMutation,create,rename,move,trashNote,restore,createProject,reload,changeDraft,composition}
}

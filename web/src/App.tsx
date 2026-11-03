import {useEffect,useLayoutEffect,useRef,useState,type CSSProperties,type UIEvent} from 'react'
import {Archive,ChevronDown,AlertCircle,Check,Eye,FileText,FolderInput,LoaderCircle,MoreHorizontal,PenLine,Plus,RefreshCw,Search,Trash2,Undo2,X} from 'lucide-react'
import {Button} from '@/components/ui/button'
import {Checkbox} from '@/components/ui/checkbox'
import {SidebarProvider,SidebarTrigger,useSidebar} from '@/components/ui/sidebar'
import {Tabs,TabsContent,TabsList,TabsTrigger} from '@/components/ui/tabs'
import {DropdownMenu,DropdownMenuContent,DropdownMenuItem,DropdownMenuSeparator,DropdownMenuTrigger} from '@/components/ui/dropdown-menu'
import {ProjectNavigation} from '@/components/project-navigation'
import {NoteContent} from '@/components/note-content'
import {DeleteDialog,ProjectDeleteDialog,Highlight,NoteDialog,SettingsDialog,MoveDialog,type NoteDialogRequest} from '@/components/workspace-dialogs'
import {useMemoNotes} from '@/hooks/use-memo-notes'
import {api,bootstrap,errorMessage,readPreference,writePreference,type Note,type Project,type ProjectSnapshot,type ProjectTrashEntry,type TrashEntry} from '@/lib/memo-api'
import {useI18n} from '@/lib/i18n'
import {noteName,nextUntitledName,projectForNote} from '@/lib/workspace-model'
import {fullDate,parentPath,relativeDate} from '@/lib/note-view'
import {addTab,noteTab,removeTab,replaceNoteTab,tabPath,pruneNoteTabs,restoreTreeNavigation,isUtilityPage,toggleTreeBranch,revealTreeNote,type TreePage} from '@/lib/workspace-tabs'

function WorkspaceNavigationToggle(){const {state,isMobile}=useSidebar(),{t}=useI18n();return isMobile||state==='collapsed'?<div className="workspace-brand"><SidebarTrigger aria-label={t('expandNav')}/></div>:null}

type Page=TreePage
function restoredNavigation(){return restoreTreeNavigation(readPreference('tabs',[]),readPreference('selected',''),readPreference('workspace-page',''))}
export default function App(){
 const {t,language}=useI18n(),memo=useMemoNotes(),{editor,library}=memo
 const [page,setPage]=useState<Page>(()=>restoredNavigation().page)
 const [tabs,setTabs]=useState<string[]>(()=>restoredNavigation().tabs)
 const tabStrip=useRef<HTMLDivElement>(null)
 const pane=useRef<HTMLDivElement>(null),restoringScroll=useRef(false),[readerReady,setReaderReady]=useState(false)
 const [expanded,setExpanded]=useState<Record<string,boolean>>(()=>readPreference('folders',{}))
 const [request,setRequest]=useState<NoteDialogRequest|null>(null),[deleting,setDeleting]=useState<Note|null>(null),[settings,setSettings]=useState(false),[moving,setMoving]=useState<Note|null>(null)
 const [deletingProject,setDeletingProject]=useState<ProjectSnapshot|null>(null),[projectBusy,setProjectBusy]=useState(false),[includeArchived,setIncludeArchived]=useState(false)
 const [preview,setPreview]=useState(true),[query,setQuery]=useState(''),[results,setResults]=useState<Note[]>([]),[searching,setSearching]=useState(false),[searchError,setSearchError]=useState('')
 const current=editor.base
 const currentProject=page==='document'&&current?projectForNote(library,current.path):undefined
 const busy=memo.mutating||memo.saving||memo.attachmentBusy||projectBusy
 const activeTab=page==='document'?noteTab(current?.path??readPreference('selected','')):page
 const visibleTabs=tabs
 const error=memo.error

 function reveal(path:string){
  const project=library.projects.find(value=>path.startsWith(value.path+'/'))
  setExpanded(previous=>{const next=revealTreeNote(previous,project?.path);writePreference('folders',next);return next})
 }
 async function openNote(note:Pick<Note,'path'>){
  if(page==='document'&&current?.path===note.path){reveal(note.path);return true}
  if(!await memo.load(note.path))return false
  setTabs(previous=>addTab(previous,noteTab(note.path)));setPage('document');setPreview(true);reveal(note.path);return true
 }
 async function changePage(next:Page){if(await memo.ensureSaved()){setPage(next);setPreview(true);memo.setError('')}}
 async function leaveUtility(){if(current)await openNote(current);else if(await memo.ensureSaved()){setPage('empty');memo.setError('')}}
 async function selectTab(id:string){if(id===activeTab)return;const path=tabPath(id);if(path)await openNote({path});else if(isUtilityPage(id))await changePage(id)}
 async function closeTab(id:string){
  if(page!=='document'&&tabPath(id)===current?.path&&!await memo.clearSelection())return
  const next=removeTab(tabs,activeTab,id)
  if(id===activeTab){
   if(!await memo.ensureSaved())return
   const path=tabPath(next.active)
   if(path){if(!await openNote({path}))return}
   else if(next.active==='empty'){if(!await memo.clearSelection())return;setPage('empty')}
   else if(isUtilityPage(next.active)){if(!await memo.clearSelection())return;setPage(next.active)}
   setPreview(true);memo.setError('')
  }
  setTabs(previous=>previous.filter(tab=>tab!==id))
 }
 async function newItem(value:NoteDialogRequest){if(value.kind==='copy'||await memo.ensureSaved()){memo.setError('');setRequest(value)}}
 async function prepareRename(note:Note){const value=await memo.readForMutation(note);if(value){memo.setError('');setRequest({kind:'rename',note:value})}}
 async function prepareDelete(note:Note){const value=await memo.readForMutation(note);if(value){memo.setError('');setDeleting(value)}}
 async function prepareMove(note:Note){const value=await memo.readForMutation(note);if(value){memo.setError('');setMoving(value)}}
 async function confirmMove(directory:string){if(!moving)return false;const value=await memo.move(moving,directory);if(!value)return false;setTabs(previous=>replaceNoteTab(previous,moving.path,value.path));await openNote(value);return true}
 async function submit(name:string,template:boolean){
  if(!request)return false
  if(request.kind==='project'){
   const project=await memo.createProject(name,language,template)
   if(!project)return false
   setExpanded(previous=>{const next=revealTreeNote(previous,project.path);writePreference('folders',next);return next})
   if(project.nodes[0]&&await memo.load(project.nodes[0].path)){setPage('document');setPreview(true)}else{await memo.clearSelection();setPage('empty')}
   return true
  }
  if(request.kind==='rename'){const renamed=await memo.rename(request.note,name);if(renamed)setTabs(previous=>replaceNoteTab(previous,request.note.path,renamed.path));return !!renamed}
  const directory=request.kind==='node'?request.project.path:request.kind==='copy'&&editor.base&&request.project?parentPath(editor.base.path):'Inbox'
  const peers=memo.notes.filter(note=>parentPath(note.path)===directory)
  const base=name||nextUntitledName(t('untitled'),peers)
  const note=await memo.create((directory?directory+'/':'')+base+'.md',request.kind==='copy'?editor.draft:'')
  if(!note)return false
  setPage('document');setPreview(request.kind==='copy')
  setExpanded(previous=>{const project=library.projects.find(p=>directory===p.path||directory.startsWith(p.path+'/'));const next=revealTreeNote(previous,project?.path);writePreference('folders',next);return next})
  return true
 }
 async function confirmDelete(){if(!deleting)return;const wasCurrent=page==='document'&&current?.path===deleting.path;if(await memo.trashNote(deleting)){setTabs(previous=>{const next=previous.filter(id=>id!==noteTab(deleting.path));return wasCurrent?addTab(next,'trash'):next});setDeleting(null);if(wasCurrent)setPage('trash')}}
 async function restore(entry:TrashEntry){const note=await memo.restore(entry);if(note){setPage('document');setPreview(true);reveal(note.path)}}
 async function projectAction(project:Project,action:'archive'|'trash'){
  if(busy||!await memo.ensureSaved())return
  setProjectBusy(true);memo.setError('')
  try{const snap=await api<ProjectSnapshot>('/api/project?path='+encodeURIComponent(project.path));if(action==='trash')setDeletingProject(snap);else await memo.projectMutation('/api/projects/archive',{path:snap.path,revision:snap.revision,archived:!snap.archived})}
  catch(error){memo.setError(errorMessage(error))}finally{setProjectBusy(false)}
 }
 async function confirmProjectDelete(){
  if(!deletingProject)return
  const path=deletingProject.path,wasCurrent=page==='document'&&current?.path.startsWith(path+'/')
  if(await memo.projectMutation<ProjectTrashEntry>('/api/projects/trash',{path,revision:deletingProject.revision},path)){
   setTabs(previous=>previous.filter(id=>!tabPath(id)?.startsWith(path+'/')));setDeletingProject(null)
   if(wasCurrent)setPage('empty')
  }
 }
 async function restoreProject(entry:ProjectTrashEntry){await memo.projectMutation<ProjectSnapshot>('/api/project-trash/restore',{id:entry.id,revision:entry.revision})}
 async function copyNoteContent(note:Note){
  if(current?.path===note.path){await copy(editor.draft);return}
  try{const value=await api<Note>('/api/notes?path='+encodeURIComponent(note.path));await copy(value.content)}catch(error){memo.setError(errorMessage(error))}
 }
 async function copy(value:string){try{await navigator.clipboard.writeText(value)}catch{memo.setError(t('clipboardError'))}}

 useEffect(()=>{let active=true;void bootstrap().then(()=>{if(!active)return;setTabs(restoredNavigation().tabs);setPage(restoredNavigation().page);setReaderReady(true)}).catch(()=>{});return()=>{active=false}},[])
 useEffect(()=>{if(readerReady)writePreference('tabs',tabs)},[tabs,readerReady])
 useEffect(()=>{if(readerReady)writePreference('workspace-page',page)},[page,readerReady])
 useLayoutEffect(()=>{
  if(!readerReady||editor.loading||page==='document'&&!current)return
  restoringScroll.current=true
  const frame=requestAnimationFrame(()=>{if(pane.current)pane.current.scrollTop=readPreference<Record<string,number>>('scrolls',{})[activeTab]??0;requestAnimationFrame(()=>{restoringScroll.current=false})})
  return()=>cancelAnimationFrame(frame)
 },[readerReady,activeTab,editor.loading,preview])
 useEffect(()=>{if(page==='document'&&current)setTabs(previous=>addTab(previous,noteTab(current.path)))},[page,current?.path])
 useEffect(()=>{
  if(!readerReady||!memo.ready||memo.mutating||editor.loading)return
  const paths=new Set(memo.notes.map(note=>note.path))
  setTabs(previous=>{const next=pruneNoteTabs(previous,paths,memo.dirty||editor.conflict?current?.path:undefined,new Set(library.projects.map(project=>project.path)));return next.length===previous.length?previous:next})
  if(page==='document'&&!current)setPage('empty')
 },[readerReady,memo.ready,memo.mutating,library,editor.loading,current?.path,memo.dirty,editor.conflict,page])

 useEffect(()=>{
  const strip=tabStrip.current;if(!strip)return
  const reveal=()=>strip.querySelector('.workspace-tab[data-active="true"]')?.scrollIntoView({block:'nearest',inline:'nearest'})
  reveal();const observer=new ResizeObserver(reveal);observer.observe(strip);return()=>observer.disconnect()
 },[activeTab,tabs.length])

 useLayoutEffect(()=>{
  const element=pane.current;if(!element)return
  const workspace=element.closest<HTMLElement>('.memo-workspace-tabs')
  const align=()=>workspace?.style.setProperty('--workspace-scrollbar-width',Math.max(0,element.offsetWidth-element.clientWidth)+'px')
  align();const observer=new ResizeObserver(align);observer.observe(element)
  return()=>observer.disconnect()
 },[page,activeTab,memo.ready])

 useEffect(()=>{
  const shortcut=(event:KeyboardEvent)=>{
   if(event.isComposing||!(event.ctrlKey||event.metaKey)||request||deleting||moving||settings||deletingProject)return
   if(event.key.toLowerCase()==='s'){event.preventDefault();void memo.save()}
   if(event.key.toLowerCase()==='n'){event.preventDefault();void newItem(currentProject?{kind:'node',project:currentProject}:{kind:'inbox'})}
  };window.addEventListener('keydown',shortcut);return()=>window.removeEventListener('keydown',shortcut)
 })
 useEffect(()=>{
  if(page!=='search'||!query.trim()){setResults([]);setSearching(false);setSearchError('');return}
  const controller=new AbortController();setSearching(true);setSearchError('')
  const timer=setTimeout(()=>{api<Note[]>('/api/search?q='+encodeURIComponent(query)+(includeArchived?'&includeArchived=1':''),'GET',undefined,controller.signal).then(value=>{setResults(value);setSearching(false)}).catch(error=>{if(!controller.signal.aborted){setSearchError(errorMessage(error));setSearching(false)}})},180)
  return()=>{clearTimeout(timer);controller.abort()}
 },[page,query,library,includeArchived])

 const empty=<div className="empty-state"><FileText className="size-7 text-muted-foreground/60" strokeWidth={1.4}/><h1>{t(memo.ready?(memo.session?'noSelection':'connectLibrary'):'opening')}</h1><p>{t(memo.session?'noSelectionHint':'connectLibraryHint')}</p>{memo.session&&<Button variant="outline" size="sm" onClick={()=>void newItem({kind:'inbox'})}><Plus/>{t('newNote')}</Button>}</div>
 const noteRows=(rows:Note[],highlight='',showLocation=false)=><div className="plain-note-list">{rows.map(note=>{
  const project=projectForNote(library,note.path)
  return <div key={note.path} className="plain-note-row note-list-row">
   <button type="button" className="note-row-open" onClick={()=>void openNote(note)}><strong><Highlight text={noteName(note)} query={highlight}/></strong>{showLocation&&<small>{project?.name??t('inbox')}</small>}{note.preview&&<p><Highlight text={note.preview} query={highlight}/></p>}</button>
   <div className="note-row-actions"><time dateTime={note.modified}>{relativeDate(note.modified)}</time><DropdownMenu><DropdownMenuTrigger asChild><Button variant="ghost" size="icon-xs" className="note-row-more" disabled={busy||!memo.session} aria-label={t('moreFor',{name:noteName(note)})}><MoreHorizontal/></Button></DropdownMenuTrigger><DropdownMenuContent align="end"><DropdownMenuItem disabled={busy} onSelect={()=>void prepareRename(note)}><PenLine/>{t('rename')}</DropdownMenuItem><DropdownMenuItem disabled={busy} onSelect={()=>void prepareMove(note)}><FolderInput/>{t('moveNote')}</DropdownMenuItem><DropdownMenuSeparator/><DropdownMenuItem variant="destructive" disabled={busy} onSelect={()=>void prepareDelete(note)}><Trash2/>{t('deleteNote')}</DropdownMenuItem></DropdownMenuContent></DropdownMenu></div>
  </div>
 })}</div>
 const document=current?<div className="document-page" role="region" aria-label={noteName(current)}>
  <NoteContent key={current.path} path={current.path} draft={editor.draft} preview={preview} disabled={memo.mutating} canUpload={!!memo.session&&memo.connected&&!editor.conflict} version={current.revision+(library.assetVersion??" ")} onChange={memo.changeDraft} onComposition={memo.composition} onBusy={value=>{memo.setAttachmentBusy(value);if(!value)void memo.refresh()}} onError={memo.setError}/>
 </div>:empty
 let content:React.ReactNode=empty
 if(page==='document')content=editor.loading?<div className="empty-state"><LoaderCircle className="size-5 animate-spin"/><p>{t('loadingNote')}</p></div>:document
 if(page==='search')content=<section className="collection-page search-page"><h1 className="sr-only">{t('searchLibrary')}</h1><p className="search-scope">{t('searchScopeHint')}</p><label className="search-archive-toggle"><Checkbox checked={includeArchived} onCheckedChange={value=>setIncludeArchived(value===true)} aria-label={t('includeArchived')}/>{t('includeArchived')}</label>{searchError?<p role="alert" className="text-sm text-destructive">{searchError}</p>:query.trim()?<><p className="collection-caption">{searching?t('searching'):t('resultsCount',{count:results.length})}</p>{!searching&&!results.length?<div className="collection-empty"><p>{t('noResults')}</p></div>:noteRows(results,query,true)}</>:<div className="collection-empty"><p>{t('searchPrompt')}</p></div>}</section>
 if(page==='archive')content=<section className="collection-page"><h1 className="sr-only">{t('archivedProjects')}</h1><p className="collection-caption">{t('archivedHint')}</p>{library.projects.filter(p=>p.archived).map(project=><details key={project.path} className="archived-project"><summary><FolderInput className="size-4"/><strong>{project.name}</strong><span>{t('noteCount',{count:project.nodes.length})}</span></summary><div className="archive-project-actions"><Button size="sm" variant="outline" disabled={busy} onClick={()=>void projectAction(project,'archive')}><Undo2/>{t('unarchiveProject')}</Button><Button size="sm" variant="ghost" disabled={busy} onClick={()=>void projectAction(project,'trash')}><Trash2/>{t('deleteProject')}</Button></div>{noteRows(project.nodes)}</details>)}{!library.projects.some(p=>p.archived)&&<p className="collection-empty">{t('archiveEmpty')}</p>}</section>
 if(page==='trash')content=<section className="collection-page"><h1 className="sr-only">{t('trash')}</h1><p className="collection-caption">{t('trashHint')}</p>{memo.projectTrash.length>0&&<><h2 className="collection-subtitle">{t('projects')}</h2><div className="plain-note-list">{memo.projectTrash.map(entry=><div className="trash-row" key={entry.id}><FolderInput className="size-4"/><div><strong>{entry.name}</strong><small>{t('noteCount',{count:entry.noteCount})} · {entry.originalPath}</small><small>{t('deletedAt')} · {fullDate(entry.deletedAt)}</small></div><Button variant="outline" size="sm" disabled={busy} onClick={()=>void restoreProject(entry)}><Undo2/>{t('restore')}</Button></div>)}</div></>}{memo.trash.length>0&&<><h2 className="collection-subtitle">{t('notesLabel')}</h2><div className="plain-note-list">{memo.trash.map(entry=><div className="trash-row" key={entry.id}><FileText className="size-4 text-muted-foreground"/><div><strong>{entry.name}</strong><small>{entry.originalPath}</small><small>{t('deletedAt')} · {fullDate(entry.deletedAt)}</small></div><Button variant="outline" size="sm" disabled={busy} onClick={()=>void restore(entry)}><Undo2/>{t('restore')}</Button></div>)}</div></>}{!memo.trash.length&&!memo.projectTrash.length&&<div className="collection-empty"><Trash2 className="size-6"/><p>{t('trashEmpty')}</p></div>}</section>

 function tabLabel(id:string){
  const path=tabPath(id)
  if(path){
   const name=noteName(memo.notes.find(note=>note.path===path)??{path})
   const duplicate=tabs.some(other=>other!==id&&tabPath(other)&&noteName(memo.notes.find(note=>note.path===tabPath(other))??{path:tabPath(other)!})===name)
   return duplicate?(projectForNote(library,path)?.name??t('inbox'))+' / '+name:name
  }
  return t(id==='inbox'?'inbox':id==='trash'?'trash':id==='search'?'search':id==='archive'?'archivedProjects':id==='projects'?'projects':'workspace')
 }
 const paneProps={ref:pane,className:'workspace-pane',onScroll:(event:UIEvent<HTMLDivElement>)=>{if(readerReady&&!editor.loading&&!restoringScroll.current&&activeTab!=='empty')writePreference('scrolls',{...readPreference<Record<string,number>>('scrolls',{}),[activeTab]:event.currentTarget.scrollTop})}}

 return <SidebarProvider className="memo-app" style={{'--sidebar-width':'224px'} as CSSProperties}>
  <ProjectNavigation onNewNote={()=>void newItem({kind:'inbox'})} query={query} onQuery={setQuery} onSearch={()=>{if(page!=='search')void changePage('search')}} searchDisabled={!memo.session||!!(request||deleting||moving||settings||deletingProject)} library={library} selectedPath={page==='document'?current?.path:undefined} page={page} expanded={expanded} busy={busy||!memo.session} onToggle={path=>setExpanded(previous=>{const value=toggleTreeBranch(previous,path);writePreference('folders',value);return value})} onProjectAction={(project,action)=>void projectAction(project,action)} onSelect={note=>void openNote(note)} onPage={next=>void changePage(next)} onNewProject={()=>void newItem({kind:'project'})} onAddNode={project=>void newItem({kind:'node',project})} onRename={note=>void prepareRename(note)} onMove={note=>void prepareMove(note)} onCopyContent={note=>void copyNoteContent(note)} onCopyPath={note=>void copy((memo.session?.vaultPath??'')+'/'+note.path)} onTrash={note=>void prepareDelete(note)} onSettings={()=>setSettings(true)}/>
  <main className="memo-workspace"><Tabs value={activeTab} onValueChange={id=>void selectTab(id)} activationMode="manual" className="memo-workspace-tabs"><header className="workspace-tabbar"><WorkspaceNavigationToggle/><div className="workspace-tabstrip" ref={tabStrip}><TabsList aria-label={t('workspace')}>{visibleTabs.map(id=>{
   const path=tabPath(id),label=tabLabel(id),Icon=id==='trash'?Trash2:id==='search'?Search:id==='archive'?Archive:FileText
   return <div className="workspace-tab" key={id} data-active={id===activeTab}><TabsTrigger value={id} title={path??label} disabled={memo.mutating||memo.attachmentBusy}><Icon className="size-3.5"/><span>{label}</span></TabsTrigger>{id!=='empty'&&<Button variant="ghost" size="icon-xs" disabled={busy} aria-label={t('closeTab')+' · '+label} onClick={()=>void closeTab(id)}><X/></Button>}</div>
  })}</TabsList></div>{tabs.length>2&&<DropdownMenu><DropdownMenuTrigger asChild><Button variant="ghost" size="icon-sm" className="tab-overflow" aria-label={t('openTabs')} title={t('openTabs')}><ChevronDown/></Button></DropdownMenuTrigger><DropdownMenuContent align="end" className="max-h-80 overflow-y-auto">{tabs.map(id=>{const path=tabPath(id);const label=tabLabel(id);return <DropdownMenuItem key={id} disabled={busy} onSelect={()=>void selectTab(id)}><span className="min-w-0"><span className="block truncate">{label}</span>{path&&<span className="block text-xs text-muted-foreground">{path}</span>}</span>{id===activeTab&&<Check className="ml-auto"/>}</DropdownMenuItem>})}</DropdownMenuContent></DropdownMenu>}</header>
   {memo.ready&&memo.session&&!memo.connected&&<div className="workspace-feedback" role="alert"><AlertCircle className="size-4 shrink-0"/><span>{t('serviceDisconnectedHint')}</span></div>}
   {error&&<div className="workspace-feedback" role="alert"><AlertCircle className="size-4 shrink-0"/><span>{error}</span>{editor.saveError&&!editor.conflict&&<Button variant="outline" size="sm" onClick={()=>void memo.save()}>{t('retrySave')}</Button>}<Button variant="ghost" size="icon-xs" aria-label={t('close')} onClick={()=>memo.setError('')}><X/></Button></div>}
   {editor.conflict&&current&&<div className="conflict-banner" role="alert"><div><strong>{t('externalChange')}</strong><p>{t('draftSafe')}</p></div><div><Button size="sm" variant="outline" onClick={()=>void memo.reload()}>{t('reloadDisk')}</Button><Button size="sm" onClick={()=>void newItem({kind:'copy',project:currentProject})}>{t('saveCopy')}</Button></div></div>}
   {page==='document'&&current&&<div className="workspace-content-toolbar">
    {page==='document'&&current&&<div className="workspace-document-actions">
     <span className={'save-state '+(!memo.connected||editor.conflict||editor.saveError?'has-error':'')} role="status" title={t('modified')+' · '+fullDate(current.modified)}>{!memo.connected?t('disconnected'):memo.saving?<><LoaderCircle className="size-3 animate-spin"/>{t('saving')}</>:editor.conflict?t('conflict'):editor.saveError?t('saveFailed'):memo.dirty?t('unsaved'):<><Check className="size-3.5"/>{t('saved')}</>}</span>
     <Button variant="ghost" size="sm" disabled={busy||editor.loading} aria-label={t(preview?'write':'preview')} onClick={()=>{void memo.ensureSaved().then(ok=>{if(ok)setPreview(value=>!value)})}}>{preview?<PenLine/>:<Eye/>}{t(preview?'write':'preview')}</Button>
    </div>}
   </div>}
   {isUtilityPage(page)&&<div className="utility-toolbar"><span>{t(page==='search'?'search':page==='trash'?'trash':'archivedProjects')}</span><div className="utility-toolbar-actions">{page==='trash'&&<Button variant="ghost" size="icon-xs" aria-label={t('refresh')} onClick={()=>void memo.refresh()}><RefreshCw/></Button>}<Button variant="ghost" size="icon-xs" onClick={()=>void leaveUtility()} aria-label={t('close')}><X/></Button></div></div>}
   {page==='document'?<TabsContent {...paneProps} value={activeTab}>{content}</TabsContent>:<div {...paneProps} role="region" aria-label={t(page==='search'?'search':page==='trash'?'trash':page==='archive'?'archivedProjects':'workspace')}>{content}</div>}
  </Tabs></main>
  <ProjectDeleteDialog project={deletingProject} onClose={()=>setDeletingProject(null)} onConfirm={confirmProjectDelete} busy={busy} error={error}/><NoteDialog request={request} onClose={()=>setRequest(null)} onSubmit={submit} busy={busy} error={error}/><DeleteDialog note={deleting} onClose={()=>setDeleting(null)} onConfirm={confirmDelete} busy={busy} error={error}/><MoveDialog note={moving} library={library} onClose={()=>setMoving(null)} onConfirm={confirmMove} busy={busy} error={error}/><SettingsDialog open={settings} onOpenChange={setSettings} session={memo.session} connected={memo.connected}/>
 </SidebarProvider>
}

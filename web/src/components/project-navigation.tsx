import {useEffect,useRef,useState} from 'react'
import {Archive,ChevronRight,Copy,FileText,FolderInput,MoreHorizontal,PenLine,Plus,Search,Trash2} from 'lucide-react'
import {Sidebar,SidebarContent,SidebarFooter,SidebarGroup,SidebarGroupContent,SidebarHeader,SidebarMenu,SidebarMenuButton,SidebarMenuItem,SidebarTrigger,useSidebar} from '@/components/ui/sidebar'
import {Button} from '@/components/ui/button'
import {Input} from '@/components/ui/input'
import {ContextMenu,ContextMenuContent,ContextMenuItem,ContextMenuSeparator,ContextMenuTrigger} from '@/components/ui/context-menu'
import {DropdownMenu,DropdownMenuContent,DropdownMenuItem,DropdownMenuSeparator,DropdownMenuTrigger} from '@/components/ui/dropdown-menu'
import {useI18n} from '@/lib/i18n'
import {ThoughtsBrand} from '@/components/thoughts-brand'
import {ProjectsIcon,QuickNotesIcon,SettingsIcon} from '@/components/navigation-icons'
import {treeBranchOpen} from '@/lib/workspace-tabs'
import {nodeLabel,noteName} from '@/lib/workspace-model'
import type {Library,Note,Project} from '@/lib/memo-api'

type Props={
 library:Library;selectedPath?:string;page:string;expanded:Record<string,boolean>;busy:boolean
 query:string;searchDisabled:boolean;onQuery:(value:string)=>void;onSearch:()=>void
 onToggle:(path:string)=>void;onSelect:(note:Note)=>void;onPage:(page:'trash'|'archive')=>void
 onProjectAction:(project:Project,action:'archive'|'trash')=>void;onNewProject:()=>void;onNewNote:()=>void;onAddNode:(project:Project)=>void
 onRename:(note:Note)=>void;onMove:(note:Note)=>void;onCopyContent:(note:Note)=>void;onCopyPath:(note:Note)=>void;onTrash:(note:Note)=>void;onSettings:()=>void
}
export function ProjectNavigation(props:Props){
 const {t}=useI18n(),{isMobile,open,openMobile,setOpen,setOpenMobile}=useSidebar()
 const searchInput=useRef<HTMLInputElement>(null),[focusSearch,setFocusSearch]=useState(false)
 const inboxOpen=treeBranchOpen(props.expanded,'ui:inbox'),projectsOpen=treeBranchOpen(props.expanded,'ui:projects')
 const inboxActive=props.library.inbox.some(note=>note.path===props.selectedPath)
 const archivedSelection=props.library.projects.find(project=>project.archived&&project.nodes.some(note=>note.path===props.selectedPath))
 const select=(note:Note)=>{props.onSelect(note);if(isMobile)setOpenMobile(false)}
 const go=(page:'trash'|'archive')=>{props.onPage(page);if(isMobile)setOpenMobile(false)}

 useEffect(()=>{
  const shortcut=(event:KeyboardEvent)=>{
   if(event.defaultPrevented||event.isComposing||props.searchDisabled||!(event.ctrlKey||event.metaKey)||event.key.toLowerCase()!=='k')return
   event.preventDefault();props.onSearch();setFocusSearch(true)
   if(isMobile)setOpenMobile(true);else setOpen(true)
  }
  window.addEventListener('keydown',shortcut);return()=>window.removeEventListener('keydown',shortcut)
 },[props.onSearch,props.searchDisabled,isMobile,setOpen,setOpenMobile])
 useEffect(()=>{
  if(!focusSearch||!(isMobile?openMobile:open))return
  const frame=requestAnimationFrame(()=>{searchInput.current?.focus();searchInput.current?.select();setFocusSearch(false)})
  return()=>cancelAnimationFrame(frame)
 },[focusSearch,isMobile,open,openMobile])

 function actions(note:Note,context:boolean){
  const Item=context?ContextMenuItem:DropdownMenuItem,Separator=context?ContextMenuSeparator:DropdownMenuSeparator
  return <>
   <Item disabled={props.busy} onSelect={()=>props.onRename(note)}><PenLine/>{t('rename')}</Item>
   <Item disabled={props.busy} onSelect={()=>props.onMove(note)}><FolderInput/>{t('moveNote')}</Item>
   <Separator/>
   <Item onSelect={()=>props.onCopyContent(note)}><Copy/>{t('copyContent')}</Item>
   <Item onSelect={()=>props.onCopyPath(note)}><FileText/>{t('copyPath')}</Item>
   <Separator/>
   <Item variant="destructive" disabled={props.busy} onSelect={()=>props.onTrash(note)}><Trash2/>{t('deleteNote')}</Item>
  </>
 }
 function noteRow(note:Note,label=noteName(note)){
  return <ContextMenu key={note.path}><ContextMenuTrigger asChild>
   <div className={'navigation-node '+(props.selectedPath===note.path?'is-selected':'')}>
    <Button variant="ghost" className="node-select" onClick={()=>select(note)} disabled={props.busy} aria-current={props.selectedPath===note.path?'page':undefined} title={note.path}><span>{label}</span></Button>
    <DropdownMenu><DropdownMenuTrigger asChild><Button variant="ghost" size="icon-sm" className="node-more" aria-label={t('moreFor',{name:noteName(note)})}><MoreHorizontal/></Button></DropdownMenuTrigger><DropdownMenuContent align="start">{actions(note,false)}</DropdownMenuContent></DropdownMenu>
   </div>
  </ContextMenuTrigger><ContextMenuContent>{actions(note,true)}</ContextMenuContent></ContextMenu>
 }
 return <Sidebar collapsible="offcanvas" className="memo-navigation">
  <div className="memo-navigation-top">
  <SidebarHeader className="memo-sidebar-header"><ThoughtsBrand/><SidebarTrigger aria-label={t('collapseNav')}/></SidebarHeader>
  <form role="search" className="sidebar-search" onSubmit={event=>{event.preventDefault();props.onSearch();if(isMobile)setOpenMobile(false)}}>
   <Search className="size-3.5"/>
   <Input ref={searchInput} type="search" value={props.query} disabled={props.searchDisabled} onFocus={()=>{if(props.page!=='search')props.onSearch()}} onChange={event=>props.onQuery(event.target.value)} placeholder={t('searchLibrary')} aria-label={t('searchLibrary')} aria-keyshortcuts="Control+k Meta+k"/>
   {!props.query&&<kbd aria-hidden="true">Ctrl K</kbd>}
  </form>
  </div>
  <SidebarContent className="memo-navigation-scroll">
   <SidebarGroup className="memo-shortcuts"><div className="navigation-section-heading"><SidebarMenuButton className="navigation-section-toggle" isActive={inboxActive} onClick={()=>props.onToggle('ui:inbox')} aria-expanded={inboxOpen}><ChevronRight className={inboxOpen?'rotate-90':''}/><QuickNotesIcon/><span>{t('inbox')}</span></SidebarMenuButton><Button variant="ghost" size="icon-sm" className="navigation-section-action" disabled={props.busy} onClick={props.onNewNote} aria-label={t('newNote')}><Plus/></Button></div>
    {inboxOpen&&<div className="quick-note-nodes">{props.library.inbox.length?props.library.inbox.map(note=>noteRow(note)):<p className="project-empty-label">{t('projectEmpty')}</p>}</div>}
   </SidebarGroup>
   <SidebarGroup className="memo-project-group"><div className="navigation-section-heading"><SidebarMenuButton className="navigation-section-toggle" onClick={()=>props.onToggle('ui:projects')} aria-expanded={projectsOpen}><ChevronRight className={projectsOpen?'rotate-90':''}/><ProjectsIcon/><span>{t('projects')}</span></SidebarMenuButton><Button variant="ghost" size="icon-sm" className="navigation-section-action" aria-label={t('newProject')} disabled={props.busy} onClick={props.onNewProject}><Plus/></Button></div>
    {projectsOpen&&<SidebarGroupContent><div className="project-navigation-tree">
     {props.library.projects.filter(project=>!project.archived).map(project=>{
      const expanded=props.expanded[project.path]===true
      return <div key={project.path} className="navigation-project"><div className="navigation-project-row">
       <Button variant="ghost" className="project-toggle" onClick={()=>props.onToggle(project.path)} aria-expanded={expanded} aria-label={t(expanded?'collapse':'expand',{name:project.name})}><ChevronRight className={'size-3.5 transition-transform '+(expanded?'rotate-90':'')}/><span>{project.name}</span></Button>
       <DropdownMenu><DropdownMenuTrigger asChild><Button variant="ghost" size="icon-sm" className="project-more" disabled={props.busy} aria-label={t('moreFor',{name:project.name})}><MoreHorizontal/></Button></DropdownMenuTrigger><DropdownMenuContent align="start"><DropdownMenuItem onSelect={()=>props.onAddNode(project)}><Plus/>{t('addNode')}</DropdownMenuItem><DropdownMenuItem onSelect={()=>props.onProjectAction(project,'archive')}><Archive/>{t('archiveProject')}</DropdownMenuItem><DropdownMenuItem variant="destructive" onSelect={()=>props.onProjectAction(project,'trash')}><Trash2/>{t('deleteProject')}</DropdownMenuItem></DropdownMenuContent></DropdownMenu>
      </div>{expanded&&<div className="project-nodes">{project.nodes.length?project.nodes.map(note=>noteRow(note,nodeLabel(project,note))):<p className="project-empty-label">{t('projectEmpty')}</p>}</div>}</div>
     })}
    </div></SidebarGroupContent>}
   </SidebarGroup>
   {props.library.projects.some(project=>project.archived)&&<SidebarGroup className="memo-shortcuts"><SidebarMenu><SidebarMenuItem><SidebarMenuButton isActive={props.page==='archive'||!!archivedSelection} onClick={()=>go('archive')}><Archive/><span>{t('archivedProjects')}</span><small className="ml-auto text-muted-foreground">{props.library.projects.filter(p=>p.archived).length}</small></SidebarMenuButton></SidebarMenuItem></SidebarMenu>
    {archivedSelection&&<div className="quick-note-nodes">{archivedSelection.nodes.filter(note=>note.path===props.selectedPath).map(note=>noteRow(note))}</div>}
   </SidebarGroup>}
  </SidebarContent>
  <SidebarFooter className="memo-navigation-footer"><SidebarMenu><SidebarMenuItem><SidebarMenuButton isActive={props.page==='trash'} onClick={()=>go('trash')}><Trash2/><span>{t('trash')}</span></SidebarMenuButton></SidebarMenuItem><SidebarMenuItem><SidebarMenuButton onClick={props.onSettings}><SettingsIcon/><span>{t('settings')}</span></SidebarMenuButton></SidebarMenuItem></SidebarMenu></SidebarFooter>
 </Sidebar>
}

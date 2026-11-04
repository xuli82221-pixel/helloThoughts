import {useEffect,useState} from 'react'
import {Check,LoaderCircle,ExternalLink} from 'lucide-react'
import {Button} from '@/components/ui/button'
import {Input} from '@/components/ui/input'
import {Checkbox} from '@/components/ui/checkbox'
import {Dialog,DialogContent,DialogDescription,DialogFooter,DialogHeader,DialogTitle} from '@/components/ui/dialog'
import {AlertDialog,AlertDialogAction,AlertDialogCancel,AlertDialogContent,AlertDialogDescription,AlertDialogFooter,AlertDialogHeader,AlertDialogTitle} from '@/components/ui/alert-dialog'
import {Select,SelectContent,SelectItem,SelectTrigger,SelectValue} from '@/components/ui/select'
import {api,errorMessage,type Note,type Project,type ProjectSnapshot,type Session,type Library} from '@/lib/thoughts-api'
import {applyLocale,isLocale,languageNames,locales,useI18n,type Locale} from '@/lib/i18n'
import {noteName,validName} from '@/lib/workspace-model'
import {PRODUCT_NAME} from '@/lib/product'
import {AppearanceSettings} from '@/components/appearance-settings'
import {SupportSettings} from '@/components/support-settings'
import {privacyPath} from '@/lib/privacy-path'
import {parentPath} from '@/lib/note-view'

export function MoveDialog({note,library,onClose,onConfirm,busy,error}:{note:Note|null;library:Library;onClose:()=>void;onConfirm:(directory:string)=>Promise<boolean>;busy:boolean;error:string}){
 const {t}=useI18n(),[destination,setDestination]=useState('')
 const choices=[{path:'Inbox',name:t('inbox')},...library.projects].filter(project=>project.path!==parentPath(note?.path??''))
 useEffect(()=>{setDestination(choices[0]?.path??'')},[note?.path])
 return <Dialog open={!!note} onOpenChange={open=>{if(!open&&!busy)onClose()}}><DialogContent><DialogHeader><DialogTitle>{t('moveTitle',{name:note?noteName(note):''})}</DialogTitle><DialogDescription>{t('moveHint')}</DialogDescription></DialogHeader>{choices.length?<div className="space-y-2"><label className="text-sm font-medium" id="move-destination">{t('moveDestination')}</label><Select value={destination} onValueChange={setDestination} disabled={busy}><SelectTrigger aria-labelledby="move-destination" className="w-full"><SelectValue/></SelectTrigger><SelectContent>{choices.map(project=><SelectItem key={project.path} value={project.path}>{project.name}</SelectItem>)}</SelectContent></Select></div>:<p className="text-sm text-muted-foreground">{t('noMoveDestination')}</p>}{error&&<p role="alert" className="text-sm text-destructive">{error}</p>}<DialogFooter><Button variant="outline" disabled={busy} onClick={onClose}>{t('cancel')}</Button><Button disabled={busy||!destination} onClick={()=>void onConfirm(destination).then(ok=>{if(ok)onClose()})}>{busy&&<LoaderCircle className="animate-spin"/>}{t('moveNote')}</Button></DialogFooter></DialogContent></Dialog>
}

export type NoteDialogRequest={kind:'node';project:Project}|{kind:'inbox'}|{kind:'project'}|{kind:'rename';note:Note}|{kind:'copy';project?:Project}
export function NoteDialog({request,onClose,onSubmit,busy,error}:{request:NoteDialogRequest|null;onClose:()=>void;onSubmit:(name:string,template:boolean)=>Promise<boolean>;busy:boolean;error:string}){
 const {t}=useI18n(),[name,setName]=useState(''),[template,setTemplate]=useState(false),[invalid,setInvalid]=useState(false)
 useEffect(()=>{setName(request?.kind==='rename'?noteName(request.note):'');setTemplate(false);setInvalid(false)},[request])
 const title=t(request?.kind==='rename'?'rename':request?.kind==='project'?'newProject':request?.kind==='inbox'?'newNote':request?.kind==='copy'?'saveCopy':'addNode')
 const description=request?.kind==='project'?t('projectHint'):request?.kind==='rename'?t('renameHint'):request?.kind==='node'?t('noteInProject',{project:request.project.name}):request?.kind==='copy'?t('copyDraftHint'):t('quickNoteHint')
 async function submit(){const value=name.trim();if((request?.kind==='project'||request?.kind==='rename')&&!value||value&&!validName(value)){setInvalid(true);return}setInvalid(false);if(await onSubmit(value,template))onClose()}
 return <Dialog open={!!request} onOpenChange={open=>{if(!open&&!busy)onClose()}}><DialogContent><DialogHeader><DialogTitle>{title}</DialogTitle><DialogDescription>{description}</DialogDescription></DialogHeader><form onSubmit={event=>{event.preventDefault();void submit()}} className="space-y-5">
  <div className="space-y-2"><label className="text-sm font-medium" htmlFor="thoughts-item-name">{t(request?.kind==='project'?'projectName':'nodeName')}</label><Input id="thoughts-item-name" value={name} onChange={event=>setName(event.target.value)} placeholder={request?.kind==='project'?'':t('untitled')} maxLength={80} disabled={busy} autoFocus aria-invalid={invalid}/></div>
  {request?.kind==='project'&&<div className="space-y-2"><label className="flex items-center gap-3 text-sm"><Checkbox checked={template} onCheckedChange={value=>setTemplate(value===true)} disabled={busy}/>{t('starterNodes')}</label><p className="text-xs leading-relaxed text-muted-foreground">{t('starterHint')}</p></div>}
  {(invalid||error)&&<p role="alert" className="text-sm text-destructive">{invalid?t('invalidName'):error}</p>}
  <DialogFooter><Button type="button" variant="outline" disabled={busy} onClick={onClose}>{t('cancel')}</Button><Button type="submit" disabled={busy}>{busy&&<LoaderCircle className="animate-spin"/>}{t(request?.kind==='rename'?'rename':request?.kind==='project'?'newProject':request?.kind==='copy'?'saveCopy':'create')}</Button></DialogFooter>
 </form></DialogContent></Dialog>
}
export function DeleteDialog({note,onClose,onConfirm,busy,error}:{note:Note|null;onClose:()=>void;onConfirm:()=>Promise<void>;busy:boolean;error:string}){
 const {t}=useI18n()
 return <AlertDialog open={!!note} onOpenChange={open=>{if(!open&&!busy)onClose()}}><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{t('deleteTitle',{name:note?noteName(note):''})}</AlertDialogTitle><AlertDialogDescription>{t('deleteHint')}</AlertDialogDescription></AlertDialogHeader>{error&&<p role="alert" className="text-sm text-destructive">{error}</p>}<AlertDialogFooter><AlertDialogCancel disabled={busy}>{t('cancel')}</AlertDialogCancel><AlertDialogAction className="bg-destructive text-white hover:bg-destructive/90" disabled={busy} onClick={event=>{event.preventDefault();void onConfirm()}}>{busy&&<LoaderCircle className="animate-spin"/>}{t('deleteNote')}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
}
export function ProjectDeleteDialog({project,onClose,onConfirm,busy,error}:{project:ProjectSnapshot|null;onClose:()=>void;onConfirm:()=>Promise<void>;busy:boolean;error:string}){
 const {t}=useI18n()
 return <AlertDialog open={!!project} onOpenChange={open=>{if(!open&&!busy)onClose()}}><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{t('deleteTitle',{name:project?.name??''})}</AlertDialogTitle><AlertDialogDescription>{t('deleteProjectHint',{count:project?.noteCount??0})}</AlertDialogDescription></AlertDialogHeader>{error&&<p role="alert" className="text-sm text-destructive">{error}</p>}<AlertDialogFooter><AlertDialogCancel disabled={busy}>{t('cancel')}</AlertDialogCancel><AlertDialogAction disabled={busy} className="bg-destructive text-white hover:bg-destructive/90" onClick={event=>{event.preventDefault();void onConfirm()}}>{busy&&<LoaderCircle className="animate-spin"/>}{t('deleteProject')}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
}
export function SettingsDialog({open,onOpenChange,session,connected}:{open:boolean;onOpenChange:(value:boolean)=>void;session:Session|null;connected:boolean}){
 const {t,language}=useI18n(),[saving,setSaving]=useState(false),[error,setError]=useState(''),[copied,setCopied]=useState(false),[openingBrowser,setOpeningBrowser]=useState(false),[browserOpened,setBrowserOpened]=useState(false)
 async function chooseLanguage(value:Locale){if(saving||!session)return;setSaving(true);setError('');try{const preferences=await api<{language:string}>('/api/preferences','PUT',{language:value});if(isLocale(preferences.language))applyLocale(preferences.language)}catch(error){setError(errorMessage(error))}finally{setSaving(false)}}
 return <Dialog open={open} onOpenChange={onOpenChange}><DialogContent className="max-h-[calc(100dvh-3rem)] overflow-y-auto"><DialogHeader><DialogTitle>{t('settings')}</DialogTitle><DialogDescription>{t('settingsHint')}</DialogDescription></DialogHeader>
  <div className="space-y-6">{session&&<AppearanceSettings/>}<section className="space-y-3 border-t pt-5"><label className="text-sm font-medium" id="language-label">{t('language')}</label><Select value={language} onValueChange={value=>{if(isLocale(value))void chooseLanguage(value)}} disabled={!session||saving}><SelectTrigger aria-labelledby="language-label" className="w-full"><SelectValue/></SelectTrigger><SelectContent>{locales.map(value=><SelectItem key={value} value={value}>{languageNames[value]}</SelectItem>)}</SelectContent></Select><p className="text-xs leading-relaxed text-muted-foreground">{t('languageHint')}</p></section>
  <section className="space-y-3 border-t pt-5"><h3 className="text-sm font-medium">{t('library')}</h3><p className="break-all rounded-md border bg-muted/40 p-3 font-mono text-xs leading-relaxed">{session?.vaultPath??t('notConnected')}</p><div className="flex items-center justify-between gap-3"><span className="text-xs text-muted-foreground">{t(connected?'connected':'disconnected')}</span><Button size="sm" variant="outline" disabled={!session} onClick={()=>{if(session)void navigator.clipboard.writeText(session.vaultPath).then(()=>setCopied(true)).catch(()=>setError(t('clipboardError')))}}>{copied?<><Check/>{t('copied')}</>:t('copyFolder')}</Button></div></section>
  <section className="space-y-3 border-t pt-5"><Button variant="outline" disabled={!session||!connected||openingBrowser} onClick={()=>{setOpeningBrowser(true);setError('');void api('/api/open-browser','POST',{}).then(()=>setBrowserOpened(true)).catch(error=>setError(errorMessage(error))).finally(()=>setOpeningBrowser(false))}}>{openingBrowser?<LoaderCircle className="animate-spin"/>:<ExternalLink/>}{t('openExternalBrowser')}</Button><p className="text-xs leading-relaxed text-muted-foreground">{t(browserOpened?'externalBrowserOpened':'browserSessionHint')}</p></section>
  <section className="space-y-2 border-t pt-5"><h3 className="text-sm font-medium">{t('privacyTitle')}</h3><p className="text-xs leading-relaxed text-muted-foreground">{t('privacyHint')}</p><p className="text-xs leading-relaxed text-muted-foreground">{t('privacyAgent')}</p><div className="flex flex-wrap gap-x-5 gap-y-2 pt-2 text-xs"><a href={privacyPath(language)} target="_blank" rel="noopener noreferrer" className="underline underline-offset-4">{t('privacyPolicy')}</a><a href={privacyPath(language)+'#usage'} target="_blank" rel="noopener noreferrer" className="underline underline-offset-4">{t('usageNotes')}</a><a href="mailto:hello@spxipo.com" className="underline underline-offset-4">{t('contactSupport')}</a></div></section>
  <SupportSettings/>
  {error&&<p role="alert" className="text-sm text-destructive">{error}</p>}<p className="border-t pt-4 text-xs text-muted-foreground">{PRODUCT_NAME} {session?.version??''}</p>
  </div>
 </DialogContent></Dialog>
}
export function Highlight({text,query}:{text:string;query:string}){const value=query.trim();if(!value)return <>{text}</>;const escaped=value.replace(/[.*+?^${}()|[\]\\]/g,'\\$&');return <>{text.split(new RegExp('('+escaped+')','gi')).map((part,index)=>part.toLocaleLowerCase()===value.toLocaleLowerCase()?<mark key={index}>{part}</mark>:part)}</>}

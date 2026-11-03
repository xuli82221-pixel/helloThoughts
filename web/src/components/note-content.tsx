import {useRef,useState,type ChangeEvent,type ClipboardEvent,type DragEvent} from 'react'
import ReactMarkdown,{defaultUrlTransform} from 'react-markdown'
import remarkGfm from 'remark-gfm'
import {ImagePlus,ImageOff,LoaderCircle} from 'lucide-react'
import {Button} from '@/components/ui/button'
import {Textarea} from '@/components/ui/textarea'
import {Dialog,DialogContent,DialogHeader,DialogTitle,DialogDescription} from '@/components/ui/dialog'
import {ApiError,errorMessage,uploadBinary} from '@/lib/memo-api'
import {imageMarkdown,insertImageMarkdown,localImageURL,type UploadedImage} from '@/lib/image-content'
import {useI18n} from '@/lib/i18n'
import {readingBody} from '@/lib/note-view'

function LocalImage({note,source,alt,version,onOpen}:{note:string;source:string;alt:string;version:string;onOpen:(url:string,alt:string)=>void}){
 const {t}=useI18n(),[failed,setFailed]=useState(false),url=localImageURL(note,source,version)
 if(!url||failed)return <span className="missing-image" role="img" aria-label={t('imageMissing')}><ImageOff className="size-4"/>{alt||t('imageMissing')}</span>
 return <span className="local-image"><button type="button" aria-label={t('openImage',{name:alt||t('image')})} onClick={()=>onOpen(url,alt)}><img src={url} alt={alt} loading="lazy" onError={()=>setFailed(true)}/></button></span>
}
type Props={path:string;draft:string;preview:boolean;disabled:boolean;canUpload:boolean;version:string;onChange:(content:string)=>void;onComposition:(value:boolean)=>void;onBusy:(value:boolean)=>void;onError:(value:string)=>void}
export function NoteContent(props:Props){
 const {t}=useI18n(),editor=useRef<HTMLTextAreaElement>(null),input=useRef<HTMLInputElement>(null),uploadingRef=useRef(false)
 const [uploading,setUploading]=useState(false),[dragging,setDragging]=useState(false),[added,setAdded]=useState(false),[lightbox,setLightbox]=useState<{url:string;alt:string}|null>(null)
 const latest=useRef(props);latest.current=props
 async function add(files:File[]){
  if(!files.length||uploadingRef.current||!props.canUpload||props.disabled)return
  if(files.length>8){props.onError(t('tooManyImages'));return}
  if(files.some(file=>file.size>10*1024*1024)){props.onError(t('error_image_too_large'));return}
  if(files.some(file=>!['image/png','image/jpeg','image/webp'].includes(file.type)&&!(file.type===''&&/\.(png|jpe?g|webp)$/i.test(file.name)))){props.onError(t('error_invalid_image'));return}
  const textarea=editor.current;if(!textarea)return
  const initial=textarea.value,start=textarea.selectionStart,end=textarea.selectionEnd,notePath=props.path
  uploadingRef.current=true;setUploading(true);props.onBusy(true);props.onError('')
  const inserted:string[]=[]
  try{for(const file of files){const image=await uploadBinary<UploadedImage>('/api/images?note='+encodeURIComponent(notePath),file);inserted.push(imageMarkdown(file.name||t('image'),image.source))}}
  catch(error){props.onError(error instanceof ApiError?errorMessage(error):t('imageUploadFailed'))}
  finally{
   if(inserted.length&&latest.current.path===notePath&&editor.current){const target=editor.current,result=insertImageMarkdown(initial,target.value,start,end,target.selectionEnd,inserted.join('\n\n'));latest.current.onChange(result.content);setAdded(true);requestAnimationFrame(()=>{if(editor.current===target){target.focus();target.setSelectionRange(result.cursor,result.cursor)}})}
   uploadingRef.current=false;setUploading(false);latest.current.onBusy(false)
  }
 }
 function paste(event:ClipboardEvent<HTMLTextAreaElement>){const files=Array.from(event.clipboardData.files);if(files.length){event.preventDefault();void add(files)}}
 function drop(event:DragEvent<HTMLTextAreaElement>){if(event.dataTransfer.files.length){event.preventDefault();setDragging(false);void add(Array.from(event.dataTransfer.files))}}
 function choose(event:ChangeEvent<HTMLInputElement>){const files=Array.from(event.target.files??[]);event.target.value='';void add(files)}
 return <>
  {props.preview?<article className="markdown"><ReactMarkdown remarkPlugins={[remarkGfm]} skipHtml urlTransform={(url,key)=>key==='src'?url:defaultUrlTransform(url)} components={{img:({src,alt})=><LocalImage key={(src??'')+props.version} note={props.path} source={src??''} alt={alt??''} version={props.version} onOpen={(url,alt)=>setLightbox({url,alt})}/>,a:({href,children})=>href?.startsWith('#')?<a href={href}>{children}</a>:href&&/^https?:\/\//i.test(href)?<a href={href} target="_blank" rel="noopener noreferrer" referrerPolicy="no-referrer">{children}</a>:<span>{children}</span>}}>{readingBody(props.draft)}</ReactMarkdown></article>:<>
   <div className="editor-media-tools"><Button variant="outline" size="sm" disabled={props.disabled||!props.canUpload||uploading} onClick={()=>input.current?.click()}>{uploading?<LoaderCircle className="animate-spin"/>:<ImagePlus/>}{t(uploading?'imageUploading':'addImage')}</Button><span role={uploading?'status':undefined}>{t(added?'imageSavedHint':'imageHint')}</span><input ref={input} type="file" accept="image/png,image/jpeg,image/webp" multiple hidden onChange={choose} aria-label={t('addImage')}/></div>
   <label className="sr-only" htmlFor="note-editor">{t('markdownBody')}</label><Textarea ref={editor} id="note-editor" aria-label={t('markdownBody')} value={props.draft} onChange={event=>props.onChange(event.target.value)} onCompositionStart={()=>props.onComposition(true)} onCompositionEnd={()=>props.onComposition(false)} onPaste={paste} onDragOver={event=>{if(event.dataTransfer.types.includes('Files')){event.preventDefault();setDragging(true)}}} onDragLeave={()=>setDragging(false)} onDrop={drop} placeholder={t('writePlaceholder')} disabled={props.disabled} spellCheck={false} className={'document-editor '+(dragging?'is-dragging':'')}/>
  </>}
  <Dialog open={!!lightbox} onOpenChange={open=>{if(!open)setLightbox(null)}}><DialogContent className="image-lightbox"><DialogHeader><DialogTitle>{lightbox?.alt||t('image')}</DialogTitle><DialogDescription className="sr-only">{t('imageLightboxHint')}</DialogDescription></DialogHeader>{lightbox&&<img src={lightbox.url} alt={lightbox.alt}/>}</DialogContent></Dialog>
 </>
}

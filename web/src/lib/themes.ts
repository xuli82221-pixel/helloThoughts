import {useSyncExternalStore} from 'react'
import {api,errorMessage,uploadBinary} from './thoughts-api'
import {applyLocale,isLocale} from './i18n'
import {isLightTheme} from './theme-policy'

export type Theme={id:string;name:string;version:string;themeSchema:number;tokens:Record<string,string>;builtIn:boolean;preview:boolean}
const colors='background foreground card card-foreground popover popover-foreground primary primary-foreground secondary secondary-foreground muted muted-foreground accent accent-foreground destructive border input ring sidebar sidebar-foreground sidebar-primary sidebar-primary-foreground sidebar-accent sidebar-accent-foreground sidebar-border sidebar-ring workspace-tab-background editor-background code-background'.split(' ')
type Appearance={themes:Theme[];selected:string;busy:boolean;error:string}
let state:Appearance={themes:[],selected:'classic',busy:false,error:''}
const listeners=new Set<()=>void>()
let sequence=0
function publish(update:Partial<Appearance>){state={...state,...update};listeners.forEach(listener=>listener())}
function valid(theme:Theme){return theme?.themeSchema===1&&isLightTheme(theme.tokens)&&colors.every(key=>/^#[0-9a-f]{6}([0-9a-f]{2})?$/i.test(theme.tokens?.[key]))&&/^(?:[0-9]|1[0-6])px$/.test(theme.tokens?.radius)}
export function applyTheme(theme:Theme){
 if(!valid(theme))return false
 for(const key of [...colors,'radius'])document.documentElement.style.setProperty('--'+key,theme.tokens[key])
 document.documentElement.style.colorScheme='light'
 document.documentElement.dataset.theme=theme.id
 try{localStorage.setItem('memo.ui.appearance',JSON.stringify(theme))}catch{/* Optional cache. */}
 publish({selected:theme.id});return true
}
export function restoreCachedTheme(){try{const theme=JSON.parse(localStorage.getItem('memo.ui.appearance')??'null');if(theme)applyTheme(theme)}catch{/* Start with Classic CSS. */}}
export async function initializeThemes(selected?:string){
 const current=++sequence
 try{const [themes,preferences]=await Promise.all([api<Theme[]>('/api/themes'),selected?Promise.resolve({theme:selected,language:undefined}):api<{theme:string;language:string}>('/api/preferences')]);if(current!==sequence)return;if(isLocale(preferences.language))applyLocale(preferences.language);publish({themes,error:''});const value=themes.find(theme=>theme.id===preferences.theme)??themes.find(theme=>theme.id==='classic');if(value)applyTheme(value)}catch(error){if(current!==sequence)return;for(const key of [...colors,'radius'])document.documentElement.style.removeProperty('--'+key);document.documentElement.style.colorScheme='light';document.documentElement.dataset.theme='classic';publish({selected:'classic',error:errorMessage(error)})}
}
export async function selectTheme(id:string){
 if(state.busy)return;const current=++sequence;publish({busy:true,error:''})
 try{const result=await api<{theme:string}>('/api/preferences','PUT',{theme:id});const value=state.themes.find(theme=>theme.id===result.theme);if(current===sequence&&value)applyTheme(value)}catch(error){publish({error:errorMessage(error)})}finally{publish({busy:false})}
}
export async function importTheme(file:File){
 if(state.busy)return;publish({busy:true,error:''})
 try{if(file.size>3*1024*1024)throw new Error('package');const theme=await uploadBinary<Theme>('/api/themes',file);const result=await api<{theme:string}>('/api/preferences','PUT',{theme:theme.id});await initializeThemes(result.theme)}catch(error){publish({error:file.size>3*1024*1024?'THEME_TOO_LARGE':errorMessage(error)})}finally{publish({busy:false})}
}
export async function removeTheme(id:string){
 if(state.busy)return;publish({busy:true,error:''})
 try{await api('/api/themes/'+encodeURIComponent(id),'DELETE');await initializeThemes()}catch(error){publish({error:errorMessage(error)})}finally{publish({busy:false})}
}
export function useAppearance(){return useSyncExternalStore(listener=>{listeners.add(listener);return()=>listeners.delete(listener)},()=>state)}

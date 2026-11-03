import {useSyncExternalStore} from 'react'
import {locales,messages,singularMessages,type Locale,type MessageKey} from './locales'
import {PRODUCT_NAME} from './product'
export {languageNames,locales,type Locale} from './locales'

function cachedLanguage():Locale {
 try { const value=JSON.parse(localStorage.getItem('memo.ui.language')??'"en"');return isLocale(value)?value:'en' } catch { return 'en' }
}
export function isLocale(value:unknown):value is Locale {return typeof value==='string'&&(locales as readonly string[]).includes(value)}
let current:Locale=cachedLanguage()
const listeners=new Set<()=>void>()
export function locale(){return current}
export function applyLocale(value:Locale){
 current=value
 document.documentElement.lang=value
 document.title=PRODUCT_NAME+' · '+t('library')
 try {localStorage.setItem('memo.ui.language',JSON.stringify(value))}catch{/* Optional browser cache. */}
 listeners.forEach(listener=>listener())
}
export function t(key:MessageKey,values:Record<string,string|number>={}) {
 let text:string=messages[key][locales.indexOf(current)]
 if(values.count===1)text=singularMessages[current]?.[key]??text
 for(const [name,value] of Object.entries(values))text=text.replaceAll('{'+name+'}',String(value))
 return text
}
export function useI18n(){const language=useSyncExternalStore(listener=>{listeners.add(listener);return()=>{listeners.delete(listener)}},()=>current);return {language,t}}
export function localeTag(){return {en:'en-US',ko:'ko-KR',ja:'ja-JP','zh-Hans':'zh-CN','zh-Hant':'zh-TW',de:'de-DE',fr:'fr-FR',es:'es-ES'}[current]}

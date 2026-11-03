import {restoreTabs,isCollectionTab,type CollectionTab} from './workspace-tabs.ts'
export type ReaderState={schema:1;tabs:string[];selected:string;page:'empty'|'document'|CollectionTab;scrolls:Record<string,number>}
export function normalizeReaderState(value:unknown):ReaderState{
 const p=(value&&typeof value==='object'?value:{}) as Partial<ReaderState>
 const selected=typeof p.selected==='string'&&restoreTabs([],p.selected).length?p.selected:''
 const tabs=restoreTabs(p.tabs,selected).slice(-200)
 const page=p.page==='document'||typeof p.page==='string'&&isCollectionTab(p.page)?p.page:'empty'
 const scrolls:Record<string,number>={}
 if(p.scrolls&&typeof p.scrolls==='object')for(const [id,pos] of Object.entries(p.scrolls).slice(0,500)){
  if(restoreTabs([id],'').length&&typeof pos==='number'&&Number.isFinite(pos)&&pos>=0&&pos<=1e8)scrolls[id]=pos
 }
 return {schema:1,tabs,selected,page,scrolls}
}

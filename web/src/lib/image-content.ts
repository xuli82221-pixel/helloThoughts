export type UploadedImage={path:string;source:string;width:number;height:number;mime:string}
export function imageMarkdown(name:string,source:string){const alt=name.replace(/\.[^.]+$/,'').replace(/[\r\n]/g,' ').replace(/[\\\[\]]/g,'\\$&');return '!['+alt+']('+source.split('/').map(encodeURIComponent).join('/')+')'}
export function insertImageMarkdown(initial:string,latest:string,start:number,end:number,currentEnd:number,markdown:string){
 const unchanged=initial===latest,from=Math.max(0,Math.min(latest.length,unchanged?start:currentEnd)),to=unchanged?Math.max(from,Math.min(latest.length,end)):from
 const addition='\n\n'+markdown+'\n\n';return {content:latest.slice(0,from)+addition+latest.slice(to),cursor:from+addition.length}
}
export function localImageURL(note:string,source:string,version:string){
 if(!source||/^[a-z][a-z0-9+.-]*:/i.test(source)||source.startsWith('/')||/[\\\x00-\x1f?#]/.test(source))return undefined
 return '/api/images?'+new URLSearchParams({note,src:source,v:version}).toString()
}

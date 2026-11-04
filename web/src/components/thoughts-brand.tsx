import {PRODUCT_NAME} from '@/lib/product'

export function ThoughtsBrand({compact=false}:{compact?:boolean}){
 return <div className={'thoughts-brand'+(compact?' thoughts-brand-compact':'')}><img src="/assets/icon.svg" alt="" width="24" height="24" className="shrink-0 rounded-md"/><span>{PRODUCT_NAME}</span></div>
}

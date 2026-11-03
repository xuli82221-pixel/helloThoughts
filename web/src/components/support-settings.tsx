import {useState} from 'react'
import {Check,Copy} from 'lucide-react'
import {Button} from '@/components/ui/button'
import {useI18n} from '@/lib/i18n'
import config from '../../../release/sponsorship.json'

export function SupportSettings(){
 const {t}=useI18n(),[copied,setCopied]=useState(false),[error,setError]=useState(false)
 if(!config.enabled||config.network!=='Base'||config.chainId!==8453||!/^0x[0-9a-fA-F]{40}$/.test(config.address))return null
 return <section className="space-y-3 border-t pt-5"><h3 className="text-sm font-medium">{t('supportThoughts')}</h3><p className="text-xs leading-relaxed text-muted-foreground">{t('supportOptional')}</p><div className="support-layout"><img src="/assets/support-base-address-qr.png" width="160" height="160" className="rounded-lg border bg-white" alt={t('supportQrAlt')}/><div className="min-w-0 space-y-3"><p className="text-sm font-medium">{t('baseMainnet')}</p><p className="break-all font-mono text-xs leading-relaxed">{config.address}</p><Button variant="outline" size="sm" onClick={()=>{setError(false);void navigator.clipboard.writeText(config.address).then(()=>setCopied(true)).catch(()=>setError(true))}}>{copied?<Check/>:<Copy/>}{t(copied?'copied':'copyAddress')}</Button></div></div><p className="text-xs leading-relaxed text-muted-foreground">{t('supportNetworkHint')}</p>{error&&<p role="alert" className="text-xs text-destructive">{t('clipboardError')}</p>}</section>
}

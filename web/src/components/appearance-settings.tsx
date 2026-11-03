import {useRef} from 'react'
import {LoaderCircle,Upload} from 'lucide-react'
import {Button} from '@/components/ui/button'
import {importTheme,useAppearance} from '@/lib/themes'
import {useI18n} from '@/lib/i18n'

export function AppearanceSettings(){
 const {t}=useI18n(),appearance=useAppearance(),input=useRef<HTMLInputElement>(null)
 return <section className="space-y-3">
  <Button variant="outline" size="sm" disabled={appearance.busy} onClick={()=>input.current?.click()}>{appearance.busy?<LoaderCircle className="animate-spin"/>:<Upload/>}{t('importTheme')}</Button>
  <input ref={input} type="file" accept=".thoughts-theme,.memo-theme,.zip" hidden aria-label={t('importTheme')} onChange={event=>{const file=event.target.files?.[0];event.target.value='';if(file)void importTheme(file)}}/>
  {appearance.error&&<p role="alert" className="text-sm text-destructive">{appearance.error==='THEME_TOO_LARGE'?t('themeTooLarge'):appearance.error}</p>}
 </section>
}

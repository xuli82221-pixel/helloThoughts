import type {Locale} from './locales'

export function privacyPath(language:Locale){return language==='en'?'/privacy.html':'/privacy-'+language+'.html'}

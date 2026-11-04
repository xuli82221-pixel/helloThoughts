import React from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import { setNonce } from 'get-nonce'
import '@fontsource/geist/latin-400.css'
import '@fontsource/geist/latin-500.css'
import '@fontsource/geist/latin-600.css'
import '@fontsource-variable/noto-sans-sc'
import '@fontsource-variable/noto-sans-jp'
import '@fontsource-variable/noto-sans-kr'
import '@fontsource-variable/noto-sans-tc'
import '@fontsource/geist-mono/latin-400.css'
import './style.css'
import {applyLocale,locale} from './lib/i18n'
import {restoreCachedTheme} from './lib/themes'

const styleNonce=document.querySelector('meta[name="thoughts-style-nonce"]')?.getAttribute('content')
if(styleNonce)setNonce(styleNonce)
applyLocale(locale())
restoreCachedTheme()
createRoot(document.getElementById('root')!).render(<React.StrictMode><App/></React.StrictMode>)

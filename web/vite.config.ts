import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({resolve:{alias:{'@':decodeURIComponent(new URL('./src',import.meta.url).pathname).replace(/^\/(?=[A-Za-z]:)/,'')}},plugins:[react(),tailwindcss()],build:{outDir:'../internal/webui/dist',emptyOutDir:true,sourcemap:false}})

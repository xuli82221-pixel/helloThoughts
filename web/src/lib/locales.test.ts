import {test} from 'node:test'
import assert from 'node:assert/strict'
import {messages,locales,languageNames,singularMessages} from './locales.ts'

test('all eight supported languages have a nonempty translation for every key',()=>{
 assert.deepEqual(locales,['en','ko','ja','zh-Hans','zh-Hant','de','fr','es'])
 for(const [key,values] of Object.entries(messages)){
  assert.equal(values.length,locales.length,key)
  for(const value of values)assert.ok(value.trim(),key)
 }
 assert.equal(languageNames.en,'English')
 assert.deepEqual(locales.slice(-3).map(locale=>languageNames[locale]),['Deutsch','Français','Español'])
})

test('European count messages include the correct singular forms',()=>{
 assert.equal(singularMessages.de?.noteCount,'{count} Notiz')
 assert.equal(singularMessages.fr?.resultsCount,'{count} résultat')
 assert.equal(singularMessages.es?.characters,'{count} carácter')
 for(const language of ['de','fr','es'] as const){
  for(const key of ['projectCount','noteCount','folderCount','resultsCount','characters'] as const){
   assert.ok(singularMessages[language]?.[key]?.includes('{count}'),language+': '+key)
  }
 }
})
test('translated messages preserve all interpolation parameters',()=>{
 const parameters=(text:string)=>Array.from(text.matchAll(/\{(\w+)\}/g),match=>match[1]).sort()
 for(const [key,values] of Object.entries(messages))for(const value of values)assert.deepEqual(parameters(value),parameters(values[0]),key)
})
test('project navigation and quick notes have distinct English labels',()=>{
 assert.equal(new Set([messages.projects[0],messages.inbox[0],messages.trash[0],messages.search[0]]).size,4)
})

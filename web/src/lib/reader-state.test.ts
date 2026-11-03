import {test} from 'node:test'
import assert from 'node:assert/strict'
import {normalizeReaderState} from './reader-state.ts'

test('project parent view and scroll restore without storing note content',()=>{
 const state=normalizeReaderState({tabs:['projects','project:Projects/工作'],selected:'',page:'project:Projects/工作',scrolls:{'project:Projects/工作':128},draft:'not saved'})
 assert.equal(state.page,'project:Projects/工作');assert.deepEqual(state.scrolls,{'project:Projects/工作':128});assert.equal('draft' in state,false)
 assert.equal(normalizeReaderState({page:'project:../outside',tabs:['project:../outside']}).page,'empty')
})

test('archive is a restorable workspace page without storing project contents',()=>{
 const state=normalizeReaderState({tabs:['archive','note:Projects/P/a.md'],selected:'',page:'archive',scrolls:{archive:220},projectBody:'not persisted'})
 assert.equal(state.page,'archive');assert.deepEqual(state.tabs,['archive','note:Projects/P/a.md']);assert.deepEqual(state.scrolls,{archive:220});assert.equal('projectBody' in state,false)
})
test('reading state keeps tabs, selection and scroll without copying drafts or credentials',()=>{
 const state=normalizeReaderState({tabs:['note:Inbox/a.md','search'],selected:'Inbox/a.md',page:'document',scrolls:{'note:Inbox/a.md':428},draft:'private draft',csrf:'private token'})
 assert.deepEqual(state,{schema:1,tabs:['note:Inbox/a.md','search'],selected:'Inbox/a.md',page:'document',scrolls:{'note:Inbox/a.md':428}})
})
test('malformed reader paths and scroll values do not survive state restoration',()=>{
 const state=normalizeReaderState({selected:'../outside.md',page:'invalid',tabs:['note:C:/outside.md','trash'],scrolls:{trash:Infinity,search:-2,'note:../outside.md':50}})
 assert.deepEqual(state,{schema:1,tabs:['trash'],selected:'',page:'empty',scrolls:{}})
})

import {test} from 'node:test'
import assert from 'node:assert/strict'
import {nextUntitledName,nodeLabel,noteName,projectForNote,validName} from './workspace-model.ts'
import type {Library,Note} from './thoughts-api.ts'
const note=(path:string,name:string,title=name):Note=>({path,name,title,content:'',revision:'version',modified:'2026-10-27T00:00:00Z'})
test('renamed file name stays authoritative even when its original heading is preserved',()=>{
 const renamed=note('Projects/SPXIPO/更细.md','更细','细节')
 assert.equal(noteName(renamed),'更细')
 assert.equal(noteName({path:'旧项目/规则.md'}),'规则')
})
test('project selection follows the real returned tree and legacy nested file paths',()=>{
 const current=note('项目/Memo/旧目录/笔记.md','笔记')
 const library:Library={projects:[{path:'项目/Memo',name:'Memo',nodes:[current]}],inbox:[]}
 assert.equal(projectForNote(library,current.path)?.name,'Memo')
 assert.equal(nodeLabel(library.projects[0],current),'旧目录 / 笔记')
 assert.equal(projectForNote(library,'项目/Memo/已删除.md'),undefined)
})
test('unsafe or reserved names are rejected while all supported scripts remain valid',()=>{
 for(const name of ['../other','a/b','a\\b','CON','con.txt','NUL','COM1','.hidden',' bad','bad.','line\nbreak',''])assert.equal(validName(name),false,name)
 for(const name of ['规则','更細','한국어','日本語','New idea'])assert.equal(validName(name),true,name)
})
test('unnamed quick notes do not overwrite a case variant of an existing file',()=>{
 assert.equal(nextUntitledName('Untitled note',[note('Inbox/untitled note.md','untitled note'),note('Inbox/Untitled note 2.md','Untitled note 2')]),'Untitled note 3')
})

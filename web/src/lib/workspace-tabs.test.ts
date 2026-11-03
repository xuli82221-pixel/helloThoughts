import {test} from 'node:test'
import assert from 'node:assert/strict'
import {addTab,noteTab,removeTab,replaceNoteTab,restoreTabs,tabPath,pruneNoteTabs,projectTab,projectTabPath,restoreTreeNavigation,treeBranchOpen,toggleTreeBranch,revealTreeNote} from './workspace-tabs.ts'

test('only note tabs survive legacy migration; utility views remain available without tabs',()=>{
 assert.deepEqual(restoreTreeNavigation(['inbox','projects','project:Projects/A','note:Inbox/a.md','search'],'Inbox/a.md','inbox'),{tabs:['note:Inbox/a.md'],page:'document'})
 assert.deepEqual(restoreTreeNavigation(['inbox','project:Projects/A'],'','project:Projects/A'),{tabs:[],page:'empty'})
 assert.equal(restoreTreeNavigation(['inbox','trash'],'Inbox/a.md','trash').page,'trash')
 assert.deepEqual(restoreTreeNavigation(['search','archive','trash'],'','search'),{tabs:[],page:'search'})
})
test('tree expansion is independent of the content page and other branches',()=>{
 assert.equal(treeBranchOpen({},'ui:inbox'),true);assert.equal(treeBranchOpen({},'ui:projects'),true)
 const expanded={'Projects/A':true,'Projects/B':true,'ui:inbox':false}
 assert.deepEqual(toggleTreeBranch(expanded,'Projects/A'),{...expanded,'Projects/A':false})
 assert.deepEqual(revealTreeNote(expanded,'Projects/C'),{...expanded,'Projects/C':true,'ui:projects':true})
 assert.deepEqual(revealTreeNote(expanded),{...expanded,'ui:inbox':true})
})

test('project parent tabs retain their real paths and reject unsafe or container paths',()=>{
 for(const path of ['Projects/工作','项目/Thoughts','旧项目'])assert.equal(projectTabPath(projectTab(path)),path)
 for(const path of ['../outside','C:/private','Projects','Inbox','assets','Projects/../escape','.memo','Projects/CON','Projects/a/b'])assert.equal(projectTabPath(projectTab(path)),null)
 assert.deepEqual(pruneNoteTabs(['projects','project:Projects/gone','note:Projects/gone/draft.md'],new Set(),'Projects/gone/draft.md',new Set()),['projects','note:Projects/gone/draft.md'])
})

test('external deletion closes missing tabs but retains an unresolved draft for recovery',()=>{
 const tabs=['note:Projects/P/a.md','note:Inbox/live.md','trash','note:Projects/P/draft.md']
 assert.deepEqual(pruneNoteTabs(tabs,new Set(['Inbox/live.md']),'Projects/P/draft.md'),['note:Inbox/live.md','trash','note:Projects/P/draft.md'])
})

test('notes open separately, repeated selection reuses a tab, same names in different projects stay separate',()=>{
 const tasks=noteTab('Projects/First/待办.md'),ideas=noteTab('Projects/First/想法.md'),other=noteTab('Projects/Second/待办.md')
 let tabs=addTab([],tasks);tabs=addTab(tabs,ideas);tabs=addTab(tabs,tasks);tabs=addTab(tabs,other)
 assert.deepEqual(tabs,[tasks,ideas,other]);assert.equal(tabPath(ideas),'Projects/First/想法.md')
})
test('closing an inactive tab preserves selection; closing active selects right then left then empty',()=>{
 const tabs=['note:Inbox/a.md','note:Inbox/b.md','note:Inbox/c.md']
 assert.deepEqual(removeTab(tabs,tabs[1],tabs[0]),{tabs:tabs.slice(1),active:tabs[1]})
 assert.equal(removeTab(tabs,tabs[1],tabs[1]).active,tabs[2])
 assert.equal(removeTab(tabs,tabs[2],tabs[2]).active,tabs[1])
 assert.equal(removeTab([tabs[0]],tabs[0],tabs[0]).active,'empty')
})
test('rename/move retains tab order and deduplicates destination; cache rejects malformed entries and never stores content',()=>{
 assert.deepEqual(replaceNoteTab(['note:Inbox/a.md','search','note:Projects/P/b.md'],'Inbox/a.md','Projects/P/b.md'),['note:Projects/P/b.md','search'])
 assert.deepEqual(restoreTabs([null,'document','note:../outside.md','note:C:/outside.md','note:Inbox/a.md','note:Inbox/a.md','trash'],'Inbox/b.md'),['note:Inbox/a.md','trash','note:Inbox/b.md'])
 assert.deepEqual(restoreTabs({},''),[])
})

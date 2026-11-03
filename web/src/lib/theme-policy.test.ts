import {test} from 'node:test'
import assert from 'node:assert/strict'
import {isLightTheme} from './theme-policy.ts'

test('old dark caches and mixed dark surfaces cannot restore a dark interface',()=>{
 const keys='background card popover sidebar workspace-tab-background editor-background code-background'.split(' ')
 const light=Object.fromEntries(keys.map(key=>[key,'#fafafa']))
 assert.equal(isLightTheme(light),true)
 assert.equal(isLightTheme({...light,background:'#14212b'}),false)
 assert.equal(isLightTheme({...light,sidebar:'#18181b'}),false)
 assert.equal(isLightTheme({...light,background:'#ffffff00'}),false)
})

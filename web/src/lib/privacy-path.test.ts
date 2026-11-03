import {test} from 'node:test'
import assert from 'node:assert/strict'
import {locales} from './locales.ts'
import {privacyPath} from './privacy-path.ts'

test('each supported UI language links to its own local notice',()=>{
 assert.deepEqual(locales.map(privacyPath),['/privacy.html','/privacy-ko.html','/privacy-ja.html','/privacy-zh-Hans.html','/privacy-zh-Hant.html','/privacy-de.html','/privacy-fr.html','/privacy-es.html'])
 assert.equal(new Set(locales.map(privacyPath)).size,locales.length)
})

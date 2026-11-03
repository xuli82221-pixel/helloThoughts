import test from 'node:test'
import assert from 'node:assert/strict'
import {imageMarkdown,insertImageMarkdown,localImageURL} from './image-content.ts'

test('image insertion uses the selected range without changing surrounding text',()=>{const value=insertImageMarkdown('before selected after','before selected after',7,15,15,'![photo](assets/x.png)');assert.equal(value.content,'before \n\n![photo](assets/x.png)\n\n after')})
test('typing during upload is preserved and insertion follows the current caret',()=>{const value=insertImageMarkdown('old','old NEW TEXT',1,2,12,'![photo](assets/x.png)');assert.equal(value.content,'old NEW TEXT\n\n![photo](assets/x.png)\n\n')})
test('image filenames cannot inject Markdown or break relative URL syntax',()=>{assert.equal(imageMarkdown('a]b.png','../assets/a b.png'),'![a\\]b](../assets/a%20b.png)')})
test('remote and absolute images never generate external image requests',()=>{for(const value of ['https://example.com/a.png','//example.com/a.png','data:image/png;base64,x','file:///a.png','C:\\a.png','/a.png'])assert.equal(localImageURL('Inbox/A.md',value,'v'),undefined);assert.match(localImageURL('Inbox/A.md','assets/a.png','revision')!,/^\/api\/images\?/)})

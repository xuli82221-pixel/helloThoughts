const surfaces='background card popover sidebar workspace-tab-background editor-background code-background'.split(' ')
export function isLightTheme(tokens:Record<string,string>){
 return surfaces.every(key=>{
  const color=tokens[key]
  if(!/^#[0-9a-f]{6}(ff)?$/i.test(color??''))return false
  const brightness=parseInt(color.slice(1,3),16)*.2126+parseInt(color.slice(3,5),16)*.7152+parseInt(color.slice(5,7),16)*.0722
  return brightness>=180
 })
}

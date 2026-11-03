import type {SVGProps} from 'react'

type IconProps=SVGProps<SVGSVGElement>
function Icon({children,...props}:IconProps){
 return <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.7" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true" focusable="false" {...props}>{children}</svg>
}

/** Local outline icons inherit the current theme color. */
export function ProjectsIcon(props:IconProps){
 return <Icon {...props}><path d="M3 8V6a2 2 0 0 1 2-2h4l2 3h8a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8Z"/></Icon>
}

export function QuickNotesIcon(props:IconProps){
 return <Icon {...props}><path d="M12 6.1C8.9 3.6 5.4 3.3 2 5v14c3.4-1.7 6.9-1.4 10 1.1 3.1-2.5 6.6-2.8 10-1.1V5c-3.4-1.7-6.9-1.4-10 1.1Z"/><path d="M12 6.1v14"/></Icon>
}

export function SettingsIcon(props:IconProps){
 return <Icon {...props}><path d="m9.5 2-.6 2.7-1.6.9-2.6-.8-2.4 4.1 2 1.9v2.4l-2 1.9 2.4 4.1 2.6-.8 1.6.9.6 2.7h5l.6-2.7 1.6-.9 2.6.8 2.4-4.1-2-1.9v-2.4l2-1.9-2.4-4.1-2.6.8-1.6-.9-.6-2.7Z"/><circle cx="12" cy="12" r="3.6"/></Icon>
}

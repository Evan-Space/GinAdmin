import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/userList')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>
    <p>Hello "/userList"!</p>
    
  </div>
}

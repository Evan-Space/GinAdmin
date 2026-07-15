import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/log/errorLog')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/log/errorLog"!</div>
}

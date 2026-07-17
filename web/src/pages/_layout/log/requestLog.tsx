import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_layout/log/requestLog')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/log"!</div>
}

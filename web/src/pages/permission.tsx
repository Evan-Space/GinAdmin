import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/permission')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/permission"!</div>
}

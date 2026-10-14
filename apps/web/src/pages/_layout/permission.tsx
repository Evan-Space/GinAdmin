import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_layout/permission')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/permission"!</div>
}

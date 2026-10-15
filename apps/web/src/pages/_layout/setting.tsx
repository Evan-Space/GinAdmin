import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_layout/setting')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/setting"!</div>
}

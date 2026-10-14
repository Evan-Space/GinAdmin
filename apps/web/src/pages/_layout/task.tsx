import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_layout/task')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/task"!</div>
}

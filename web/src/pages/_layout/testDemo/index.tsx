import { createFileRoute } from '@tanstack/react-router'
import { UploadFile } from './-components/uploadFile'

export const Route = createFileRoute('/_layout/testDemo/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>
    <UploadFile />
  </div>
}

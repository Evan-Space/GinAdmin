import { createFileRoute } from '@tanstack/react-router'
import { CustomUploadFile } from './-components/customDemo/CustomUploadFile'

export const Route = createFileRoute('/_layout/UploadFile/')({
    component: RouteComponent,
})

function RouteComponent() {
    return <div>
        <CustomUploadFile />
    </div>
}

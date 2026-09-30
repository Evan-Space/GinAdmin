import { createFileRoute } from '@tanstack/react-router'
import { Index } from './-components/UploadFile'

export const Route = createFileRoute('/_layout/UploadFile/')({
    component: RouteComponent,
})

function RouteComponent() {
    return <div>
        <Index />
    </div>
}

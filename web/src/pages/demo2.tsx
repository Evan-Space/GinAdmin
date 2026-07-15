import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/demo2')({
  component: Demo2,
})

function Demo2() {
    return <div>Demo2</div>
}
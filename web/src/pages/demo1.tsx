import { createFileRoute } from '@tanstack/react-router'


export const Route = createFileRoute('/demo1')({
    component: Demo1,
  })

  

function Demo1() {
    return <div>Demo1</div>
}



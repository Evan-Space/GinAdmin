import { createRoot } from 'react-dom/client'
import '@src/style/global.css'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'

const router = createRouter({
    routeTree: routeTree,
})

// 注册路由类型，Link/useNavigate 才有类型提示
declare module '@tanstack/react-router' {
    interface Register {
        router: typeof router
    }
}


createRoot(document.getElementById('root') as HTMLElement).render(<RouterProvider router={router} />)
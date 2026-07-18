import { createFileRoute, redirect } from '@tanstack/react-router'
import Layout from '@src/layout'
import { getToken } from '@src/request/request'

export const Route = createFileRoute('/_layout')({
    component: Layout,
    beforeLoad: async () => {
        const token = getToken()
        if (!token) {
            throw redirect({ to: '/login' })
        }
    }
})

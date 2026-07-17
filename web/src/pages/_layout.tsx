import { createFileRoute } from '@tanstack/react-router'
import Layout from '@src/layout'

export const Route = createFileRoute('/_layout')({
    component: Layout,
})

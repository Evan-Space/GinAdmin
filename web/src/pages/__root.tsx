import { createRootRoute } from '@tanstack/react-router'
import { ConfigProvider } from 'antd'
import { appTheme } from '@src/style/antdTheme'
import Layout from '@src/layout'

export const Route = createRootRoute({
    component: () => (
        <ConfigProvider theme={appTheme}>
            <div className="w-screen h-screen">
                <Layout />
            </div>
        </ConfigProvider>
    ),
})

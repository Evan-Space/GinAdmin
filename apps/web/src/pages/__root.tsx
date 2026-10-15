import { createRootRoute, Outlet } from '@tanstack/react-router'
import { ConfigProvider, App } from 'antd'
import { appTheme } from '@src/style/antdTheme'

export const Route = createRootRoute({
    component: () => (
        <ConfigProvider theme={appTheme}>
            <App className="w-screen h-screen">
                <Outlet />
            </App>
        </ConfigProvider>
    ),
})

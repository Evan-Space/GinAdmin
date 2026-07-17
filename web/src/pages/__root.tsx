import { createRootRoute, Outlet } from '@tanstack/react-router'
import { ConfigProvider } from 'antd'
import { appTheme } from '@src/style/antdTheme'

export const Route = createRootRoute({
    component: () => (
        <ConfigProvider theme={appTheme}>
            <div className="w-screen h-screen">
                <Outlet />
            </div>
        </ConfigProvider>
    ),
})

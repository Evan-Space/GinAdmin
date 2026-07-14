import Layout from './layout/index.tsx'
import { ConfigProvider } from 'antd'
import { appTheme } from './style/antdTheme.ts'

export default function App() {
    return (
        <ConfigProvider theme={appTheme}>
            <div className="w-screen h-screen">
                <Layout />
            </div>
        </ConfigProvider>
    )
}

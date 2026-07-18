import { Form } from 'antd'
import { App } from 'antd'
import { useNavigate } from '@tanstack/react-router'
import { loginAPI } from '@src/request/login'


type LoginFormValues = {
    username: string
    password: string
}


export const useLogin = () => {
    const navigate = useNavigate()
    const [form] = Form.useForm<LoginFormValues>()
    const { message: messageApi } = App.useApp()
    const handleSubmit = async (values: LoginFormValues) => {
        try {

            const response = await loginAPI(values)
            if (response.code !== 0) {
                messageApi.error(response.msg || '登录失败')
                return
            }

            localStorage.setItem('auth_token', response.data.access_token)
            messageApi.success('登录成功')
            navigate({ to: '/' })
        } catch {
            messageApi.error('网络异常，请稍后重试')
        }
    }

    return {
        form,
        messageApi,
        handleSubmit,
    }
}
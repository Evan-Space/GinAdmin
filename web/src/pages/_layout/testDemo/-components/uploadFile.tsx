import { InboxOutlined } from '@ant-design/icons'
import { Button, Progress, Space, Upload } from 'antd'
import { UploadStatus, useChunkUpload } from '../hooks'

const STATUS_TEXT: Record<UploadStatus, string> = {
    idle: '等待选择文件',
    hashing: '正在计算文件指纹',
    uploading: '正在上传',
    merging: '服务端合并中',
    success: '上传完成',
    paused: '已暂停',
    error: '上传失败',
}

export const UploadFile = () => {
    const { status, progress, hashProgress, fileName, start, pause, resume, cancel } =
        useChunkUpload()
    const busy = status === 'hashing' || status === 'uploading' || status === 'merging'
    return (
        <div className="w-[520px] rounded-lg border border-gray-200 p-4">
            <Upload.Dragger
                multiple={false}
                showUploadList={false}
                disabled={busy}
                beforeUpload={(file) => {
                    start(file)
                    return false // 阻止 antd 自己上传，只借它的选择文件能力
                }}
            >
                <p className="text-3xl">
                    <InboxOutlined />
                </p>
                <p>点击或拖拽文件到此处上传</p>
            </Upload.Dragger>
            <div className="mt-4">
                <div className="text-sm text-gray-500">
                    {fileName || '未选择文件'} · {STATUS_TEXT[status]}
                </div>
                {status === 'hashing' ? (
                    <Progress percent={hashProgress} status="active" />
                ) : (
                    <Progress
                        percent={progress}
                        status={status === 'error' ? 'exception' : busy ? 'active' : undefined}
                    />
                )}
            </div>
            <Space className="mt-2">
                <Button onClick={pause} disabled={status !== 'uploading'}>
                    暂停
                </Button>
                <Button onClick={resume} disabled={status !== 'paused'}>
                    继续
                </Button>
                <Button
                    danger
                    onClick={cancel}
                    disabled={status === 'idle' || status === 'success'}
                >
                    取消
                </Button>
            </Space>
        </div>
    )
}

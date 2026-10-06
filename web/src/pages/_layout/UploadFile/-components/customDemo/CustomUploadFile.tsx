import { useCustomUploadFile } from './hooks.tsx'
import { Progress, Button } from 'antd'
export const CustomUploadFile = () => {
    const { status, handleFileChange, onProgress, handlePause, handleResume } = useCustomUploadFile()
    return (
        <div className="my-10 min-h-25 border border-solid border-#000">
            <input
                className={'border border-solid border-[#ccc]'}
                type="file"
                placeholder={'sss'}
                onChange={(event) => handleFileChange(event.target.files?.[0])}
            />
            <div className={'w-[80%] mx-auto'}>
                <Progress percent={onProgress} />
            </div>

            {status === 'uploading' && (
                <Button type="primary" onClick={handlePause}>
                    暂停
                </Button>
            )}
            {status === 'paused' && (
                <Button type="primary" onClick={handleResume}>
                    继续
                </Button>
            )}

            {status}
        </div>
    )
}

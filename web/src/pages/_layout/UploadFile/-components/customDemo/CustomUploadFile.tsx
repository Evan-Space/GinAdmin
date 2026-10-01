import { useCustomUploadFile } from './hooks.tsx'
import { Progress } from 'antd'
export const CustomUploadFile = () => {
    const { handleFileChange, onProgress } = useCustomUploadFile()
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
        </div>
    )
}

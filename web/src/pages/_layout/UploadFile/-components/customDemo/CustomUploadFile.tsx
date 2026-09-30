import { useCustomUploadFile } from './hooks.tsx'

export const CustomUploadFile = () => {
    const { handleFileChange } = useCustomUploadFile()
    return (
        <div className="my-10 min-h-25 border border-solid border-#000">
            <input
                className={'border border-solid border-[#ccc]'}
                type="file"
                multiple
                placeholder={'sss'}
                onChange={(event) => handleFileChange(event.target.files || [])}
            />
        </div>
    )
}

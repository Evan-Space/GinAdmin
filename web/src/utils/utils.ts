type PlainObject = Record<string, unknown>

/** 判断单个值是否为空 */
export function isEmptyValue(value: unknown): boolean {
    if (value === null || value === undefined) return true
    if (typeof value === 'string' && value.trim() === '') return true
    if (Array.isArray(value) && value.length === 0) return true
    return false
}

/** 删除对象中空值的 key，返回新对象（不修改原对象） */
export function omitEmptyValues<T extends PlainObject>(obj: T): Partial<T> {
    const result: Partial<T> = {}

    for (const key in obj) {
        if (!Object.prototype.hasOwnProperty.call(obj, key)) continue
        const value = obj[key]
        if (!isEmptyValue(value)) {
            result[key] = value as T[Extract<keyof T, string>]
        }
    }

    return result
}

/**
 * 处理请求错误
 */
export function handleRequestError(error: unknown) {
    if (error instanceof Error) {
        return error.message
    }
    return '请求错误'
}

/**
 * 控制请求并发数量
 * */
export const limitPromise = <T>(taskArr: Array<() => Promise<T>>, limit: number): Promise<T[]> => {
    return new Promise((resolve) => {
        const results = new Array(taskArr.length)
        let activeCount = 0
        let nextTaskINdex = 0

        const runNext = () => {
            if (nextTaskINdex >= taskArr.length) {
                if (activeCount === 0) {
                    resolve(results)
                }
                return
            }

            const currentTaskIndex = nextTaskINdex++
            activeCount++

            Promise.resolve()
                .then(taskArr[currentTaskIndex])
                .then(
                    (res) => {
                        results[currentTaskIndex] = res
                    },
                    (err) => {
                        results[currentTaskIndex] = err
                    },
                )
                .finally(() => {
                    activeCount--
                    runNext()
                })

            if (activeCount < limit) {
                runNext()
            }
        }

        runNext()
    })
}

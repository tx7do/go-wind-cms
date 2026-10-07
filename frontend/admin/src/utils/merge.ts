/**
 * 深度合并对象
 * @param target - 目标对象
 * @param source - 源对象
 * @returns 合并后的新对象
 */
export function merge<T extends Record<string, any>, U extends Record<string, any>>(
    target: T,
    source: U
): T & U {
    const result: any = {...target};

    for (const key in source) {
        if (Object.prototype.hasOwnProperty.call(source, key)) {
            const sourceValue = source[key];
            const targetValue = target[key as keyof T];

            // 如果两个值都是普通对象，则递归合并
            if (
                isPlainObject(targetValue) &&
                isPlainObject(sourceValue)
            ) {
                result[key] = merge(
                    targetValue as Record<string, any>,
                    sourceValue as Record<string, any>
                );
            }
            // 否则直接使用源值覆盖
            else {
                result[key] = sourceValue;
            }
        }
    }

    return result;
}

/**
 * 创建自定义合并函数
 * @param customizer - 自定义处理函数
 * @returns 合并函数
 */
export function createMerge(
    customizer?: (
        obj: Record<string, any>,
        key: string,
        value: any
    ) => boolean | void
) {
    return function <T extends Record<string, any>, U extends Record<string, any>>(
        target: T,
        source: U
    ): T & U {
        const result: any = {...target};

        for (const key in source) {
            if (Object.prototype.hasOwnProperty.call(source, key)) {
                const sourceValue = source[key];
                const targetValue = target[key as keyof T];

                // 如果有自定义处理函数，先尝试调用
                if (customizer) {
                    const handled = customizer(result, key, sourceValue);
                    if (handled === true) {
                        continue; // 如果自定义函数返回true，跳过默认处理
                    }
                }

                // 如果两个值都是普通对象，则递归合并
                if (
                    isPlainObject(targetValue) &&
                    isPlainObject(sourceValue)
                ) {
                    result[key] = createMerge(customizer)(
                        targetValue as Record<string, any>,
                        sourceValue as Record<string, any>
                    );
                }
                // 否则直接使用源值覆盖
                else {
                    result[key] = sourceValue;
                }
            }
        }

        return result;
    };
}

/**
 * 带数组覆盖功能的合并函数
 */
export const mergeWithArrayOverride = createMerge((obj, key, value) => {
    if (Array.isArray(obj[key]) && Array.isArray(value)) {
        obj[key] = value;
        return true; // 表示已处理，跳过默认逻辑
    }
    return false; // 表示未处理，继续默认逻辑
});

/**
 * 判断是否为普通对象
 * @param value - 要判断的值
 * @returns 是否为普通对象
 */
function isPlainObject(value: any): value is Record<string, any> {
    return (
        value !== null &&
        typeof value === 'object' &&
        Object.prototype.toString.call(value) === '[object Object]' &&
        !Array.isArray(value)
    );
}

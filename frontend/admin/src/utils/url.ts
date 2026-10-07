/**
 * 获取url中的全部参数
 * @param search - url参数
 */

export function getUrlAllParam(search: string): Record<string, string> {
    if (!search) return {};
    // 去除首个字符串问号
    if (search?.[0] === '?') search = search.substring(1, search.length);

    const arr = search.split('&'); // 分割数组
    const result: Record<string, string> = {};

    for (let i = 0; i < arr.length; i++) {
        const value = arr[i]?.split('=');
        if (value?.length === 2) {
            result[value[0]] = value[1];
        }
    }

    return result;
}

/**
 * 获取url中参数某个值
 * @param search - url参数
 * @param key - 搜索值
 */
export function getUrlParam(search: string, key: string) {
    if (!search || !key) return '';
    // 去除首个字符串问号
    if (search?.[0] === '?') search = search.substring(1, search.length);
    const urlParams = getUrlAllParam(search);
    return urlParams?.[key] || '';
}

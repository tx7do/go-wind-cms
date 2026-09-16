'use client';

import {useLocale} from 'next-intl';

import {resolveApiErrorMsg} from '@/core/transport/rest/utils';
import {allMessages, defaultLocale, validateLocale} from '@/i18n/config';

/**
 * API 错误文案解析：不再使用后端 message 字段，
 * 用 reason 查询 i18n 错误文案（error.<REASON> 命名空间），
 * 查不到时回退到默认错误文案。
 */
export function useApiError() {
    const locale = useLocale();

    return (error: unknown): string => {
        const dict = (allMessages[validateLocale(locale)] ?? allMessages[defaultLocale])?.error as Record<string, string> | undefined;
        return resolveApiErrorMsg(error, dict);
    };
}

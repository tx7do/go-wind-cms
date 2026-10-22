import Taro from '@tarojs/taro';
import {AxiosError, AxiosHeaders, type AxiosAdapter, type AxiosResponse, type InternalAxiosRequestConfig} from 'axios';

/**
 * axios ↔ Taro.request 适配器
 *
 * 小程序运行时没有 XHR / Node http，axios 默认适配器直接派不出请求。
 * 这里用 Taro.request 复刻默认适配器的契约：
 * - 拼 baseURL + url + params 序列化 query
 * - 响应按 config.validateStatus 决议 resolve/reject（AxiosError 带 response）
 * - 支持 config.signal 中止、config.timeout
 * - 仅在「无 XHR」的运行时启用（见 request-client.ts），H5 仍走 axios 自带 XHR 适配器
 */

const isAbsoluteURL = (url: string) => /^([a-z][a-z\d+\-.]*:)?\/\//i.test(url);

function combineURLs(baseURL: string, relativeURL: string) {
    return relativeURL
        ? baseURL.replace(/\/+$/, '') + '/' + relativeURL.replace(/^\/+/, '')
        : baseURL;
}

function buildFullPath(baseURL: string | undefined, requestedURL: string) {
    if (baseURL && !isAbsoluteURL(requestedURL)) {
        return combineURLs(baseURL, requestedURL);
    }
    return requestedURL;
}

function encodeParam(value: unknown): string {
    if (value === null || value === undefined) return '';
    if (typeof value === 'object') return JSON.stringify(value);
    return String(value);
}

function buildURL(url: string, params?: Record<string, unknown> | null) {
    if (!params) return url;
    const parts: string[] = [];
    for (const [key, value] of Object.entries(params)) {
        if (value === null || value === undefined) continue;
        if (Array.isArray(value)) {
            for (const item of value) {
                if (item !== null && item !== undefined) parts.push(`${encodeURIComponent(key)}=${encodeURIComponent(encodeParam(item))}`);
            }
        } else {
            parts.push(`${encodeURIComponent(key)}=${encodeURIComponent(encodeParam(value))}`);
        }
    }
    if (!parts.length) return url;
    return url + (url.includes('?') ? '&' : '?') + parts.join('&');
}

const KNOWN_STATUS_ERROR = (status: number): 'ERR_BAD_REQUEST' | 'ERR_BAD_RESPONSE' =>
    status >= 500 ? 'ERR_BAD_RESPONSE' : 'ERR_BAD_REQUEST';

export const taroRequestAdapter: AxiosAdapter = (config: InternalAxiosRequestConfig) =>
    new Promise((resolve, reject) => {
        const fullPath = buildFullPath(config.baseURL, config.url ?? '');
        const url = buildURL(fullPath, config.params as Record<string, unknown> | null);

        // AxiosHeaders 归一化后摊平成 Taro 要的普通键值对
        const headers: Record<string, string> = {};
        AxiosHeaders.from(config.headers).forEach((value, key) => {
            if (value !== null && value !== undefined) headers[key] = String(value);
        });

        let task: Taro.RequestTask<unknown> | null = null;
        const onCanceled = () => {
            task?.abort?.();
            reject(new AxiosError('canceled', AxiosError.ERR_CANCELED, config));
        };
        // 最小信号接口：本项目 tsconfig 未开 DOM lib，axios 引用的 AbortSignal
        // 全局类型在这里不可靠，按结构收窄
        const signal = config.signal as {aborted: boolean; addEventListener(type: string, fn: () => void): void; removeEventListener(type: string, fn: () => void): void} | undefined;
        if (signal) {
            if (signal.aborted) {
                onCanceled();
                return;
            }
            signal.addEventListener('abort', onCanceled);
        }

        // 小程序仅支持 text / arrayBuffer；JSON 解析交给 Taro 的 dataType='json'
        const wantsArrayBuffer = config.responseType === 'arraybuffer' || config.responseType === 'blob';

        task = Taro.request({
            url,
            method: (config.method ?? 'get').toUpperCase() as keyof Taro.request.Method,
            data: config.data as string | Record<string, unknown> | undefined,
            header: headers,
            timeout: config.timeout,
            dataType: wantsArrayBuffer ? 'other' : 'json',
            responseType: wantsArrayBuffer ? 'arraybuffer' : 'text',
            enableHttp2: true,
            success: (res) => {
                const response: AxiosResponse = {
                    data: res.data as AxiosResponse['data'],
                    status: res.statusCode,
                    statusText: '',
                    headers: (res.header ?? {}) as AxiosResponse['headers'],
                    config,
                    request: task,
                };
                const validate = config.validateStatus ?? ((status: number) => status >= 200 && status < 300);
                if (!res.statusCode || validate(res.statusCode)) {
                    resolve(response);
                } else {
                    reject(new AxiosError(
                        `Request failed with status code ${res.statusCode}`,
                        KNOWN_STATUS_ERROR(res.statusCode),
                        config,
                        task,
                        response,
                    ));
                }
            },
            fail: (err) => {
                if (signal?.aborted) {
                    reject(new AxiosError('canceled', AxiosError.ERR_CANCELED, config));
                    return;
                }
                reject(new AxiosError(
                    err?.errMsg || 'Network Error',
                    AxiosError.ERR_NETWORK,
                    config,
                    task,
                    undefined,
                ));
            },
            complete: () => {
                signal?.removeEventListener('abort', onCanceled);
            },
        });
    });

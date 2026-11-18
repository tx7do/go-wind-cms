/**
 * ApiClient 单例
 *
 * 生成的代码提供了统一的 ApiClient，它通过 ClientTransport 接口发送请求。
 * 这里将已有的 requestApi（基于 axios 的 RequestClient）适配为 ClientTransport，
 * 保留 token 注入、错误拦截、自动刷新等全部已有逻辑。
 */
import {
    type ClientTransport,
    createApiClient,
} from '@/api/generated/app/service/v1';
import {requestApi} from '@/core/transport/rest';

// 渐进式登录验证码的 header 状态（模块级）：登录页声明支持 + 填入验证码 id/value，
// 只在后续请求上携带。空值不进 header——恒传空值会把空头发给后端（空 header
// 即校验失败，CAPTCHA_REQUIRED 点亮流程会被打断）。
const captchaHeaders: Record<string, string> = {};

/** 声明本客户端支持验证码 UI（后端仅对声明方启用验证码 gate）。 */
export function declareCaptchaSupport(): void {
    captchaHeaders['X-Captcha-Supported'] = '1';
}

/** 设置/清除验证码 id 与用户输入（id/value 齐备才携带）。 */
export function setCaptchaValues(id: string | undefined, value: string | undefined): void {
    if (id && value) {
        captchaHeaders['X-Captcha-Id'] = id;
        captchaHeaders['X-Captcha-Value'] = value;
    } else {
        delete captchaHeaders['X-Captcha-Id'];
        delete captchaHeaders['X-Captcha-Value'];
    }
}

const transport: ClientTransport = {
    unary(path, method, body, _meta) {
        return requestApi({
            body,
            method,
            path,
            headers: Object.keys(captchaHeaders).length > 0 ? {...captchaHeaders} : undefined,
        });
    },
    serverStream(path, _meta) {
        // SSE 流由 transport/sse 模块独立管理，此处不应被调用
        throw new Error(`serverStream not supported via ApiClient: ${path}`);
    },
    byteStream(path, _meta) {
        // HttpBody 字节流不经过 ApiClient，此处不应被调用
        throw new Error(`byteStream not supported via ApiClient: ${path}`);
    },
    duplexStream(path, _meta) {
        throw new Error(`duplexStream not supported via ApiClient: ${path}`);
    },
};

export const apiClient = createApiClient(transport);

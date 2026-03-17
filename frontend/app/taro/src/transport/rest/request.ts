import {createRequestClient, requestClient} from "@/transport/rest/rest-client";
import {useAccessStore} from "@/store/core/access/hooks";
import {useLanguageStore} from "@/store/core/language/hooks";


/**
 * 通用请求处理器 Hook
 */
export function useRequest() {

    const accessStore = useAccessStore();
    const languageStore = useLanguageStore();

    const requestClient = createRequestClient(
        API_BASE_URL,
        () => languageStore.language.locale,
        () => accessStore.access.accessToken?.value || ''
    );

    /**
     * 发起请求
     */
    function request({path, method, body}: { path: string; method: string; body: null | string }) {
        return requestClient.request(path, {
            method,
            data: body,
        } as never);
    }

    return request;
}

export function requestApi({path, method, body}: { path: string; method: string; body: null | string }) {
    return requestClient.request(path, {
        method,
        data: body,
    } as never);
}

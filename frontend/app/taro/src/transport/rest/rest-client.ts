import {RequestClient} from "@/transport/rest/request-client";
import {HttpResponse} from "@/transport/rest/types";
import {authenticateResponseInterceptor, errorMessageResponseInterceptor} from "@/transport/rest/preset-interceptors";

import store from "@/store";


export function createRequestClient(baseURL: string, getLocale?: () => string, getToken?: () => string) {
  const client = new RequestClient({
    baseURL,
  });

  // 格式化令牌
  function formatToken(token: null | string) {
    return token ? `Bearer ${token}` : null;
  }

  // 请求头处理
  client.addRequestInterceptor({
    fulfilled: (config) => {
      if (getToken) {
        config.headers.Authorization = formatToken(getToken());
      }
      if (getLocale) {
        config.headers['Accept-Language'] = getLocale();
      }
      return config as never;
    },
  });

  // response数据解构
  client.addResponseInterceptor({
    fulfilled: (response) => {
      const {data: responseData, status} = response;

      // TODO 根据Kratos进行定制

      if (status >= 200 && status < 400) {
        return responseData;
      }

      const {code} = responseData as HttpResponse;
      if (code !== null) {
        throw Object.assign({}, responseData, {responseData});
      }

      console.error('parse HttpResponse failed!', response);
      throw Object.assign({}, response, {response});
    },
  });

  // token过期的处理
  client.addResponseInterceptor(
    authenticateResponseInterceptor({
      client,
      doReAuthenticate: async () => {
      },
      doRefreshToken: async () => '',
      enableRefreshToken: true,
      formatToken,
    }),
  );

  // 通用的错误处理,如果没有进入上面的错误处理逻辑，就会进入这里
  client.addResponseInterceptor(
    errorMessageResponseInterceptor(async (msg: string, error) => {
      const responseData = (error as unknown as {
        response?: { data?: Record<string, unknown> }
      })?.response?.data ?? {};
      const errorMessage = responseData?.error ?? responseData?.message ?? '';
      window.alert(errorMessage || msg);
    }),
  );

  return client;
}

// Taro 项目使用自定义的环境变量配置
// 通过 defineConstants 注入的全局常量
function getAccessToken() {
  return store.getState().access.accessToken?.value || '';
}

function getLocale() {
  return store.getState().language.locale || 'zh-CN';
}

export const requestClient = createRequestClient(
  API_BASE_URL,
  getLocale,
  getAccessToken
);

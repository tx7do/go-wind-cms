import type { AxiosInstance, AxiosResponse } from 'axios';

import type { RequestInterceptorConfig, ResponseInterceptorConfig } from '../types';

const defaultRequestInterceptorConfig: RequestInterceptorConfig = {
  fulfilled: (config) => config,
  rejected: (error) => Promise.reject(error),
};

const defaultResponseInterceptorConfig: ResponseInterceptorConfig = {
  fulfilled: (response: AxiosResponse) => response,
  rejected: (error) => Promise.reject(error),
};

class InterceptorManager {
  private axiosInstance: AxiosInstance;

  constructor(instance: AxiosInstance) {
    this.axiosInstance = instance;
  }

  addRequestInterceptor({
    fulfilled,
    rejected,
  }: RequestInterceptorConfig = defaultRequestInterceptorConfig) {
    this.axiosInstance.interceptors.request.use(fulfilled, rejected);
  }

  addResponseInterceptor<T = never>({
    fulfilled,
    rejected,
  }: ResponseInterceptorConfig<T> = defaultResponseInterceptorConfig as ResponseInterceptorConfig<T>) {
    this.axiosInstance.interceptors.response.use(fulfilled, rejected);
  }
}

export { InterceptorManager };

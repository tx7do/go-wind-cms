import type { AxiosInstance, AxiosResponse, InternalAxiosRequestConfig } from 'axios';

import type { RequestInterceptorConfig, ResponseInterceptorConfig } from '../types';

const defaultRequestInterceptorConfig: RequestInterceptorConfig = {
  fulfilled: (config) => config as InternalAxiosRequestConfig<never>,
  rejected: (error) => Promise.reject(error) as never,
};

const defaultResponseInterceptorConfig: ResponseInterceptorConfig = {
  fulfilled: (response: AxiosResponse) => response,
  rejected: (error) => Promise.reject(error) as never,
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
    this.axiosInstance.interceptors.request.use(
      fulfilled as never,
      rejected as never,
    );
  }

  addResponseInterceptor<T = never>(
    config: ResponseInterceptorConfig<T> = defaultResponseInterceptorConfig as ResponseInterceptorConfig<T>,
  ) {
    const {fulfilled, rejected} = config;
    this.axiosInstance.interceptors.response.use(
      fulfilled as never,
      rejected as never,
    );
  }
}

export { InterceptorManager };

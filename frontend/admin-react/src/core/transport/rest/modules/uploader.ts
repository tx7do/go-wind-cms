import type { AxiosRequestConfig, AxiosResponse } from 'axios';

import type { RequestClient } from '../request-client';

class FileUploader {
  private client: RequestClient;

  constructor(client: RequestClient) {
    this.client = client;
  }

  public async upload(
    url: string,
    data: { file: Blob | File } & Record<string, any>,
    config?: AxiosRequestConfig,
  ): Promise<AxiosResponse> {
    const formData = new FormData();

    Object.entries(data).forEach(([key, value]) => {
      formData.append(key, value);
    });

    // 必须显式声明 multipart：实例默认头是 application/json，axios 1.x 对
    // "FormData + JSON 头" 会 formDataToJSON 序列化成 JSON 体，后端解析不到
    // multipart 的 file 字段。声明 multipart 后，xhr adapter 检测到 FormData
    // 会剥掉该头，由浏览器自动生成带 boundary 的 Content-Type。
    const finalConfig: AxiosRequestConfig = {
      ...config,
      headers: {
        'Content-Type': 'multipart/form-data',
        ...config?.headers,
      },
    };

    // 使用 request 方法代替 post 方法，避免类型推断问题
    return this.client.request(url, {
      method: 'POST',
      data: formData,
      ...finalConfig,
    });
  }
}

export { FileUploader };

import { RequestClient } from './request-client';

export function requestApi({
  path,
  method,
  body,
  headers,
}: {
  path: string;
  method: string;
  body: null | string;
  headers?: Record<string, string>;
}) {
  return RequestClient.getInstance().request(path, {
    method,
    data: body,
    headers,
  } as never);
}

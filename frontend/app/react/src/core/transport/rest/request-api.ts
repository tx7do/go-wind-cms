import { RequestClient } from './request-client';

export function requestApi({
  path,
  method,
  body,
}: {
  path: string;
  method: string;
  body: null | string;
}) {
  return RequestClient.getInstance().request(path, {
    method,
    data: body,
  } as never);
}

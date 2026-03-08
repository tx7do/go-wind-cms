import {defineStore} from 'pinia';

import {
  createCommentServiceClient,
} from '@/api/generated/app/service/v1';
import {type Paging, requestClientRequestHandler} from "@/transport/rpc/request";
import {useUserStore} from '@/stores';
import {makeOrderBy, makeQueryString, makeUpdateMask} from "@/transport/rpc";

export const useCommentStore = defineStore('comment', () => {
  const service = createCommentServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询评论列表
   */
  async function listComment(
    paging?: Paging,
    formValues?: null | object,
    fieldMask?: null | string,
    orderBy?: null | string[],
  ) {
    const noPaging =
      paging?.page === undefined && paging?.pageSize === undefined;
    // @ts-ignore proto generated code is error.
    return await service.ListComment({
      fieldMask,
      orderBy: makeOrderBy(orderBy),
      query: makeQueryString(formValues, userStore.isTenantUser()),
      page: paging?.page,
      pageSize: paging?.pageSize,
      noPaging,
    });
  }

  /**
   * 获取评论
   */
  async function getComment(id: number) {
    return await service.Get({id});
  }

  /**
   * 创建评论
   */
  async function createComment(values: Record<string, any> = {}) {
    return await service.Create({
      // @ts-ignore proto generated code is error.
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新评论
   */
  async function updateComment(id: number, values: Record<string, any> = {}) {
    return await service.Update({
      id,
      // @ts-ignore proto generated code is error.
      data: {
        ...values,
      },
      updateMask: makeUpdateMask(Object.keys(values ?? [])),
    });
  }

  /**
   * 删除评论
   */
  async function deleteComment(id: number) {
    return await service.Delete({id});
  }

  function $reset() {
  }

  return {
    $reset,
    listComment,
    getComment,
    createComment,
    updateComment,
    deleteComment,
  };
});

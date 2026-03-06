import {defineStore} from 'pinia';

import {
  createPostServiceClient,
} from '@/api/generated/app/service/v1';
import {type Paging, requestClientRequestHandler} from "@/transport/rpc/request";
import {makeOrderBy, makeQueryString, makeUpdateMask} from "@/transport/rpc";
import {useUserStore} from "@/stores";
import {currentLocaleLanguageCode} from "@/locales";

export const usePostStore = defineStore('post', () => {
  const service = createPostServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询帖子列表
   */
  async function listPost(
    paging?: Paging,
    formValues?: null | object,
    fieldMask?: null | string,
    orderBy?: null | string[],
  ) {
    const locale = currentLocaleLanguageCode();
    formValues = {
      ...formValues,
      locale,
    };

    const noPaging =
      paging?.page === undefined && paging?.pageSize === undefined;
    // @ts-ignore proto generated code is error.
    return await service.ListPost({
      fieldMask,
      orderBy: makeOrderBy(orderBy),
      query: makeQueryString(formValues, userStore.isTenantUser()),
      page: paging?.page,
      pageSize: paging?.pageSize,
      noPaging,
    });
  }

  /**
   * 获取帖子
   */
  async function getPost(id: number) {
    if (!id) return null;

    const locale = currentLocaleLanguageCode();

    return await service.GetPost({id, locale});
  }

  /**
   * 创建帖子
   */
  async function createPost(values: Record<string, any> = {}) {
    return await service.CreatePost({
      // @ts-ignore proto generated code is error.
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新帖子
   */
  async function updatePost(id: number, values: Record<string, any> = {}) {
    return await service.UpdatePost({
      id,
      // @ts-ignore proto generated code is error.
      data: {
        ...values,
      },
      updateMask: makeUpdateMask(Object.keys(values ?? [])),
    });
  }

  /**
   * 删除帖子
   */
  async function deletePost(id: number) {
    return await service.DeletePost({id});
  }

  function $reset() {
  }

  return {
    $reset,
    listPost,
    getPost,
    createPost,
    updatePost,
    deletePost,
  };
});

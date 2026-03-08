import {defineStore} from 'pinia';

import {
  type contentservicev1_Post,
  type contentservicev1_PostTranslation,
  createPostServiceClient,
} from '@/api/generated/app/service/v1';
import {type Paging, requestClientRequestHandler} from "@/transport/rpc/request";
import {makeOrderBy, makeQueryString, makeUpdateMask} from "@/transport/rpc";
import {useUserStore} from '@/stores';
import {$t, currentLocaleLanguageCode} from "@/locales";

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
    return await service.List({
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

    return await service.Get({id, locale});
  }

  /**
   * 创建帖子
   */
  async function createPost(values: Record<string, any> = {}) {
    return await service.Create({
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
   * 删除帖子
   */
  async function deletePost(id: number) {
    return await service.Delete({id});
  }

  /**
   * 获取帖子的翻译
   */
  function getTranslation(post: contentservicev1_Post) {
    if (!post?.translations || post.translations.length === 0) return null;

    const locale = currentLocaleLanguageCode();
    // 优先查找当前语言的翻译
    const translation = post.translations?.find((t: contentservicev1_PostTranslation) => t.languageCode === locale);
    // 如果找不到，回退到第一个翻译
    return translation || post.translations?.[0];
  }

  function getPostTitle(post: contentservicev1_Post) {
    const translation = getTranslation(post)
    return translation?.title || $t('page.post_detail.untitled')
  }

  function getPostSummary(post: contentservicev1_Post) {
    const translation = getTranslation(post)
    return translation?.summary || ''
  }

  function getPostThumbnail(post: contentservicev1_Post) {
    const translation = getTranslation(post)
    return translation?.thumbnail || '/placeholder.jpg'
  }

  function getPostContent(post: contentservicev1_Post) {
    const translation = getTranslation(post)
    return translation?.content || ''
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
    getTranslation,
    getPostTitle,
    getPostSummary,
    getPostThumbnail,
    getPostContent,
  };
});

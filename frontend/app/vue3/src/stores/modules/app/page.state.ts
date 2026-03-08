import {defineStore} from 'pinia';

import {
  type contentservicev1_Page,
  type contentservicev1_PageTranslation,
  createPageServiceClient,
} from '@/api/generated/app/service/v1';
import {type Paging, requestClientRequestHandler} from "@/transport/rpc/request";
import {useUserStore} from "@/stores";
import {makeOrderBy, makeQueryString, makeUpdateMask} from "@/transport/rpc";
import {currentLocaleLanguageCode} from "@/locales";

export const usePageStore = defineStore('page', () => {
  const service = createPageServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询页面列表
   */
  async function listPage(
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
   * 获取页面
   */
  async function getPage(id: number) {
    return await service.Get({id});
  }

  /**
   * 创建页面
   */
  async function createPage(values: Record<string, any> = {}) {
    return await service.Create({
      // @ts-ignore proto generated code is error.
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新页面
   */
  async function updatePage(id: number, values: Record<string, any> = {}) {
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
   * 删除页面
   */
  async function deletePage(id: number) {
    return await service.Delete({id});
  }

  /**
   * 获取页面的翻译
   */
  function getTranslation(page: contentservicev1_Page) {
    if (!page?.translations || page.translations.length === 0) return null;

    const locale = currentLocaleLanguageCode();
    // 优先查找当前语言的翻译
    const translation = page.translations?.find((t: contentservicev1_PageTranslation) => t.languageCode === locale);
    // 如果找不到，回退到第一个翻译
    return translation || page.translations?.[0];
  }

  function $reset() {
  }

  return {
    $reset,
    listPage,
    getPage,
    createPage,
    updatePage,
    deletePage,
    getTranslation,
  };
});

import {defineStore} from 'pinia';

import {
  type contentservicev1_Category,
  type contentservicev1_CategoryTranslation,
  createCategoryServiceClient,
} from '@/api/generated/app/service/v1';
import {type Paging, requestClientRequestHandler} from "@/transport/rpc/request";
import {makeOrderBy, makeQueryString, makeUpdateMask} from "@/transport/rpc";
import {useUserStore} from "@/stores";
import {currentLocaleLanguageCode} from "@/locales";

export const useCategoryStore = defineStore('category', () => {
  const service = createCategoryServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询分类列表
   */
  async function listCategory(
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
    return await service.ListCategory({
      fieldMask,
      orderBy: makeOrderBy(orderBy),
      query: makeQueryString(formValues, userStore.isTenantUser()),
      page: paging?.page,
      pageSize: paging?.pageSize,
      noPaging,
    });
  }

  /**
   * 获取分类
   */
  async function getCategory(id: number) {
    return await service.GetCategory({id});
  }

  /**
   * 创建分类
   */
  async function createCategory(values: Record<string, any> = {}) {
    return await service.CreateCategory({
      // @ts-ignore proto generated code is error.
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新分类
   */
  async function updateCategory(id: number, values: Record<string, any> = {}) {
    return await service.UpdateCategory({
      id,
      // @ts-ignore proto generated code is error.
      data: {
        ...values,
      },
      updateMask: makeUpdateMask(Object.keys(values ?? [])),
    });
  }

  /**
   * 删除分类
   */
  async function deleteCategory(id: number) {
    return await service.DeleteCategory({id});
  }

  /**
   * 获取分类的翻译
   */
  function getTranslation(category: contentservicev1_Category) {
    if (!category?.translations || category.translations.length === 0) return null;

    const locale = currentLocaleLanguageCode();
    // 优先查找当前语言的翻译
    const translation = category.translations?.find((t: contentservicev1_CategoryTranslation) => t.languageCode === locale);
    // 如果找不到，回退到第一个翻译
    return translation || category.translations?.[0];
  }

  function $reset() {
  }

  return {
    $reset,
    listCategory,
    getCategory,
    createCategory,
    updateCategory,
    deleteCategory,
    getTranslation,
  };
});

import {defineStore} from 'pinia';

import {
  createNavigationServiceClient
} from '@/api/generated/app/service/v1';
import {type Paging, requestClientRequestHandler} from "@/transport/rpc/request";
import {makeOrderBy, makeQueryString, makeUpdateMask} from "@/transport/rpc";
import {useUserStore} from "@/stores";
import {currentLocaleLanguageCode} from "@/locales";

export const useNavigationStore = defineStore('navigation', () => {
  const service = createNavigationServiceClient(requestClientRequestHandler);
  const userStore = useUserStore();

  /**
   * 查询导航列表
   */
  async function listNavigation(
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
   * 获取导航
   */
  async function getNavigation(id: number) {
    return await service.Get({id});
  }

  /**
   * 创建导航
   */
  async function createNavigation(values: Record<string, any> = {}) {
    return await service.Create({
      // @ts-ignore proto generated code is error.
      data: {
        ...values,
      },
    });
  }

  /**
   * 更新导航
   */
  async function updateNavigation(
    id: number,
    values: Record<string, any> = {},
  ) {
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
   * 删除导航
   */
  async function deleteNavigation(id: number) {
    return await service.Delete({id});
  }

  function $reset() {
  }

  return {
    $reset,
    listNavigation,
    getNavigation,
    createNavigation,
    updateNavigation,
    deleteNavigation,
  };
});

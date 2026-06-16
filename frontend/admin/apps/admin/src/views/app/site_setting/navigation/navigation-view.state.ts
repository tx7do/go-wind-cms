import { defineStore } from 'pinia';

import {
  fetchListNavigationItems,
  fetchListNavigations,
  type siteservicev1_ListNavigationItemResponse as ListNavigationItemResponse,
  type siteservicev1_ListNavigationResponse as ListNavigationResponse,
  PaginationQuery,
} from '#/api';

/**
 * 导航视图状态接口
 */
interface NavigationViewState {
  loading: boolean; // 加载状态
  needReloadNavigationList: boolean; // 是否需要重新加载导航列表
  needReloadNavigationItemList: boolean; // 是否需要重新加载导航项列表

  currentNavigationId: null | number; // 当前选中的导航ID
  navigationList: ListNavigationResponse; // 导航列表
  itemList: ListNavigationItemResponse; // 导航项列表
}

/**
 * 导航视图状态
 */
export const useNavigationViewStore = defineStore('permission-view', {
  state: (): NavigationViewState => ({
    currentNavigationId: null,
    loading: false,
    needReloadNavigationList: false,
    needReloadNavigationItemList: false,
    navigationList: { items: [], total: 0 },
    itemList: { items: [], total: 0 },
  }),

  actions: {
    /**
     * 获取导航列表
     */
    async fetchNavigationList(
      currentPage: number,
      pageSize: number,
      formValues: any,
    ) {
      this.loading = true;
      try {
        this.navigationList = await fetchListNavigations(
          new PaginationQuery({
            paging: { page: currentPage, pageSize },
            formValues,
          }),
        );
        return this.navigationList;
      } catch (error) {
        console.error('获取导航列表失败:', error);
        this.resetNavigationList();
      } finally {
        this.loading = false;
      }

      return this.navigationList;
    },

    /**
     * 根据导航ID获取导航项列表
     * @param navigationId 导航ID
     * @param currentPage
     * @param pageSize
     * @param formValues
     */
    async fetchItemList(
      navigationId: null | number,
      currentPage: number,
      pageSize: number,
      formValues: any,
    ) {
      if (!navigationId) {
        this.resetItemList();
        return this.itemList;
      }

      this.loading = true;
      try {
        this.itemList = await fetchListNavigationItems(
          new PaginationQuery({
            paging: { page: currentPage, pageSize },
            formValues: {
              ...formValues,
              navigation_id: navigationId.toString(),
            },
          }),
        );
      } catch (error) {
        console.error(`获取导航[${navigationId}]的导航项失败:`, error);
        this.resetItemList();
      } finally {
        this.loading = false;
      }

      return this.itemList;
    },

    /**
     * 设置当前选中的导航ID，并联动刷新导航列表
     * @param navigationId 导航ID
     */
    async setCurrentNavigationId(navigationId: number) {
      this.currentNavigationId = navigationId; // 更新当前选中的分组ID
      this.needReloadNavigationItemList = true;
    },

    resetNavigationList() {
      this.navigationList = { items: [], total: 0 };
    },

    resetItemList() {
      this.itemList = { items: [], total: 0 };
    },

    reloadNavigationList() {
      this.navigationList = { items: [], total: 0 };
      this.itemList = { items: [], total: 0 };
      this.currentNavigationId = 0;
      this.needReloadNavigationList = true;
      this.needReloadNavigationItemList = true;
    },
  },
});

import {defineStore} from "pinia";
import {navigateTo} from "@/router";

interface NavbarState {
  /**
   * 是否展示侧边栏
   */
  isShowSidebar: boolean;

  /**
   * 是否展示顶部栏
   */
  isShowNavbar: boolean;

  /**
   * 是否展示底部栏
   */
  isShowFooter: boolean;

  /**
   * 是否展示搜索栏
   */
  isShowSearchBar: boolean;

  /**
   * 悬浮层
   */
  activeOverlay: string | undefined;
}

export const useNavbarStore = defineStore('navbar', {
  actions: {
    /**
     * 设置侧边栏是否展示
     * @param isShow
     */
    setShowSidebar(isShow: boolean) {
      this.isShowSidebar = isShow;
    },

    /**
     * 设置顶部栏是否展示
     * @param isShow
     */
    setShowNavbar(isShow: boolean) {
      this.isShowNavbar = isShow;
    },

    /**
     * 设置底部栏是否展示
     * @param isShow
     */
    setShowFooter(isShow: boolean) {
      this.isShowFooter = isShow;
    },

    /**
     * 设置搜索栏是否展示
     * @param isShow
     */
    setShowSearchBar(isShow: boolean) {
      this.isShowSearchBar = isShow;
    },

    /**
     * 设置悬浮层是否展示
     * @param tab
     */
    setActiveOverlay(tab: string | undefined) {
      this.activeOverlay = tab;
    },

    async navigateToAndHideOverlay(path: string, force: boolean = false) {
      await navigateTo(path, force);
      this.setActiveOverlay(null);
    }
  },

  getters: {
    getShowSidebar: (state: NavbarState) => state.isShowSidebar,
    getShowNavbar: (state: NavbarState) => state.isShowNavbar,
    getShowFooter: (state: NavbarState) => state.isShowFooter,
    getShowSearchBar: (state: NavbarState) => state.isShowSearchBar,
    getActiveOverlay: (state: NavbarState) => state.activeOverlay,
  },

  state: (): NavbarState => ({
    isShowSidebar: false,
    isShowNavbar: true,
    isShowFooter: true,
    isShowSearchBar: true,
    activeOverlay: null,
  }),
});

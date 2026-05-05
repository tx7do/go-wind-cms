import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router';
import { routes } from 'vue-router/auto-routes';

import type {
  Router,
  RouterHistory
} from 'vue-router';

import type {EnhancedRouteLocation} from './types';
import {useRouteCacheStore} from '@/stores';

import NProgress from '@/utils/progress';
import {RouteMode} from "./types";

class RouteView {
  // 路由对象
  private readonly router: Router | undefined = undefined;

  constructor() {
    this.router = this.createRouter();
    this.beforeRouteChange();
    this.afterRouteChange();
  }

  /**
   * 根据环境变量中的配置生成路由模式
   */
  private createHistory = (): RouterHistory => {
    if (import.meta.env.VITE_ROUTER_MODE === RouteMode.HISTORY) {
      return createWebHistory();
    } else {
      return createWebHashHistory();
    }
  }

  /**
   * 创建路由对象
   * @private
   */
  private createRouter(): Router {
    return createRouter({
      history: this.createHistory(),
      routes,
    });
  }

  /**
   * 路由守卫
   * @private
   */
  private beforeRouteChange(): void {
    if (!this.router) {
      return;
    }

    this.router.beforeEach((to: EnhancedRouteLocation, _from, next) => {
      NProgress.start();

      const routeCacheStore = useRouteCacheStore();

      // Route cache
      routeCacheStore.addRoute(to);

      next();
    });
  }

  private afterRouteChange(): void {
    if (!this.router) {
      return;
    }

    this.router.afterEach(() => {
      NProgress.done();
    });
  }

  /**
   * 获取路由对象
   */
  public getRouter(): Router {
    return this.router as Router;
  }
}

const routeView = new RouteView();

export default routeView.getRouter();

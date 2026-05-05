import {createApp} from 'vue'

import App from '@/App.vue'

import {setupRouter} from '@/router'
import {setupStore} from '@/stores'
import {setupI18n} from '@/locales'
import {setupNaive} from "@/plugins/naive";
import {setupNaiveDiscreteApi, setupCustomComponents, setupGlobalMethods} from "@/plugins";

import 'virtual:uno.css'
import '@/styles/app.less'
import '@/styles/reset.min.css'

async function bootstrap(namespace: string) {
  const app = createApp(App);

  // 注册国际化
  await setupI18n(app);

  // 挂载状态管理
  await setupStore(app, {namespace});

  // 注册全局常用的 naive-ui 组件
  setupNaive(app);

  // 挂载 naive-ui 脱离上下文的 Api
  setupNaiveDiscreteApi();

  // 注册全局自定义组件
  setupCustomComponents(app);

  // 注册全局方法，如：app.config.globalProperties.$message = message
  setupGlobalMethods(app);

  // 注册路由模块
  setupRouter(app);

  app.mount('#app');
}

async function initApplication() {
  // name用于指定项目唯一标识
  // 用于区分不同项目的偏好设置以及存储数据的key前缀以及其他一些需要隔离的数据
  const env = import.meta.env.PROD ? 'prod' : 'dev';
  const appVersion = import.meta.env.VITE_APP_VERSION;
  const namespace = `${import.meta.env.VITE_APP_NAMESPACE}-${appVersion}-${env}`;

  await bootstrap(namespace);
}

initApplication();

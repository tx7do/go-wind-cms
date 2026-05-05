import type {Pinia} from 'pinia'

import type {App} from 'vue'

import {createPinia} from 'pinia'

let pinia: Pinia

export interface InitStoreOptions {
  namespace: string
}

/**
 * @zh_CN 初始化pinia
 */
export async function setupStore(app: App, options: InitStoreOptions) {
  const {createPersistedState} = await import('pinia-plugin-persistedstate');
  pinia = createPinia();
  const {namespace} = options;
  pinia.use(
    createPersistedState({
      // key $appName-$store.id
      key: storeKey => `${namespace}-${storeKey}`,
      storage: localStorage,
    }),
  );
  app.use(pinia);
  return pinia;
}

export function resetAllStores() {
  if (!pinia) {
    console.error('Pinia is not installed')
    return
  }
  const allStores = (pinia as any)._s
  for (const [_key, store] of allStores)
    store.$reset()
}
